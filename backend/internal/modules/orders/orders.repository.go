package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrders(filter types.OrdersFilterQuery) ([]data.Order, error)
	GetAllBaristaOrders(filter types.OrdersTimeZoneFilter) ([]data.Order, error)
	GetOrderById(orderID uint) (*data.Order, error)
	CreateOrder(order *data.Order) error
	UpdateOrderStatus(orderID uint, status data.OrderStatus) error
	DeleteOrder(orderID uint) error

	GetStatusesCount(filter types.OrdersTimeZoneFilter) (map[data.OrderStatus]int64, error)

	GetSubOrdersByOrderID(orderID uint) ([]data.Suborder, error)
	UpdateSubOrderStatus(subOrderID uint, status data.SubOrderStatus) error
	AddSubOrderAdditive(subOrderID uint, additive *data.SuborderAdditive) error
	CheckAllSubordersCompleted(orderID, subOrderID uint) (bool, error)
	GetOrderBySubOrderID(subOrderID uint) (*data.Order, error)

	GetOrderDetails(orderID uint, filter *contexts.StoreContextFilter) (*data.Order, error)
	GetOrdersForExport(filter *types.OrdersExportFilterQuery) ([]data.Order, error)
	GetSuborderByID(suborderID uint) (*data.Suborder, error)

	CalculateFrozenStock(storeID uint) (map[uint]float64, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *data.Order) error {
	const workingHours = 16 // Working hours for a cafe

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the store's orders for consistency
		if err := tx.Exec(`
			SELECT 1
			FROM orders
			WHERE store_id = ?
			FOR UPDATE
		`, order.StoreID).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to lock rows for store %d: %w", order.StoreID, err)
			}
		}

		var lastOrder struct {
			DisplayNumber int
			CreatedAt     time.Time
		}

		// Fetch the latest order for the store
		if err := tx.Model(&data.Order{}).
			Where("store_id = ?", order.StoreID).
			Order("created_at DESC").
			Limit(1).
			Select("display_number, created_at").
			Scan(&lastOrder).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to fetch last order: %w", err)
			}
		}

		var nextDisplayNumber int
		now := time.Now()

		// Check if the elapsed time since the last order exceeds the working hours
		if now.Sub(lastOrder.CreatedAt) > time.Duration(workingHours)*time.Hour {
			// Reset display number if more than working hours have passed
			nextDisplayNumber = 1
		} else {
			// Increment display number if within the same working period
			nextDisplayNumber = lastOrder.DisplayNumber + 1
		}

		// Ensure display number doesn't exceed the maximum
		if nextDisplayNumber > 999 {
			nextDisplayNumber = 1
		}
		order.DisplayNumber = nextDisplayNumber

		// Create the new order
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})
}

func (r *orderRepository) GetOrders(filter types.OrdersFilterQuery) ([]data.Order, error) {
	var orders []data.Order

	query := r.db.
		Preload("Suborders.StoreProductSize.ProductSize.Unit").
		Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Order("created_at DESC")

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"customer_name ILIKE ?",
			searchTerm,
		)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.StoreID != nil {
		query = query.Where("store_id = ?", *filter.StoreID)
	}

	var err error
	query, err = utils.ApplyPagination(query, filter.Pagination, &data.Order{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, nil
}

func (r *orderRepository) GetAllBaristaOrders(filter types.OrdersTimeZoneFilter) ([]data.Order, error) {
	var orders []data.Order

	// Validate StoreID first
	if filter.StoreID == nil {
		return nil, fmt.Errorf("storeID is required")
	}

	// Get correct timezone location
	location, err := getTimeZoneLocation(filter)
	if err != nil {
		return nil, err
	}

	// Get the current date in the specified timezone
	now := time.Now().In(location)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, location)

	// Convert start and end of today **back to UTC** for database querying
	startOfTodayUTC := startOfToday.UTC()
	endOfTodayUTC := endOfToday.UTC()

	// Query orders for the given store within the correct time range
	err = r.db.
		Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.StoreProductSize.ProductSize.Unit").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Where("store_id = ?", *filter.StoreID).
		Where("created_at BETWEEN ? AND ?", startOfTodayUTC, endOfTodayUTC).
		Order("created_at ASC").
		Find(&orders).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch barista orders: %w", err)
	}

	return orders, nil
}

