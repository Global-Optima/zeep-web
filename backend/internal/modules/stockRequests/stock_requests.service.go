package stockRequests

import (
	"fmt"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StockRequestService interface {
	CreateStockRequest(req types.CreateStockRequestDTO) (uint, error)
	GetStockRequests(filter types.GetStockRequestsFilter) ([]types.StockRequestResponse, error)
	GetStockRequestByID(id uint) (types.StockRequestResponse, error)

	UpdateStockRequestStatus(requestID uint, status string) error
	AcceptStockRequestWithChange(requestID uint, status string, dto types.AcceptWithChangeRequestStatusDTO) error
	RejectStockRequest(requestID uint, status string, dto types.RejectStockRequestStatusDTO) error

	GetLowStockIngredients(storeID uint) ([]types.LowStockIngredientResponse, error)
	GetAllStockMaterials(storeID uint, filter types.StockMaterialFilter) ([]types.StockMaterialDTO, error)
	UpdateStockRequestIngredients(requestID uint, items []types.StockRequestStockMaterialDTO) error
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
	for _, item := range req.StockMaterials {
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

func (s *stockRequestService) GetStockRequests(filter types.GetStockRequestsFilter) ([]types.StockRequestResponse, error) {
	requests, err := s.repo.GetStockRequests(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock requests: %w", err)
	}

	responses := []types.StockRequestResponse{}
	for _, request := range requests {
		responses = append(responses, types.ToStockRequestResponse(&request))
	}
	return responses, nil
}

func (s *stockRequestService) GetStockRequestByID(id uint) (types.StockRequestResponse, error) {
	request, err := s.repo.GetStockRequestByID(id)
	if err != nil {
		return types.StockRequestResponse{}, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	return types.ToStockRequestResponse(request), nil
}

func (s *stockRequestService) RejectStockRequest(requestID uint, status string, dto types.RejectStockRequestStatusDTO) error {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock request: %w", err)
	}

	requestStatus := data.StockRequestStatus(status)
	switch requestStatus {
	case data.StockRequestRejectedByWarehouse:
		if err := s.handleRejectedByWarehouseStatus(request, dto.Comment); err != nil {
			return err
		}

	case data.StockRequestRejectedByStore:
		if err := s.handleRejectedByStoreStatus(request, dto.Comment); err != nil {
			return err
		}
	}

	request.Status = requestStatus
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return fmt.Errorf("failed to update stock request status: %w", err)
	}

	return nil
}

func (s *stockRequestService) UpdateStockRequestStatus(requestID uint, status string) error {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock request: %w", err)
	}

	storeWarehouse, err := s.repo.GetStoreWarehouse(request.StoreID)
	if err != nil {
		return fmt.Errorf("failed to fetch store warehouse: %w", err)
	}

	requestStatus := data.StockRequestStatus(status)

	switch requestStatus {
	case data.StockRequestInDelivery:
		if err := s.handleInDeliveryStatus(request); err != nil {
			return err
		}

	case data.StockRequestCompleted:
		if err := s.handleCompletedStatus(request, storeWarehouse.ID); err != nil {
			return err
		}
	}

	request.Status = requestStatus
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return fmt.Errorf("failed to update stock request status: %w", err)
	}

	return nil
}

func (s *stockRequestService) AcceptStockRequestWithChange(requestID uint, status string, dto types.AcceptWithChangeRequestStatusDTO) error {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch stock request: %w", err)
	}

	storeWarehouse, err := s.repo.GetStoreWarehouse(request.StoreID)
	if err != nil {
		return fmt.Errorf("failed to fetch store warehouse: %w", err)
	}

	if err := s.handleAcceptedWithChange(request, storeWarehouse.ID, dto.Items, dto.Comment); err != nil {
		return err
	}

	request.Status = data.StockRequestStatus(status)
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return fmt.Errorf("failed to update stock request status: %w", err)
	}

	return nil
}

