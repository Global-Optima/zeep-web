package storeProducts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type StoreProductService interface {
	GetStoreProductById(storeID, storeProductID uint) (*types.StoreProductDetailsDTO, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]types.StoreProductDTO, error)
	CreateStoreProduct(storeID uint, dto *types.CreateStoreProductDTO) (uint, error)
	CreateMultipleStoreProducts(storeID uint, dtos []types.CreateStoreProductDTO) ([]uint, error)
	UpdateStoreProduct(storeID, storeProductID uint, dto *types.UpdateStoreProductDTO) error
	DeleteStoreProduct(storeID, storeProductID uint) error

	//GetStoreProductSizeById(storeID uint, storeProductSizeID uint) (*types.StoreProductSizeDTO, error)
	//GetStoreProductSizes(storeID uint, filter *types.StoreProductSizesFilterDTO) ([]types.StoreProductSizeDTO, error)
	//CreateStoreProductSize(storeID uint, dto *types.CreateStoreProductSizeDTO) (uint, error)
	//UpdateStoreProductSize(storeID, storeProductSizeID uint, dto *types.UpdateStoreProductSizeDTO) error
	//DeleteStoreProductSize(storeID, storeProductSizeID uint) error
}

type storeProductService struct {
	repo        StoreProductRepository
	productRepo product.ProductRepository
	logger      *zap.SugaredLogger
}

func NewStoreProductService(repo StoreProductRepository, productRepo product.ProductRepository, logger *zap.SugaredLogger) StoreProductService {
	return &storeProductService{
		repo:        repo,
		productRepo: productRepo,
		logger:      logger,
	}
}

func (s *storeProductService) GetStoreProductById(storeID uint, storeProductID uint) (*types.StoreProductDetailsDTO, error) {
	dto, err := s.repo.GetStoreProductById(storeID, storeProductID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store product by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return dto, nil
}

func (s *storeProductService) GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]types.StoreProductDTO, error) {
	storeProducts, err := s.repo.GetStoreProducts(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreProductDTO, len(storeProducts))
	for i, storeProduct := range storeProducts {
		dtos[i] = *types.MapToStoreProductDTO(&storeProduct)
	}

	return dtos, nil
}

func (s *storeProductService) CreateStoreProduct(storeID uint, dto *types.CreateStoreProductDTO) (uint, error) {

	if dto.ProductSizes != nil && len(dto.ProductSizes) > 0 {
		inputSizeIDs := make([]uint, len(dto.ProductSizes))
		for i, productSize := range dto.ProductSizes {
			inputSizeIDs[i] = productSize.ProductSizeID
		}

		if err := s.validateProductSizesByProductID(inputSizeIDs, dto.ProductID); err != nil {
			return 0, utils.WrapError("failed to create store product", err)
		}
	}

	storeProduct := types.CreateToStoreProduct(dto)
	storeProduct.StoreID = storeID
	id, err := s.repo.CreateStoreProduct(storeProduct)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store product", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
}

func (s *storeProductService) CreateMultipleStoreProducts(storeID uint, dtos []types.CreateStoreProductDTO) ([]uint, error) {

	storeProducts := make([]data.StoreProduct, len(dtos))
	for i, dto := range dtos {
		storeProducts[i] = *types.CreateToStoreProduct(&dto)
		storeProducts[i].StoreID = storeID

		if dto.ProductSizes != nil && len(dto.ProductSizes) > 0 {
			inputSizeIDs := make([]uint, len(dto.ProductSizes))
			for j, productSize := range dto.ProductSizes {
				inputSizeIDs[j] = productSize.ProductSizeID
			}

			if err := s.validateProductSizesByProductID(inputSizeIDs, dto.ProductID); err != nil {
				return nil, utils.WrapError("failed to update store product", err)
			}
		}
	}

	ids, err := s.repo.CreateMultipleStoreProducts(storeProducts)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create %d store product: %w", len(dtos), err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	return ids, nil
}

func (s *storeProductService) UpdateStoreProduct(productID, storeProductID uint, dto *types.UpdateStoreProductDTO) error {
	if dto.ProductSizes != nil && len(dto.ProductSizes) > 0 {
		inputSizeIDs := make([]uint, len(dto.ProductSizes))
		for i, productSize := range dto.ProductSizes {
			inputSizeIDs[i] = productSize.ProductSizeID
		}

		if err := s.validateProductSizesByProductID(inputSizeIDs, productID); err != nil {
			wrappedErr := utils.WrapError("failed to update store product: ", err)
			return wrappedErr
		}
	}

	storeProduct := types.UpdateToStoreProduct(dto)
	err := s.repo.UpdateStoreProductByID(productID, storeProductID, storeProduct)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update store product", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *storeProductService) DeleteStoreProduct(storeID, storeProductID uint) error {
	err := s.repo.DeleteStoreProduct(storeID, storeProductID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete store product", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *storeProductService) GetStoreProductSizeById(storeID uint, productID uint) (*types.StoreProductSizeDTO, error) {
	size, err := s.repo.GetStoreProductSizeById(storeID, productID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store product size by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dto := types.MapToStoreProductSizeDTO(*size)
	return &dto, nil
}

func (s *storeProductService) validateProductSizesByProductID(productSizeIDs []uint, productID uint) error {
	productSizes, err := s.productRepo.GetProductSizesByProductID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("%w: %w", types.ErrValidationFailed, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if err := types.ValidateStoreProductSizes(productSizeIDs, productSizes); err != nil {
		wrappedErr := fmt.Errorf("%w: %w", types.ErrValidationFailed, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

/*func (s *storeProductService) GetStoreProductSizes(storeID uint, filter *types.StoreProductSizesFilterDTO) ([]types.StoreProductSizeDTO, error) {
	sizes, err := s.repo.GetStoreProductSizes(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store product sizes", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreProductSizeDTO, len(sizes))
	for i, size := range sizes {
		dtos[i] = types.MapToStoreProductSizeDTO(size)
	}
	return dtos, nil
}*/

func (s *storeProductService) CreateStoreProductSize(storeID, storeProductID uint, dto *types.CreateStoreProductSizeDTO) (uint, error) {
	size := types.CreateToStoreProductSize(dto)
	id, err := s.repo.CreateStoreProductSize(size)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store product size", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
}

func (s *storeProductService) UpdateStoreProductSize(storeID, productSizeID uint, dto *types.UpdateStoreProductSizeDTO) error {
	size := types.UpdateToStoreProductSize(dto)
	err := s.repo.UpdateProductSize(storeID, productSizeID, size)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update store product size", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *storeProductService) DeleteStoreProductSize(storeID, productSizeID uint) error {
	err := s.repo.DeleteStoreProductSize(storeID, productSizeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete store product size", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}
