package stockMaterial

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
)

type StockMaterialService interface {
	GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error)
	GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error)
	CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error)
	UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) (*types.StockMaterialsDTO, error)
	DeleteStockMaterial(stockMaterialID uint) error
	DeactivateStockMaterial(stockMaterialID uint) error
}

type stockMaterialService struct {
	repo        StockMaterialRepository
	packageRepo stockMaterialPackage.StockMaterialPackageRepository
}

func NewStockMaterialService(repo StockMaterialRepository, packageRepo stockMaterialPackage.StockMaterialPackageRepository) StockMaterialService {
	return &stockMaterialService{
		repo:        repo,
		packageRepo: packageRepo,
	}
}

func (s *stockMaterialService) GetAllStockMaterials(filter *types.StockMaterialFilter) ([]types.StockMaterialsDTO, error) {
	stockMaterials, err := s.repo.GetAllStockMaterials(filter)
	if err != nil {
		return nil, err
	}

	stockMaterialResponses := make([]types.StockMaterialsDTO, 0)
	for _, stockMaterial := range stockMaterials {
		stockMaterialResponses = append(stockMaterialResponses, *types.ConvertStockMaterialToStockMaterialResponse(&stockMaterial))
	}

	return stockMaterialResponses, nil
}

func (s *stockMaterialService) GetStockMaterialByID(stockMaterialID uint) (*types.StockMaterialsDTO, error) {
	stockMaterial, err := s.repo.GetStockMaterialByID(stockMaterialID)
	if err != nil {
		return nil, err
	}

	if stockMaterial == nil {
		return nil, errors.New("StockMaterial not found")
	}

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)

	return stockMaterialResponse, nil
}

func (s *stockMaterialService) CreateStockMaterial(req *types.CreateStockMaterialDTO) (*types.StockMaterialsDTO, error) {
	stockMaterial := types.ConvertCreateStockMaterialRequestToStockMaterial(req)

	err := s.repo.CreateStockMaterial(stockMaterial)
	if err != nil {
		return nil, err
	}

	if len(req.Packages) != 0 {
		packages := make([]data.StockMaterialPackage, len(req.Packages))
		for i, pkg := range req.Packages {
			packages[i] = *types.ConvertPackageDTOToModel(stockMaterial.ID, &pkg)
		}

		err = s.packageRepo.CreateMultiplePackages(packages)
		if err != nil {
			return nil, err
		}
	}

	stockMaterialResponse := types.ConvertStockMaterialToStockMaterialResponse(stockMaterial)
	return stockMaterialResponse, nil
}

func (s *stockMaterialService) UpdateStockMaterial(stockMaterialID uint, req *types.UpdateStockMaterialDTO) (*types.StockMaterialsDTO, error) {
	updated, err := s.repo.UpdateStockMaterialFields(stockMaterialID, *req)
	if err != nil {
		return nil, err
	}
	updatedStockMaterial := types.ConvertStockMaterialToStockMaterialResponse(updated)
	return updatedStockMaterial, nil
}

func (s *stockMaterialService) DeleteStockMaterial(stockMaterialID uint) error {
	return s.repo.DeleteStockMaterial(stockMaterialID)
}

func (s *stockMaterialService) DeactivateStockMaterial(stockMaterialID uint) error {
	return s.repo.DeactivateStockMaterial(stockMaterialID)
}
