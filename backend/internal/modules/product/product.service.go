package product

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProductService interface {
	GetProductByID(productID uint) (*types.ProductDetailsDTO, error)
	GetProducts(filter *types.ProductsFilterDto) ([]types.ProductDetailsDTO, error)
	CreateProduct(dto *types.CreateProductDTO, img, vid *multipart.FileHeader) (uint, error)
	UpdateProduct(productID uint, dto *types.UpdateProductDTO, img, vid *multipart.FileHeader) error
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
		key := fmt.Sprintf("%s/%s", storage.IMAGES_CONVERTED_STORAGE_REPO_KEY, product.ImageURL)
		imageUrl, err := s.storageRepo.GetFileURL(key)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to retrieve product image url for productID = %d: %w", product.ID, err)
			s.logger.Error(wrappedErr)
		}
		productDTOs[i] = *types.MapToProductDetailsDTO(&product, imageUrl, product.VideoURL)
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

	key := fmt.Sprintf("%s/%s", storage.IMAGES_CONVERTED_STORAGE_REPO_KEY, product.ImageURL)
	imageUrl, err := s.storageRepo.GetFileURL(key)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve product image url for productID = %d: %w", productID, err)
		s.logger.Error(wrappedErr)
	}

	return types.MapToProductDetailsDTO(product, imageUrl, product.VideoURL), nil
}

func (s *productService) CreateProduct(
	dto *types.CreateProductDTO,
	img *multipart.FileHeader,
	vid *multipart.FileHeader,
) (uint, error) {
	startTime := time.Now()
	fmt.Println("[TIMER] Processing started...")

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

	imageUrl, videoUrl, err := s.storageRepo.ConvertAndUploadMedia(img, vid)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to convert and upload media for productID = %d: %w", product.ID, err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	product.ImageURL = imageUrl
	product.VideoURL = videoUrl

	productID, err := s.repo.CreateProduct(product)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create product: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	fmt.Printf("[TIMER] Total process finished. Took %v\n", time.Since(startTime))

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

func (s *productService) UpdateProduct(productID uint, dto *types.UpdateProductDTO, img, vid *multipart.FileHeader) error {
	if strings.TrimSpace(dto.Name) != "" {
		exists, err := s.repo.CheckProductExists(dto.Name)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to check product: %w", err)
			s.logger.Error(wrappedErr)
			return wrappedErr
		}
		if exists {
			wrappedErr := fmt.Errorf("%w: product with the name %s already exists", types.ErrProductAlreadyExists, dto.Name)
			s.logger.Error(wrappedErr)
			return wrappedErr
		}
	}

	product := types.UpdateProductToModel(dto)

	productBefore, err := s.repo.GetProductByID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch product: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	imageUrl, videoUrl, err := s.storageRepo.ConvertAndUploadMedia(img, vid)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to convert and upload media for productID = %d: %w", productID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	product.ImageURL = imageUrl
	product.VideoURL = videoUrl

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

	err = s.notificationService.NotifyCentralCatalogUpdate(notificationDetails)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to send notification: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
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
		for _, productSizeAdditive := range productSize.Additives {
			key := fmt.Sprintf("%s/%s", storage.IMAGES_CONVERTED_STORAGE_REPO_KEY, productSizeAdditive.Additive.ImageURL)
			imageUrl, err := s.storageRepo.GetFileURL(key)
			if err != nil {
				wrappedErr := fmt.Errorf("failed to retrieve product image url for productID = %d: %w", productID, err)
				s.logger.Error(wrappedErr)
			}
			productSizeAdditive.Additive.ImageURL = imageUrl
		}

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

func (s *productService) getPresignedURLsConcurrently(products []data.Product) []string {
	numProducts := len(products)
	imageUrls := make([]string, numProducts)

	var wg sync.WaitGroup
	errChan := make(chan error, numProducts)

	for i, product := range products {
		wg.Add(1)

		go func(idx int, imageKey string) {
			defer wg.Done()

			presignedURL, err := s.storageRepo.GetPresignedURL(imageKey)
			if err != nil {
				s.logger.Errorf("Failed to get presigned URL for image %s: %v", imageKey, err)
				errChan <- err // Send error to channel
				return
			}

			imageUrls[idx] = presignedURL
		}(i, product.ImageURL)
	}

	wg.Wait()
	close(errChan)

	// Check for errors (if needed, handle them based on your use case)
	if len(errChan) > 0 {
		s.logger.Warnf("Some presigned URLs failed to generate.")
	}

	return imageUrls
}
