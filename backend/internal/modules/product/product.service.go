package product

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProductService interface {
	GetStoreProducts(filter types.ProductFilterDao) ([]types.StoreProductDTO, error)
	GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error)
	CreateProduct(product *types.CreateStoreProduct) error
	UpdateProduct(product *types.UpdateStoreProduct) error
	DeleteProduct(productID uint) error
}

type productService struct {
	repo   ProductRepository
	logger *zap.SugaredLogger
}

func NewProductService(repo ProductRepository, logger *zap.SugaredLogger) ProductService {
	return &productService{
		repo:   repo,
		logger: logger,
	}
}

func (s *productService) GetStoreProducts(filter types.ProductFilterDao) ([]types.StoreProductDTO, error) {
	products, err := s.repo.GetStoreProducts(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	productDTOs := make([]types.StoreProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = types.MapToStoreProductDTO(product)
	}

	return productDTOs, nil
}

func (s *productService) GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error) {
	productDetails, err := s.repo.GetStoreProductDetails(storeID, productID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve product", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	if productDetails == nil {
		return nil, nil
	}

	return productDetails, nil
}

func (s *productService) CreateProduct(dto *types.CreateStoreProduct) error {
	product := types.CreateToProductModel(dto)

	productID, err := s.repo.CreateProduct(product)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create product", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	productSizes := types.ToProductSizesModels(dto.ProductSizes, productID)
	if err := s.repo.CreateProductSizes(productID, productSizes); err != nil {
		wrappedErr := utils.WrapError("failed to create product sizes", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if err := AssignAdditives(productID, productSizes, dto.Additives, s.repo); err != nil {
		return err
	}

	return nil
}

func (s *productService) UpdateProduct(dto *types.UpdateStoreProduct) error {
	product := types.UpdateToProductModel(dto)

	if err := s.repo.UpdateProduct(product); err != nil {
		wrappedErr := utils.WrapError("failed to update product", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *productService) DeleteProduct(productID uint) error {
	return s.repo.DeleteProduct(productID)
}

func AssignAdditives(
	productID uint,
	productSizes []data.ProductSize,
	additives []types.SelectedAdditiveDTO,
	repo ProductRepository,
) error {
	var defaultAdditives []data.DefaultProductAdditive
	productAdditives := make(map[uint][]data.ProductAdditive)

	for _, additive := range additives {
		if additive.IsDefault {
			defaultAdditives = append(defaultAdditives, data.DefaultProductAdditive{
				ProductID:  productID,
				AdditiveID: additive.AdditiveID,
			})
		} else {
			for _, size := range productSizes {
				productAdditives[size.ID] = append(productAdditives[size.ID], data.ProductAdditive{
					ProductSizeID: size.ID,
					AdditiveID:    additive.AdditiveID,
				})
			}
		}
	}

	if len(defaultAdditives) > 0 {
		if err := repo.AssignDefaultAdditives(productID, defaultAdditives); err != nil {
			return err
		}
	}

	for sizeID, additives := range productAdditives {
		if len(additives) > 0 {
			if err := repo.AssignProductAdditives(sizeID, additives); err != nil {
				return err
			}
		}
	}

	return nil
}
