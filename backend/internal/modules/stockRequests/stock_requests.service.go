package stockRequests

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

const maxRequestsPerDay int64 = 10

type StockRequestService interface {
	CreateStockRequest(storeID uint, req types.CreateStockRequestDTO) (uint, string, error)
	GetStockRequests(filter types.GetStockRequestsFilter) ([]types.StockRequestResponse, error)
	GetStockRequestByID(id uint) (types.StockRequestResponse, error)

	// statuses
	RejectStockRequestByStore(requestID uint, dto types.RejectStockRequestStatusDTO) (*data.StockRequest, error)
	RejectStockRequestByWarehouse(requestID uint, dto types.RejectStockRequestStatusDTO) (*data.StockRequest, error)
	SetProcessedStatus(requestID uint) (*data.StockRequest, error)
	SetInDeliveryStatus(requestID uint) (*data.StockRequest, error)
	SetCompletedStatus(requestID uint) (*data.StockRequest, error)
	AcceptStockRequestWithChange(requestID uint, dto types.AcceptWithChangeRequestStatusDTO) (*data.StockRequest, error)

	UpdateStockRequest(requestID uint, items []types.StockRequestStockMaterialDTO) (*data.StockRequest, error)

	DeleteStockRequest(requestID uint) (*data.StockRequest, error)
	GetLastCreatedStockRequest(storeID uint) (*types.StockRequestResponse, error)
	AddStockMaterialToCart(storeID uint, dto types.StockRequestStockMaterialDTO) (*data.StockRequest, error)
}

type stockRequestService struct {
	repo                StockRequestRepository
	stockMaterialRepo   stockMaterial.StockMaterialRepository
	transactionManager  TransactionManager
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewStockRequestService(
	repo StockRequestRepository,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
	transactionManager TransactionManager,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) StockRequestService {
	return &stockRequestService{
		repo:                repo,
		stockMaterialRepo:   stockMaterialRepo,
		transactionManager:  transactionManager,
		notificationService: notificationService,
		logger:              logger,
	}
}

func (s *stockRequestService) CreateStockRequest(storeID uint, req types.CreateStockRequestDTO) (uint, string, error) {
	existingRequest, err := s.repo.GetOpenCartByStoreID(storeID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, "", fmt.Errorf("failed to check for open cart: %w", err)
	}
	if existingRequest != nil {
		return 0, "", types.ErrExistingRequest
	}

	count, err := s.repo.CountFinalizedStockRequestsInLast24Hours(storeID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to count today's stock requests: %w", err)
	}
	if count >= maxRequestsPerDay {
		return 0, "", types.ErrOneRequestPerDay
	}

	store, err := s.repo.GetStoreWarehouse(storeID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to fetch store warehouse for store ID %d: %w", storeID, err)
	}

	stockRequest := &data.StockRequest{
		StoreID:     storeID,
		WarehouseID: store.WarehouseID,
		Status:      data.StockRequestCreated,
	}

	if err := s.repo.CreateStockRequest(stockRequest); err != nil {
		return 0, "", fmt.Errorf("failed to create stock request: %w", err)
	}

	ingredients := []data.StockRequestIngredient{}
	for _, item := range req.StockMaterials {
		var material data.StockMaterial
		if err := s.stockMaterialRepo.PopulateStockMaterial(item.StockMaterialID, &material); err != nil {
			return 0, "", fmt.Errorf("failed to fetch stock material for ID %d: %w", item.StockMaterialID, err)
		}

		ingredients = append(ingredients, data.StockRequestIngredient{
			StockRequestID:  stockRequest.ID,
			StockMaterialID: item.StockMaterialID,
			Quantity:        item.Quantity,
		})
	}

	if err := s.repo.AddIngredientsToStockRequest(ingredients); err != nil {
		return 0, "", fmt.Errorf("failed to add ingredients to stock request ID %d: %w", stockRequest.ID, err)
	}

	requestDetails := &details.NewStockRequestDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           store.WarehouseID,
			FacilityName: store.Name,
		},
		RequesterName: store.Name,
		RequestID:     stockRequest.ID,
	}

	err = s.notificationService.NotifyNewStockRequest(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return stockRequest.ID, store.Name, nil
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

func (s *stockRequestService) RejectStockRequestByStore(requestID uint, dto types.RejectStockRequestStatusDTO) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	if !types.IsValidTransition(request.Status, data.StockRequestRejectedByStore) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", request.Status, data.StockRequestRejectedByStore)
	}

	if err := s.handleRejectedByStoreStatus(request, dto.Comment); err != nil {
		return nil, err
	}

	request.Status = data.StockRequestRejectedByStore
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return nil, fmt.Errorf("failed to update stock request status: %w", err)
	}

	requestDetails := &details.StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           request.WarehouseID,
			FacilityName: request.Warehouse.Name,
		},
		StockRequestID: request.ID,
		RequestStatus:  request.Status,
	}
	err = s.notificationService.NotifyStockRequestStatusUpdated(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return request, nil
}

