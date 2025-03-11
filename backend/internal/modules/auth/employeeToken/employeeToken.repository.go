package employeeToken

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type EmployeeTokenRepository interface {
	CreateToken(token *data.EmployeeTokens) error

	GetTokenByEmployeeID(employeeID uint) (*data.EmployeeTokens, error)

	DeleteToken(token *data.EmployeeTokens) error

	DeleteTokenByEmployeeID(employeeID uint) error
}

type employeeTokenRepository struct {
	db *gorm.DB
}

func NewEmployeeTokenRepository(db *gorm.DB) EmployeeTokenRepository {
	return &employeeTokenRepository{db: db}
}

func (r *employeeTokenRepository) CreateToken(token *data.EmployeeTokens) error {
	return r.db.Create(token).Error
}

func (r *employeeTokenRepository) GetTokenByEmployeeID(employeeID uint) (*data.EmployeeTokens, error) {
	var token data.EmployeeTokens
	if err := r.db.Where("employee_id = ?", employeeID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *employeeTokenRepository) DeleteToken(token *data.EmployeeTokens) error {
	return r.db.Unscoped().Delete(token).Error
}

func (r *employeeTokenRepository) DeleteTokenByEmployeeID(employeeID uint) error {
	return r.db.Unscoped().Where("employee_id = ?", employeeID).Delete(&data.EmployeeTokens{}).Error
}
