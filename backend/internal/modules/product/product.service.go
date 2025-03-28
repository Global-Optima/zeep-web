package product

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProductService interface {
	GetProductByID(productID uint) (*types.ProductDetailsDTO, error)
	GetProducts(filter *types.ProductsFilterDto) ([]types.ProductDetailsDTO, error)
	CreateProduct(dto *types.CreateProductDTO) (uint, error)
	UpdateProduct(productID uint, dto *types.UpdateProductDTO) (*types.ProductDTO, error)
	DeleteProduct(productID uint) (*data.Product, error)

	GetProductSizesByProductID(productID uint) ([]types.ProductSizeDetailsDTO, error)
	GetProductSizeDetailsByID(productID uint) (*types.ProductSizeDetailsDTO, error)
	CreateProductSize(dto *types.CreateProductSizeDTO) (uint, error)
	UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error
	DeleteProductSize(productSizeID uint) error
}

type productService struct {
	repo                ProductRepository
	notificationService notifications.NotificationService
	storageRepo         storage.StorageRepository
	logger              *zap.SugaredLogger
}

func NewProductService(repo ProductRepository, notificationService notifications.NotificationService, storageRepo storage.StorageRepository, logger *zap.SugaredLogger) ProductService {
	return &productService{
		repo:                repo,
		logger:              logger,
		storageRepo:         storageRepo,
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

	exists, err := s.repo.CheckProductExists(dto.Name)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check product: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if exists {
		wrappedErr := fmt.Errorf("%w: product with the name %s already exists", types.ErrProductAlreadyExists, dto.Name)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	if dto.Image != nil || dto.Video != nil {
		imageKey, videoKey, err := s.storageRepo.ConvertAndUploadMedia(dto.Image, dto.Video)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to convert and upload media for productID = %d: %w", product.ID, err)
			s.logger.Error(wrappedErr)
			return 0, wrappedErr
		}
		if imageKey != "" {
			newImageKey := data.StorageImageKey(imageKey)
			product.ImageKey = &newImageKey
		}
		if videoKey != "" {
			newVideoKey := data.StorageVideoKey(videoKey)
			product.VideoKey = &newVideoKey
		}
	}

	productID, err := s.repo.CreateProduct(product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create product: %w", err)
		s.logger.Error(wrappedErr)
		go func() {
			if product.ImageKey != nil {
				err := s.storageRepo.DeleteImageFiles(*product.ImageKey)
				if err != nil {
					wrappedErr := fmt.Errorf("failed to delete image files: %w", err)
					s.logger.Error(wrappedErr)
				}
			}

			if product.VideoKey != nil {
				err := s.storageRepo.DeleteFile(product.VideoKey.GetConvertedVideoObjectKey())
				if err != nil {
					wrappedErr := fmt.Errorf("failed to delete video file: %w", err)
					s.logger.Error(wrappedErr)
				}
			}
		}()
		return 0, wrappedErr
	}

	notificationDetails := &details.NewProductDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID: productID,
		},
		ProductName: product.Name,
	}
	err = s.notificationService.NotifyNewProductAdded(notificationDetails)
	if err != nil {
		return 0, fmt.Errorf("failed to notify new product added: %w", err)
	}

	return productID, nil
}

func (s *productService) CreateProductSize(dto *types.CreateProductSizeDTO) (uint, error) {
	productSize := types.CreateToProductSizeModel(dto)

	productSizeID, err := s.repo.CreateProductSize(productSize)
	if err != nil {
		if errors.Is(err, types.ErrProductSizeUniqueName) {
			return 0, types.ErrProductSizeUniqueName
		}

		wrappedErr := fmt.Errorf("failed to create product size: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	product, err := s.repo.GetProductByID(productSize.ProductID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product by productSizeID: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	notificationDetails := &details.NewProductSizeDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID: productSizeID,
		},
		ProductName:     product.Name,
		ProductSizeName: productSize.Name,
		Size:            productSize.Size,
	}
	err = s.notificationService.NotifyNewProductSizeAdded(notificationDetails)
	if err != nil {
		return 0, fmt.Errorf("failed to notify new product size added: %w", err)
	}

	return productSizeID, nil
}

