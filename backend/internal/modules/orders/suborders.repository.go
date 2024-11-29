package orders

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type SubOrderRepository interface {
	CreateSubOrder(subOrder *data.OrderProduct) error
	UpdateSubOrder(subOrder *data.OrderProduct) error
	GetSubOrdersByOrderID(orderID uint) ([]data.OrderProduct, error)
	GetSubOrderByID(subOrderID uint) (*data.OrderProduct, error)
	GetSubOrders(orderID uint) ([]data.OrderProduct, error)
	DeleteSubOrder(subOrderID uint) error
	GetSubOrderCount(orderID uint) (int64, error)

	CreateSubOrderAdditive(orderProductAdditive *data.OrderProductAdditive) error
	GetSubOrderAdditivesBySubOrderID(subOrderID uint) ([]data.OrderProductAdditive, error)
	DeleteSubOrderAdditivesBySubOrderID(subOrderID uint) error
}

type subOrderRepository struct {
	db *gorm.DB
}

func NewSubOrderRepository(db *gorm.DB) SubOrderRepository {
	return &subOrderRepository{db: db}
}

func (r *subOrderRepository) GetSubOrders(orderID uint) ([]data.OrderProduct, error) {
	var subOrders []data.OrderProduct
	err := r.db.Preload("ProductSize").
		Where("order_id = ?", orderID).Find(&subOrders).Error
	if err != nil {
		return nil, err
	}

	return subOrders, nil
}

func (r *subOrderRepository) CreateSubOrder(orderProduct *data.OrderProduct) error {
	return r.db.Create(orderProduct).Error
}

func (r *subOrderRepository) UpdateSubOrder(orderProduct *data.OrderProduct) error {
	return r.db.Save(orderProduct).Error
}

func (r *subOrderRepository) GetSubOrdersByOrderID(orderID uint) ([]data.OrderProduct, error) {
	var orderProducts []data.OrderProduct
	err := r.db.Preload("ProductSize").
		Preload("Additives.Additive").
		Where("order_id = ?", orderID).
		Find(&orderProducts).Error
	return orderProducts, err
}

func (r *subOrderRepository) DeleteSubOrder(subOrderID uint) error {
	return r.db.Delete(&data.OrderProduct{}, subOrderID).Error
}

func (r *subOrderRepository) CreateSubOrderAdditive(orderProductAdditive *data.OrderProductAdditive) error {
	return r.db.Create(orderProductAdditive).Error
}

func (r *subOrderRepository) GetSubOrderAdditivesBySubOrderID(subOrderID uint) ([]data.OrderProductAdditive, error) {
	var additives []data.OrderProductAdditive
	err := r.db.Preload("Additive").
		Where("order_product_id = ?", subOrderID).
		Find(&additives).Error
	return additives, err
}

func (r *subOrderRepository) DeleteSubOrderAdditivesBySubOrderID(subOrderID uint) error {
	return r.db.Where("order_product_id = ?", subOrderID).Delete(&data.OrderProductAdditive{}).Error
}

func (r *subOrderRepository) GetSubOrderByID(subOrderID uint) (*data.OrderProduct, error) {
	var orderProduct data.OrderProduct
	err := r.db.Where("id = ?", subOrderID).First(&orderProduct).Error
	return &orderProduct, err
}

func (r *subOrderRepository) GetStatusesCount() (map[string]int64, error) {
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

func (r *subOrderRepository) GetSubOrderCount(orderID uint) (int64, error) {
	var count int64
	err := r.db.Model(&data.OrderProduct{}).
		Where("order_id = ?", orderID).
		Count(&count).Error
	return count, err
}
