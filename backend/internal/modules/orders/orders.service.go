package orders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/taskqueue"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
	"go.uber.org/zap"
)

const (
	OrderPaymentFailure = "order-payment-failure"
)

type OrderService interface {
	GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error)
	GetAllBaristaOrders(filter types.OrdersTimeZoneFilter) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.SuborderDTO, error)
	GetStatusesCount(filter types.OrdersTimeZoneFilter) (types.OrderStatusesCountDTO, error)
	CreateOrder(storeId uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error)
	CompleteSubOrder(orderID, subOrderID uint) error
	CompleteSubOrderByBarcode(subOrderID uint) (*types.SuborderDTO, error)
	GeneratePDFReceipt(orderID uint) ([]byte, error)
	GetOrderBySubOrder(subOrderID uint) (*data.Order, error)
	GetOrderById(orderId uint) (types.OrderDTO, error)

	GetOrderDetails(orderID uint, filter *contexts.StoreContextFilter) (*types.OrderDetailsDTO, error)
	ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error)
	GenerateSuborderBarcodePDF(suborderID uint) ([]byte, error)

	AcceptSubOrder(subOrderID uint) error
	AdvanceSubOrderStatus(subOrderID uint, options *types.ToggleNextSuborderStatusOptions) (*types.SuborderDTO, error)

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
	logger              *zap.SugaredLogger
}

func NewOrderService(
	taskQueue taskqueue.TaskQueue,
	orderRepo OrderRepository,
	storeProductRepo storeProducts.StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) OrderService {
	return &orderService{
		taskQueue:           taskQueue,
		orderRepo:           orderRepo,
		storeProductRepo:    storeProductRepo,
		storeAdditiveRepo:   storeAdditiveRepo,
		storeStockRepo:      storeStockRepo,
		notificationService: notificationService,
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

	// This map will count how many additives are selected per category.
	categoryCount := make(map[uint]int)

	for _, addID := range storeAdditiveIDs {
		storeAdd, err := repo.GetStoreAdditiveByID(addID, &contexts.StoreContextFilter{StoreID: &storeID})
		if err != nil {
			return nil, nil, fmt.Errorf("error with store additive: %w", err)
		}
		if storeAdd == nil {
			return nil, nil, fmt.Errorf("store additive with ID %d is nil", addID)
		}
		if storeAdd.Additive.Name == "" {
			return nil, nil, fmt.Errorf("store additive with ID %d has an empty name", addID)
		}

		// Increase the count for the additive's category.
		categoryID := storeAdd.Additive.AdditiveCategoryID
		categoryCount[categoryID]++

		// If the additive category does NOT allow multiple selection,
		// then more than one additive in this category is an error.
		if !storeAdd.Additive.Category.IsMultipleSelect && categoryCount[categoryID] > 1 {
			return nil, nil, types.ErrMultipleSelect
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

func ValidateStoreProductSizes(
	storeID uint,
	storeProductSizeIDs []uint,
	repo storeProducts.StoreProductRepository,
	frozenMap map[uint]float64,
) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	productNames := make(map[uint]string)

	for _, psID := range storeProductSizeIDs {
		storePS, err := repo.GetStoreProductSizeById(storeID, psID)
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
				"insufficient stock for store product size %d: %w",
				sub.StoreProductSizeID, err,
			))
			return types.ErrInsufficientStock
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
					"insufficient stock for store additive %d: %w",
					addID, err,
				))
				return types.ErrInsufficientStock
			}
			// freeze additive usage
			for _, ingrUsage := range sa.Additive.Ingredients {
				frozenMap[ingrUsage.IngredientID] += ingrUsage.Quantity
			}
		}
	}

	return nil
}

