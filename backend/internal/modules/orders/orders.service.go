package orders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/kafka"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
)

type OrderService interface {
	GetAllOrders(status *string) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.OrderProductDTO, error)
	GetStatusesCount() (map[string]int64, error)
	GetSubOrderCount(orderID uint) (int64, error)

	CreateOrder(createOrderDTO *types.CreateOrderDTO) error
	CompleteSubOrder(subOrderID uint) error
	GeneratePDFReceipt(orderID uint) ([]byte, error)

	GetLowStockIngredients(threshold float64) ([]data.Ingredient, error)
}

type orderService struct {
	orderRepo      OrderRepository
	subOrderRepo   SubOrderRepository
	kafkaManager   *kafka.KafkaManager
	ordersNotifier *OrdersNotifier
}

func NewOrderService(orderRepo OrderRepository, subOrderRepo SubOrderRepository, kafkaManager *kafka.KafkaManager, ordersNotifier *OrdersNotifier) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		subOrderRepo:   subOrderRepo,
		kafkaManager:   kafkaManager,
		ordersNotifier: ordersNotifier,
	}
}

func (s *orderService) CreateOrder(createOrderDTO *types.CreateOrderDTO) error {
	var productSizeIDs []uint
	var additiveIDs []uint
	for _, product := range createOrderDTO.Products {
		productSizeIDs = append(productSizeIDs, product.ProductSizeID)
		additiveIDs = append(additiveIDs, product.Additives...)
	}

	productPrices, err := ValidateProductSizes(productSizeIDs, s.orderRepo)
	if err != nil {
		return err
	}
	additivePrices, err := ValidateAdditives(additiveIDs, s.orderRepo)
	if err != nil {
		return err
	}

	order, total := types.ConvertCreateOrderDTOToOrder(createOrderDTO, productPrices, additivePrices)
	order.Total = total

	Logger.Debug(fmt.Sprintf("%+v", order))
	err = s.orderRepo.CreateOrder(&order)
	if err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	orderEvent := types.OrderEvent{
		OrderID:   order.ID,
		StoreID:   order.StoreID,
		Status:    types.OrderStatusPending,
		Timestamp: time.Now(),
	}
	for _, product := range order.OrderProducts {
		orderEvent.Items = append(orderEvent.Items, types.SubOrderEvent{
			SubOrderID: product.ID,
			Status:     types.OrderStatusPending,
		})
	}

	eventData, err := json.Marshal(orderEvent)
	if err != nil {
		return fmt.Errorf("failed to serialize order event: %w", err)
	}

	err = s.kafkaManager.PublishEvent(s.kafkaManager.Topics.ActiveOrders, fmt.Sprintf("%d", order.ID), eventData)
	if err != nil {
		return fmt.Errorf("failed to publish order to Kafka: %w", err)
	}

	s.ordersNotifier.NotifyNewOrder(order.ID, orderEvent)

	return nil
}

func (s *orderService) CompleteSubOrder(subOrderID uint) error {
	subOrder, err := s.subOrderRepo.GetSubOrderByID(subOrderID)
	if err != nil {
		return fmt.Errorf("suborder not found: %w", err)
	}

	orderEvent, err := s.kafkaManager.GetOrderEvent(s.kafkaManager.Topics.ActiveOrders, subOrder.OrderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order event from Kafka: %w", err)
	}

	for i, item := range orderEvent.Items {
		if item.SubOrderID == subOrderID {
			orderEvent.Items[i].Status = types.OrderStatusCompleted
			break
		}
	}

	allCompleted := true
	for _, item := range orderEvent.Items {
		if item.Status != types.OrderStatusCompleted {
			allCompleted = false
			break
		}
	}

	orderEvent.Timestamp = time.Now()
	var topic string
	if allCompleted {
		orderEvent.Status = types.OrderStatusCompleted
		topic = s.kafkaManager.Topics.CompletedOrders

		err = s.orderRepo.UpdateOrderStatus(subOrder.OrderID, string(types.OrderStatusCompleted))
		if err != nil {
			return fmt.Errorf("failed to update order status in database: %w", err)
		}

		s.ordersNotifier.NotifyOrderCompleted(subOrder.OrderID, orderEvent)
	} else {
		topic = s.kafkaManager.Topics.ActiveOrders
		s.ordersNotifier.NotifySubOrderCompleted(subOrder.OrderID, subOrderID, orderEvent)
	}

	err = s.kafkaManager.PublishEvent(topic, fmt.Sprintf("%d", subOrder.OrderID), orderEvent)
	if err != nil {
		return fmt.Errorf("failed to publish order to Kafka: %w", err)
	}

	return nil
}

func (s *orderService) GetAllOrders(status *string) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetAllOrders(status)
	if err != nil {
		return nil, err
	}

	var orderDTOs []types.OrderDTO
	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}
	return orderDTOs, nil
}

func (s *orderService) GetSubOrders(orderID uint) ([]types.OrderProductDTO, error) {
	subOrders, err := s.subOrderRepo.GetSubOrders(orderID)
	if err != nil {
		return nil, err
	}

	var subOrderDTOs []types.OrderProductDTO
	for _, subOrder := range subOrders {
		subOrderDTOs = append(subOrderDTOs, types.ConvertOrderProductToDTO(&subOrder))
	}
	return subOrderDTOs, nil
}

func (s *orderService) GetStatusesCount() (map[string]int64, error) {
	return s.orderRepo.GetStatusesCount()
}

func (s *orderService) GetSubOrderCount(orderID uint) (int64, error) {
	return s.subOrderRepo.GetSubOrderCount(orderID)
}

func (s *orderService) GeneratePDFReceipt(orderID uint) ([]byte, error) {
	order, err := s.orderRepo.GetOrderById(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}

	details := pdf.PDFReceiptDetails{
		OrderID:   order.ID,
		StoreID:   *order.StoreID,
		OrderDate: order.OrderDate.Format("2006-01-02 15:04:05"),
		Total:     order.Total,
	}

	for _, product := range order.OrderProducts {
		productSizeLabel, err := s.orderRepo.GetProductSizeLabel(product.ProductSizeID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch size label: %w", err)
		}

		pdfProduct := pdf.PDFProduct{
			ProductName: fmt.Sprintf("Product #%d", product.ProductSizeID),
			Size:        productSizeLabel,
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

func (s *orderService) GetLowStockIngredients(threshold float64) ([]data.Ingredient, error) {
	return s.orderRepo.GetLowStockIngredients(threshold)
}

func ValidateProductSizes(productSizeIDs []uint, repo OrderRepository) (map[uint]float64, error) {
	prices := make(map[uint]float64)
	for _, id := range productSizeIDs {
		productSize, err := repo.GetProductSizeByID(id)
		if err != nil {
			return nil, fmt.Errorf("invalid product size ID: %d", id)
		}
		prices[id] = productSize.BasePrice
	}
	return prices, nil
}

func ValidateAdditives(additiveIDs []uint, repo OrderRepository) (map[uint]float64, error) {
	prices := make(map[uint]float64)
	for _, id := range additiveIDs {
		additive, err := repo.GetAdditiveByID(id)
		if err != nil {
			return nil, fmt.Errorf("invalid additive ID: %d", id)
		}
		prices[id] = additive.BasePrice
	}
	return prices, nil
}
