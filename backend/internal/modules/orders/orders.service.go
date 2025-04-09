package orders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/taskqueue"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"go.uber.org/zap"
)

const (
	OrderPaymentFailure = "order-payment-failure"
)

type OrderService interface {
	GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error)
	GetAllBaristaOrders(filter types.OrdersTimeZoneFilter) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.SuborderDTO, error)
	CreateOrder(storeId uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error)
	GetOrderBySubOrder(subOrderID uint) (*data.Order, error)
	GetOrderById(orderId uint) (types.OrderDTO, error)

	GetOrderDetails(orderID uint, filter *contexts.StoreContextFilter) (*types.OrderDetailsDTO, error)
	ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error)

	SetNextSubOrderStatus(subOrderID uint, options *types.ToggleNextSuborderStatusOptions) (*types.SuborderDTO, error)

	SuccessOrderPayment(orderID uint, dto *types.TransactionDTO) error
	FailOrderPayment(orderID uint) error
}

type orderValidationResults struct {
	ProductPrices  map[uint]float64
	ProductNames   map[uint]string
	AdditivePrices map[uint]float64
	AdditiveNames  map[uint]string
}

type orderService struct {
	taskQueue           taskqueue.TaskQueue
	orderRepo           OrderRepository
	storeProductRepo    storeProducts.StoreProductRepository
	storeAdditiveRepo   storeAdditives.StoreAdditiveRepository
	storeStockRepo      storeStocks.StoreStockRepository
	notificationService notifications.NotificationService
	transactionManager  TransactionManager
	logger              *zap.SugaredLogger
}

func NewOrderService(
	taskQueue taskqueue.TaskQueue,
	orderRepo OrderRepository,
	storeProductRepo storeProducts.StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	notificationService notifications.NotificationService,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) OrderService {
	return &orderService{
		taskQueue:           taskQueue,
		orderRepo:           orderRepo,
		storeProductRepo:    storeProductRepo,
		storeAdditiveRepo:   storeAdditiveRepo,
		storeStockRepo:      storeStockRepo,
		notificationService: notificationService,
		transactionManager:  transactionManager,
		logger:              logger,
	}
}

func (s *orderService) GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetOrders(filter)
	if err != nil {
		return nil, err
	}

	orderDTOs := make([]types.OrderDTO, 0)
	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}

	return orderDTOs, nil
}

func (s *orderService) GetAllBaristaOrders(filter types.OrdersTimeZoneFilter) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetAllBaristaOrders(filter)
	if err != nil {
		wrappedErr := fmt.Errorf("error getting barista orders: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	orderDTOs := make([]types.OrderDTO, 0)

	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}

	return orderDTOs, nil
}

