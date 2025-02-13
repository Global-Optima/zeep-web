package product

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProductService interface {
	GetProductByID(productID uint) (*types.ProductDetailsDTO, error)
	GetProducts(filter *types.ProductsFilterDto) ([]types.ProductDetailsDTO, error)
	CreateProduct(product *types.CreateProductDTO) (uint, error)
	UpdateProduct(productID uint, dto *types.UpdateProductDTO) error
	DeleteProduct(productID uint) error

	GetProductSizesByProductID(productID uint) ([]types.ProductSizeDetailsDTO, error)
	GetProductSizeDetailsByID(productID uint) (*types.ProductSizeDetailsDTO, error)
	CreateProductSize(dto *types.CreateProductSizeDTO) (uint, error)
	UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error
	DeleteProductSize(productSizeID uint) error
}

type productService struct {
	repo                ProductRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewProductService(repo ProductRepository, notificationService notifications.NotificationService, logger *zap.SugaredLogger) ProductService {
	return &productService{
		repo:                repo,
		logger:              logger,
		notificationService: notificationService,
	}
}

func (s *productService) GetProducts(filter *types.ProductsFilterDto) ([]types.ProductDetailsDTO, error) {
	products, err := s.repo.GetProducts(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	productDTOs := make([]types.ProductDetailsDTO, len(products))
	for i, product := range products {
		productDTOs[i] = *types.MapToProductDetailsDTO(&product)
	}

	return productDTOs, nil
}

func (s *productService) GetProductByID(productID uint) (*types.ProductDetailsDTO, error) {
	product, err := s.repo.GetProductByID(productID)
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
	productBefore, err := s.repo.GetProductByID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch product: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	product := types.UpdateProductToModel(dto, productBefore)

	err = s.repo.UpdateProduct(productID, product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	changes := types.GenerateProductChanges(productBefore, dto)

	notificationDetails := &details.CentralCatalogUpdateDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           productID,
			FacilityName: "Central Catalog",
		},
		Changes: changes,
	}

	if len(changes) != 0 {
		err = s.notificationService.NotifyCentralCatalogUpdate(notificationDetails)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to send notification: %w", err)
			s.logger.Error(wrappedErr)
		}
	}

	return nil
}

func (s *productService) UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error {
	if dto.IsDefault != nil && !*dto.IsDefault {
		wrappedErr := fmt.Errorf("failed to update product size: cannot set isDefault to false")
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	updateModels := types.UpdateProductSizeToModels(dto)

	productSize, err := s.repo.GetProductSizeById(productSizeID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch product: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.repo.UpdateProductSizeWithAssociations(productSizeID, updateModels)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product size: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if dto.BasePrice != nil && productSize.BasePrice != *dto.BasePrice {
		notificationDetails := &details.PriceChangeNotificationDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID: productSizeID,
			},
			ProductName: productSize.Product.Name,
			OldPrice:    productSize.BasePrice,
			NewPrice:    *dto.BasePrice,
		}
		err = s.notificationService.NotifyPriceChange(notificationDetails)
		if err != nil {
			return fmt.Errorf("failed to notify price change: %w", err)
		}
	}

	return nil
}

func (s *productService) GetProductSizesByProductID(productID uint) ([]types.ProductSizeDetailsDTO, error) {
	productSizes, err := s.repo.GetProductSizesByProductID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product sizes: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.ProductSizeDetailsDTO, len(productSizes))
	for i, productSize := range productSizes {
		dtos[i] = types.MapToProductSizeDetails(productSize)
	}

	return dtos, nil
}

func (s *productService) GetProductSizeDetailsByID(productID uint) (*types.ProductSizeDetailsDTO, error) {
	productSize, err := s.repo.GetProductSizeDetailsByID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product sizes: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dto := types.MapToProductSizeDetails(*productSize)

	return &dto, nil
}

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
