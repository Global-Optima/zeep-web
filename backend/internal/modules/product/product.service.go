package product

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProductService interface {
	GetProductDetails(productID uint) (*types.ProductDetailsDTO, error)

	GetProducts(filter *types.ProductsFilterDto) ([]types.ProductDTO, error)
	CreateProduct(product *types.CreateProductDTO) (uint, error)
	UpdateProduct(productID uint, dto *types.UpdateProductDTO) error
	DeleteProduct(productID uint) error

	GetProductSizesByProductID(productID uint) ([]types.ProductSizeDTO, error)
	CreateProductSize(dto *types.CreateProductSizeDTO) (uint, error)
	UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error
	DeleteProductSize(productSizeID uint) error
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

func (s *productService) GetProducts(filter *types.ProductsFilterDto) ([]types.ProductDTO, error) {
	products, err := s.repo.GetProducts(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	productDTOs := make([]types.ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = types.MapToProductDTO(product)
	}

	return productDTOs, nil
}

func (s *productService) GetProductDetails(productID uint) (*types.ProductDetailsDTO, error) {
	product, err := s.repo.GetProductDetails(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve product for productID = %d: %w", productID, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.MapToProductDetailsDTO(product), nil
}

func (s *productService) CreateProduct(dto *types.CreateProductDTO) (uint, error) {
	product := types.CreateToProductModel(dto)

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

	err := s.repo.UpdateProduct(productID, product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product additives: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *productService) UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error {
	updateModels := types.UpdateProductSizeToModels(dto)

	err := s.repo.UpdateProductSizeWithAssociations(productSizeID, updateModels)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product additives: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *productService) GetProductSizesByProductID(productID uint) ([]types.ProductSizeDTO, error) {
	productSizes, err := s.repo.GetProductSizesByProductID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product sizes: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.ProductSizeDTO, len(productSizes))
	for i, productSize := range productSizes {
		dtos[i] = types.MapToProductSizeDTO(productSize)
	}

	return dtos, nil
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
	err := s.repo.DeleteProduct(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete product: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *productService) DeleteProductSize(productSizeID uint) error {
	err := s.repo.DeleteProductSize(productSizeID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete product size: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}
