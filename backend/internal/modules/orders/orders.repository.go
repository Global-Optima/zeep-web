package orders

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAllOrders(status *string) ([]data.Order, error)
	CreateOrder(order *data.Order) error
	UpdateOrderStatus(orderID uint, status string) error
	GetOrderById(orderID uint) (*data.Order, error)
	DeleteOrder(orderID uint) error

	CreateSubOrder(subOrder *data.OrderProduct) error
	UpdateSubOrder(subOrder *data.OrderProduct) error
	GetSubOrdersByOrderID(orderID uint) ([]data.OrderProduct, error)
	GetSubOrderByID(subOrderID uint) (*data.OrderProduct, error)
	GetSubOrders(orderID uint) ([]data.OrderProduct, error)
	DeleteSubOrder(subOrderID uint) error

	CreateSubOrderAdditive(orderProductAdditive *data.OrderProductAdditive) error
	GetSubOrderAdditivesBySubOrderID(subOrderID uint) ([]data.OrderProductAdditive, error)
	DeleteSubOrderAdditivesBySubOrderID(subOrderID uint) error

	GetStatusesCount() (map[string]int64, error)
	GetSubOrderCount(orderID uint) (int64, error)
	UpdateTime(subOrderID uint, updatedTime time.Time) error
	UpdateInventory(productID uint, quantity int) error
	GetLowStockIngredients(threshold float64) ([]data.Ingredient, error)
	GetProductSizeLabel(productSizeID uint) (string, error)
	GetAdditiveByID(additiveID uint) (*data.Additive, error)
	GetProductSizeByID(productSizeID uint) (*data.ProductSize, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetAllOrders(status *string) ([]data.Order, error) {
	var orders []data.Order
	query := r.db.Preload("OrderProducts")

	if status != nil && *status != "" {
		query = query.Where("order_status = ?", *status)
	}

	err := query.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) GetSubOrders(orderID uint) ([]data.OrderProduct, error) {
	var subOrders []data.OrderProduct
	err := r.db.Preload("ProductSize").
		Where("order_id = ?", orderID).Find(&subOrders).Error
	if err != nil {
		return nil, err
	}

	return subOrders, nil
}

func (r *orderRepository) CreateOrder(order *data.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		for _, subOrder := range order.OrderProducts {
			var productSize data.ProductSize
			if err := tx.Preload("Product").
				First(&productSize, subOrder.ProductSizeID).Error; err != nil {
				return fmt.Errorf("failed to fetch product size ID %d: %w", subOrder.ProductSizeID, err)
			}

			var product data.Product
			if err := tx.Preload("ProductIngredients.ItemIngredient.Ingredient").
				First(&product, productSize.ProductID).Error; err != nil {
				return fmt.Errorf("failed to fetch product ID %d: %w", productSize.ProductID, err)
			}

			for _, pi := range product.ProductIngredients {
				totalQuantity := float64(subOrder.Quantity) * pi.ItemIngredient.Quantity

				var stock data.StoreWarehouseStock
				if err := tx.Where("ingredient_id = ?", pi.ItemIngredient.IngredientID).
					First(&stock).Error; err != nil {
					return fmt.Errorf("failed to find stock for ingredient %d: %w", pi.ItemIngredient.IngredientID, err)
				}

				if stock.Quantity < totalQuantity {
					return fmt.Errorf("insufficient stock for ingredient %d: available %f, required %f",
						pi.ItemIngredient.IngredientID, stock.Quantity, totalQuantity)
				}

				if err := tx.Model(&data.StoreWarehouseStock{}).
					Where("ingredient_id = ?", pi.ItemIngredient.IngredientID).
					UpdateColumn("quantity", gorm.Expr("quantity - ?", totalQuantity)).Error; err != nil {
					return fmt.Errorf("failed to update inventory for ingredient %d: %w", pi.ItemIngredient.IngredientID, err)
				}
			}
		}

		// commit the transaction
		return nil
	})
}

func (r *orderRepository) UpdateOrderStatus(orderID uint, status string) error {
	return r.db.Model(&data.Order{}).Where("id = ?", orderID).Update("order_status", status).Error
}

func (r *orderRepository) GetOrderById(orderID uint) (*data.Order, error) {
	var order data.Order
	err := r.db.Preload("OrderProducts.Additives").
		Where("id = ?", orderID).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) DeleteOrder(orderID uint) error {
	return r.db.Delete(&data.Order{}, orderID).Error
}

