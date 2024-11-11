package product

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
)

type ProductService interface {
	GetStoreProducts(storeID, categoryID uint, searchQuery string, limit, offset int) ([]types.StoreProductDTO, error)
	GetStoreProductDetails(storeID, productID uint) (*types.StoreProductDetailsDTO, error)
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetStoreProducts(storeID, categoryID uint, searchQuery string, limit, offset int) ([]types.StoreProductDTO, error) {
	products, err := s.repo.GetStoreProducts(storeID, categoryID, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]types.StoreProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = mapToStoreProductDTO(product)
	}

	return productDTOs, nil
}

func (s *productService) GetStoreProductDetails(storeID, productID uint) (*types.StoreProductDetailsDTO, error) {
	product, err := s.repo.GetStoreProductDetails(storeID, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	productDetailsDTO := mapToStoreProductDetailsDTO(product)

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
