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
	CreateOrder(order *data.Order) (uint, error)
	UpdateOrderStatus(suborderID uint, updateData types.UpdateOrderDTO) error
	DeleteOrder(orderID uint) error

	GetSubOrdersByOrderID(orderID uint) ([]data.Suborder, error)
	UpdateSubOrderStatus(suborderID uint, updateData types.UpdateSubOrderDTO) error
	AddSubOrderAdditive(subOrderID uint, additive *data.SuborderAdditive) error
	CheckAllSubordersCompleted(orderID, subOrderID uint) (bool, error)
	GetOrderBySubOrderID(subOrderID uint) (*data.Order, error)

	GetOrderDetails(orderID uint, filter *contexts.StoreContextFilter) (*data.Order, error)
	GetOrdersForExport(filter *types.OrdersExportFilterQuery) ([]data.Order, error)
	GetSuborderByID(suborderID uint) (*data.Suborder, error)

	CalculateFrozenStock(storeID uint) (map[uint]float64, error)
	HandlePaymentSuccess(orderID uint, paymentTransaction *data.Transaction) (*data.Order, error)
	HandlePaymentFailure(orderID uint) error
	CloneWithTransaction(tx *gorm.DB) orderRepository
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CloneWithTransaction(tx *gorm.DB) orderRepository {
	return orderRepository{
		db: tx,
	}
}

func (r *orderRepository) CreateOrder(order *data.Order) (uint, error) {
	const workingHours = 16 // Working hours for a cafe

	err := r.db.Transaction(func(tx *gorm.DB) error {
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

		orderIngredients, err := r.getOrderIngredients(tx, order.ID)
		if err != nil {
			return fmt.Errorf("failed to get order ingredients: %w", err)
		}

		if err := data.RecalculateOutOfStock(tx, order.StoreID, orderIngredients, nil, nil); err != nil {
			return fmt.Errorf("failed to recalculate out of stock: %w", err)
		}

		return nil
	})

	return order.ID, err
}

func (r *orderRepository) getOrderIngredients(tx *gorm.DB, orderID uint) ([]uint, error) {
	ingredientIDsFromProductSizes, err := r.getOrderProductSizeIngredients(tx, orderID)
	if err != nil {
		return nil, err
	}

	ingredientIDsFromAdditives, err := r.getOrderAdditiveIngredients(tx, orderID)
	if err != nil {
		return nil, err
	}

	return utils.UnionSlices(ingredientIDsFromProductSizes, ingredientIDsFromAdditives), nil
}

func (r *orderRepository) getOrderProductSizeIngredients(tx *gorm.DB, orderID uint) ([]uint, error) {
	var ingredientIDsFromProductSizes []uint

	err := tx.Model(&data.ProductSizeIngredient{}).
		Distinct("product_size_ingredients.ingredient_id").
		Joins("JOIN store_product_sizes ON store_product_sizes.product_size_id = product_size_ingredients.product_size_id").
		Joins("JOIN suborders ON suborders.store_product_size_id = store_product_sizes.id").
		Where("suborders.order_id = ?", orderID).
		Pluck("product_size_ingredients.ingredient_id", &ingredientIDsFromProductSizes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients from product sizes: %w", err)
	}

	return ingredientIDsFromProductSizes, nil
}

