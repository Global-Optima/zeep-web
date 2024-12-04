package orders

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAllOrders(storeID uint, status *string) ([]data.Order, error)
	GetStoreOrderById(storeID, orderID uint) (*data.Order, error)
	GetOrderByOrderId(orderID uint) (*data.Order, error)
	CreateOrder(order *data.Order) error
	UpdateOrderStatus(orderID uint, storeID uint, status data.OrderStatus) error
	DeleteOrder(orderID uint) error

	GetStatusesCount(storeID uint) (map[string]int64, error)
	GetLowStockIngredients(threshold float64) ([]data.Ingredient, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetAllOrders(storeID uint, status *string) ([]data.Order, error) {
	var orders []data.Order
	query := r.db.Preload("OrderProducts").Where("store_id = ?", storeID)

	if status != nil && *status != "" {
		query = query.Where("order_status = ?", *status)
	}

	err := query.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) GetStoreOrderById(storeID, orderID uint) (*data.Order, error) {
	var order data.Order
	err := r.db.Preload("OrderProducts.Additives").
		Where("id = ? AND store_id = ?", orderID, storeID).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) GetOrderByOrderId(orderID uint) (*data.Order, error) {
	var order data.Order
	err := r.db.Preload("OrderProducts.Additives").
		Where("id = ?", orderID).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) CreateOrder(order *data.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		productSizeIDs := extractProductSizeIDs(order.OrderProducts)
		productSizes, err := fetchProductSizes(tx, productSizeIDs)
		if err != nil {
			return fmt.Errorf("failed to fetch product sizes: %w", err)
		}

		for _, subOrder := range order.OrderProducts {
			productSize, ok := productSizes[subOrder.ProductSizeID]
			if !ok {
				return fmt.Errorf("product size ID %d not found", subOrder.ProductSizeID)
			}

			if err := updateStockForProductSize(tx, subOrder, productSize); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *orderRepository) UpdateOrderStatus(orderID uint, storeID uint, status data.OrderStatus) error {
	return r.db.Model(&data.Order{}).
		Where("id = ? AND store_id = ?", orderID, storeID).
		Update("order_status", status).Error
}

func (r *orderRepository) DeleteOrder(orderID uint) error {
	return r.db.Delete(&data.Order{}, orderID).Error
}

func (r *orderRepository) GetStatusesCount(storeID uint) (map[string]int64, error) {
	var counts []struct {
		OrderStatus string
		Count       int64
	}
	err := r.db.Model(&data.Order{}).
		Select("order_status, COUNT(*) as count").
		Where("store_id = ?", storeID).
		Group("order_status").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	statusCount := make(map[string]int64)
	for _, c := range counts {
		statusCount[c.OrderStatus] = c.Count
	}
	return statusCount, nil
}

func (r *orderRepository) GetLowStockIngredients(threshold float64) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	err := r.db.Model(&data.Ingredient{}).
		Joins("JOIN store_warehouse_stocks ON ingredients.id = store_warehouse_stocks.ingredient_id").
		Where("store_warehouse_stocks.quantity <= ?", threshold).
		Group("ingredients.id").
		Find(&ingredients).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch low stock ingredients: %w", err)
	}

	return ingredients, nil
}

func extractProductSizeIDs(subOrders []data.OrderProduct) []uint {
	var ids []uint
	for _, subOrder := range subOrders {
		ids = append(ids, subOrder.ProductSizeID)
	}
	return ids
}

func fetchProductSizes(tx *gorm.DB, productSizeIDs []uint) (map[uint]data.ProductSize, error) {
	var productSizes []data.ProductSize
	if err := tx.Where("id IN ?", productSizeIDs).Find(&productSizes).Error; err != nil {
		return nil, err
	}

	productSizeMap := make(map[uint]data.ProductSize)
	for _, ps := range productSizes {
		productSizeMap[ps.ID] = ps
	}
	return productSizeMap, nil
}

func updateStockForProductSize(tx *gorm.DB, subOrder data.OrderProduct, productSize data.ProductSize) error {
	for _, pi := range productSize.ProductIngredients {
		totalQuantity := float64(subOrder.Quantity) * pi.ItemIngredient.Quantity

		if err := updateIngredientStock(tx, pi.ItemIngredient.IngredientID, totalQuantity); err != nil {
			return err
		}
	}
	return nil
}

func updateIngredientStock(tx *gorm.DB, ingredientID uint, requiredQuantity float64) error {
	var stock data.StoreWarehouseStock
	if err := tx.Where("ingredient_id = ?", ingredientID).First(&stock).Error; err != nil {
		return fmt.Errorf("failed to find stock for ingredient %d: %w", ingredientID, err)
	}

	if stock.Quantity < requiredQuantity {
		return fmt.Errorf("insufficient stock for ingredient %d: available %f, required %f",
			ingredientID, stock.Quantity, requiredQuantity)
	}

	if err := tx.Model(&data.StoreWarehouseStock{}).
		Where("ingredient_id = ?", ingredientID).
		UpdateColumn("quantity", gorm.Expr("quantity - ?", requiredQuantity)).Error; err != nil {
		return fmt.Errorf("failed to update inventory for ingredient %d: %w", ingredientID, err)
	}

	return nil
}
