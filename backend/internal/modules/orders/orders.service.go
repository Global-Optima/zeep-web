package orders

import (
	"fmt"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"

	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
)

type OrderService interface {
	GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error)
	GetAllBaristaOrders(storeID uint) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.SuborderDTO, error)
	GetStatusesCount(storeID uint) (types.OrderStatusesCountDTO, error)
	CreateOrder(storeId uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error)
	CompleteSubOrder(subOrderID uint) error
	GeneratePDFReceipt(orderID uint) ([]byte, error)
	GetOrderBySubOrder(subOrderID uint) (*data.Order, error)
	GetOrderById(orderId uint) (types.OrderDTO, error)

	GetOrderDetails(orderID uint) (*types.OrderDetailsDTO, error)
	ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error)
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
	storeWarehouseRepo  storeWarehouses.StoreWarehouseRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewOrderService(orderRepo OrderRepository, storeProductRepo storeProducts.StoreProductRepository, storeAdditiveRepo storeAdditives.StoreAdditiveRepository, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository, notificationService notifications.NotificationService, logger *zap.SugaredLogger) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		storeProductRepo:    storeProductRepo,
		storeAdditiveRepo:   storeAdditiveRepo,
		storeWarehouseRepo:  storeWarehouseRepo,
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

func (s *orderService) GetAllBaristaOrders(storeID uint) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetAllBaristaOrders(storeID)

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

func (s *orderService) CreateOrder(storeID uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error) {
	storeProductSizeIDs, storeAdditiveIDs := RetrieveIDs(*createOrderDTO)

	validations, err := s.ValidationResults(storeID, storeProductSizeIDs, storeAdditiveIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("validation failed: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	createOrderDTO.StoreID = storeID
	order, total := types.ConvertCreateOrderDTOToOrder(createOrderDTO,
		validations.ProductPrices, validations.AdditivePrices)
	order.Total = total

	err = s.orderRepo.CreateOrder(&order)
	if err != nil {
		wrappedErr := fmt.Errorf("error creating order: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	details := &details.NewOrderNotificationDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           order.StoreID,
			FacilityName: order.Store.Name,
		},
		CustomerName: createOrderDTO.CustomerName,
		OrderID:      order.ID,
	}
	err = s.notificationService.NotifyNewOrder(details)
	if err != nil {
		return nil, fmt.Errorf("failed to notify new order: %w", err)
	}

	return &order, nil
}

func (s *orderService) CompleteSubOrder(subOrderID uint) error {
	err := s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusCompleted)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	orderID, allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(subOrderID)
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

		stockMap := make(map[uint]*data.StoreWarehouseStock)

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

func (s *orderService) deductProductSizeIngredientsFromStock(order *data.Order, stockMap map[uint]*data.StoreWarehouseStock) error {
	for _, suborder := range order.Suborders {
		updatedStocks, err := s.storeWarehouseRepo.DeductStockByProductSizeTechCart(order.StoreID, suborder.StoreProductSizeID)
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

func (s *orderService) deductAdditiveIngredientsFromStock(order *data.Order, stockMap map[uint]*data.StoreWarehouseStock) error {
	for _, suborder := range order.Suborders {
		for _, storeAdditive := range suborder.StoreAdditives {
			updatedStocks, err := s.storeWarehouseRepo.DeductStockByAdditiveTechCart(order.StoreID, storeAdditive.StoreAdditiveID)
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

func (s *orderService) notifyLowStockIngredients(order *data.Order, stockMap map[uint]*data.StoreWarehouseStock) {
	for _, stock := range stockMap {
		if stock.Quantity <= stock.LowStockThreshold {
			details := &details.StoreWarehouseRunOutDetails{
				BaseNotificationDetails: details.BaseNotificationDetails{
					ID:           order.StoreID,
					FacilityName: order.Store.Name,
				},
				StockItem:   stock.Ingredient.Name,
				StockItemID: stock.IngredientID,
			}
			err := s.notificationService.NotifyStoreWarehouseRunOut(details)
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

	details := pdf.PDFReceiptDetails{
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

		for _, storeAdditive := range suborder.StoreAdditives {
			pdfSuborder.Additives = append(pdfSuborder.Additives, pdf.PDFAdditive{
				Name:  fmt.Sprintf("Additive #%d", storeAdditive.StoreAdditiveID),
				Price: storeAdditive.Price,
			})
		}

		details.SubOrders = append(details.SubOrders, pdfSuborder)
	}

	return pdf.GeneratePDFReceipt(details)
}

func (s *orderService) GetStatusesCount(storeID uint) (types.OrderStatusesCountDTO, error) {
	countsMap, err := s.orderRepo.GetStatusesCount(storeID)
	if err != nil {
		wrappedErr := fmt.Errorf("error couting statuses: %w", err)
		s.logger.Error(wrappedErr.Error())
		return types.OrderStatusesCountDTO{}, wrappedErr
	}

	dto := types.OrderStatusesCountDTO{
		PREPARING:   countsMap[data.OrderStatusPreparing],
		COMPLETED:   countsMap[data.OrderStatusCompleted],
		IN_DELIVERY: countsMap[data.OrderStatusInDelivery],
		DELIVERED:   countsMap[data.OrderStatusDelivered],
		CANCELLED:   countsMap[data.OrderStatusCancelled],
	}

	dto.ALL = dto.PREPARING + dto.COMPLETED + dto.IN_DELIVERY + dto.DELIVERED + dto.CANCELLED

	return dto, nil
}

func (s *orderService) ValidationResults(storeID uint, storeProductSizeIDs, storeAdditiveIDs []uint) (*orderValidationResults, error) {
	productPrices, productNames, err := ValidateProductSizes(storeID, storeProductSizeIDs, s.storeProductRepo)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	additivePrices, additiveNames, err := ValidateAdditives(storeID, storeAdditiveIDs, s.storeAdditiveRepo)
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

func ValidateAdditives(storeID uint, storeAdditiveIDs []uint, repo storeAdditives.StoreAdditiveRepository) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	additiveNames := make(map[uint]string)
	for _, id := range storeAdditiveIDs {
		storeAdditive, err := repo.GetStoreAdditiveByID(id, storeID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid additive ID: %d", id)
		}
		if storeAdditive == nil {
			return nil, nil, fmt.Errorf("additive with ID %d is nil", id)
		}
		if storeAdditive.Additive.Name == "" {
			return nil, nil, fmt.Errorf("additive with ID %d has an empty name", id)
		}

		price := storeAdditive.Additive.BasePrice
		if storeAdditive.StorePrice == nil {
			price = *storeAdditive.StorePrice
		}
		prices[id] = price
		additiveNames[id] = storeAdditive.Additive.Name
	}
	return prices, additiveNames, nil
}

func ValidateProductSizes(storeID uint, storeProductSizeIDs []uint, repo storeProducts.StoreProductRepository) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	productNames := make(map[uint]string)
	for _, id := range storeProductSizeIDs {
		storeProductSize, err := repo.GetStoreProductSizeById(storeID, id)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid store product size ID: %d", id)
		}
		if storeProductSize == nil {
			return nil, nil, fmt.Errorf("store product size with ID %d is nil", id)
		}
		if storeProductSize.StoreProduct.Product.Name == "" {
			return nil, nil, fmt.Errorf("product size with ID %d has an associated product with an empty name", id)
		}

		price := storeProductSize.ProductSize.BasePrice
		if storeProductSize.StorePrice == nil {
			price = *storeProductSize.StorePrice
		}

		prices[id] = price
		productNames[id] = storeProductSize.StoreProduct.Product.Name

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
