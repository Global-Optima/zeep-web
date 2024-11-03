package product

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"gorm.io/gorm"
)

type ProductService interface {
	GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductDTO, error)
	SearchStoreProducts(storeID uint, searchQuery, category string, offset, limit int) ([]types.ProductDTO, error)
	GetStoreProductDetails(storeID, productID uint) (*types.ProductDTO, error)
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductDTO, error) {

	if offset < 0 || limit <= 0 {
		return nil, errors.New("offset must be >= 0 and limit must be > 0")
	}

	products, err := s.repo.GetStoreProducts(storeID, category, offset, limit)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no products found for store ID %d", storeID)
		}
		return nil, fmt.Errorf("error retrieving products: %w", err)
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no available products found for store ID %d and category '%s'", storeID, category)
	}

	return products, nil
}

func (s *productService) SearchStoreProducts(storeID uint, searchQuery, category string, offset, limit int) ([]types.ProductDTO, error) {

	if searchQuery == "" {
		return nil, errors.New("search query cannot be empty")
	}
	if offset < 0 || limit <= 0 {
		return nil, errors.New("offset must be >= 0 and limit must be > 0")
	}

	products, err := s.repo.SearchStoreProducts(storeID, searchQuery, category, offset, limit)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no products found matching '%s' for store ID %d", searchQuery, storeID)
		}
		return nil, fmt.Errorf("error searching products: %w", err)
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no products found for search query '%s' in store ID %d and category '%s'", searchQuery, storeID, category)
	}

	return products, nil
}

func (s *productService) GetStoreProductDetails(storeID, productID uint) (*types.ProductDTO, error) {

	if storeID == 0 || productID == 0 {
		return nil, errors.New("storeID and productID must be greater than 0")
	}

	product, err := s.repo.GetStoreProductDetails(storeID, productID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product ID %d not found for store ID %d", productID, storeID)
		}
		return nil, fmt.Errorf("error retrieving product details: %w", err)
	}

	return product, nil
}