func (s *orderService) CompleteSubOrder(orderID, subOrderID uint) error {
	completedAt := time.Now()
	updateSubOrder := types.UpdateSubOrderDTO{
		Status:      data.SubOrderStatusCompleted,
		CompletedAt: &completedAt,
	}
	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, updateSubOrder); err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	order, err := s.orderRepo.GetOrderById(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order by id: %w", err)
	}

	var completedSuborder *data.Suborder
	for i := range order.Suborders {
		if order.Suborders[i].ID == subOrderID {
			completedSuborder = &order.Suborders[i]
			break
		}
	}
	if completedSuborder == nil {
		return fmt.Errorf("suborder %d not found in order", subOrderID)
	}

	stockMap := make(map[uint]*data.StoreStock)
	if err := s.deductSuborderIngredientsFromStock(order.StoreID, completedSuborder, stockMap); err != nil {
		return fmt.Errorf("failed to deduct ingredients for suborder: %w", err)
	}

	s.notifyLowStockIngredients(order, stockMap)

	allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(orderID, subOrderID)
	if err != nil {
		return fmt.Errorf("failed to check suborder completion: %w", err)
	}
	if allCompleted {
		completedAt = time.Now()
		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}

		updateOrder := types.UpdateOrderDTO{
			Status:      newStatus,
			CompletedAt: &completedAt,
		}
		if err := s.orderRepo.UpdateOrderStatus(orderID, updateOrder); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return nil
}

func (s *orderService) GenerateSuborderBarcodePDF(suborderID uint) ([]byte, error) {
	suborder, err := s.orderRepo.GetSuborderByID(suborderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suborder (id=%d): %w", suborderID, err)
	}
	if suborder == nil {
		return nil, fmt.Errorf("no suborder found for ID %d", suborderID)
	}
	// if suborder.Status != data.SubOrderStatusPending {
	// 	return nil, fmt.Errorf("barcode generation allowed only for pending suborders")
	// }

	barcodeData := fmt.Sprintf("suborder-%d", suborder.ID)

	return utils.GenerateBarcodePDF(barcodeData)
}

func (s *orderService) CompleteSubOrderByBarcode(subOrderID uint) (*types.SuborderDTO, error) {
	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order for suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	completedAt := time.Now()
	updateSubOrder := types.UpdateSubOrderDTO{
		Status:      data.SubOrderStatusCompleted,
		CompletedAt: &completedAt,
	}
	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, updateSubOrder); err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	var completedSuborder *data.Suborder
	for i := range order.Suborders {
		if order.Suborders[i].ID == subOrderID {
			completedSuborder = &order.Suborders[i]
			break
		}
	}
	if completedSuborder == nil {
		return nil, fmt.Errorf("suborder %d not found in order", subOrderID)
	}

	stockMap := make(map[uint]*data.StoreStock)
	if err := s.deductSuborderIngredientsFromStock(order.StoreID, completedSuborder, stockMap); err != nil {
		return nil, fmt.Errorf("failed to deduct ingredients for additives: %w", err)
	}
	s.notifyLowStockIngredients(order, stockMap)

	allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(order.ID, subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}
	if allCompleted {
		completedAt = time.Now()
		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}
		updateOrder := types.UpdateOrderDTO{
			Status:      newStatus,
			CompletedAt: &completedAt,
		}
		if err := s.orderRepo.UpdateOrderStatus(order.ID, updateOrder); err != nil {
			wrappedErr := fmt.Errorf("failed to update order status: %w", err)
			s.logger.Error(wrappedErr.Error())
			return nil, wrappedErr
		}
	}

	subOrder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, err
	}
	response := types.ConvertSuborderToDTO(subOrder)
	return &response, nil
}

func (s *orderService) deductSuborderIngredientsFromStock(storeID uint, suborder *data.Suborder, stockMap map[uint]*data.StoreStock) error {
	updatedStocks, err := s.storeStockRepo.DeductStockByProductSizeTechCart(storeID, suborder.StoreProductSizeID)
	if err != nil {
		return fmt.Errorf("failed to deduct product size ingredients: %w", err)
	}
	for _, stock := range updatedStocks {
		if existingStock, exists := stockMap[stock.IngredientID]; exists {
			existingStock.Quantity = stock.Quantity
		} else {
			stockMap[stock.IngredientID] = &stock
		}
	}

	for _, subAdditive := range suborder.SuborderAdditives {
		updatedStocks, err := s.storeStockRepo.DeductStockByAdditiveTechCart(storeID, subAdditive.StoreAdditiveID)
		if err != nil {
			return fmt.Errorf("failed to deduct additive ingredients: %w", err)
		}
		for _, stock := range updatedStocks {
			if existingStock, exists := stockMap[stock.IngredientID]; exists {
				existingStock.Quantity = stock.Quantity
			} else {
				stockMap[stock.IngredientID] = &stock
			}
		}
	}

	return nil
}