func (s *orderService) GetSubOrders(orderID uint) ([]types.SuborderDTO, error) {
	suborders, err := s.orderRepo.GetSubOrdersByOrderID(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("error getting suborders by orderID: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	var subOrderDTOs []types.SuborderDTO
	for _, suborder := range suborders {
		subOrderDTOs = append(subOrderDTOs, types.ConvertSuborderToDTO(&suborder))
	}
	return subOrderDTOs, nil
}

func ExpandSuborders(suborders []types.CreateSubOrderDTO) []types.CreateSubOrderDTO {
	var expanded []types.CreateSubOrderDTO
	for _, s := range suborders {
		// Repeat each suborder 'quantity' times, but each with quantity=1
		for range s.Quantity {
			expanded = append(expanded, types.CreateSubOrderDTO{
				StoreProductSizeID: s.StoreProductSizeID,
				Quantity:           1,
				StoreAdditivesIDs:  s.StoreAdditivesIDs,
			})
		}
	}
	return expanded
}

func (s *orderService) CreateOrder(storeID uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error) {
	censorValidator := censor.GetCensorValidator()

	if err := censorValidator.ValidateText(createOrderDTO.CustomerName); err != nil {
		s.logger.Error(err)
		return nil, types.ErrInvalidCustomerNameCensor
	}

	if len(createOrderDTO.Suborders) == 0 {
		return nil, fmt.Errorf("order can not be empty")
	}

	if err := ValidateMultipleSelect(storeID, *createOrderDTO, s.storeAdditiveRepo); err != nil {
		s.logger.Error(err)
		return nil, types.ErrMultipleSelect // TODO: test updates
	}

	frozenMap, err := s.orderRepo.CalculateFrozenStock(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate frozen stock: %w", err)
	}

	if err := s.CheckAndAccumulateSuborders(storeID, createOrderDTO.Suborders, frozenMap); err != nil {
		return nil, err
	}

	storeProductSizeIDs, storeAdditiveIDs := RetrieveIDs(*createOrderDTO)
	validations, err := s.StockAndPriceValidationResults(
		storeID,
		storeProductSizeIDs,
		storeAdditiveIDs,
		frozenMap,
	)
	if err != nil {
		s.logger.Error("validation failed: %w", err)
		return nil, err
	}

	createOrderDTO.StoreID = storeID
	order, total := types.ConvertCreateOrderDTOToOrder(
		createOrderDTO,
		validations.ProductPrices,
		validations.AdditivePrices,
	)
	order.Status = data.OrderStatusWaitingForPayment
	order.Total = total

	id, err := s.orderRepo.CreateOrder(&order)
	if err != nil {
		wrappedErr := fmt.Errorf("error creating order: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	payload, err := json.Marshal(types.WaitingOrderPayload{OrderID: id})
	if err != nil {
		return &order, err
	}

	err = s.taskQueue.EnqueueTask(OrderPaymentFailure, payload, config.GetConfig().Payment.WaitingTime)
	if err != nil {
		return &order, err
	}

	return &order, nil
}

func (s *orderService) StockAndPriceValidationResults(
	storeID uint,
	storeProductSizeIDs, storeAdditiveIDs []uint,
	frozenMap map[uint]float64,
) (*orderValidationResults, error) {
	productPrices, productNames, err := ValidateStoreProductSizes(storeID, storeProductSizeIDs, s.storeProductRepo, frozenMap)
	if err != nil {
		s.logger.Error(fmt.Errorf("product validation failed: %w", err))
		return nil, err
	}

	additivePrices, additiveNames, err := ValidateStoreAdditives(storeID, storeAdditiveIDs, s.storeAdditiveRepo, frozenMap)
	if err != nil {
		s.logger.Error(fmt.Errorf("additive validation failed: %w", err))
		return nil, err
	}

	return &orderValidationResults{
		ProductPrices:  productPrices,
		ProductNames:   productNames,
		AdditivePrices: additivePrices,
		AdditiveNames:  additiveNames,
	}, nil
}

func ValidateStoreAdditives(
	storeID uint,
	storeAdditiveIDs []uint,
	repo storeAdditives.StoreAdditiveRepository,
	frozenMap map[uint]float64,
) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	additiveNames := make(map[uint]string)

	for _, addID := range storeAdditiveIDs {
		storeAdd, err := repo.GetStoreAdditiveWithDetailsByID(addID, &contexts.StoreContextFilter{StoreID: &storeID})
		if err != nil {
			return nil, nil, fmt.Errorf("error with store additive: %w", err)
		}
		if storeAdd == nil {
			return nil, nil, fmt.Errorf("store additive with ID %d is nil", addID)
		}
		if storeAdd.Additive.Name == "" {
			return nil, nil, fmt.Errorf("store additive with ID %d has an empty name", addID)
		}

		// Get the effective price (store-specific price overrides the base price if available).
		price := storeAdd.Additive.BasePrice
		if storeAdd.StorePrice != nil {
			price = *storeAdd.StorePrice
		}
		prices[addID] = price
		additiveNames[addID] = storeAdd.Additive.Name
	}

	return prices, additiveNames, nil
}

func ValidateMultipleSelect(storeID uint, createOrderDTO types.CreateOrderDTO, repo storeAdditives.StoreAdditiveRepository) error {
	for _, suborder := range createOrderDTO.Suborders {
		categoryCount := make(map[uint]int)
		for _, addID := range suborder.StoreAdditivesIDs {
			storeAdd, err := repo.GetStoreAdditiveWithDetailsByID(addID, &contexts.StoreContextFilter{StoreID: &storeID})
			if err != nil {
				return err
			}
			categoryID := storeAdd.Additive.AdditiveCategoryID
			categoryCount[categoryID]++
			if !storeAdd.Additive.Category.IsMultipleSelect && categoryCount[categoryID] > 1 {
				return types.ErrMultipleSelect
			}
		}
	}

	return nil
}

func ValidateStoreProductSizes(
	storeID uint,
	storeProductSizeIDs []uint,
	repo storeProducts.StoreProductRepository,
	frozenMap map[uint]float64,
) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	productNames := make(map[uint]string)

	for _, psID := range storeProductSizeIDs {
		storePS, err := repo.GetStoreProductSizeWithDetailsByID(storeID, psID)
		if err != nil {
			return nil, nil, fmt.Errorf("error with store product size: %w", err)
		}
		if storePS == nil {
			return nil, nil, fmt.Errorf("store product size with ID %d is nil", psID)
		}
		if storePS.ProductSize.Product.Name == "" {
			return nil, nil, fmt.Errorf("product size with ID %d has an associated product with an empty name", psID)
		}

		price := storePS.ProductSize.BasePrice
		if storePS.StorePrice != nil {
			price = *storePS.StorePrice
		}

		prices[psID] = price
		productNames[psID] = storePS.StoreProduct.Product.Name
	}

	return prices, productNames, nil
}

func RetrieveIDs(createOrderDTO types.CreateOrderDTO) ([]uint, []uint) {
	var storeProductSizeIDs []uint
	var storeAdditiveIDs []uint
	for _, product := range createOrderDTO.Suborders {
		storeProductSizeIDs = append(storeProductSizeIDs, product.StoreProductSizeID)
		storeAdditiveIDs = append(storeAdditiveIDs, product.StoreAdditivesIDs...)
	}
	return storeProductSizeIDs, storeAdditiveIDs
}

func (s *orderService) CheckAndAccumulateSuborders(
	storeID uint,
	suborders []types.CreateSubOrderDTO,
	frozenMap map[uint]float64,
) error {
	// Expand any suborders with quantity > 1 into separate single-quantity items
	expanded := ExpandSuborders(suborders)

	// For each single suborder, call existing checks
	for _, sub := range expanded {
		// Product Size
		sps, err := s.storeProductRepo.GetSufficientStoreProductSizeById(storeID, sub.StoreProductSizeID, frozenMap)
		if err != nil {
			s.logger.Error(fmt.Errorf(
				"error occured while trying to get sufficient store product size %d: %w",
				sub.StoreProductSizeID, err,
			))
			return err
		}

		// If success, we 'freeze' that usage. We add the usage from the product size to the frozenMap.
		for _, usage := range sps.ProductSize.ProductSizeIngredients {
			frozenMap[usage.IngredientID] += usage.Quantity
		}

		// Additives
		for _, addID := range sub.StoreAdditivesIDs {
			sa, err := s.storeAdditiveRepo.GetSufficientStoreAdditiveByID(storeID, addID, frozenMap)
			if err != nil {
				s.logger.Error(fmt.Errorf(
					"error occured while trying to get sufficient store additive %d: %w",
					addID, err,
				))
				return err
			}
			// freeze additive usage
			for _, ingrUsage := range sa.Additive.Ingredients {
				frozenMap[ingrUsage.IngredientID] += ingrUsage.Quantity
			}
		}
	}

	return nil
}

func (s *orderService) GetOrderBySubOrder(subOrderID uint) (*data.Order, error) {
	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order for suborder %d: %w", subOrderID, err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	return order, nil
}

func (s *orderService) GetOrderById(orderId uint) (types.OrderDTO, error) {
	order, err := s.orderRepo.GetOrderById(orderId)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order with ID %d: %w", orderId, err)
		s.logger.Error(wrappedErr.Error())
		return types.OrderDTO{}, wrappedErr
	}

	orderDTO := types.ConvertOrderToDTO(order)
	return orderDTO, nil
}

func (s *orderService) GetOrderDetails(orderID uint, filter *contexts.StoreContextFilter) (*types.OrderDetailsDTO, error) {
	order, err := s.orderRepo.GetOrderDetails(orderID, filter)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, nil
	}

	return types.MapToOrderDetailsDTO(order), nil
}

