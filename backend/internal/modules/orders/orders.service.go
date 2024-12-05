package orders

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/kafka"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
	"github.com/google/uuid"
)

var Logger = logger.GetInstance()

type OrderService interface {
	GetAllOrders(storeID uint, status *string, limit int, offset int) ([]types.OrderDTO, error)
	GetSubOrders(storeID, orderID uint) ([]types.SubOrderEvent, error)
	GetStatusesCount(storeID uint) (map[string]int64, error)
	GetSubOrderCount(orderID, storeID uint) (int64, error)

	CreateOrder(createOrderDTO *types.CreateOrderDTO) (*uint, error)
	CompleteSubOrder(subOrderID, orderID, storeID uint) error
	GeneratePDFReceipt(orderID uint) ([]byte, error)
	GetActiveOrderEvent(orderID uint, storeID uint) (*types.OrderEvent, error)

	GetLowStockIngredients(threshold float64) ([]data.Ingredient, error)
}

type orderService struct {
	orderRepo      OrderRepository
	subOrderRepo   SubOrderRepository
	productRepo    product.ProductRepository
	additiveRepo   additives.AdditiveRepository
	kafkaManager   *kafka.KafkaManager
	ordersNotifier *OrdersNotifier
}

func NewOrderService(orderRepo OrderRepository, subOrderRepo SubOrderRepository, productRepo product.ProductRepository, additiveRepo additives.AdditiveRepository, kafkaManager *kafka.KafkaManager, ordersNotifier *OrdersNotifier) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		subOrderRepo:   subOrderRepo,
		productRepo:    productRepo,
		additiveRepo:   additiveRepo,
		kafkaManager:   kafkaManager,
		ordersNotifier: ordersNotifier,
	}
}

func (s *orderService) GetAllOrders(storeID uint, status *string, limit int, offset int) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetAllOrders(storeID, status, limit, offset)
	if err != nil {
		return nil, err
	}

	var orderDTOs []types.OrderDTO
	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}
	return orderDTOs, nil
}

func (s *orderService) CreateOrder(createOrderDTO *types.CreateOrderDTO) (*uint, error) {
	productSizeIDs, additiveIDs := RetrieveIDs(*createOrderDTO)

	productPrices, productNames, err := ValidateProductSizes(productSizeIDs, s.productRepo)
	if err != nil {
		return nil, err
	}
	additivePrices, additiveNames, err := ValidateAdditives(additiveIDs, s.additiveRepo)
	if err != nil {
		return nil, err
	}

	order, total := types.ConvertCreateOrderDTOToOrder(createOrderDTO, productPrices, additivePrices)
	order.Total = total

	err = s.orderRepo.CreateOrder(&order)
	if err != nil {
		return nil, fmt.Errorf("failed to save order to database: %w", err)
	}

	orderEvent := types.OrderEvent{
		ID:                order.ID,
		CustomerName:      createOrderDTO.CustomerName,
		ETA:               calculateETA(createOrderDTO.OrderType),
		OrderType:         createOrderDTO.OrderType,
		StoreID:           order.StoreID,
		SubOrdersQuantity: 0,
		Status:            string(data.OrderStatusPending),
		Timestamp:         time.Now(),
	}

	for _, product := range order.OrderProducts {
		for i := 0; i < product.Quantity; i++ {
			ID := uuid.New().ID()

			subOrderAdditives := []types.AdditiveDetail{}
			for _, additive := range product.Additives {
				subOrderAdditives = append(subOrderAdditives, types.AdditiveDetail{
					AdditiveID: additive.AdditiveID,
					Name:       additiveNames[additive.AdditiveID],
					Price:      additivePrices[additive.AdditiveID],
				})
			}

			subOrder := types.SubOrderEvent{
				ID:            uint(ID),
				SubOrderID:    product.ID,
				ProductSizeID: product.ProductSizeID,
				ProductName:   productNames[product.ProductSizeID],
				Additives:     subOrderAdditives,
				ETA:           calculateETA(createOrderDTO.OrderType),
				Status:        data.OrderStatusPending,
			}
			orderEvent.Items = append(orderEvent.Items, subOrder)
			orderEvent.SubOrdersQuantity++
		}
	}

	err = s.kafkaManager.PublishOrderEvent(s.kafkaManager.Topics.ActiveOrders, orderEvent.StoreID, orderEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to publish order to Kafka: %w", err)
	}

	s.ordersNotifier.NotifyNewOrder(order.ID, order.StoreID, orderEvent)

	return &order.ID, nil
}

