package product

import (
	"encoding/json"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ProductService interface {
	GetStoreProducts(c *gin.Context, storeID, categoryID uint, searchQuery string, limit, offset int) ([]types.StoreProductDTO, error)
	GetStoreProductDetails(c *gin.Context, storeID, productID uint) (*types.StoreProductDetailsDTO, error)
}

type productService struct {
	repo        ProductRepository
	redisClient *redis.Client
}

func NewProductService(repo ProductRepository, redisClient *redis.Client) ProductService {
	return &productService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *productService) GetStoreProducts(c *gin.Context, storeID, categoryID uint, searchQuery string, limit, offset int) ([]types.StoreProductDTO, error) {
	cacheKey := utils.GenerateCacheKey("products:store", storeID, "category", categoryID, "search", searchQuery, "limit", limit, "offset", offset)

	cachedData, err := s.redisClient.Get(c, cacheKey).Result()
	if err == nil {
		var products []types.StoreProductDTO
		if err := json.Unmarshal([]byte(cachedData), &products); err == nil && len(products) > 0 {
			return products, nil
		}
	}

	products, err := s.repo.GetStoreProducts(storeID, categoryID, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]types.StoreProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = mapToStoreProductDTO(product)
	}

	data, _ := json.Marshal(productDTOs)
	ttl := utils.GetTTL("warm")
	s.redisClient.Set(c, cacheKey, data, ttl)

	return productDTOs, nil
}

func (s *productService) GetStoreProductDetails(c *gin.Context, storeID, productID uint) (*types.StoreProductDetailsDTO, error) {
	cacheKey := utils.GenerateCacheKey("product:store", storeID, "product", productID)

	cachedData, err := s.redisClient.Get(c, cacheKey).Result()
	if err == nil {
		var productDetails types.StoreProductDetailsDTO
		if err := json.Unmarshal([]byte(cachedData), &productDetails); err == nil {
			return &productDetails, nil
		}
	}

	product, err := s.repo.GetStoreProductDetails(storeID, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	productDetailsDTO := mapToStoreProductDetailsDTO(product)
	data, _ := json.Marshal(productDetailsDTO)
	ttl := utils.GetTTL("warm")
	s.redisClient.Set(c, cacheKey, data, ttl)

	return productDetailsDTO, nil
}

func mapToStoreProductDetailsDTO(product *data.Product) *types.StoreProductDetailsDTO {
	sizes := make([]types.ProductSizeDTO, len(product.ProductSizes))
	for i, size := range product.ProductSizes {
		sizes[i] = types.ProductSizeDTO{
			ID:        size.ID,
			Name:      size.Name,
			BasePrice: size.BasePrice,
			Measure:   size.Measure,
		}
	}

	defaultAdditives := make([]types.ProductAdditiveDTO, len(product.DefaultAdditives))
	for i, pa := range product.DefaultAdditives {
		additive := pa.Additive
		defaultAdditives[i] = types.ProductAdditiveDTO{
			ID:          additive.ID,
			Name:        additive.Name,
			Description: additive.Description,
			ImageURL:    additive.ImageURL,
		}
	}

	recipeSteps := make([]types.RecipeStepDTO, len(product.RecipeSteps))
	for i, step := range product.RecipeSteps {
		recipeSteps[i] = types.RecipeStepDTO{
			ID:          step.ID,
			Description: step.Description,
			ImageURL:    step.ImageURL,
			Step:        step.Step,
		}
	}

	return &types.StoreProductDetailsDTO{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		ImageURL:         product.ImageURL,
		Sizes:            sizes,
		DefaultAdditives: defaultAdditives,
		RecipeSteps:      recipeSteps,
	}
}

func mapToStoreProductDTO(product data.Product) types.StoreProductDTO {
	var basePrice float64 = 0
	if len(product.ProductSizes) > 0 {
		basePrice = product.ProductSizes[0].BasePrice
	}

	return types.StoreProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		BasePrice:   basePrice,
	}
}