func (s *productService) UpdateProduct(productID uint, dto *types.UpdateProductDTO) (*types.ProductDTO, error) {
	product, err := s.repo.GetRawProductByID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch product: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	oldImageKey := product.ImageKey
	oldVideoKey := product.VideoKey

	err = types.UpdateProductToModel(dto, product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if dto.DeleteImage && dto.Image == nil {
		product.ImageKey = nil
	}

	if dto.DeleteVideo && dto.Video == nil {
		product.VideoKey = nil
	}

	if dto.Image != nil || dto.Video != nil {
		imageKey, videoKey, err := s.storageRepo.ConvertAndUploadMedia(dto.Image, dto.Video)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to upload media for productID = %d: %w", productID, err)
			s.logger.Error(wrappedErr)
			return nil, wrappedErr
		}
		if imageKey != "" {
			newImageKey := data.StorageImageKey(imageKey)
			product.ImageKey = &newImageKey
		}
		if videoKey != "" {
			newVideoKey := data.StorageVideoKey(videoKey)
			product.VideoKey = &newVideoKey
		}
	}

	err = s.repo.SaveProduct(product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update product: %w", err)
		s.logger.Error(wrappedErr)

		go func() {
			if product.ImageKey != nil && product.ImageKey != oldImageKey {
				if err := s.storageRepo.DeleteImageFiles(*product.ImageKey); err != nil {
					s.logger.Errorf("rollback failed to delete new image: %v", err)
				}
			}
			if product.VideoKey != nil && product.VideoKey != oldVideoKey {
				if err := s.storageRepo.DeleteFile(product.VideoKey.GetConvertedVideoObjectKey()); err != nil {
					s.logger.Errorf("rollback failed to delete new video: %v", err)
				}
			}
		}()

		return nil, wrappedErr
	}

	if oldImageKey != nil && oldImageKey != product.ImageKey {
		go func() {
			if err := s.storageRepo.MarkImagesAsDeleted(*oldImageKey); err != nil {
				s.logger.Errorf("failed to delete old image for productID = %d: %v", productID, err)
			}
		}()
	}

	if oldVideoKey != nil && oldVideoKey != product.VideoKey {
		go func() {
			if err := s.storageRepo.MarkFileAsDeleted(oldVideoKey.GetConvertedVideoObjectKey()); err != nil {
				s.logger.Errorf("failed to delete old video for productID = %d: %v", productID, err)
			}
		}()
	}

	changes := types.GenerateProductChanges(product, dto, product.ImageKey)

	if len(changes) != 0 {
		notificationDetails := &details.CentralCatalogUpdateDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           productID,
				FacilityName: "Central Catalog",
			},
			Changes: changes,
		}

		err = s.notificationService.NotifyCentralCatalogUpdate(notificationDetails)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to send notification: %w", err)
			s.logger.Error(wrappedErr)
		}
	}

	oldProductDto := types.MapToProductDTO(*product)
	return &oldProductDto, nil
}

func (s *productService) UpdateProductSize(productSizeID uint, dto *types.UpdateProductSizeDTO) error {
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

func (s *productService) DeleteProduct(productID uint) (*data.Product, error) {
	product, err := s.repo.DeleteProduct(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete product: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	go func() {
		if product.ImageKey != nil {
			err := s.storageRepo.MarkImagesAsDeleted(*product.ImageKey)
			if err != nil {
				wrappedErr := fmt.Errorf("failed to mark images as deleted: %w", err)
				s.logger.Error(wrappedErr)
			}
		}

		if product.VideoKey != nil {
			err = s.storageRepo.MarkFileAsDeleted(product.VideoKey.GetConvertedVideoObjectKey())
			if err != nil {
				wrappedErr := fmt.Errorf("failed to mark file as deleted: %w", err)
				s.logger.Error(wrappedErr)
			}
		}
	}()
	return product, nil
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
