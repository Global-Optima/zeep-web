package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/kafka"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/websockets"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
	"github.com/IBM/sarama"
)

type OrderService interface {
	GetAllOrders(status *string) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.OrderProductDTO, error)
	GetStatusesCount() (map[string]int64, error)
	GetSubOrderCount(orderID uint) (int64, error)

	CreateOrder(createOrderDTO *types.CreateOrderDTO) error
	CompleteSubOrder(subOrderID uint) error
	GeneratePDFReceipt(orderID uint) ([]byte, error)

	UpdateInventory(productID uint, quantity int) error
	GetLowStockProducts(threshold float64) ([]data.Product, error)
}

type orderService struct {
	repo          OrderRepository
	kafkaProducer *kafka.KafkaProducer
}

func NewOrderService(repo OrderRepository, kafkaProducer *kafka.KafkaProducer) OrderService {
	return &orderService{
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

func (s *orderService) CreateOrder(createOrderDTO *types.CreateOrderDTO) error {

	var productSizeIDs []uint
	var additiveIDs []uint
	for _, product := range createOrderDTO.Products {
		productSizeIDs = append(productSizeIDs, product.ProductSizeID)
		additiveIDs = append(additiveIDs, product.Additives...)
	}

	productPrices, err := ValidateProductSizes(productSizeIDs, s.repo)
	if err != nil {
		return err
	}
	additivePrices, err := ValidateAdditives(additiveIDs, s.repo)
	if err != nil {
		return err
	}

	order, total := types.ConvertCreateOrderDTOToOrder(createOrderDTO, productPrices, additivePrices)
	order.Total = total

	err = s.repo.CreateOrder(&order)
	if err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	orderEvent := types.OrderEvent{
		OrderID:   order.ID,
		StoreID:   order.StoreID,
		Status:    "pending",
		Timestamp: time.Now(),
	}
	for _, product := range order.OrderProducts {
		orderEvent.Items = append(orderEvent.Items, types.SubOrderEvent{
			SubOrderID: product.ID,
			Status:     "pending",
		})
	}

	eventData, err := json.Marshal(orderEvent)
	if err != nil {
		return fmt.Errorf("failed to serialize order event: %w", err)
	}

	kafkaConfig := config.GetConfig().Kafka
	err = s.kafkaProducer.SendMessage(kafkaConfig.Topics.ActiveOrders, fmt.Sprintf("%d", order.ID), string(eventData))
	if err != nil {
		return fmt.Errorf("failed to publish order to Kafka: %w", err)
	}

	websockets.GetHubInstance().Broadcast("orders", "order_created", eventData)

	return nil
}

func (s *orderService) CompleteSubOrder(subOrderID uint) error {
	subOrder, err := s.repo.GetSubOrderByID(subOrderID)
	if err != nil {
		return fmt.Errorf("suborder not found: %w", err)
	}

	orderEvent, err := s.GetActiveOrderEvent(subOrder.OrderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order event from Kafka: %w", err)
	}

	for i, item := range orderEvent.Items {
		if item.SubOrderID == subOrderID {
			orderEvent.Items[i].Status = "completed"
			break
		}
	}

	eventData, err := json.Marshal(orderEvent)
	if err != nil {
		return fmt.Errorf("failed to serialize updated order event: %w", err)
	}

	kafkaConfig := config.GetConfig().Kafka
	err = s.kafkaProducer.SendMessage(kafkaConfig.Topics.ActiveOrders, fmt.Sprintf("%d", subOrder.OrderID), string(eventData))
	if err != nil {
		return fmt.Errorf("failed to publish updated order event to Kafka: %w", err)
	}

	allCompleted := true
	for _, item := range orderEvent.Items {
		if item.Status != "completed" {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		orderEvent.Status = "completed"
		orderEvent.Timestamp = time.Now()

		eventData, err = json.Marshal(orderEvent)
		if err != nil {
			return fmt.Errorf("failed to serialize completed order event: %w", err)
		}

		err = s.kafkaProducer.SendMessage(kafkaConfig.Topics.ActiveOrders, fmt.Sprintf("%d", subOrder.OrderID), string(eventData))
		if err != nil {
			return fmt.Errorf("failed to publish completed order event to Kafka: %w", err)
		}
	}

	websockets.GetHubInstance().Broadcast("orders", "suborder_completed", eventData)

	return nil
}

func (s *orderService) GetActiveOrderEvent(orderID uint) (*types.OrderEvent, error) {
	kafkaConfig := config.GetConfig().Kafka

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(kafkaConfig.Brokers, kafkaConfig.ConsumerGroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	defer consumer.Close()

	resultChan := make(chan *types.OrderEvent, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	handler := kafka.NewKafkaHandler(func(topic, key, value string) error {
		var event types.OrderEvent
		if err := json.Unmarshal([]byte(value), &event); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		if event.OrderID == orderID {
			resultChan <- &event
			cancel()
		}
		return nil
	})

	go func() {
		if err := consumer.Consume(ctx, []string{kafkaConfig.Topics.ActiveOrders}, handler); err != nil {
			fmt.Printf("Error consuming Kafka topic: %v\n", err)
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout fetching order event for order ID %d", orderID)
	}
}

func (s *orderService) GetAllOrders(status *string) ([]types.OrderDTO, error) {
	orders, err := s.repo.GetAllOrders(status)
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
	subOrders, err := s.repo.GetSubOrders(orderID)
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
	return s.repo.GetStatusesCount()
}

func (s *orderService) GetSubOrderCount(orderID uint) (int64, error) {
	return s.repo.GetSubOrderCount(orderID)
}

func (s *orderService) GeneratePDFReceipt(orderID uint) ([]byte, error) {
	order, err := s.repo.GetOrderById(orderID)
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
		productSizeLabel, err := s.repo.GetProductSizeLabel(product.ProductSizeID)
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

func (s *orderService) UpdateInventory(productID uint, quantity int) error {
	return s.repo.UpdateInventory(productID, quantity)
}

func (s *orderService) GetLowStockProducts(threshold float64) ([]data.Product, error) {
	return s.repo.GetLowStockProducts(threshold)
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