func (s *orderService) notifyLowStockIngredients(order *data.Order, stockMap map[uint]*data.StoreStock) {
	for _, stock := range stockMap {
		if stock.Quantity <= stock.LowStockThreshold {
			notificationDetails := &details.StoreWarehouseRunOutDetails{
				BaseNotificationDetails: details.BaseNotificationDetails{
					ID:           order.StoreID,
					FacilityName: order.Store.Name,
				},
				StockItem:   stock.Ingredient.Name,
				StockItemID: stock.IngredientID,
			}
			err := s.notificationService.NotifyStoreWarehouseRunOut(notificationDetails)
			if err != nil {
				s.logger.Errorf("failed to send notification: %v", err)
			}
		}
	}
}

func (s *orderService) GeneratePDFReceipt(orderID uint) ([]byte, error) {
	order, err := s.orderRepo.GetOrderById(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order details: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	detailsPDF := pdf.PDFReceiptDetails{
		OrderID:   order.ID,
		StoreID:   order.StoreID,
		OrderDate: order.CreatedAt.Format("2006-01-02 15:04:05"),
		Total:     order.Total,
	}

	for _, suborder := range order.Suborders {
		pdfSuborder := pdf.PDFSubOrder{
			ProductName: fmt.Sprintf("Product #%d", suborder.StoreProductSizeID),
			Price:       suborder.Price,
			Status:      string(suborder.Status),
		}

		for _, storeAdditive := range suborder.SuborderAdditives {
			pdfSuborder.Additives = append(pdfSuborder.Additives, pdf.PDFAdditive{
				Name:  fmt.Sprintf("Additive #%d", storeAdditive.StoreAdditiveID),
				Price: storeAdditive.Price,
			})
		}

		detailsPDF.SubOrders = append(detailsPDF.SubOrders, pdfSuborder)
	}

	return pdf.GeneratePDFReceipt(detailsPDF)
}

func (s *orderService) GetStatusesCount(filter types.OrdersTimeZoneFilter) (types.OrderStatusesCountDTO, error) {
	countsMap, err := s.orderRepo.GetStatusesCount(filter)
	if err != nil {
		wrappedErr := fmt.Errorf("error couting statuses: %w", err)
		s.logger.Error(wrappedErr.Error())
		return types.OrderStatusesCountDTO{}, wrappedErr
	}

	dto := types.OrderStatusesCountDTO{
		PENDING:     countsMap[data.OrderStatusPending],
		PREPARING:   countsMap[data.OrderStatusPreparing],
		COMPLETED:   countsMap[data.OrderStatusCompleted],
		IN_DELIVERY: countsMap[data.OrderStatusInDelivery],
		DELIVERED:   countsMap[data.OrderStatusDelivered],
		CANCELLED:   countsMap[data.OrderStatusCancelled],
	}

	dto.ALL = dto.PENDING + dto.PREPARING + dto.COMPLETED + dto.IN_DELIVERY + dto.DELIVERED + dto.CANCELLED

	return dto, nil
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

func (s *orderService) AcceptSubOrder(subOrderID uint) error {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve suborder: %w", err)
	}
	if suborder == nil {
		return fmt.Errorf("suborder %d not found", subOrderID)
	}

	if suborder.Status != data.SubOrderStatusPending {
		return fmt.Errorf("suborder %d is not pending", subOrderID)
	}

	updateSubOrder := types.UpdateSubOrderDTO{
		Status: data.SubOrderStatusPreparing,
	}
	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, updateSubOrder); err != nil {
		return fmt.Errorf("failed to update suborder status: %w", err)
	}

	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	suborders, err := s.orderRepo.GetSubOrdersByOrderID(order.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch suborders for order %d: %w", order.ID, err)
	}

	allAccepted := true
	for _, so := range suborders {
		if so.Status == data.SubOrderStatusPending {
			allAccepted = false
			break
		}
	}

	updateOrder := types.UpdateOrderDTO{
		Status: data.OrderStatusPreparing,
	}
	// If all suborders have been accepted (i.e. transitioned to PREPARING or later)
	// and the order is still in the PENDING state, update the order's status.
	if allAccepted && order.Status == data.OrderStatusPending {
		if err := s.orderRepo.UpdateOrderStatus(order.ID, updateOrder); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}
	return nil
}

