package orders

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAllBaristaOrders(storeID uint, status *string) ([]data.Order, error)
	GetOrderById(orderID uint) (*data.Order, error)
	CreateOrder(order *data.Order) error
	UpdateOrderStatus(orderID uint, status data.OrderStatus) error
	DeleteOrder(orderID uint) error

	GetStatusesCount(storeID uint) (map[data.OrderStatus]int64, error)

	GetSubOrdersByOrderID(orderID uint) ([]data.Suborder, error)
	UpdateSubOrderStatus(subOrderID uint, status data.SubOrderStatus) error
	AddSubOrderAdditive(subOrderID uint, additive *data.SuborderAdditive) error
	CheckAllSubordersCompleted(subOrderID uint) (uint, bool, error)
	GetOrderBySubOrderID(subOrderID uint) (*data.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetAllBaristaOrders(storeID uint, status *string) ([]data.Order, error) {
	var orders []data.Order

	now := time.Now().UTC()

	// Calculate the start and end of the range (yesterday to tomorrow)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)

	query := r.db.Preload("Suborders.ProductSize").
		Preload("Suborders.ProductSize.Product").
		Preload("Suborders.Additives.Additive").
		Where("store_id = ?", storeID).
		Where("created_at >= ? AND created_at <= ?", startOfToday, endOfToday)

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at asc").Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch barista orders: %w", err)
	}

	return orders, nil
}

// Get a specific order for a store
func (r *orderRepository) GetOrderById(orderId uint) (*data.Order, error) {
	var order data.Order
	err := r.db.Preload("Suborders.ProductSize").
		Preload("Suborders.ProductSize.Product").
		Preload("Suborders.Additives.Additive").
		Where("id = ?", orderId).
		First(&order).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch order with ID %d: %w", orderId, err)
	}

	return &order, nil
}

// Create a new order along with its suborders and additives
func (r *orderRepository) CreateOrder(order *data.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})
}

// Update the status of an order
func (r *orderRepository) UpdateOrderStatus(orderID uint, status data.OrderStatus) error {
	return r.db.Model(&data.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

// Delete an order
func (r *orderRepository) DeleteOrder(orderID uint) error {
	return r.db.Delete(&data.Order{}, orderID).Error
}

// Get count of orders grouped by status
func (r *orderRepository) GetStatusesCount(storeID uint) (map[data.OrderStatus]int64, error) {
	var counts []struct {
		Status data.OrderStatus
		Count  int64
	}

	// Ensure the query groups and counts correctly
	err := r.db.Model(&data.Order{}).
		Select("status, COUNT(*) as count").
		Where("store_id = ?", storeID).
		Group("status").
		Find(&counts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch statuses count: %w", err)
	}

	// Initialize statusCount with all possible statuses to avoid missing keys
	statusCount := map[data.OrderStatus]int64{
		data.OrderStatusPreparing:  0,
		data.OrderStatusCompleted:  0,
		data.OrderStatusInDelivery: 0,
		data.OrderStatusDelivered:  0,
		data.OrderStatusCancelled:  0,
	}

	// Populate statusCount with actual counts
	for _, c := range counts {
		statusCount[c.Status] = c.Count
	}

	return statusCount, nil
}

// Get all suborders for a specific order
func (r *orderRepository) GetSubOrdersByOrderID(orderID uint) ([]data.Suborder, error) {
	var suborders []data.Suborder
	err := r.db.Preload("Additives").Where("order_id = ?", orderID).Find(&suborders).Error
	if err != nil {
		return nil, err
	}
	return suborders, nil
}

// Update the status of a specific suborder
func (r *orderRepository) UpdateSubOrderStatus(suborderID uint, status data.SubOrderStatus) error {
	return r.db.Model(&data.Suborder{}).
		Where("id = ?", suborderID).
		Update("status", status).Error
}

// Add an additive to a suborder
func (r *orderRepository) AddSubOrderAdditive(suborderID uint, additive *data.SuborderAdditive) error {
	additive.SuborderID = suborderID
	return r.db.Create(additive).Error
}

func (r *orderRepository) CheckAllSubordersCompleted(subOrderID uint) (uint, bool, error) {
	var suborder data.Suborder
	err := r.db.Select("order_id").Where("id = ?", subOrderID).First(&suborder).Error
	if err != nil {
		return 0, false, fmt.Errorf("failed to fetch suborder: %w", err)
	}

	orderID := suborder.OrderID

	var count int64
	err = r.db.Model(&data.Suborder{}).
		Where("order_id = ? AND status != ?", orderID, data.SubOrderStatusCompleted).
		Count(&count).Error
	if err != nil {
		return 0, false, fmt.Errorf("failed to count suborders: %w", err)
	}

	return orderID, count == 0, nil
}

func (r *orderRepository) GetOrderBySubOrderID(subOrderID uint) (*data.Order, error) {
	var order data.Order

	err := r.db.
		Joins("JOIN suborders ON suborders.order_id = orders.id").
		Where("suborders.id = ?", subOrderID).
		Preload("Suborders.ProductSize").
		Preload("Suborders.ProductSize.Product").
		Preload("Suborders.Additives.Additive").
		First(&order).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch order for suborder %d: %w", subOrderID, err)
	}

	return &order, nil
}
