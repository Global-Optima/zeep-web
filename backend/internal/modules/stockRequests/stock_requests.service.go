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
	GetAllStockMaterials(storeID uint, filter types.StockMaterialFilter) ([]types.StockMaterialDTO, error)
	UpdateStockRequestIngredients(requestID uint, items []types.StockRequestItemDTO) error
	GetAvailableStockMaterialsByIngredient(ingredientID uint, warehouseID *uint) ([]types.StockMaterialAvailabilityDTO, error)
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
		var stockMaterial data.StockMaterial
		if err := s.repo.GetStockMaterialByID(item.StockMaterialID, &stockMaterial); err != nil {
			return 0, fmt.Errorf("failed to fetch stock material for ID %d: %w", item.StockMaterialID, err)
		}

		ingredients = append(ingredients, data.StockRequestIngredient{
			StockRequestID:  stockRequest.ID,
			IngredientID:    stockMaterial.IngredientID,
			StockMaterialID: item.StockMaterialID,
			Quantity:        item.Quantity,
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

	switch status.Status {
	case data.StockRequestInDelivery:
		for _, ingredient := range request.Ingredients {
			if err := s.repo.DeductWarehouseStock(ingredient.StockMaterialID, request.WarehouseID, ingredient.Quantity); err != nil {
				return fmt.Errorf("failed to deduct warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
			}
		}

	case data.StockRequestCompleted:
		for _, ingredient := range request.Ingredients {
			dates := types.UpdateIngredientDates{
				DeliveredDate:  time.Now(),
				ExpirationDate: time.Now().AddDate(0, 0, ingredient.StockMaterial.ExpirationPeriodInDays),
			}

			if err := s.repo.UpdateStockRequestIngredientDates(ingredient.ID, &dates); err != nil {
				return fmt.Errorf("failed to update ingredient dates for stock material ID %d: %w", ingredient.StockMaterialID, err)
			}

			if err := s.repo.AddToStoreWarehouseStock(storeWarehouse.ID, ingredient.StockMaterialID, ingredient.Quantity); err != nil {
				return fmt.Errorf("failed to update store warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
			}
		}

	case data.StockRequestRejected:
		fmt.Printf("Stock request rejected, ID: %d\n", requestID)
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

func (s *stockRequestService) GetAllStockMaterials(storeID uint, filter types.StockMaterialFilter) ([]types.StockMaterialDTO, error) {
	return s.repo.GetAllStockMaterials(storeID, filter)
}

func (s *stockRequestService) UpdateStockRequestIngredients(requestID uint, items []types.StockRequestItemDTO) error {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock request: %w", err)
	}

	ingredients := []data.StockRequestIngredient{}
	for _, item := range items {
		var stockMaterial data.StockMaterial
		if err := s.repo.GetStockMaterialByID(item.StockMaterialID, &stockMaterial); err != nil {
			return fmt.Errorf("failed to fetch stock material for ID %d: %w", item.StockMaterialID, err)
		}

		ingredients = append(ingredients, data.StockRequestIngredient{
			StockRequestID:  request.ID,
			IngredientID:    stockMaterial.IngredientID,
			StockMaterialID: item.StockMaterialID,
			Quantity:        item.Quantity,
		})
	}

	if err := s.repo.ReplaceStockRequestIngredients(request.ID, ingredients); err != nil {
		return fmt.Errorf("failed to replace ingredients for stock request ID %d: %w", requestID, err)
	}

	return nil
}

func (s *stockRequestService) GetAvailableStockMaterialsByIngredient(ingredientID uint, warehouseID *uint) ([]types.StockMaterialAvailabilityDTO, error) {
	stocks, err := s.repo.GetStockMaterialsByIngredient(ingredientID, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock materials: %w", err)
	}

	availability := make([]types.StockMaterialAvailabilityDTO, len(stocks))
	for i, stock := range stocks {
		availability[i] = types.StockMaterialAvailabilityDTO{
			StockMaterialID: stock.StockMaterialID,
			Name:            stock.StockMaterial.Name,
			Category:        stock.StockMaterial.Category,
			AvailableQty:    stock.Quantity,
			WarehouseID:     stock.WarehouseID,
			WarehouseName:   stock.Warehouse.Name,
			Unit:            stock.StockMaterial.Unit.Name,
		}
	}

	return availability, nil
}