func RetrieveIDs(createOrderDTO types.CreateOrderDTO) ([]uint, []uint) {
	var productSizeIDs []uint
	var additiveIDs []uint
	for _, product := range createOrderDTO.OrderItems {
		productSizeIDs = append(productSizeIDs, product.ProductSizeID)
		additiveIDs = append(additiveIDs, product.AdditivesIDs...)
	}
	return productSizeIDs, additiveIDs
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

func calculateETA(orderType string) time.Time {
	if orderType == "Доставка" {
		return time.Now().Add(30 * time.Minute)
	}
	return time.Now().Add(15 * time.Minute)
}

func GetCachedProductSize(productSizeID uint, repo product.ProductRepository) (*data.ProductSize, error) {
	cache := utils.GetCacheInstance()
	cacheKey := fmt.Sprintf("product_size:%d", productSizeID)

	var productSize data.ProductSize
	if err := cache.Get(cacheKey, &productSize); err == nil {
		return &productSize, nil
	}

	productSizePtr, err := repo.GetProductSizeWithProduct(productSizeID)
	if err != nil {
		return nil, fmt.Errorf("product size with ID %d not found", productSizeID)
	}

	_ = cache.Set(cacheKey, productSizePtr, 30*time.Minute)
	return productSizePtr, nil
}

func GetCachedAdditive(additiveID uint, repo additives.AdditiveRepository) (*data.Additive, error) {
	cache := utils.GetCacheInstance()
	cacheKey := fmt.Sprintf("additive:%d", additiveID)

	var additive data.Additive
	if err := cache.Get(cacheKey, &additive); err == nil {
		return &additive, nil
	}

	additivePtr, err := repo.GetAdditiveByID(additiveID)
	if err != nil {
		return nil, fmt.Errorf("additive with ID %d not found", additiveID)
	}

	_ = cache.Set(cacheKey, additivePtr, 10*time.Minute)
	return additivePtr, nil
}

func (s *orderService) CompleteSubOrder(subOrderID, orderID, storeID uint) error {
	orderEvent, err := s.kafkaManager.FetchOrderEvent(orderID, storeID, string(s.kafkaManager.Topics.ActiveOrders))
	if err != nil {
		return fmt.Errorf("failed to fetch order event from Kafka: %w", err)
	}

	subOrderFound := false
	for i, item := range orderEvent.Items {
		if item.ID == subOrderID {
			orderEvent.Items[i].Status = data.OrderStatusCompleted
			subOrderFound = true
			break
		}
	}

	if !subOrderFound {
		return fmt.Errorf("suborder %d not found in order %d", subOrderID, orderID)
	}

	allCompleted := true
	for _, item := range orderEvent.Items {
		if item.Status != data.OrderStatusCompleted {
			allCompleted = false
			break
		}
	}

	order, err := s.orderRepo.GetOrderByOrderId(orderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order from db %s", err.Error())
	}

	var topic kafka.Topic
	if allCompleted {
		orderEvent.Status = string(data.OrderStatusCompleted)
		topic = s.kafkaManager.Topics.CompletedOrders

		err = s.orderRepo.UpdateOrderStatus(order.ID, storeID, data.OrderStatusCompleted)
		if err != nil {
			return fmt.Errorf("failed to update order status in database: %w", err)
		}

		s.ordersNotifier.NotifyOrderCompleted(order.ID, storeID, orderEvent)
	} else {
		topic = s.kafkaManager.Topics.ActiveOrders
		s.ordersNotifier.NotifySubOrderCompleted(order.ID, subOrderID, storeID, orderEvent)
	}

	err = s.kafkaManager.PublishOrderEvent(topic, storeID, *orderEvent)
	if err != nil {
		return fmt.Errorf("failed to publish order to Kafka: %w", err)
	}

	return nil
}

func (s *orderService) GetSubOrders(storeID, orderID uint) ([]types.SubOrderEvent, error) {
	orderEvent, err := s.kafkaManager.FetchOrderEvent(orderID, storeID, string(s.kafkaManager.Topics.ActiveOrders))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OrderEvent from Kafka: %w", err)
	}

	if orderEvent.StoreID != storeID {
		return nil, fmt.Errorf("order %d does not belong to store %d", orderID, storeID)
	}

	var subOrders []types.SubOrderEvent
	subOrders = append(subOrders, orderEvent.Items...)
	return subOrders, nil
}

func (s *orderService) GeneratePDFReceipt(orderID uint) ([]byte, error) {
	order, err := s.orderRepo.GetOrderByOrderId(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}

	details := pdf.PDFReceiptDetails{
		OrderID:   order.ID,
		StoreID:   order.StoreID,
		OrderDate: order.CreatedAt.Format("2006-01-02 15:04:05"),
		Total:     order.Total,
	}

	for _, product := range order.OrderProducts {
		productSize, err := s.productRepo.GetProductSizeWithProduct(product.ProductSizeID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch size label: %w", err)
		}

		pdfProduct := pdf.PDFProduct{
			ProductName: fmt.Sprintf("Product #%d", product.ProductSizeID),
			Size:        productSize.Name,
			Quantity:    product.Quantity,
			Price:       product.Price,
		}

		for _, additive := range product.Additives {
			pdfProduct.Additives = append(pdfProduct.Additives, pdf.PDFAdditive{
				Name:  fmt.Sprintf("Additive #%d", additive.AdditiveID),
				Price: additive.Price,
			})
		}

		details.Products = append(details.Products, pdfProduct)
	}

	return pdf.GeneratePDFReceipt(details)
}

func (s *orderService) GetActiveOrderEvent(orderID, storeID uint) (*types.OrderEvent, error) {
	event, err := s.kafkaManager.FetchOrderEvent(orderID, storeID, string(s.kafkaManager.Topics.ActiveOrders))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active order event: %w", err)
	}
	return event, nil
}

func (s *orderService) GetStatusesCount(storeID uint) (map[string]int64, error) {
	return s.orderRepo.GetStatusesCount(storeID)
}

func (s *orderService) GetSubOrderCount(orderID, storeID uint) (int64, error) {
	return s.subOrderRepo.GetSubOrderCount(orderID)
}

func (s *orderService) GetLowStockIngredients(threshold float64) ([]data.Ingredient, error) {
	return s.orderRepo.GetLowStockIngredients(threshold)
}
