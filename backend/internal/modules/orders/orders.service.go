package orders

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

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

	GetOrderDetails(orderID uint) (*types.OrderDetailsDTO, error)
	ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error)
	GenerateSuborderBarcodePDF(suborderID uint) ([]byte, error)

	AcceptSubOrder(subOrderID uint) error
	AdvanceSubOrderStatus(subOrderID uint) (*types.SuborderDTO, error)
}

type orderValidationResults struct {
	ProductPrices  map[uint]float64
	ProductNames   map[uint]string
	AdditivePrices map[uint]float64
	AdditiveNames  map[uint]string
}

type orderService struct {
	orderRepo           OrderRepository
	storeProductRepo    storeProducts.StoreProductRepository
	storeAdditiveRepo   storeAdditives.StoreAdditiveRepository
	storeStockRepo      storeStocks.StoreStockRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewOrderService(
	orderRepo OrderRepository,
	storeProductRepo storeProducts.StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) OrderService {
	return &orderService{
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
		return nil, err
	}

	if len(createOrderDTO.Suborders) == 0 {
		return nil, fmt.Errorf("order can not be empty")
	}

	// 1) Calculate how much is already frozen by existing PENDING/PREPARING orders
	frozenMap, err := s.orderRepo.CalculateFrozenStock(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate frozen stock: %w", err)
	}

	// 2) Attempt to create + freeze new suborders usage
	//    If we can't fulfill them all, it returns an error early
	if err := s.CheckAndAccumulateSuborders(storeID, createOrderDTO.Suborders, frozenMap); err != nil {
		// e.g. \"not enough stock for ingredient X\"
		return nil, err
	}

	// 3) If we get here, we have effectively 'reserved' usage in frozenMap
	//    Next, do your usual 'ValidationResults' to figure out item names & prices
	storeProductSizeIDs, storeAdditiveIDs := RetrieveIDs(*createOrderDTO)
	validations, err := s.ValidationResults(
		storeID,
		storeProductSizeIDs,
		storeAdditiveIDs, // we no longer need to pass frozenMap here if it only calculates price,
		frozenMap,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("validation failed: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	// 4) Build + save the order
	createOrderDTO.StoreID = storeID
	order, total := types.ConvertCreateOrderDTOToOrder(
		createOrderDTO,
		validations.ProductPrices,
		validations.AdditivePrices,
	)
	order.Status = data.OrderStatusPending
	order.Total = total

	if err := s.orderRepo.CreateOrder(&order); err != nil {
		wrappedErr := fmt.Errorf("error creating order: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	// optional notification
	notificationDetails := &details.NewOrderNotificationDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           order.StoreID,
			FacilityName: order.Store.Name,
		},
		CustomerName: createOrderDTO.CustomerName,
		OrderID:      order.ID,
	}
	if err := s.notificationService.NotifyNewOrder(notificationDetails); err != nil {
		s.logger.Errorf("failed to notify new order: %w", err)
	}

	return &order, nil
}

func (s *orderService) ValidationResults(
	storeID uint,
	storeProductSizeIDs, storeAdditiveIDs []uint,
	frozenMap map[uint]float64,
) (*orderValidationResults, error) {
	productPrices, productNames, err := ValidateStoreProductSizes(storeID, storeProductSizeIDs, s.storeProductRepo, frozenMap)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	additivePrices, additiveNames, err := ValidateStoreAdditives(storeID, storeAdditiveIDs, s.storeAdditiveRepo, frozenMap)
	if err != nil {
		return nil, fmt.Errorf("additive validation failed: %w", err)
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
		storeAdd, err := repo.GetSufficientStoreAdditiveByID(storeID, addID, frozenMap)
		if err != nil {
			return nil, nil, fmt.Errorf("error with store additive: %w", err)
		}
		if storeAdd == nil {
			return nil, nil, fmt.Errorf("store additive with ID %d is nil", addID)
		}
		if storeAdd.Additive.Name == "" {
			return nil, nil, fmt.Errorf("store additive with ID %d has an empty name", addID)
		}

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
		storePS, err := repo.GetSufficientStoreProductSizeById(storeID, psID, frozenMap)
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
	// 1) Expand any suborders with quantity > 1 into separate single-quantity items
	expanded := ExpandSuborders(suborders)

	// 2) For each single suborder, call existing checks
	for _, sub := range expanded {
		// A) Product Size
		sps, err := s.storeProductRepo.GetSufficientStoreProductSizeById(storeID, sub.StoreProductSizeID, frozenMap)
		if err != nil {
			return fmt.Errorf(
				"insufficient stock for store product size %d: %w",
				sub.StoreProductSizeID, err,
			)
		}

		// If success, we 'freeze' that usage. We add the usage from the product size to the frozenMap.
		for _, usage := range sps.ProductSize.ProductSizeIngredients {
			frozenMap[usage.IngredientID] += usage.Quantity
		}

		// B) Additives
		for _, addID := range sub.StoreAdditivesIDs {
			sa, err := s.storeAdditiveRepo.GetSufficientStoreAdditiveByID(storeID, addID, frozenMap)
			if err != nil {
				return fmt.Errorf(
					"insufficient stock for store additive %d: %w",
					addID, err,
				)
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
	err := s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusCompleted)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(orderID, subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	if allCompleted {
		order, err := s.orderRepo.GetOrderById(orderID)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get order by id: %w", err)
			s.logger.Error(wrappedErr.Error())
			return wrappedErr
		}

		stockMap := make(map[uint]*data.StoreStock)

		err = s.deductProductSizeIngredientsFromStock(order, stockMap)
		if err != nil {
			return fmt.Errorf("failed to deduct ingredients for products: %w", err)
		}

		err = s.deductAdditiveIngredientsFromStock(order, stockMap)
		if err != nil {
			return fmt.Errorf("failed to deduct ingredients for additives: %w", err)
		}

		s.notifyLowStockIngredients(order, stockMap)

		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}

		err = s.orderRepo.UpdateOrderStatus(orderID, newStatus)

		if err != nil {
			wrappedErr := fmt.Errorf("failed to update order status: %w", err)
			s.logger.Error(wrappedErr.Error())
			return wrappedErr
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

	err = s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusCompleted)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(order.ID, subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	if allCompleted {
		order, err := s.orderRepo.GetOrderById(order.ID)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get order by id: %w", err)
			s.logger.Error(wrappedErr.Error())
			return nil, wrappedErr
		}

		stockMap := make(map[uint]*data.StoreStock)

		err = s.deductProductSizeIngredientsFromStock(order, stockMap)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct ingredients for products: %w", err)
		}

		err = s.deductAdditiveIngredientsFromStock(order, stockMap)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct ingredients for additives: %w", err)
		}

		s.notifyLowStockIngredients(order, stockMap)

		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}

		err = s.orderRepo.UpdateOrderStatus(order.ID, newStatus)

		if err != nil {
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

func (s *orderService) deductProductSizeIngredientsFromStock(order *data.Order, stockMap map[uint]*data.StoreStock) error {
	for _, suborder := range order.Suborders {
		updatedStocks, err := s.storeStockRepo.DeductStockByProductSizeTechCart(order.StoreID, suborder.StoreProductSizeID)
		if err != nil {
			return fmt.Errorf("failed to deduct from store stock: %w", err)
		}

		for _, stock := range updatedStocks {
			if existingStock, exists := stockMap[stock.IngredientID]; exists {
				existingStock.Quantity = stock.Quantity // Update latest quantity
			} else {
				stockMap[stock.IngredientID] = &stock
			}
		}
	}

	return nil
}

func (s *orderService) deductAdditiveIngredientsFromStock(order *data.Order, stockMap map[uint]*data.StoreStock) error {
	for _, suborder := range order.Suborders {
		for _, storeAdditive := range suborder.SuborderAdditives {
			updatedStocks, err := s.storeStockRepo.DeductStockByAdditiveTechCart(order.StoreID, storeAdditive.StoreAdditiveID)
			if err != nil {
				return fmt.Errorf("failed to deduct from store stock: %w", err)
			}

			for _, stock := range updatedStocks {
				if existingStock, exists := stockMap[stock.IngredientID]; exists {
					existingStock.Quantity = stock.Quantity // Update latest quantity
				} else {
					stockMap[stock.IngredientID] = &stock
				}
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

func (s *orderService) GetOrderDetails(orderID uint) (*types.OrderDetailsDTO, error) {
	order, err := s.orderRepo.GetOrderDetails(orderID)
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

	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusPreparing); err != nil {
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

	// If all suborders have been accepted (i.e. transitioned to PREPARING or later)
	// and the order is still in the PENDING state, update the order's status.
	if allAccepted && order.Status == data.OrderStatusPending {
		if err := s.orderRepo.UpdateOrderStatus(order.ID, data.OrderStatusPreparing); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}
	return nil
}

var allowedTransitions = map[data.SubOrderStatus]data.SubOrderStatus{
	data.SubOrderStatusPending:   data.SubOrderStatusPreparing,
	data.SubOrderStatusPreparing: data.SubOrderStatusCompleted,
}

func (s *orderService) AdvanceSubOrderStatus(subOrderID uint) (*types.SuborderDTO, error) {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suborder: %w", err)
	}
	if suborder == nil {
		return nil, fmt.Errorf("suborder %d not found", subOrderID)
	}

	currentStatus := suborder.Status
	nextStatus, ok := allowedTransitions[currentStatus]
	if !ok {
		return nil, fmt.Errorf("no allowed transition from status %s", currentStatus)
	}

	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, nextStatus); err != nil {
		return nil, fmt.Errorf("failed to update suborder status: %w", err)
	}

	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	suborders, err := s.orderRepo.GetSubOrdersByOrderID(order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch suborders for order %d: %w", order.ID, err)
	}

	hasPreparing, allCompleted := evaluateSuborderStatuses(suborders)

	if hasPreparing {
		if order.Status != data.OrderStatusPreparing {
			if err := s.orderRepo.UpdateOrderStatus(order.ID, data.OrderStatusPreparing); err != nil {
				return nil, fmt.Errorf("failed to update order status to pending: %w", err)
			}
		}
	} else if allCompleted {
		var newOrderStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newOrderStatus = data.OrderStatusInDelivery
		} else {
			newOrderStatus = data.OrderStatusCompleted
		}
		if err := s.orderRepo.UpdateOrderStatus(order.ID, newOrderStatus); err != nil {
			return nil, fmt.Errorf("failed to update order status to completed: %w", err)
		}

		stockMap := make(map[uint]*data.StoreStock)

		err = s.deductProductSizeIngredientsFromStock(order, stockMap)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct ingredients for products: %w", err)
		}

		err = s.deductAdditiveIngredientsFromStock(order, stockMap)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct ingredients for additives: %w", err)
		}

		s.notifyLowStockIngredients(order, stockMap)

	} else {
		if order.Status != data.OrderStatusPreparing {
			if err := s.orderRepo.UpdateOrderStatus(order.ID, data.OrderStatusPreparing); err != nil {
				return nil, fmt.Errorf("failed to update order status to preparing: %w", err)
			}
		}
	}

	updatedSuborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated suborder: %w", err)
	}

	dto := types.ConvertSuborderToDTO(updatedSuborder)
	return &dto, nil
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
