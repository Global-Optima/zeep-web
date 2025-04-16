package orders

import (
	"encoding/json"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	storeAdditivesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	storeStocksTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/sirupsen/logrus"

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

type storeAdditiveValidationResults struct {
	storeAdditivesList []data.StoreAdditive
	prices             map[uint]float64
	names              map[uint]string
}

type suborderAdditivesContext struct {
	storeAddKeys []storeAdditivesTypes.StorePSAdditiveKey
	storeAddIDs  []uint
}

type OrderService interface {
	GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error)
	GetAllBaristaOrders(filter types.OrdersTimeZoneFilter) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.SuborderDTO, error)
	CreateOrder(createOrderDTO *types.CreateOrderDTO) (*data.Order, error)
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
	taskQueue                 taskqueue.TaskQueue
	orderRepo                 OrderRepository
	storeProductRepo          storeProducts.StoreProductRepository
	storeAdditiveRepo         storeAdditives.StoreAdditiveRepository
	storeStockRepo            storeStocks.StoreStockRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
	storeProductService       storeProducts.StoreProductService
	storeAdditiveService      storeAdditives.StoreAdditiveService
	notificationService       notifications.NotificationService
	transactionManager        TransactionManager
	logger                    *zap.SugaredLogger
}

func NewOrderService(
	taskQueue taskqueue.TaskQueue,
	orderRepo OrderRepository,
	storeProductRepo storeProducts.StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	storeProductService storeProducts.StoreProductService,
	storeAdditiveService storeAdditives.StoreAdditiveService,
	notificationService notifications.NotificationService,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) OrderService {
	return &orderService{
		taskQueue:                 taskQueue,
		orderRepo:                 orderRepo,
		storeProductRepo:          storeProductRepo,
		storeAdditiveRepo:         storeAdditiveRepo,
		storeStockRepo:            storeStockRepo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
		storeProductService:       storeProductService,
		storeAdditiveService:      storeAdditiveService,
		notificationService:       notificationService,
		transactionManager:        transactionManager,
		logger:                    logger,
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

func (s *orderService) CreateOrder(createOrderDTO *types.CreateOrderDTO) (*data.Order, error) {
	start := time.Now()
	censorValidator := censor.GetCensorValidator()

	if err := censorValidator.ValidateText(createOrderDTO.CustomerName); err != nil {
		s.logger.Error(err)
		return nil, types.ErrInvalidCustomerNameCensor
	}

	if len(createOrderDTO.Suborders) == 0 {
		return nil, fmt.Errorf("order can not be empty")
	}

	validations, err := validateSuborders(createOrderDTO, s.storeProductRepo, s.storeAdditiveRepo)
	if err != nil {
		wrappedErr := fmt.Errorf("suborders validation failed: %v", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	logrus.Infof("Order: %v", createOrderDTO)
	logrus.Infof("storeID: %v, createOrderDTO: %v, storeAdditiveRepo: %v", createOrderDTO.StoreID, createOrderDTO, s.storeAdditiveRepo)

	frozenInventory, err := s.storeInventoryManagerRepo.CalculateFrozenInventory(createOrderDTO.StoreID, nil)
	if err != nil {
		wrappedErr := fmt.Errorf("suborders inventory check failed: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	logrus.Infof("Frozen Stock BEFORE creating order: %v", frozenInventory)
	if err := s.CheckAndAccumulateSuborders(createOrderDTO.StoreID, createOrderDTO.Suborders, frozenInventory); err != nil {
		wrappedErr := fmt.Errorf("suborders inventory check failed: %v", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	order, total := types.ConvertCreateOrderDTOToOrder(
		createOrderDTO,
		validations.ProductPrices,
		validations.AdditivePrices,
	)
	order.Status = data.OrderStatusWaitingForPayment
	order.Total = total

	id, err := s.orderRepo.CreateOrder(&order)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create order: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	go func() {
		inventoryLists, err := s.orderRepo.GetOrderInventory(id)
		if err != nil {
			s.logger.Error("failed to get ingredients", zap.Error(err))
			return
		}

		err = s.storeInventoryManagerRepo.RecalculateStoreInventory(order.StoreID, &storeInventoryManagersTypes.RecalculateInput{
			IngredientIDs:   inventoryLists.IngredientIDs,
			ProvisionIDs:    inventoryLists.ProvisionIDs,
			FrozenInventory: frozenInventory,
		})
		if err != nil {
			s.logger.Error("failed to recalculate out of stock", zap.Error(err))
			return
		}
	}()

	go func() {
		payload, err := json.Marshal(types.WaitingOrderPayload{OrderID: id})
		if err != nil {
			s.logger.Error("failed to marshal payload", zap.Error(err))
			return
		}

		err = s.taskQueue.EnqueueTask(OrderPaymentFailure, payload, config.GetConfig().Payment.WaitingTime)
		if err != nil {
			s.logger.Error("failed to enqueue payment failure", zap.Error(err))
			return
		}
	}()

	logrus.Infof("________________order done in %v_________________", time.Since(start))
	return &order, nil
}

func validateSuborders(
	order *types.CreateOrderDTO,
	storeProductRepo storeProducts.StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
) (*orderValidationResults, error) {
	storeProductSizeIDs, subAddCtx := retrieveIDs(order.Suborders)

	productPrices, productNames, err := validateStoreProductSizes(order.StoreID, storeProductSizeIDs, storeProductRepo)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	saValRes, err := validateStoreAdditives(
		order.StoreID,
		subAddCtx,
		storeAdditiveRepo,
	)
	if err != nil {
		return nil, fmt.Errorf("additive validation failed: %w", err)
	}

	if err := validateMultipleSelect(order.Suborders, saValRes.storeAdditivesList); err != nil {
		return nil, err
	}

	return &orderValidationResults{
		ProductPrices:  productPrices,
		ProductNames:   productNames,
		AdditivePrices: saValRes.prices,
		AdditiveNames:  saValRes.names,
	}, nil
}

func validateStoreAdditives(
	storeID uint,
	ctx *suborderAdditivesContext,
	repo storeAdditives.StoreAdditiveRepository,
) (*storeAdditiveValidationResults, error) {
	// 1) Batch fetch StoreAdditives by the collected additive IDs.
	storeAdditivesList, err := repo.GetStoreAdditivesByIDs(storeID, ctx.storeAddIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch storeAdditives: %w", err)
	}
	// Build a lookup map: StoreAdditiveID -> *StoreAdditive
	storeAdditivesMap := make(map[uint]*data.StoreAdditive, len(storeAdditivesList))
	for i := range storeAdditivesList {
		// Use index-based addressing so that the pointer is stable.
		sa := &storeAdditivesList[i]
		storeAdditivesMap[sa.ID] = sa
	}

	// 2) Batch fetch ProductSizeAdditives using the aggregated keys.
	psaMap, err := repo.GetPSAByProductSizeAndAdditive(ctx.storeAddKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch productSizeAdditives: %w", err)
	}

	// 3) Validate each aggregated key and build results.
	prices := make(map[uint]float64)
	additiveNames := make(map[uint]string)

	// Iterate over the unique keys from context
	for _, key := range ctx.storeAddKeys {
		sa := storeAdditivesMap[key.StoreAdditiveID]
		if sa == nil {
			return nil, fmt.Errorf("no storeAdditive found for ID=%d", key.StoreAdditiveID)
		}
		// Look up the corresponding ProductSizeAdditive
		psa, ok := psaMap[key]
		if !ok || psa == nil {
			return nil, fmt.Errorf(
				"additive %d not linked to storeProductSize %d in store %d",
				key.StoreAdditiveID, key.StoreProductSizeID, storeID,
			)
		}

		// Check stock
		if sa.IsOutOfStock {
			return nil, fmt.Errorf(
				"%w: additive %s (ID=%d) is out of stock",
				storeStocksTypes.ErrInsufficientStock, sa.Additive.Name, key.StoreAdditiveID,
			)
		}

		// Compute the final price:
		var price float64
		if psa.IsDefault {
			return nil, fmt.Errorf("additive with ID %d is a default additive for productSize with ID %d", psa.AdditiveID, psa.ProductSizeID)
		} else if sa.StorePrice != nil {
			price = *sa.StorePrice
		} else {
			price = sa.Additive.BasePrice
		}

		// We store the results keyed by StoreAdditiveID.
		prices[key.StoreAdditiveID] = price
		additiveNames[key.StoreAdditiveID] = sa.Additive.Name
	}

	return &storeAdditiveValidationResults{
		storeAdditivesList: storeAdditivesList,
		prices:             prices,
		names:              additiveNames,
	}, nil
}

func validateMultipleSelect(
	suborders []types.CreateSubOrderDTO,
	storeAdditivesList []data.StoreAdditive,
) error {
	// Build a map: storeAdditiveID → *StoreAdditive
	storeAdditivesMap := make(map[uint]*data.StoreAdditive, len(storeAdditivesList))
	for i := range storeAdditivesList {
		sa := &storeAdditivesList[i]
		storeAdditivesMap[sa.ID] = sa
	}

	// For each suborder, ensure non-multiple-select categories aren’t picked >1 time
	for _, suborder := range suborders {
		categoryCount := make(map[uint]int)
		for _, addID := range suborder.StoreAdditivesIDs {
			sa := storeAdditivesMap[addID]
			if sa == nil {
				return fmt.Errorf("additive %d not found in the storeAdditivesList map", addID)
			}

			categoryID := sa.Additive.AdditiveCategoryID
			categoryCount[categoryID]++

			if !sa.Additive.Category.IsMultipleSelect && categoryCount[categoryID] > 1 {
				return types.ErrMultipleSelect
			}
		}
	}

	return nil
}

func validateStoreProductSizes(
	storeID uint,
	storeProductSizeIDs []uint,
	repo storeProducts.StoreProductRepository,
) (map[uint]float64, map[uint]string, error) {

	// 1) Fetch all needed storeProductSizes at once
	storePSList, err := repo.GetStoreProductSizesWithDetailsByIDs(storeID, storeProductSizeIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching store product sizes: %w", err)
	}

	// 2) Convert slice → map for quick lookups by ID
	storePSMap := make(map[uint]*data.StoreProductSize, len(storePSList))
	for i := range storePSList {
		ps := &storePSList[i]
		storePSMap[ps.ID] = ps
	}

	// 3) Validate & collect results
	prices := make(map[uint]float64, len(storeProductSizeIDs))
	productNames := make(map[uint]string, len(storeProductSizeIDs))

	for _, psID := range storeProductSizeIDs {
		storePS := storePSMap[psID]
		if storePS == nil {
			return nil, nil, fmt.Errorf("storeProductSize with ID %d not found in store %d", psID, storeID)
		}

		// Check for missing product name
		// (We assume storePS.ProductSize is preloaded)
		if storePS.ProductSize.Product.Name == "" {
			return nil, nil, fmt.Errorf(
				"product size with ID %d has an associated product with an empty name",
				psID,
			)
		}

		// Compute the final price
		price := storePS.ProductSize.BasePrice
		if storePS.StorePrice != nil {
			price = *storePS.StorePrice
		}

		prices[psID] = price
		// storePS.StoreProduct.Product.Name is the "actual" product name
		productNames[psID] = storePS.StoreProduct.Product.Name
	}

	return prices, productNames, nil
}

func retrieveIDs(suborders []types.CreateSubOrderDTO) ([]uint, *suborderAdditivesContext) {
	var storeProductSizeIDs []uint

	// Use maps as sets for uniqueness.
	addKeySet := make(map[storeAdditivesTypes.StorePSAdditiveKey]struct{})
	additiveIDSet := make(map[uint]struct{})

	for _, sub := range suborders {
		storeProductSizeIDs = append(storeProductSizeIDs, sub.StoreProductSizeID)
		for _, addID := range sub.StoreAdditivesIDs {
			key := storeAdditivesTypes.StorePSAdditiveKey{
				StoreProductSizeID: sub.StoreProductSizeID,
				StoreAdditiveID:    addID,
			}
			addKeySet[key] = struct{}{}
			additiveIDSet[addID] = struct{}{}
		}
	}

	// Convert set to slice.
	addKeys := make([]storeAdditivesTypes.StorePSAdditiveKey, 0, len(addKeySet))
	for k := range addKeySet {
		addKeys = append(addKeys, k)
	}
	additiveIDs := make([]uint, 0, len(additiveIDSet))
	for id := range additiveIDSet {
		additiveIDs = append(additiveIDs, id)
	}

	return storeProductSizeIDs, &suborderAdditivesContext{
		storeAddKeys: addKeys,
		storeAddIDs:  additiveIDs,
	}
}

func (s *orderService) CheckAndAccumulateSuborders(
	storeID uint,
	suborders []types.CreateSubOrderDTO,
	frozenInventory *storeInventoryManagersTypes.FrozenInventory,
) error {
	// Expand any suborders with quantity > 1 into separate single-quantity items
	expanded := ExpandSuborders(suborders)

	// For each single suborder, call existing checks
	for _, sub := range expanded {
		// Product Size

		logrus.Infof("storeID: %v, storeProductSizeID: %v, frozenInventory: %v", storeID, sub.StoreProductSizeID, frozenInventory)
		err := s.storeProductService.CheckSufficientStoreProductSizeByID(storeID, sub.StoreProductSizeID, frozenInventory)
		if err != nil {
			s.logger.Error(fmt.Errorf(
				"error occured while trying to get sufficient store product size %d: %w",
				sub.StoreProductSizeID, err,
			))
			return err
		}
		logrus.Infof("Frozen Stock AFTER 1st suborder: %v", frozenInventory)

		// Additives
		for _, addID := range sub.StoreAdditivesIDs {
			err := s.storeAdditiveService.CheckSufficientStoreAdditiveByID(storeID, addID, frozenInventory)
			if err != nil {
				s.logger.Error(fmt.Errorf(
					"error occured while trying to get sufficient store additive %d: %w",
					addID, err,
				))
				return err
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
	order, err := s.orderRepo.GetRawOrderById(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete the order %d after payment refuse: failed to get order data: %w", orderID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	inventoryLists, err := s.orderRepo.GetOrderInventory(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("could not get inventory lists for order id %d: %w", orderID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if order.Status != data.OrderStatusWaitingForPayment {
		return types.ErrInappropriateOrderStatus
	}

	err = s.orderRepo.HardDeleteOrderByID(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete the order %d after payment refuse: %w", orderID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	go func() {
		err = s.storeInventoryManagerRepo.RecalculateStoreInventory(order.StoreID, &storeInventoryManagersTypes.RecalculateInput{
			IngredientIDs: inventoryLists.IngredientIDs,
			ProvisionIDs:  inventoryLists.ProvisionIDs,
		})
		if err != nil {
			s.logger.Errorf("could not recalculate stock after order deletion: %v", err)
		}
	}()

	return nil
}