func (r *orderRepository) getOrderAdditiveIngredients(tx *gorm.DB, orderID uint) ([]uint, error) {
	var ingredientIDsFromAdditives []uint

	err := tx.Model(&data.AdditiveIngredient{}).
		Distinct("additive_ingredients.ingredient_id").
		Joins("JOIN store_additives ON store_additives.additive_id = additive_ingredients.additive_id").
		Joins("JOIN suborder_additives ON suborder_additives.store_additive_id = store_additives.id").
		Joins("JOIN suborders ON suborders.id = suborder_additives.suborder_id").
		Where("suborders.order_id = ?", orderID).
		Pluck("additive_ingredients.ingredient_id", &ingredientIDsFromAdditives).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients from additives: %w", err)
	}

	return ingredientIDsFromAdditives, nil
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

	// Get the correct timezone location
	location, err := getTimeZoneLocation(filter)
	if err != nil {
		return nil, err
	}

	// Get the current date in the specified timezone
	now := time.Now().In(location)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	if filter.IncludeYesterdayOrders != nil && *filter.IncludeYesterdayOrders {
		yesterday := now.AddDate(0, 0, -1)
		startOfToday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, location)
	}

	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, location)

	// Convert start and end of today to UTC for database querying
	startOfTodayUTC := startOfToday.UTC()
	endOfTodayUTC := endOfToday.UTC()

	// Base query with preloads and business logic
	query := r.db.
		Preload("Suborders.StoreProductSize.ProductSize.Product").
		Preload("Suborders.StoreProductSize.ProductSize.Unit").
		Preload("Suborders.SuborderAdditives.StoreAdditive.Additive").
		Where("store_id = ?", *filter.StoreID).
		Where("status NOT IN (?)", []data.OrderStatus{data.OrderStatusWaitingForPayment}).
		Where("created_at BETWEEN ? AND ?", startOfTodayUTC, endOfTodayUTC).
		Order("created_at ASC")

	// Optionally filter by an array of statuses if provided
	if len(filter.Statuses) > 0 {
		query = query.Where("status IN (?)", filter.Statuses)
	}

	// Apply additional time gap filtering if needed
	if filter.TimeGapMinutes != nil && *filter.TimeGapMinutes > 0 {
		cutoffTime := now.Add(-time.Duration(*filter.TimeGapMinutes) * time.Minute)
		cutoffTimeUTC := cutoffTime.UTC()
		query = query.Where("(status != ? OR completed_at >= ?)", data.OrderStatusCompleted, cutoffTimeUTC)
	}

	// Execute the query
	err = query.Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch barista orders: %w", err)
	}

	return orders, nil
}

// Helper function to get the timezone location
func getTimeZoneLocation(filter types.OrdersTimeZoneFilter) (*time.Location, error) {
	// Default to UTC
	location := time.UTC

	// If timezone location is specified, use that
	if filter.TimeZoneLocation != nil && *filter.TimeZoneLocation != "" {
		var err error
		location, err = time.LoadLocation(*filter.TimeZoneLocation)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone location: %w", err)
		}
	} else if filter.TimeZoneOffset != nil {
		// If timezone offset is specified, create a fixed timezone
		hours := int(*filter.TimeZoneOffset) / 60
		minutes := int(*filter.TimeZoneOffset) % 60
		location = time.FixedZone(fmt.Sprintf("UTC%+d:%02d", hours, minutes), int(*filter.TimeZoneOffset)*60)
	}

	return location, nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to fetch order with ID %d: %w", orderId, err)
	}

	return &order, nil
}

func (r *orderRepository) GetRawOrderById(orderId uint) (*data.Order, error) {
	var order data.Order
	err := r.db.Where("id = ?", orderId).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to fetch raw order with ID %d: %w", orderId, err)
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(orderID uint, updateData types.UpdateOrderDTO) error {
	return r.db.Model(&data.Order{}).
		Where("id = ?", orderID).
		Updates(updateData).Error
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

func (r *orderRepository) UpdateSubOrderStatus(suborderID uint, updateData types.UpdateSubOrderDTO) error {
	return r.db.Model(&data.Suborder{}).
		Where("id = ?", suborderID).
		Updates(updateData).Error
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
			data.OrderStatusWaitingForPayment,
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

func (r *orderRepository) HandlePaymentSuccess(orderID uint, paymentTransaction *data.Transaction) (*data.Order, error) {
	order, err := r.GetRawOrderById(orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != data.OrderStatusWaitingForPayment {
		return nil, types.ErrInappropriateOrderStatus
	}
	order.Status = data.OrderStatusPending

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&data.Order{}).
			Where(&data.Order{BaseEntity: data.BaseEntity{ID: orderID}}).
			Save(order).Error
		if err != nil {
			return err
		}

		if err := tx.Create(paymentTransaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) HandlePaymentFailure(orderID uint) error {
	order, err := r.GetRawOrderById(orderID)
	if err != nil {
		return err
	}

	if order.Status != data.OrderStatusWaitingForPayment {
		return types.ErrInappropriateOrderStatus
	}

	orderIngredients, err := r.getOrderIngredients(r.db, orderID)
	if err != nil {
		return err
	}

	if err := data.RecalculateOutOfStock(r.db, order.StoreID, orderIngredients, nil, nil); err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&order).Error; err != nil {
		return err
	}

	return nil
}