func (s *stockRequestService) RejectStockRequestByWarehouse(requestID uint, dto types.RejectStockRequestStatusDTO) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	if !types.IsValidTransition(request.Status, data.StockRequestRejectedByWarehouse) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", request.Status, data.StockRequestRejectedByWarehouse)
	}

	if err := s.handleRejectedByWarehouseStatus(request, dto.Comment); err != nil {
		return nil, err
	}

	request.Status = data.StockRequestRejectedByWarehouse
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return nil, fmt.Errorf("failed to update stock request status: %w", err)
	}

	requestDetails := &details.StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           request.StoreID,
			FacilityName: request.Store.Name,
		},
		StockRequestID: request.ID,
		RequestStatus:  request.Status,
	}
	err = s.notificationService.NotifyStockRequestStatusUpdated(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return request, nil
}

func (s *stockRequestService) SetProcessedStatus(requestID uint) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	if !types.IsValidTransition(request.Status, data.StockRequestProcessed) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", request.Status, data.StockRequestProcessed)
	}

	if request.Status == data.StockRequestCreated {
		count, err := s.repo.CountStockRequestsInLast24Hours(request.StoreID)
		if err != nil {
			return nil, fmt.Errorf("failed to count today's stock requests: %w", err)
		}
		if count >= maxRequestsPerDay {
			return nil, types.ErrOneRequestPerDay
		}
	}

	request.Status = data.StockRequestProcessed
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return nil, fmt.Errorf("failed to update stock request status: %w", err)
	}

	requestDetails := &details.StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           request.WarehouseID,
			FacilityName: request.Warehouse.Name,
		},
		StockRequestID: request.ID,
		RequestStatus:  request.Status,
	}

	err = s.notificationService.NotifyStockRequestStatusUpdated(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return request, nil
}

func (s *stockRequestService) SetInDeliveryStatus(requestID uint) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	if !types.IsValidTransition(request.Status, data.StockRequestInDelivery) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", request.Status, data.StockRequestInDelivery)
	}

	if err := s.handleInDeliveryStatus(request); err != nil {
		return nil, err
	}

	request.Status = data.StockRequestInDelivery
	if err := s.repo.UpdateStockRequestStatus(request); err != nil {
		return nil, fmt.Errorf("failed to update stock request status: %w", err)
	}

	requestDetails := &details.StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           request.StoreID,
			FacilityName: request.Store.Name,
		},
		StockRequestID: request.ID,
		RequestStatus:  request.Status,
	}
	err = s.notificationService.NotifyStockRequestStatusUpdated(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return request, nil
}

func (s *stockRequestService) SetCompletedStatus(requestID uint) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	if !types.IsValidTransition(request.Status, data.StockRequestCompleted) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", request.Status, data.StockRequestCompleted)
	}

	err = s.transactionManager.HandleCompleteStockRequest(request)
	if err != nil {
		return nil, err
	}

	requestDetails := &details.StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           request.WarehouseID,
			FacilityName: request.Warehouse.Name,
		},
		StockRequestID: request.ID,
		RequestStatus:  request.Status,
	}
	err = s.notificationService.NotifyStockRequestStatusUpdated(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return request, nil
}

func (s *stockRequestService) AcceptStockRequestWithChange(requestID uint, dto types.AcceptWithChangeRequestStatusDTO) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	store, err := s.repo.GetStoreWarehouse(request.StoreID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store warehouse: %w", err)
	}

	if !types.IsValidTransition(request.Status, data.StockRequestAcceptedWithChange) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", request.Status, data.StockRequestAcceptedWithChange)
	}

	if err := s.transactionManager.HandleAcceptedWithChange(request, store.ID, dto.Items, dto.Comment); err != nil {
		return nil, err
	}

	requestDetails := &details.StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           request.WarehouseID,
			FacilityName: request.Warehouse.Name,
		},
		StockRequestID: request.ID,
		RequestStatus:  request.Status,
	}
	err = s.notificationService.NotifyStockRequestStatusUpdated(requestDetails)
	if err != nil {
		s.logger.Errorf("failed to send notification: %v", err)
	}

	return request, nil
}