func getTimeZoneLocation(filter types.OrdersTimeZoneFilter) (*time.Location, error) {
	if filter.TimeZoneLocation != nil && *filter.TimeZoneLocation != "" {
		location, err := time.LoadLocation(*filter.TimeZoneLocation)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone location: %w", err)
		}

		return location, nil
	}

	if filter.TimeZoneOffset != nil {
		offsetSeconds := int(*filter.TimeZoneOffset) * 60
		return time.FixedZone("Custom Offset", offsetSeconds), nil
	}

	// Default to UTC if no timezone is provided
	return time.UTC, nil
}

func (r *orderRepository) GetStatusesCount(filter types.OrdersTimeZoneFilter) (map[data.OrderStatus]int64, error) {
	var counts []struct {
		Status data.OrderStatus
		Count  int64
	}

	var location *time.Location
	var now time.Time

	// Determine which location to use.
	if filter.TimeZoneLocation != nil && *filter.TimeZoneLocation != "" {
		var err error
		location, err = time.LoadLocation(*filter.TimeZoneLocation)
		if err != nil {
			return nil, fmt.Errorf("failed to load timezone: %w", err)
		}
		now = time.Now().In(location)
	} else if filter.TimeZoneOffset != nil {
		offsetSeconds := int(*filter.TimeZoneOffset) * 60
		location = time.FixedZone("offset", offsetSeconds)
		now = time.Now().UTC().In(location)
	} else {
		location = time.UTC
		now = time.Now().UTC()
	}

	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, location)

	startOfTodayUTC := startOfToday.UTC()
	endOfTodayUTC := endOfToday.UTC()

	if filter.StoreID == nil {
		return nil, fmt.Errorf("storeID is required")
	}

	query := r.db.Model(&data.Order{}).
		Select("status, COUNT(*) as count").
		Where("store_id = ?", filter.StoreID).
		Where("created_at >= ? AND created_at <= ?", startOfTodayUTC, endOfTodayUTC).
		Group("status")

	err := query.Find(&counts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch statuses count: %w", err)
	}

	statusCount := map[data.OrderStatus]int64{
		data.OrderStatusPending:    0,
		data.OrderStatusPreparing:  0,
		data.OrderStatusCompleted:  0,
		data.OrderStatusInDelivery: 0,
		data.OrderStatusDelivered:  0,
		data.OrderStatusCancelled:  0,
	}

	for _, c := range counts {
		statusCount[c.Status] = c.Count
	}

	return statusCount, nil
}

func (r *orderRepository) GetOrderById(orderId uint) (*data.Order, error) {
	var order data.Order
	err := r.db.Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Preload("Suborders.StoreProductSize.ProductSize.Unit").
		Preload("Store").
		Where("id = ?", orderId).
		First(&order).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch order with ID %d: %w", orderId, err)
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(orderID uint, status data.OrderStatus) error {
	return r.db.Model(&data.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *orderRepository) DeleteOrder(orderID uint) error {
	return r.db.Delete(&data.Order{}, orderID).Error
}

func (r *orderRepository) GetSubOrdersByOrderID(orderID uint) ([]data.Suborder, error) {
	var suborders []data.Suborder
	err := r.db.Preload("SuborderAdditives").Where("order_id = ?", orderID).Find(&suborders).Error
	if err != nil {
		return nil, err
	}
	return suborders, nil
}

func (r *orderRepository) UpdateSubOrderStatus(suborderID uint, status data.SubOrderStatus) error {
	return r.db.Model(&data.Suborder{}).
		Where("id = ?", suborderID).
		Update("status", status).Error
}

func (r *orderRepository) AddSubOrderAdditive(suborderID uint, additive *data.SuborderAdditive) error {
	additive.SuborderID = suborderID
	return r.db.Create(additive).Error
}

func (r *orderRepository) CheckAllSubordersCompleted(orderID, subOrderID uint) (bool, error) {
	var suborder data.Suborder
	err := r.db.Select("order_id").Where("id = ?", subOrderID).First(&suborder).Error
	if err != nil {
		return false, fmt.Errorf("failed to fetch suborder: %w", err)
	}

	var count int64
	err = r.db.Model(&data.Suborder{}).
		Where("order_id = ? AND status != ?", orderID, data.SubOrderStatusCompleted).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to count suborders: %w", err)
	}

	return count == 0, nil
}

func (r *orderRepository) GetOrderBySubOrderID(subOrderID uint) (*data.Order, error) {
	var order data.Order

	err := r.db.
		Preload("Suborders").
		Preload("Suborders.StoreProductSize").
		Preload("Suborders.StoreProductSize.ProductSize").
		Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.StoreProductSize.ProductSize.Unit").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Joins("JOIN suborders ON suborders.order_id = orders.id").
		Where("suborders.id = ?", subOrderID).
		First(&order).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch order for suborder %d: %w", subOrderID, err)
	}

	return &order, nil
}