func (s *stockRequestService) handleInDeliveryStatus(request *data.StockRequest) error {
	for _, ingredient := range request.Ingredients {
		if err := s.repo.DeductWarehouseStock(ingredient.StockMaterialID, request.WarehouseID, ingredient.Quantity); err != nil {
			return fmt.Errorf("failed to deduct warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
		}
	}
	return nil
}

func (s *stockRequestService) handleCompletedStatus(request *data.StockRequest, storeWarehouseID uint) error {
	for _, ingredient := range request.Ingredients {
		if ingredient.StockMaterial.Package == nil {
			return utils.WrapError("package is not presenent for stock material", fmt.Errorf("stock material ID %d", ingredient.StockMaterialID))
		}

		dates := types.UpdateIngredientDates{
			DeliveredDate:  time.Now(),
			ExpirationDate: time.Now().AddDate(0, 0, ingredient.StockMaterial.ExpirationPeriodInDays),
		}

		if err := s.repo.UpdateStockRequestIngredientDates(ingredient.ID, &dates); err != nil {
			return fmt.Errorf("failed to update ingredient dates for stock material ID %d: %w", ingredient.StockMaterialID, err)
		}

		if err := s.repo.AddToStoreWarehouseStock(storeWarehouseID, ingredient.StockMaterialID, ingredient.Quantity); err != nil {
			return fmt.Errorf("failed to update store warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
		}
	}
	return nil
}

func (s *stockRequestService) handleAcceptedWithChange(request *data.StockRequest, storeWarehouseID uint, items []types.StockRequestStockMaterialDTO, comment *string) error {
	var autoGeneratedComments []string
	updatedIngredients := []data.StockRequestIngredient{}

	for _, item := range items {
		var stockMaterial data.StockMaterial
		if err := s.repo.GetStockMaterialByID(item.StockMaterialID, &stockMaterial); err != nil {
			return fmt.Errorf("failed to fetch stock material for ID %d: %w", item.StockMaterialID, err)
		}

		// Find the original ingredient in the request
		originalIngredient := findOriginalIngredient(request.Ingredients, item.StockMaterialID)

		if originalIngredient != nil {
			if originalIngredient.Quantity != item.Quantity {
				autoGeneratedComments = append(autoGeneratedComments, fmt.Sprintf(
					"Material '%s' (ID %d): Expected %.2f, Received %.2f",
					stockMaterial.Name,
					originalIngredient.StockMaterialID,
					originalIngredient.Quantity,
					item.Quantity,
				))
			}
		} else {
			autoGeneratedComments = append(autoGeneratedComments, fmt.Sprintf(
				"Material '%s' (ID %d): Unexpected material received with quantity %.2f",
				stockMaterial.Name,
				item.StockMaterialID,
				item.Quantity,
			))
		}

		if originalIngredient != nil && originalIngredient.Quantity > 0 {
			if err := s.repo.DeductWarehouseStock(originalIngredient.StockMaterialID, request.WarehouseID, originalIngredient.Quantity); err != nil {
				return fmt.Errorf("failed to deduct warehouse stock for stock material ID %d: %w", originalIngredient.StockMaterialID, err)
			}
		}

		if item.Quantity > 0 {
			if err := s.repo.AddToStoreWarehouseStock(storeWarehouseID, item.StockMaterialID, item.Quantity); err != nil {
				return fmt.Errorf("failed to add stock to store warehouse for stock material ID %d: %w", item.StockMaterialID, err)
			}
		}

		updatedIngredients = append(updatedIngredients, data.StockRequestIngredient{
			StockRequestID:  request.ID,
			IngredientID:    stockMaterial.IngredientID,
			StockMaterialID: item.StockMaterialID,
			Quantity:        item.Quantity,
		})
	}

	if err := s.repo.ReplaceStockRequestIngredients(request.ID, updatedIngredients); err != nil {
		return fmt.Errorf("failed to replace ingredients for stock request ID %d: %w", request.ID, err)
	}

	if comment != nil {
		storeComment := fmt.Sprintf("Comment from store: %s", *comment)
		autoGeneratedComments = append(autoGeneratedComments, storeComment)
	}

	// Save auto-generated comments
	if len(autoGeneratedComments) > 0 {
		combinedComments := strings.Join(autoGeneratedComments, "\n")
		if err := s.repo.AddStoreComment(request.ID, combinedComments); err != nil {
			return fmt.Errorf("failed to add store comment for request ID %d: %w", request.ID, err)
		}
	}

	return nil
}

// Helper function to find the original ingredient in the stock request by StockMaterialID
func findOriginalIngredient(ingredients []data.StockRequestIngredient, stockMaterialID uint) *data.StockRequestIngredient {
	for _, ingredient := range ingredients {
		if ingredient.StockMaterialID == stockMaterialID {
			return &ingredient
		}
	}
	return nil
}

func (s *stockRequestService) handleRejectedByWarehouseStatus(request *data.StockRequest, comment *string) error {
	if comment != nil {
		if err := s.repo.AddWarehouseComment(request.ID, *comment); err != nil {
			return fmt.Errorf("failed to add rejection comment for request ID %d: %w", request.ID, err)
		}
	}
	fmt.Printf("Stock request rejected, ID: %d, Comment: %s\n", request.ID, utils.StringOrEmpty(comment))
	return nil
}

func (s *stockRequestService) handleRejectedByStoreStatus(request *data.StockRequest, comment *string) error {
	if comment != nil {
		if err := s.repo.AddStoreComment(request.ID, *comment); err != nil {
			return fmt.Errorf("failed to add rejection comment for request ID %d: %w", request.ID, err)
		}
	}
	fmt.Printf("Stock request rejected, ID: %d, Comment: %s\n", request.ID, utils.StringOrEmpty(comment))
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

func (s *stockRequestService) UpdateStockRequestIngredients(requestID uint, items []types.StockRequestStockMaterialDTO) error {
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
			StockMaterialID:   stock.StockMaterialID,
			Name:              stock.StockMaterial.Name,
			Category:          stock.StockMaterial.StockMaterialCategory.Name,
			AvailableQuantity: stock.Quantity,
			Unit:              stock.StockMaterial.Unit.Name,
			Warehouse: types.WarehouseDTO{
				ID:   stock.WarehouseID,
				Name: stock.Warehouse.Name,
			},
		}
	}

	return availability, nil
}
