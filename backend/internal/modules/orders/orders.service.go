package orders

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
)

type OrderService interface {
	GetAllBaristaOrders(storeID uint, status *string) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.SuborderDTO, error)
	GetStatusesCount(storeID uint) (types.OrderStatusesCountDTO, error)
	CreateOrder(createOrderDTO *types.CreateOrderDTO) (*data.Order, error)
	CompleteSubOrder(subOrderID uint) error
	GeneratePDFReceipt(orderID uint) ([]byte, error)
	GetOrderBySubOrder(subOrderID uint) (*data.Order, error)
	GetOrderById(orderId uint) (types.OrderDTO, error)
}

type orderValidationResults struct {
	ProductPrices  map[uint]float64
	ProductNames   map[uint]string
	AdditivePrices map[uint]float64
	AdditiveNames  map[uint]string
}

type orderService struct {
	orderRepo    OrderRepository
	productRepo  product.ProductRepository
	additiveRepo additives.AdditiveRepository
	logger       *zap.SugaredLogger
}

func NewOrderService(orderRepo OrderRepository, productRepo product.ProductRepository, additiveRepo additives.AdditiveRepository, logger *zap.SugaredLogger) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		additiveRepo: additiveRepo,
		logger:       logger,
	}
}

// Get all orders with optional filtering by status
func (s *orderService) GetAllBaristaOrders(storeID uint, status *string) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetAllBaristaOrders(storeID, status)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	var orderDTOs []types.OrderDTO
	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}
	return orderDTOs, nil
}

// Get suborders for a specific order
func (s *orderService) GetSubOrders(orderID uint) ([]types.SuborderDTO, error) {
	suborders, err := s.orderRepo.GetSubOrdersByOrderID(orderID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	var subOrderDTOs []types.SuborderDTO
	for _, suborder := range suborders {
		subOrderDTOs = append(subOrderDTOs, types.ConvertSuborderToDTO(&suborder))
	}
	return subOrderDTOs, nil
}

// Create a new order
func (s *orderService) CreateOrder(createOrderDTO *types.CreateOrderDTO) (*data.Order, error) {
	productSizeIDs, additiveIDs := RetrieveIDs(*createOrderDTO)

	// Validate product sizes and additives
	validations, err := s.ValidationResults(productSizeIDs, additiveIDs)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert DTO to Order
	order, total := types.ConvertCreateOrderDTOToOrder(createOrderDTO, validations.ProductPrices, validations.AdditivePrices)
	order.Total = total

	// Save the order and related data to the database
	err = s.orderRepo.CreateOrder(&order)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &order, nil
}

// Complete a suborder
func (s *orderService) CompleteSubOrder(subOrderID uint) error {
	// Update suborder status
	err := s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusCompleted)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("failed to complete suborder: %w", err)
	}

	// Check if all suborders for the parent order are completed
	orderID, allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(subOrderID)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("failed to check suborders: %w", err)
	}

	// If all suborders are completed, determine the order status
	if allCompleted {
		order, err := s.orderRepo.GetOrderById(orderID)
		if err != nil {
			s.logger.Error(err.Error())
			return fmt.Errorf("failed to fetch order: %w", err)
		}

		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}

		err = s.orderRepo.UpdateOrderStatus(orderID, newStatus)

		if err != nil {
			s.logger.Error(err.Error())
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return nil
}

// Generate a PDF receipt for an order
func (s *orderService) GeneratePDFReceipt(orderID uint) ([]byte, error) {
	order, err := s.orderRepo.GetOrderById(orderID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
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

// Get statuses count for orders in a store
func (s *orderService) GetStatusesCount(storeID uint) (types.OrderStatusesCountDTO, error) {
	countsMap, err := s.orderRepo.GetStatusesCount(storeID)
	if err != nil {
		s.logger.Error(err.Error())
		return types.OrderStatusesCountDTO{}, err
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

// Validate product sizes and additives
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
		productSize, err := repo.GetProductSizeWithProduct(id)
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
		return nil, fmt.Errorf("failed to fetch order for suborder %d: %w", subOrderID, err)
	}

	return order, nil
}

func (s *orderService) GetOrderById(orderId uint) (types.OrderDTO, error) {
	order, err := s.orderRepo.GetOrderById(orderId)
	if err != nil {
		return types.OrderDTO{}, fmt.Errorf("failed to fetch order with ID %d: %w", orderId, err)
	}

	orderDTO := types.ConvertOrderToDTO(order)
	return orderDTO, nil
}
