package product

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
)

type ProductService interface {
	GetStoreProducts(filter types.ProductFilterDao) ([]types.StoreProductDTO, error)
	GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error)
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetStoreProducts(filter types.ProductFilterDao) ([]types.StoreProductDTO, error) {
	products, err := s.repo.GetStoreProducts(filter)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]types.StoreProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = mapToStoreProductDTO(product)
	}

	return productDTOs, nil
}

func (s *productService) GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error) {
	productDetails, err := s.repo.GetStoreProductDetails(storeID, productID)
	if err != nil {
		return nil, err
	}
	if productDetails == nil {
		return nil, nil
	}

	return productDetails, nil
}

func mapToStoreProductDetailsDTO(product *data.Product, defaultAdditives []data.DefaultProductAdditive) *types.StoreProductDetailsDTO {
	dto := &types.StoreProductDetailsDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
	}

	for _, size := range product.ProductSizes {
		sizeDTO := types.ProductSizeDTO{
			ID:        size.ID,
			Name:      size.Name,
			BasePrice: size.BasePrice,
			Measure:   size.Measure,
		}

		additiveCategoryMap := make(map[uint]*additiveTypes.AdditiveCategoryDTO)
		for _, pa := range size.Additives {
			additive := pa.Additive
			categoryID := additive.AdditiveCategoryID

			if _, exists := additiveCategoryMap[categoryID]; !exists {
				additiveCategoryMap[categoryID] = &additiveTypes.AdditiveCategoryDTO{
					ID:        categoryID,
					Name:      additive.Category.Name,
					Additives: []additiveTypes.AdditiveDTO{},
				}
			}

			additiveDTO := additiveTypes.AdditiveDTO{
				ID:          additive.ID,
				Name:        additive.Name,
				Description: additive.Description,
				Price:       additive.BasePrice,
				ImageURL:    additive.ImageURL,
			}

			additiveCategoryMap[categoryID].Additives = append(additiveCategoryMap[categoryID].Additives, additiveDTO)
		}

		dto.Sizes = append(dto.Sizes, sizeDTO)
	}

	for _, da := range defaultAdditives {
		additive := da.Additive
		additiveDTO := additiveTypes.AdditiveDTO{
			ID:          additive.ID,
			Name:        additive.Name,
			Description: additive.Description,
			Price:       additive.BasePrice,
			ImageURL:    additive.ImageURL,
		}
		dto.DefaultAdditives = append(dto.DefaultAdditives, additiveDTO)
	}

	return dto
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