func (s *orderService) ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error) {
	orders, err := s.orderRepo.GetOrdersForExport(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for export: %w", err)
	}

	exports := make([]types.OrderExportDTO, len(orders))
	for i, order := range orders {
		exports[i] = types.ToOrderExportDTO(&order, order.Store.Name)
	}

	return exports, nil
}

var allowedTransitions = map[data.SubOrderStatus]data.SubOrderStatus{
	data.SubOrderStatusPending:   data.SubOrderStatusPreparing,
	data.SubOrderStatusPreparing: data.SubOrderStatusCompleted,
}

func (s *orderService) SetNextSubOrderStatus(subOrderID uint, options *types.ToggleNextSuborderStatusOptions) (*types.SuborderDTO, error) {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil || suborder == nil {
		return nil, fmt.Errorf("failed to retrieve suborder %d: %w", subOrderID, err)
	}

	// Handle fallback if suborder is already completed within time gap
	if dto := s.returnIfRecentlyCompleted(suborder, options); dto != nil {
		return dto, nil
	}

	if err := s.transactionManager.SetNextSubOrderStatus(suborder); err != nil {
		wrappedErr := fmt.Errorf("failed to set next suborder status suborder %d: %w", subOrderID, err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	// Return updated suborder
	updatedSuborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated suborder: %w", err)
	}
	dto := types.ConvertSuborderToDTO(updatedSuborder)
	return &dto, nil
}

func (s *orderService) returnIfRecentlyCompleted(suborder *data.Suborder, options *types.ToggleNextSuborderStatusOptions) *types.SuborderDTO {
	if options == nil || options.IncludeIfCompletedGapMinutes == nil {
		return nil
	}

	if suborder.Status == data.SubOrderStatusCompleted && suborder.CompletedAt != nil {
		minutes := *options.IncludeIfCompletedGapMinutes
		if time.Since(*suborder.CompletedAt) <= time.Duration(minutes)*time.Minute {
			dto := types.ConvertSuborderToDTO(suborder)
			return &dto
		}
	}
	return nil
}

func (s *orderService) SuccessOrderPayment(orderID uint, dto *types.TransactionDTO) error {
	paymentTransaction := types.ToTransactionModel(dto, orderID, data.TransactionTypePayment)
	order, err := s.orderRepo.HandlePaymentSuccess(orderID, paymentTransaction)
	if err != nil {
		s.logger.Errorf("failed to handle the order %d success: %v", orderID, err)
		return err
	}

	notificationDetails := &details.NewOrderNotificationDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           order.StoreID,
			FacilityName: order.Store.Name,
		},
		CustomerName: order.CustomerName,
		OrderID:      order.ID,
	}

	if err := s.notificationService.NotifyNewOrder(notificationDetails); err != nil {
		s.logger.Errorf("failed to notify new order: %w", err)
	}

	return nil
}

func (s *orderService) FailOrderPayment(orderID uint) error {
	err := s.orderRepo.HandlePaymentFailure(orderID)
	if err != nil {
		s.logger.Errorf("failed to delete the order %d after payment refuse", err)
		return err
	}

	return nil
}
