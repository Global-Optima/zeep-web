package product

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
)

type ProductService interface {
	GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductCatalogDTO, error)
	SearchStoreProducts(storeID uint, searchQuery, category string, offset, limit int) ([]types.ProductCatalogDTO, error)
	GetStoreProductDetails(storeID, productID uint) (*types.ProductDTO, error)
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductCatalogDTO, error) {
	results, err := s.repo.GetStoreProducts(storeID, category, offset, limit)
	if err != nil {
		return nil, err
	}

	var productDTOs []types.ProductCatalogDTO
	for _, result := range results {
		productDTO := types.ProductCatalogDTO{
			ProductID:          result.ProductID,
			ProductName:        result.ProductName,
			ProductDescription: result.ProductDescription,
			Category:           result.CategoryName,
			ProductImageURL:    result.ProductImageURL,
			Price:              result.BasePrice,
			IsAvailable:        result.IsAvailable,
			IsOutOfStock:       result.IsOutOfStock,
		}
		productDTOs = append(productDTOs, productDTO)
	}

	return productDTOs, nil
}

func (s *productService) SearchStoreProducts(storeID uint, searchQuery, category string, offset, limit int) ([]types.ProductCatalogDTO, error) {
	results, err := s.repo.SearchStoreProducts(storeID, searchQuery, category, offset, limit)
	if err != nil {
		return nil, err
	}

	var productDTOs []types.ProductCatalogDTO
	for _, result := range results {
		productDTO := types.ProductCatalogDTO{
			ProductID:          result.ProductID,
			ProductName:        result.ProductName,
			ProductDescription: result.ProductDescription,
			Category:           result.CategoryName,
			ProductImageURL:    result.ProductImageURL,
			Price:              result.BasePrice,
			IsAvailable:        result.IsAvailable,
			IsOutOfStock:       result.IsOutOfStock,
		}
		productDTOs = append(productDTOs, productDTO)
	}

	return productDTOs, nil
}

func (s *productService) GetStoreProductDetails(storeID, productID uint) (*types.ProductDTO, error) {
	result, err := s.repo.GetStoreProductDetails(storeID, productID)
	if err != nil {
		return nil, err
	}

	productDTO := mapProductDAOToDTO(result)

	return productDTO, nil
}

func mapProductDAOToDTO(result *types.ProductDAO) *types.ProductDTO {
	productDTO := &types.ProductDTO{
		ProductID:          result.ProductID,
		ProductName:        result.ProductName,
		ProductDescription: result.ProductDescription,
		Category:           result.CategoryName,
		ProductImageURL:    result.ProductImageURL,
		ProductVideoURL:    result.ProductVideoURL,
		Price:              result.BasePrice,
		IsAvailable:        result.IsAvailable,
		IsOutOfStock:       result.IsOutOfStock,
		Nutrition: types.NutritionDTO{
			Calories:      result.Nutrition.Calories,
			Fat:           result.Nutrition.Fat,
			Carbohydrates: result.Nutrition.Carbohydrates,
			Proteins:      result.Nutrition.Proteins,
		},
	}

	for _, additive := range result.Additives {
		productDTO.Additives = append(productDTO.Additives, types.AdditivesDTO(additive))
	}

	for _, defaultAdditive := range result.DefaultAdditives {
		productDTO.DefaultAdditives = append(productDTO.DefaultAdditives, types.AdditivesDTO(defaultAdditive))
	}

	for _, size := range result.Sizes {
		productDTO.Sizes = append(productDTO.Sizes, types.SizeDTO(size))
	}

	return productDTO
}
