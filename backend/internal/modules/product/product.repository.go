package product

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductSizeWithProduct(productSizeID uint) (*data.ProductSize, error)
	GetStoreProducts(filter types.ProductsFilterDto) ([]data.Product, error)
	GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error)
	DeleteProduct(productID uint) error

	UpdateProduct(product *data.Product) error
	CreateProduct(product *data.Product) (uint, error)
	CreateProductSizes(productID uint, sizes []data.ProductSize) error
	AssignProductAdditives(productSizeID uint, additives []data.ProductAdditive) error
	AssignDefaultAdditives(productID uint, additives []data.DefaultProductAdditive) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProductSizeWithProduct(productSizeID uint) (*data.ProductSize, error) {
	var productSize data.ProductSize
	err := r.db.Preload("Product").First(&productSize, productSizeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch ProductSize ID %d: %w", productSizeID, err)
	}

	if productSize.Product.ID == 0 || productSize.Product.Name == "" {
		return nil, fmt.Errorf("product size with ID %d has no valid associated product", productSizeID)
	}

	return &productSize, nil
}

func (r *productRepository) GetStoreProducts(filter types.ProductsFilterDto) ([]data.Product, error) {
	var products []data.Product

	query := r.db.
		Model(&data.Product{}).
		Joins("JOIN store_products ON store_products.product_id = products.id").
		Preload("ProductSizes", "is_default = TRUE")

	if filter.StoreID != nil {
		query = query.Where("store_products.store_id = ? AND store_products.is_available = TRUE", filter.StoreID)
	}

	if filter.CategoryID != nil {
		query = query.Where("products.category_id = ?", *filter.CategoryID)
	}

	if filter.SearchTerm != nil {
		searchPattern := "%" + *filter.SearchTerm + "%"
		query = query.Where("(products.name ILIKE ? OR products.description ILIKE ?)", searchPattern, searchPattern)
	}

	err := query.
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	if products == nil {
		products = []data.Product{}
	}

	return products, nil
}

func (r *productRepository) GetStoreProductDetails(storeID uint, productID uint) (*types.StoreProductDetailsDTO, error) {
	var product data.Product

	// Preload only ProductSizes and DefaultAdditives, exclude additive categories and their additives
	err := r.db.
		Preload("ProductSizes").
		Preload("DefaultAdditives.Additive").
		Preload("Category").
		Where("id = ?", productID).
		First(&product).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Verify that the product is available in the specified store
	var storeProduct data.StoreProduct
	err = r.db.
		Where("product_id = ? AND store_id = ? AND is_available = ?", productID, storeID, true).
		First(&storeProduct).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Fetch default additives
	var defaultAdditives []data.DefaultProductAdditive
	err = r.db.
		Preload("Additive").
		Where("product_id = ?", productID).
		Find(&defaultAdditives).
		Error

	if err != nil {
		return nil, err
	}

	// Map models to DTO
	productDTO := types.MapToStoreProductDetailsDTO(&product, defaultAdditives)

	return productDTO, nil
}

func (r *productRepository) CreateProduct(product *data.Product) (uint, error) {
	if err := r.db.Create(product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (r *productRepository) CreateProductSizes(productID uint, sizes []data.ProductSize) error {
	for i := range sizes {
		sizes[i].ProductID = productID
	}
	return r.db.Create(&sizes).Error
}

func (r *productRepository) AssignProductAdditives(productSizeID uint, additives []data.ProductAdditive) error {
	for i := range additives {
		additives[i].ProductSizeID = productSizeID
	}
	return r.db.Create(&additives).Error
}

func (r *productRepository) AssignDefaultAdditives(productID uint, additives []data.DefaultProductAdditive) error {
	for i := range additives {
		additives[i].ProductID = productID
	}
	return r.db.Create(&additives).Error
}

func (r *productRepository) UpdateProduct(product *data.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateProductSize(productSize *data.ProductSize) error {
	return r.db.Save(productSize).Error
}

func (r *productRepository) DeleteProduct(productID uint) error {
	return r.db.Where("id = ?", productID).Delete(&data.Product{}).Error
}

func (r *productRepository) DeleteProductSize(productSizeID uint) error {
	return r.db.Where("id = ?", productSizeID).Delete(&data.ProductSize{}).Error
}
