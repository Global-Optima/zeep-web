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

	GetStatusesCount() (map[string]int64, error)
	UpdateTime(subOrderID uint, updatedTime time.Time) error
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

func (r *orderRepository) CreateOrder(order *data.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		productSizeIDs := extractProductSizeIDs(order.OrderProducts)
		productSizes, err := FetchProductSizes(tx, productSizeIDs)
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

func (r *orderRepository) UpdateTime(subOrderID uint, updatedTime time.Time) error {
	return r.db.Model(&data.OrderProduct{}).
		Where("id = ?", subOrderID).
		Update("updated_at", updatedTime).Error
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

func extractProductSizeIDs(subOrders []data.OrderProduct) []uint {
	var ids []uint
	for _, subOrder := range subOrders {
		ids = append(ids, subOrder.ProductSizeID)
	}
	return ids
}

func FetchProductSizes(tx *gorm.DB, productSizeIDs []uint) (map[uint]data.ProductSize, error) {
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