func (r *orderRepository) CreateSubOrder(orderProduct *data.OrderProduct) error {
	return r.db.Create(orderProduct).Error
}

func (r *orderRepository) UpdateSubOrder(orderProduct *data.OrderProduct) error {
	return r.db.Save(orderProduct).Error
}

func (r *orderRepository) GetSubOrdersByOrderID(orderID uint) ([]data.OrderProduct, error) {
	var orderProducts []data.OrderProduct
	err := r.db.Preload("ProductSize").
		Preload("Additives.Additive").
		Where("order_id = ?", orderID).
		Find(&orderProducts).Error
	return orderProducts, err
}

func (r *orderRepository) DeleteSubOrder(subOrderID uint) error {
	return r.db.Delete(&data.OrderProduct{}, subOrderID).Error
}

func (r *orderRepository) CreateSubOrderAdditive(orderProductAdditive *data.OrderProductAdditive) error {
	return r.db.Create(orderProductAdditive).Error
}

func (r *orderRepository) GetSubOrderAdditivesBySubOrderID(subOrderID uint) ([]data.OrderProductAdditive, error) {
	var additives []data.OrderProductAdditive
	err := r.db.Preload("Additive").
		Where("order_product_id = ?", subOrderID).
		Find(&additives).Error
	return additives, err
}

func (r *orderRepository) DeleteSubOrderAdditivesBySubOrderID(subOrderID uint) error {
	return r.db.Where("order_product_id = ?", subOrderID).Delete(&data.OrderProductAdditive{}).Error
}

func (r *orderRepository) GetStatusesCount() (map[string]int64, error) {
	var counts []struct {
		OrderStatus string
		Count       int64
	}
	err := r.db.Model(&data.Order{}).
		Select("order_status, COUNT(*) as count").
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

func (r *orderRepository) GetSubOrderCount(orderID uint) (int64, error) {
	var count int64
	err := r.db.Model(&data.OrderProduct{}).
		Where("order_id = ?", orderID).
		Count(&count).Error
	return count, err
}

func (r *orderRepository) UpdateTime(subOrderID uint, updatedTime time.Time) error {
	return r.db.Model(&data.OrderProduct{}).
		Where("id = ?", subOrderID).
		Update("updated_at", updatedTime).Error
}

func (r *orderRepository) GetSubOrderByID(subOrderID uint) (*data.OrderProduct, error) {
	var orderProduct data.OrderProduct
	err := r.db.Where("id = ?", subOrderID).First(&orderProduct).Error
	return &orderProduct, err
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

func (r *orderRepository) UpdateInventory(productID uint, quantity int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var product data.Product
		if err := tx.Preload("ProductIngredients.ItemIngredient.Ingredient").
			First(&product, productID).Error; err != nil {
			return fmt.Errorf("failed to fetch product: %w", err)
		}

		for _, pi := range product.ProductIngredients {
			totalQuantity := float64(quantity) * pi.ItemIngredient.Quantity

			if err := updateIngredientStock(tx, pi.ItemIngredient.IngredientID, totalQuantity); err != nil {
				return err
			}
		}

		return nil // Transaction commits automatically
	})
}

func (r *orderRepository) GetLowStockIngredients(threshold float64) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	err := r.db.Model(&data.Ingredient{}).
		Joins("JOIN store_warehouse_stock ON ingredients.id = store_warehouse_stock.ingredient_id").
		Where("store_warehouse_stock.quantity <= ?", threshold).
		Group("ingredients.id").
		Find(&ingredients).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch low stock ingredients: %w", err)
	}

	return ingredients, nil
}

func (r *orderRepository) GetProductSizeLabel(productSizeID uint) (string, error) {
	var productSize data.ProductSize
	err := r.db.Where("id = ?", productSizeID).First(&productSize).Error
	return productSize.Name, err
}

func (r *orderRepository) GetProductSizeByID(productSizeID uint) (*data.ProductSize, error) {
	var productSize data.ProductSize
	err := r.db.Where("id = ?", productSizeID).First(&productSize).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size with ID %d: %w", productSizeID, err)
	}
	return &productSize, nil
}

func (r *orderRepository) GetAdditiveByID(additiveID uint) (*data.Additive, error) {
	var additive data.Additive
	err := r.db.Where("id = ?", additiveID).First(&additive).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive with ID %d: %w", additiveID, err)
	}
	return &additive, nil
}
