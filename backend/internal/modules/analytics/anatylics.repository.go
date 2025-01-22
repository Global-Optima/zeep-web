package analytics

import (
	"time"

	models "github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type AnalyticsRepo interface {
	GetOrdersForSummary(startDate, endDate *time.Time, storeID *uint) (totalSales float64, totalOrders, totalProductsSold, totalAdditivesSold int, err error)
	GetOrdersForMonthlySales(startDate, endDate *time.Time, storeID *uint) ([]MonthlySalesData, error)
	GetPopularProducts(startDate, endDate *time.Time, storeID *uint) ([]PopularProductData, error)
	GetProductsSold(startDate, endDate *time.Time, storeID *uint) ([]ProductSoldData, error)
}

type analyticsRepo struct {
	db *gorm.DB
}

func NewAnalyticsRepo(db *gorm.DB) AnalyticsRepo {
	return &analyticsRepo{db: db}
}

type MonthlySalesData struct {
	Month  string  `json:"month"`
	Year   int     `json:"year"`
	Sales  float64 `json:"sales"`
	Orders int     `json:"orders"`
}

type PopularProductData struct {
	ProductName string  `json:"productName"`
	TotalSold   int     `json:"totalSold"`
	Revenue     float64 `json:"revenue"`
}

type ProductSoldData struct {
	ProductName  string             `json:"productName"`
	TotalSold    int                `json:"totalSold"`
	TotalRevenue float64            `json:"totalRevenue"`
	Additives    []AdditiveSoldData `json:"additives"`
}

type AdditiveSoldData struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

func completedOrDeliveredScope(db *gorm.DB) *gorm.DB {
	return db.Model(&models.Order{}).Where("status IN ?", []models.OrderStatus{models.OrderStatusCompleted, models.OrderStatusDelivered})
}

func dateRangeScope(startDate, endDate *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if startDate != nil {
			db = db.Where("created_at >= ?", startDate)
		}
		if endDate != nil {
			db = db.Where("created_at <= ?", endDate)
		}
		return db
	}
}

func storeScope(storeID *uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if storeID != nil {
			db = db.Where("store_id = ?", storeID)
		}
		return db
	}
}

func (r *analyticsRepo) GetOrdersForSummary(startDate, endDate *time.Time, storeID *uint) (float64, int, int, int, error) {
	var order models.Order
	var result struct {
		TotalSales         float64
		TotalOrders        int64
		TotalProductsSold  int64
		TotalAdditivesSold int64
	}

	err := r.db.Model(&order).
		Scopes(
			completedOrDeliveredScope,
			dateRangeScope(startDate, endDate),
			storeScope(storeID),
		).
		Select(`
            COALESCE(SUM(total), 0) as total_sales,
            COUNT(*) as total_orders,
            (SELECT COUNT(*) FROM suborders WHERE order_id IN (SELECT id FROM orders)) as total_products_sold,
            (SELECT COUNT(*) FROM suborder_additives WHERE suborder_id IN 
                (SELECT id FROM suborders WHERE order_id IN (SELECT id FROM orders))) as total_additives_sold
        `).
		Scan(&result).Error

	return result.TotalSales, int(result.TotalOrders), int(result.TotalProductsSold), int(result.TotalAdditivesSold), err
}

func (r *analyticsRepo) GetPopularProducts(startDate, endDate *time.Time, storeID *uint) ([]PopularProductData, error) {
	var results []PopularProductData

	ordersQuery := r.db.Model(&models.Order{}).
		Scopes(
			completedOrDeliveredScope,
			dateRangeScope(startDate, endDate),
			storeScope(storeID),
		).
		Select("id")

	err := r.db.Model(&models.Suborder{}).
		Joins("JOIN product_sizes ps ON ps.id = suborders.product_size_id").
		Joins("JOIN products p ON p.id = ps.product_id").
		Where("order_id IN (?)", ordersQuery).
		Group("p.name").
		Select(`
            p.name as product_name,
            COUNT(*) as total_sold,
            SUM(suborders.price) as revenue
        `).
		Order("total_sold DESC").
		Scan(&results).Error

	return results, err
}

func (r *analyticsRepo) GetProductsSold(startDate, endDate *time.Time, storeID *uint) ([]ProductSoldData, error) {
	var results []ProductSoldData

	query := r.db.Model(&models.Product{}).
		Joins("JOIN product_sizes ON product_sizes.product_id = products.id").
		Joins("JOIN suborders ON suborders.product_size_id = product_sizes.id").
		Joins("JOIN orders ON orders.id = suborders.order_id").
		Scopes(
			completedOrDeliveredScope,
			dateRangeScope(startDate, endDate),
			storeScope(storeID),
		).
		Select(`
			products.name AS product_name,
			COUNT(suborders.id) AS total_sold,
			SUM(suborders.price) AS total_revenue
		`).
		Group("products.name")

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for i := range results {
		var additives []AdditiveSoldData
		err := r.db.Model(&models.SuborderAdditive{}).
			Joins("JOIN suborders ON suborders.id = suborder_additives.suborder_id").
			Joins("JOIN product_sizes ON product_sizes.id = suborders.product_size_id").
			Joins("JOIN products ON products.id = product_sizes.product_id").
			Where("products.name = ?", results[i].ProductName).
			Select(`
				additives.name AS name,
				COUNT(suborder_additives.id) AS quantity,
				SUM(suborder_additives.price) AS revenue
			`).
			Joins("JOIN additives ON additives.id = suborder_additives.additive_id").
			Group("additives.name").
			Scan(&additives).Error
		if err != nil {
			return nil, err
		}

		results[i].Additives = additives
	}

	return results, nil
}

func (r *analyticsRepo) GetOrdersForMonthlySales(startDate, endDate *time.Time, storeID *uint) ([]MonthlySalesData, error) {
	var results []MonthlySalesData

	err := r.db.Model(&models.Order{}).
		Scopes(
			completedOrDeliveredScope,
			dateRangeScope(startDate, endDate),
			storeScope(storeID),
		).
		Select(`
            TO_CHAR(created_at, 'Month') as month,
            EXTRACT(YEAR FROM created_at) as year,
            COUNT(*) as orders,
            SUM(total) as sales
        `).
		Group("TO_CHAR(created_at, 'Month'), EXTRACT(YEAR FROM created_at)").
		Order("year, month").
		Scan(&results).Error

	return results, err
}