func (r *orderRepository) GetOrderDetails(orderID uint, filter *contexts.StoreContextFilter) (*data.Order, error) {
	var order data.Order

	query := r.db.Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.StoreProductSize.ProductSize.Unit").
		Preload("DeliveryAddress").
		Where(&data.Order{
			BaseEntity: data.BaseEntity{
				ID: orderID,
			},
		})

	if filter != nil {
		if filter.StoreID != nil {
			query.Where(&data.Order{StoreID: *filter.StoreID})
		}

		if filter.FranchiseeID != nil {
			query.Joins("JOIN stores ON stores.id = orders.store_id").
				Where(&data.Order{
					Store: data.Store{
						FranchiseeID: filter.FranchiseeID,
					},
				})
		}
	}

	err := query.First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}

	return &order, nil
}

func (r *orderRepository) GetOrdersForExport(filter *types.OrdersExportFilterQuery) ([]data.Order, error) {
	var orders []data.Order
	query := r.db.Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Preload("Store").
		Preload("DeliveryAddress")

	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.StoreID != nil {
		query = query.Where("store_id = ?", *filter.StoreID)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetSuborderByID(suborderID uint) (*data.Suborder, error) {
	var suborder data.Suborder

	err := r.db.
		Preload("StoreProductSize").
		Preload("StoreProductSize.ProductSize").
		Preload("StoreProductSize.ProductSize.Product").
		Preload("StoreProductSize.ProductSize.Unit").
		Preload("SuborderAdditives").
		Preload("SuborderAdditives.StoreAdditive").
		Preload("SuborderAdditives.StoreAdditive.Additive").
		Where("id = ?", suborderID).
		First(&suborder).Error

	if err != nil {
		return nil, err
	}
	return &suborder, nil
}

func (r *orderRepository) CalculateFrozenStock(storeID uint) (map[uint]float64, error) {
	frozenStock := make(map[uint]float64)

	activeOrders, err := r.loadActiveOrders(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to load active orders for store %d: %w", storeID, err)
	}

	for _, order := range activeOrders {
		for _, sub := range order.Suborders {
			if !isSuborderActive(sub) {
				continue
			}
			accumulateProductUsage(&frozenStock, sub)
			accumulateAdditiveUsage(&frozenStock, sub)
		}
	}

	return frozenStock, nil
}

func (r *orderRepository) loadActiveOrders(storeID uint) ([]data.Order, error) {
	var orders []data.Order
	err := r.db.
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive.Ingredients.Ingredient").
		Preload("Suborders.StoreProductSize.ProductSize.ProductSizeIngredients.Ingredient").
		Where("store_id = ?", storeID).
		Where("status IN ?", []data.OrderStatus{
			data.OrderStatusPending,
			data.OrderStatusPreparing,
		}).
		Find(&orders).Error
	return orders, err
}

func isSuborderActive(sub data.Suborder) bool {
	return sub.Status == data.SubOrderStatusPending || sub.Status == data.SubOrderStatusPreparing
}

func accumulateProductUsage(frozenStock *map[uint]float64, sub data.Suborder) {
	for _, usage := range sub.StoreProductSize.ProductSize.ProductSizeIngredients {
		(*frozenStock)[usage.IngredientID] += usage.Quantity
	}
}

func accumulateAdditiveUsage(frozenStock *map[uint]float64, sub data.Suborder) {
	for _, subAdd := range sub.SuborderAdditives {
		for _, ingrUsage := range subAdd.StoreAdditive.Additive.Ingredients {
			(*frozenStock)[ingrUsage.IngredientID] += ingrUsage.Quantity
		}
	}
}