func (s *stockRequestService) handleInDeliveryStatus(request *data.StockRequest) error {
	for _, ingredient := range request.Ingredients {
		stockQuantity, err := s.repo.GetWarehouseStockQuantity(request.WarehouseID, ingredient.StockMaterialID)
		if err != nil {
			return fmt.Errorf("failed to fetch warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
		}

		requestDetails := &details.OutOfStockDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           request.WarehouseID,
				FacilityName: request.Warehouse.Name,
			},
			ItemName: ingredient.StockMaterial.Name,
		}

		if stockQuantity < ingredient.Quantity {
			err := s.notificationService.NotifyOutOfStock(requestDetails)
			if err != nil {
				return fmt.Errorf("failed to send out of stock notification: %w", err)
			}

			return types.ErrInsufficientStock
		}

		updatedStock, err := s.repo.DeductWarehouseStock(ingredient.StockMaterialID, request.WarehouseID, ingredient.Quantity)
		if err != nil {
			return fmt.Errorf("failed to deduct warehouse stock for stock material ID %d: %w", ingredient.StockMaterialID, err)
		}

		if updatedStock.Quantity < ingredient.StockMaterial.SafetyStock {
			err := s.notificationService.NotifyOutOfStock(requestDetails)
			if err != nil {
				s.logger.Errorf("failed to send out of stock notification: %v", err)
			}
		}
	}
	return nil
}

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

	for _, ingredient := range request.Ingredients {
		_, err := s.repo.ReturnWarehouseStock(ingredient.StockMaterialID, request.WarehouseID, ingredient.Quantity)
		if err != nil {
			return fmt.Errorf("failed to return stock for material ID %d: %w", ingredient.StockMaterialID, err)
		}
	}

	fmt.Printf("Stock request rejected, ID: %d, Comment: %s\n", request.ID, utils.StringOrEmpty(comment))
	return nil
}

func (s *stockRequestService) UpdateStockRequest(requestID uint, items []types.StockRequestStockMaterialDTO) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	ingredients := []data.StockRequestIngredient{}
	for _, item := range items {
		var material data.StockMaterial
		if err := s.stockMaterialRepo.PopulateStockMaterial(item.StockMaterialID, &material); err != nil {
			return nil, fmt.Errorf("failed to fetch stock material for ID %d: %w", item.StockMaterialID, err)
		}
		ingredients = append(ingredients, data.StockRequestIngredient{
			StockRequestID:  request.ID,
			StockMaterialID: item.StockMaterialID,
			Quantity:        item.Quantity,
		})
	}

	if err := s.repo.ReplaceStockRequestIngredients(*request, ingredients); err != nil {
		return nil, fmt.Errorf("failed to replace ingredients for stock request ID %d: %w", requestID, err)
	}
	return request, nil
}

func (s *stockRequestService) DeleteStockRequest(requestID uint) (*data.StockRequest, error) {
	request, err := s.repo.GetStockRequestByID(requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock request not found")
		}
		return nil, fmt.Errorf("failed to fetch stock request: %w", err)
	}

	if request.Status != data.StockRequestCreated {
		return nil, fmt.Errorf("only stock requests in 'CREATED' status can be deleted")
	}

	if err := s.repo.DeleteStockRequest(requestID); err != nil {
		return nil, fmt.Errorf("failed to delete stock request: %w", err)
	}

	return request, nil
}

func (s *stockRequestService) GetLastCreatedStockRequest(storeID uint) (*types.StockRequestResponse, error) {
	request, err := s.repo.GetOpenCartByStoreID(storeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no open cart exists for store ID %d", storeID)
		}
		return nil, fmt.Errorf("failed to fetch last created stock request: %w", err)
	}

	response := types.ToStockRequestResponse(request)
	return &response, nil
}

func (s *stockRequestService) AddStockMaterialToCart(storeID uint, dto types.StockRequestStockMaterialDTO) (*data.StockRequest, error) {
	cart, err := s.repo.GetOpenCartByStoreID(storeID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to fetch open cart: %w", err)
	}

	if cart == nil {
		cart, err = s.createNewCart(storeID)
		if err != nil {
			return nil, fmt.Errorf("failed to create new cart: %w", err)
		}
	}

	for _, ingredient := range cart.Ingredients {
		if ingredient.StockMaterialID == dto.StockMaterialID {
			newQuantity := ingredient.Quantity + dto.Quantity
			err := s.repo.UpdateStockRequestIngredientQuantity(ingredient.ID, newQuantity)
			if err != nil {
				return nil, fmt.Errorf("failed to update stock material quantity in the cart: %w", err)
			}
			return cart, nil
		}
	}

	material, err := s.stockMaterialRepo.GetStockMaterialByID(dto.StockMaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock material for cart")
	}

	newIngredient := data.StockRequestIngredient{
		StockRequestID:  cart.ID,
		StockMaterialID: material.ID,
		Quantity:        dto.Quantity,
	}
	err = s.repo.AddIngredientsToStockRequest([]data.StockRequestIngredient{newIngredient})
	if err != nil {
		return nil, fmt.Errorf("failed to add stock material to cart: %w", err)
	}

	return cart, nil
}

func (s *stockRequestService) createNewCart(storeID uint) (*data.StockRequest, error) {
	store, err := s.repo.GetStoreWarehouse(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch corresponding warehouse, %s", err.Error())
	}
	newCart := &data.StockRequest{
		StoreID:     storeID,
		Store:       *store,
		WarehouseID: store.WarehouseID,
		Status:      data.StockRequestCreated,
	}

	err = s.repo.CreateStockRequest(newCart)
	if err != nil {
		return nil, fmt.Errorf("failed to create new cart: %w", err)
	}

	return newCart, nil
}
