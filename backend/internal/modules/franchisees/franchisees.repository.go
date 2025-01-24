package franchisees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type FranchiseeRepository interface {
	Create(franchisee *data.Franchisee) (uint, error)
	Update(id uint, updateData *data.Franchisee) error
	Delete(id uint) error
	GetByID(id uint) (*data.Franchisee, error)
	GetAll(filter *types.FranchiseeFilter) ([]data.Franchisee, error)
}

type franchiseeRepository struct {
	db *gorm.DB
}

func NewFranchiseeRepository(db *gorm.DB) FranchiseeRepository {
	return &franchiseeRepository{db: db}
}

func (r *franchiseeRepository) Create(franchisee *data.Franchisee) (uint, error) {
	err := r.db.Create(franchisee).Error
	if err != nil {
		return 0, err
	}
	return franchisee.ID, nil
}

func (r *franchiseeRepository) Update(id uint, updateData *data.Franchisee) error {
	return r.db.Model(&data.Franchisee{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *franchiseeRepository) Delete(id uint) error {
	return r.db.Delete(&data.Franchisee{}, id).Error
}

func (r *franchiseeRepository) GetByID(id uint) (*data.Franchisee, error) {
	var franchisee data.Franchisee
	if err := r.db.First(&franchisee, id).Error; err != nil {
		return nil, err
	}
	return &franchisee, nil
}

func (r *franchiseeRepository) GetAll(filter *types.FranchiseeFilter) ([]data.Franchisee, error) {
	var franchisees []data.Franchisee
	query := r.db.Model(&data.Franchisee{})

	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Franchisee{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&franchisees).Error; err != nil {
		return nil, err
	}
	return franchisees, nil
}
