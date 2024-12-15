package customers

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	GetCustomerByID(customerID uint) (*data.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) GetCustomerByID(customerID uint) (*data.Customer, error) {
	var customer data.Customer
	err := r.db.First(&customer, customerID).Error
	return &customer, err
}