var allowedTransitions = map[data.SubOrderStatus]data.SubOrderStatus{
	data.SubOrderStatusPending:   data.SubOrderStatusPreparing,
	data.SubOrderStatusPreparing: data.SubOrderStatusCompleted,
}

func (s *orderService) AdvanceSubOrderStatus(subOrderID uint, options *types.ToggleNextSuborderStatusOptions) (*types.SuborderDTO, error) {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil || suborder == nil {
		return nil, fmt.Errorf("failed to retrieve suborder %d: %w", subOrderID, err)
	}

	// Handle fallback if suborder is already completed within time gap
	if dto := s.returnIfRecentlyCompleted(suborder, options); dto != nil {
		return dto, nil
	}

	// Attempt to advance suborder status
	if err := s.advanceSuborder(subOrderID, suborder); err != nil {
		return nil, err
	}

	// Sync and update order status
	if err := s.updateOrderStatusBySuborder(subOrderID); err != nil {
		return nil, err
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

func (s *orderService) advanceSuborder(subOrderID uint, suborder *data.Suborder) error {
	currentStatus := suborder.Status
	nextStatus, ok := allowedTransitions[currentStatus]
	if !ok {
		return fmt.Errorf("no allowed transition from status %s", currentStatus)
	}

	completedAt := time.Now()
	update := types.UpdateSubOrderDTO{
		Status:      nextStatus,
		CompletedAt: &completedAt,
	}
	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, update); err != nil {
		return fmt.Errorf("failed to update suborder status: %w", err)
	}

	// If suborder is completed, deduct ingredients
	if nextStatus == data.SubOrderStatusCompleted {
		if err := s.handleSuborderCompletion(subOrderID); err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) handleSuborderCompletion(subOrderID uint) error {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve updated suborder: %w", err)
	}

	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	stockMap := make(map[uint]*data.StoreStock)
	if err := s.deductSuborderIngredientsFromStock(order.StoreID, suborder, stockMap); err != nil {
		return fmt.Errorf("failed to deduct ingredients: %w", err)
	}

	s.notifyLowStockIngredients(order, stockMap)
	return nil
}

func (s *orderService) updateOrderStatusBySuborder(subOrderID uint) error {
	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	suborders, err := s.orderRepo.GetSubOrdersByOrderID(order.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch suborders for order %d: %w", order.ID, err)
	}

	hasPreparing, allCompleted := evaluateSuborderStatuses(suborders)

	switch {
	case hasPreparing:
		return s.ensureOrderStatus(order, data.OrderStatusPreparing, nil)

	case allCompleted:
		newStatus := data.OrderStatusCompleted
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		}
		now := time.Now()
		return s.ensureOrderStatus(order, newStatus, &now)

	default:
		return s.ensureOrderStatus(order, data.OrderStatusPreparing, nil)
	}
}

func (s *orderService) ensureOrderStatus(order *data.Order, status data.OrderStatus, completedAt *time.Time) error {
	if order.Status == status {
		return nil
	}

	update := types.UpdateOrderDTO{
		Status:      status,
		CompletedAt: completedAt,
	}
	if err := s.orderRepo.UpdateOrderStatus(order.ID, update); err != nil {
		return fmt.Errorf("failed to update order status to %s: %w", status, err)
	}
	return nil
}

func evaluateSuborderStatuses(suborders []data.Suborder) (hasPreparing bool, allCompleted bool) {
	hasPreparing = false
	allCompleted = true
	for _, so := range suborders {
		if so.Status == data.SubOrderStatusPreparing {
			hasPreparing = true
		}
		if so.Status != data.SubOrderStatusCompleted {
			allCompleted = false
		}
	}
	return
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
