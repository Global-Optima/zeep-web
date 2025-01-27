package orders

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
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
	productRepo         product.ProductRepository
	additiveRepo        additives.AdditiveRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewOrderService(orderRepo OrderRepository, productRepo product.ProductRepository, additiveRepo additives.AdditiveRepository, notificationService notifications.NotificationService, logger *zap.SugaredLogger) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		productRepo:         productRepo,
		additiveRepo:        additiveRepo,
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
	productSizeIDs, additiveIDs := RetrieveIDs(*createOrderDTO)

	validations, err := s.ValidationResults(productSizeIDs, additiveIDs)
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
			ProductName: fmt.Sprintf("Product #%d", suborder.ProductSizeID),
			Price:       suborder.Price,
			Status:      string(suborder.Status),
		}

		for _, additive := range suborder.Additives {
			pdfSuborder.Additives = append(pdfSuborder.Additives, pdf.PDFAdditive{
				Name:  fmt.Sprintf("Additive #%d", additive.AdditiveID),
				Price: additive.Price,
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

func (s *orderService) ValidationResults(productSizeIDs, additiveIDs []uint) (*orderValidationResults, error) {
	productPrices, productNames, err := ValidateProductSizes(productSizeIDs, s.productRepo)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	additivePrices, additiveNames, err := ValidateAdditives(additiveIDs, s.additiveRepo)
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

func ValidateAdditives(additiveIDs []uint, repo additives.AdditiveRepository) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	additiveNames := make(map[uint]string)
	for _, id := range additiveIDs {
		additive, err := repo.GetAdditiveByID(id)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid additive ID: %d", id)
		}
		if additive == nil {
			return nil, nil, fmt.Errorf("additive with ID %d is nil", id)
		}
		if additive.Name == "" {
			return nil, nil, fmt.Errorf("additive with ID %d has an empty name", id)
		}

		if additive.BasePrice < 0 {
			return nil, nil, fmt.Errorf("additive with ID %d has an invalid base price", id)
		}
		prices[id] = additive.BasePrice
		additiveNames[id] = additive.Name
	}
	return prices, additiveNames, nil
}

func ValidateProductSizes(productSizeIDs []uint, repo product.ProductRepository) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	productNames := make(map[uint]string)
	for _, id := range productSizeIDs {
		productSize, err := repo.GetProductSizeById(id)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid product size ID: %d", id)
		}
		if productSize == nil {
			return nil, nil, fmt.Errorf("product size with ID %d is nil", id)
		}
		if productSize.Product.Name == "" {
			return nil, nil, fmt.Errorf("product size with ID %d has an associated product with an empty name", id)
		}

		prices[id] = productSize.BasePrice
		productNames[id] = productSize.Product.Name

	}
	return prices, productNames, nil
}

func RetrieveIDs(createOrderDTO types.CreateOrderDTO) ([]uint, []uint) {
	var productSizeIDs []uint
	var additiveIDs []uint
	for _, product := range createOrderDTO.Suborders {
		productSizeIDs = append(productSizeIDs, product.ProductSizeID)
		additiveIDs = append(additiveIDs, product.AdditivesIDs...)
	}
	return productSizeIDs, additiveIDs
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
