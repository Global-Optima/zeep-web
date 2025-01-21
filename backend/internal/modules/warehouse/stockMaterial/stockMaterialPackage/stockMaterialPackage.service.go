package stockMaterialPackage

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage/types"
	"gorm.io/gorm"
)

type StockMaterialPackageService interface {
	Create(dto types.CreateStockMaterialPackageDTO) (uint, error)
	GetByID(id uint) (*types.StockMaterialPackageResponse, error)
	Update(id uint, dto types.UpdateStockMaterialPackageDTO) error
	Delete(id uint) error
	GetAll(filter types.StockMaterialPackageFilter) ([]types.StockMaterialPackageResponse, error)
}

type stockMaterialPackageService struct {
	repo StockMaterialPackageRepository
}

func NewStockMaterialPackageService(repo StockMaterialPackageRepository) StockMaterialPackageService {
	return &stockMaterialPackageService{repo: repo}
}

func (s *stockMaterialPackageService) Create(dto types.CreateStockMaterialPackageDTO) (uint, error) {
	packageEntity := &data.StockMaterialPackage{
		StockMaterialID: dto.StockMaterialID,
		Size:            dto.Size,
		UnitID:          dto.UnitID,
	}

	if err := s.repo.Create(packageEntity); err != nil {
		return 0, fmt.Errorf("failed to create stock material package: %w", err)
	}

	return packageEntity.ID, nil
}

func (s *stockMaterialPackageService) Update(id uint, dto types.UpdateStockMaterialPackageDTO) error {
	packageEntity := data.StockMaterialPackage{}
	if dto.StockMaterialID != nil {
		packageEntity.StockMaterialID = *dto.StockMaterialID
	}
	if dto.Size != nil {
		packageEntity.Size = *dto.Size
	}
	if dto.UnitID != nil {
		packageEntity.UnitID = *dto.UnitID
	}

	if err := s.repo.Update(id, packageEntity); err != nil {
		return fmt.Errorf("failed to update stock material package: %w", err)
	}

	return nil
}

func (s *stockMaterialPackageService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete stock material package: %w", err)
	}
	return nil
}

func (s *stockMaterialPackageService) GetByID(id uint) (*types.StockMaterialPackageResponse, error) {
	packageEntity, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock material package not found")
		}
		return nil, fmt.Errorf("failed to fetch stock material package: %w", err)
	}
	response := types.ToStockMaterialPackageResponse(packageEntity)
	return &response, nil
}

func (s *stockMaterialPackageService) GetAll(filter types.StockMaterialPackageFilter) ([]types.StockMaterialPackageResponse, error) {
	packages, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock material packages: %w", err)
	}
	return types.ToStockMaterialPackageResponses(packages), nil
}
