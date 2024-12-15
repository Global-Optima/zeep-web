package auth

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AuthenticationRepository interface {
	CreateCustomer(customer *data.Customer) (uint, error)
	GetCustomerByPhone(phone string) (*data.Customer, error)
}

type authorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) AuthenticationRepository {
	return &authorizationRepository{db: db}
}

func (r *authorizationRepository) CreateCustomer(customer *data.Customer) (uint, error) {
	createdCustomer := &data.Customer{}

	err := r.db.Model(&data.Customer{}).
		Create(customer).Scan(createdCustomer).Error

	if err != nil {
		return 0, err
	}

	return createdCustomer.ID, nil
}

func (r *authorizationRepository) GetCustomerByPhone(phone string) (*data.Customer, error) {
	var customer data.Customer
	err := r.db.
		Where("phone = ?", phone).
		First(&customer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &customer, err
}
