package stockRequests

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
)

type StockRequestService interface {
	CreateStockRequest(req types.CreateStockRequestDTO) (uint, error)
	GetStockRequests(filter types.StockRequestFilter) ([]types.StockRequestResponse, error)
	UpdateStockRequestStatus(requestID uint, status types.UpdateStockRequestStatusDTO) error
	GetLowStockIngredients(storeID uint) ([]types.LowStockIngredientResponse, error)
	GetMarketplaceProducts(storeID uint, filter types.MarketplaceFilter) ([]types.ProductMarketplaceDTO, error)
	AddStockRequestIngredient(requestID uint, item types.StockRequestItemDTO) error
	DeleteStockRequestIngredient(ingredientID uint) error
}

type stockRequestService struct {
	repo StockRequestRepository
}

func NewStockRequestService(repo StockRequestRepository) StockRequestService {
	return &stockRequestService{repo: repo}
}

func (s *stockRequestService) CreateStockRequest(req types.CreateStockRequestDTO) (uint, error) {
	lastRequestDate, err := s.repo.GetLastStockRequestDate(req.StoreID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch last stock request date: %w", err)
	}
	err = types.ValidateStockRequestRate(lastRequestDate)
	if err != nil {
		return 0, err
	}

	storeWarehouse, err := s.repo.GetStoreWarehouse(req.StoreID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch store warehouse for store ID %d: %w", req.StoreID, err)
	}

	stockRequest := &data.StockRequest{
		StoreID:     req.StoreID,
		WarehouseID: storeWarehouse.WarehouseID,
		Status:      data.StockRequestCreated,
	}

	if err := s.repo.CreateStockRequest(stockRequest); err != nil {
		return 0, fmt.Errorf("failed to create stock request: %w", err)
	}

	ingredients := []data.StockRequestIngredient{}
	for _, item := range req.Items {
		var mapping data.IngredientStockMaterialMapping
		if err := s.repo.GetMappingByIngredientID(item.StockMaterialID, &mapping); err != nil {
			return 0, fmt.Errorf("failed to map ingredient for stock material ID %d: %w", item.StockMaterialID, err)
		}

		ingredients = append(ingredients, data.StockRequestIngredient{
			StockRequestID: stockRequest.ID,
			IngredientID:   mapping.IngredientID,
			Quantity:       item.Quantity,
		})
	}

	if err := s.repo.AddIngredientsToStockRequest(ingredients); err != nil {
		return 0, fmt.Errorf("failed to add ingredients to stock request ID %d: %w", stockRequest.ID, err)
	}

	return stockRequest.ID, nil
}

func (s *stockRequestService) GetStockRequests(filter types.StockRequestFilter) ([]types.StockRequestResponse, error) {
	requests, err := s.repo.GetStockRequests(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock requests: %w", err)
	}

	responses := []types.StockRequestResponse{}
	for _, request := range requests {
		responses = append(responses, types.ToStockRequestResponse(request))
	}
	return responses, nil
}

func (s *stockRequestService) UpdateStockRequestStatus(requestID uint, status types.UpdateStockRequestStatusDTO) error {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock request: %w", err)
	}

	storeWarehouse, err := s.repo.GetStoreWarehouse(request.StoreID)
	if err != nil {
		return fmt.Errorf("failed to fetch store warehouse: %w", err)
	}

	if status.Status == data.StockRequestInDelivery {
		for _, ingredient := range request.Ingredients {
			var mapping data.IngredientStockMaterialMapping
			if err := s.repo.GetMappingByIngredientID(ingredient.IngredientID, &mapping); err != nil {
				return fmt.Errorf("failed to map ingredient for stock material ID %d: %w", ingredient.IngredientID, err)
			}

			if err := s.repo.DeductWarehouseStock(mapping.StockMaterialID, request.WarehouseID, ingredient.Quantity); err != nil {
				return fmt.Errorf("failed to deduct warehouse stock for stock material ID %d: %w", mapping.StockMaterialID, err)
			}
		}
	}

	if status.Status == data.StockRequestCompleted {
		for _, ingredient := range request.Ingredients {
			var mapping data.IngredientStockMaterialMapping
			if err := s.repo.GetMappingByIngredientID(ingredient.IngredientID, &mapping); err != nil {
				return fmt.Errorf("failed to map ingredient for stock material ID %d: %w", ingredient.IngredientID, err)
			}

			dates := types.UpdateIngredientDates{
				DeliveredDate:  time.Now(),
				ExpirationDate: time.Now().AddDate(0, 0, mapping.StockMaterial.ExpirationPeriodInDays),
			}

			if err := s.repo.UpdateStockRequestIngredientDates(&dates); err != nil {
				return fmt.Errorf("failed to update ingredient dates for ingredient ID %d: %w", ingredient.IngredientID, err)
			}

			if err := s.repo.AddToStoreWarehouseStock(storeWarehouse.ID, ingredient.IngredientID, ingredient.Quantity); err != nil {
				return fmt.Errorf("failed to update store warehouse stock for ingredient ID %d: %w", ingredient.IngredientID, err)
			}
		}
	}

	request.Status = status.Status
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return fmt.Errorf("failed to update stock request status: %w", err)
	}

	return nil
}

func (s *stockRequestService) GetLowStockIngredients(storeID uint) ([]types.LowStockIngredientResponse, error) {
	storeWarehouse, err := s.repo.GetStoreWarehouse(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store warehouse for store ID %d: %w", storeID, err)
	}

	lowStockItems, err := s.repo.GetLowStockIngredients(storeWarehouse.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch low stock ingredients: %w", err)
	}

	responses := []types.LowStockIngredientResponse{}
	for _, stock := range lowStockItems {
		responses = append(responses, types.ToLowStockIngredientResponse(stock))
	}
	return responses, nil
}

func (s *stockRequestService) GetMarketplaceProducts(storeID uint, filter types.MarketplaceFilter) ([]types.ProductMarketplaceDTO, error) {
	return s.repo.GetMarketplaceStockMaterials(storeID, filter)
}

func (s *stockRequestService) AddStockRequestIngredient(requestID uint, item types.StockRequestItemDTO) error {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock request: %w", err)
	}

	var mapping data.IngredientStockMaterialMapping
	if err := s.repo.GetMappingByIngredientID(item.StockMaterialID, &mapping); err != nil {
		return fmt.Errorf("failed to map ingredient for stock material ID %d: %w", item.StockMaterialID, err)
	}

	ingredient := data.StockRequestIngredient{
		StockRequestID: request.ID,
		IngredientID:   mapping.IngredientID,
		Quantity:       item.Quantity,
	}

	if err := s.repo.AddIngredientsToStockRequest([]data.StockRequestIngredient{ingredient}); err != nil {
		return fmt.Errorf("failed to add ingredient to stock request ID %d: %w", requestID, err)
	}

	return nil
}

func (s *stockRequestService) DeleteStockRequestIngredient(ingredientID uint) error {
	if err := s.repo.DeleteStockRequestIngredient(ingredientID); err != nil {
		return fmt.Errorf("failed to delete stock request ingredient ID %d: %w", ingredientID, err)
	}
	return nil
}
