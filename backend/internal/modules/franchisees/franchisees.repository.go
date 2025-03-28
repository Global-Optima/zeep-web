package franchisees

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type FranchiseeRepository interface {
	CreateFranchisee(franchisee *data.Franchisee) (uint, error)
	UpdateFranchisee(id uint, updateData *data.Franchisee) error
	DeleteFranchisee(id uint) error
	GetFranchiseeByID(id uint) (*data.Franchisee, error)
	GetFranchisees(filter *types.FranchiseeFilter) ([]data.Franchisee, error)
	GetAllFranchisees(filter *types.FranchiseeFilter) ([]data.Franchisee, error)
	IsFranchiseeStore(franchiseeID, storeID uint) (bool, error)
}

type franchiseeRepository struct {
	db *gorm.DB
}

func NewFranchiseeRepository(db *gorm.DB) FranchiseeRepository {
	return &franchiseeRepository{db: db}
}

func (r *franchiseeRepository) CreateFranchisee(franchisee *data.Franchisee) (uint, error) {
	err := r.db.Create(franchisee).Error
	if err != nil {
		return 0, err
	}
	return franchisee.ID, nil
}

func (r *franchiseeRepository) UpdateFranchisee(id uint, updateData *data.Franchisee) error {
	return r.db.Model(&data.Franchisee{}).Where("id = ?", id).Select("name", "description").Updates(updateData).Error
}

func (r *franchiseeRepository) DeleteFranchisee(id uint) error {
	return r.db.Delete(&data.Franchisee{}, id).Error
}

func (r *franchiseeRepository) GetFranchiseeByID(id uint) (*data.Franchisee, error) {
	var franchisee data.Franchisee
	if err := r.db.First(&franchisee, id).Error; err != nil {
		return nil, err
	}
	return &franchisee, nil
}

func (r *franchiseeRepository) GetFranchisees(filter *types.FranchiseeFilter) ([]data.Franchisee, error) {
	var franchisees []data.Franchisee
	query := r.db.Model(&data.Franchisee{})

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

func (r *franchiseeRepository) GetAllFranchisees(filter *types.FranchiseeFilter) ([]data.Franchisee, error) {
	var franchisees []data.Franchisee
	query := r.db.Model(&data.Franchisee{})

	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	if err := query.Scopes(filter.Sort.SortGorm()).Find(&franchisees).Error; err != nil {
		return nil, err
	}
	return franchisees, nil
}

func (r *franchiseeRepository) IsFranchiseeStore(franchiseeID, storeID uint) (bool, error) {
	var store data.Store
	if err := r.db.Model(&data.Store{}).Where("franchisee_id = ? AND id = ?", franchiseeID, storeID).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
