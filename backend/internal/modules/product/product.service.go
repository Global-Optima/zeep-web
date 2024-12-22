package product

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProductService interface {
	GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error)

	GetProducts(filter *types.ProductsFilterDto) ([]types.StoreProductDTO, error)
	CreateProduct(product *types.CreateProductDTO) (uint, error)
	UpdateProduct(productID uint, dto *types.UpdateProductDTO) error
	//DeleteProduct(productID uint) error

	CreateProductSize(dto *types.CreateProductSizeDTO) (uint, error)
	UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error
	//DeleteProductSize(productSizeID uint) error
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

func (s *productService) GetProducts(filter *types.ProductsFilterDto) ([]types.StoreProductDTO, error) {
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
		wrappedErr := fmt.Errorf("failed to retrieve product for storeID = %d, productID = %d: %w", storeID, productID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return productDetails, nil
}

func (s *productService) CreateProduct(dto *types.CreateProductDTO) (uint, error) {
	product := types.CreateToProductModel(dto)

	ids := product.DefaultAdditives

	s.logger.Warn(ids)
	productID, err := s.repo.CreateProduct(product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create product: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return productID, nil
}

func (s *productService) CreateProductSize(dto *types.CreateProductSizeDTO) (uint, error) {
	productSize := types.CreateToProductSizeModel(dto)

	productSizeID, err := s.repo.CreateProductSize(productSize)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create product size: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return productSizeID, nil
}

func (s *productService) UpdateProduct(productID uint, dto *types.UpdateProductDTO) error {
	product := types.UpdateProductToModel(dto)

	err := s.repo.UpdateProduct(productID, product, dto.DefaultAdditives)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product additives: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *productService) UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error {
	productSize := types.UpdateProductSizeToModel(dto)

	err := s.repo.UpdateProductSizeWithAssociations(productSizeID, productSize, dto.Additives, dto.Ingredients)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product additives: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

/*func (s *productService) UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error {
	additiveIDs := dto.AdditiveIDList

	err := s.repo.UpdateProductSizeAdditives(productSizeID, additiveIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product additives: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *productService) CreateProductWithSizes(dto *types.CreateProductWithAttachesDTO) (uint, error) {
	productID, err := s.repo.CreateProductWithSizes(dto)
	if err != nil {
		wrappedError := utils.WrapError("failed to create product", err)
		s.logger.Error(wrappedError)
		return 0, wrappedError
	}

	return productID, nil
}*/

func (s *productService) DeleteProduct(productID uint) error {
	return s.repo.DeleteProduct(productID)
}
