package product

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	CreateProduct(product *data.Product) (uint, error)
	GetProducts(filter *types.ProductsFilterDto) ([]data.Product, error)
	GetProductByID(productID uint) (*data.Product, error)
	UpdateProduct(id uint, product *data.Product) error
	DeleteProduct(productID uint) error

	CreateProductSize(createModels *data.ProductSize) (uint, error)
	GetProductSizesByProductID(productID uint) ([]data.ProductSize, error)
	GetProductSizeById(productSizeID uint) (*data.ProductSize, error)
	GetProductSizeDetailsByID(productSizeID uint) (*data.ProductSize, error)
	UpdateProductSizeWithAssociations(id uint, updateModels *types.ProductSizeModels) error
	DeleteProductSize(productID uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProductSizeById(productSizeID uint) (*data.ProductSize, error) {
	var productSize data.ProductSize
	err := r.db.Preload("Product").
		Preload("Unit").
		First(&productSize, productSizeID).Error

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

func (r *productRepository) GetProductSizeDetailsByID(productSizeID uint) (*data.ProductSize, error) {
	var productSize data.ProductSize

	err := r.db.Model(&data.ProductSize{}).
		Preload("Unit").
		Preload("Additives.Additive.Category").
		Preload("Additives.Additive.Unit").
		Preload("ProductSizeIngredients.Ingredient.IngredientCategory").
		Preload("ProductSizeIngredients.Ingredient.Unit").
		First(&productSize, productSizeID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch ProductSize ID %d: %w", productSizeID, err)
	}

	return &productSize, nil
}

func (r *productRepository) GetProducts(filter *types.ProductsFilterDto) ([]data.Product, error) {
	var products []data.Product

	query := r.db.
		Model(&data.Product{}).
		Preload("Category").
		Preload("ProductSizes.Unit")

	if filter.CategoryID != nil {
		query = query.Where("products.category_id = ?", *filter.CategoryID)
	}

	if filter.Search != nil {
		searchPattern := "%" + *filter.Search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", searchPattern, searchPattern)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Product{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	if products == nil {
		products = []data.Product{}
	}

	return products, nil
}

func (r *productRepository) GetProductByID(productID uint) (*data.Product, error) {
	var product data.Product

	err := r.db.
		Preload("ProductSizes.Unit").
		Preload("Category").
		Where("id = ?", productID).
		First(&product).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) CreateProduct(product *data.Product) (uint, error) {
	if err := r.db.Create(product).Error; err != nil {
		return 0, err
	}

	return product.ID, nil
}

func (r *productRepository) CreateProductSize(productSize *data.ProductSize) (uint, error) {

	if err := r.db.Create(productSize).Error; err != nil {
		return 0, err
	}

	return productSize.ID, nil
}

func (r *productRepository) UpdateProductSizeWithAssociations(id uint, updateModels *types.ProductSizeModels) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels != nil {
			if err := r.handleDefaultSizeUpdate(tx, id, updateModels.ProductSize.IsDefault); err != nil {
				return err
			}

			if err := r.updateProductSize(tx, id, *updateModels.ProductSize); err != nil {
				return err
			}
		}

		if updateModels.Additives != nil {
			if err := r.updateAdditives(tx, id, updateModels.Additives); err != nil {
				return err
			}
		}

		if updateModels.Ingredients != nil {
			if err := r.updateIngredients(tx, id, updateModels.Ingredients); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *productRepository) handleDefaultSizeUpdate(tx *gorm.DB, id uint, isDefault bool) error {
	if !isDefault {
		return nil
	}

	var productSize data.ProductSize
	if err := tx.Select("product_id").First(&productSize, id).Error; err != nil {
		return fmt.Errorf("failed to get product size: %w", err)
	}

	if err := tx.Model(&data.ProductSize{}).
		Where("product_id = ? AND id != ?", productSize.ProductID, id).
		Update("is_default", false).Error; err != nil {
		return fmt.Errorf("failed to update other product sizes: %w", err)
	}

	return nil
}

func (r *productRepository) updateProductSize(tx *gorm.DB, id uint, productSize data.ProductSize) error {
	if err := tx.Model(&data.ProductSize{}).
		Where("id = ?", id).
		Updates(productSize).Error; err != nil {
		return fmt.Errorf("failed to update product size: %w", err)
	}
	return nil
}

func (r *productRepository) updateAdditives(tx *gorm.DB, productSizeID uint, additives []data.ProductSizeAdditive) error {
	if err := tx.Where("product_size_id = ?", productSizeID).Delete(&data.ProductSizeAdditive{}).Error; err != nil {
		return fmt.Errorf("failed to delete additives: %w", err)
	}

	newAdditives := make([]data.ProductSizeAdditive, len(additives))
	for i, additive := range additives {
		newAdditives[i] = data.ProductSizeAdditive{
			ProductSizeID: productSizeID,
			AdditiveID:    additive.AdditiveID,
			IsDefault:     additive.IsDefault,
		}
		logrus.Info(newAdditives[i].IsDefault)
	}

	if err := tx.Create(newAdditives).Error; err != nil {
		return fmt.Errorf("failed to create additives: %w", err)
	}

	return nil
}

func (r *productRepository) updateIngredients(tx *gorm.DB, productSizeID uint, ingredients []data.ProductSizeIngredient) error {
	if err := tx.Where("product_size_id = ?", productSizeID).Delete(&data.ProductSizeIngredient{}).Error; err != nil {
		return fmt.Errorf("failed to delete ingredients: %w", err)
	}

	newIngredients := make([]data.ProductSizeIngredient, len(ingredients))
	for i, ingredient := range ingredients {
		newIngredients[i] = data.ProductSizeIngredient{
			ProductSizeID: productSizeID,
			IngredientID:  ingredient.IngredientID,
			Quantity:      ingredient.Quantity,
		}
	}

	if err := tx.Create(newIngredients).Error; err != nil {
		return fmt.Errorf("failed to create ingredients: %w", err)
	}

	return nil
}

func (r *productRepository) UpdateProduct(productID uint, product *data.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&data.Product{}).Where("id = ?", productID).Updates(product).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *productRepository) DeleteProduct(productID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", productID).
			Delete(&data.Product{}).Error; err != nil {
			return err
		}

		if err := tx.Where("product_id = ?", productID).
			Delete(&data.ProductSize{}).Error; err != nil {
			return err
		}

		var storeProduct data.StoreProduct
		if err := tx.Clauses(clause.Returning{}).
			Where("product_id = ?", productID).
			Delete(&storeProduct).Error; err != nil {
			return err
		}

		if err := tx.Where("store_product_id = ?", storeProduct.ID).
			Delete(&data.StoreProductSize{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetProductSizesByProductID(productID uint) ([]data.ProductSize, error) {
	var productSizes []data.ProductSize

	// Optimize this
	err := r.db.
		Preload("Unit").
		Preload("Additives").
		Preload("Additives.Additive.Category").
		Preload("ProductSizeIngredients.Ingredient").
		Preload("ProductSizeIngredients.Ingredient.IngredientCategory").
		Preload("ProductSizeIngredients.Ingredient.Unit").
		Where("product_id = ?", productID).
		Find(&productSizes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product sizes: %w", err)
	}

	return productSizes, nil
}

func (r *productRepository) DeleteProductSize(productSizeID uint) error {
	return r.db.Where("id = ?", productSizeID).Delete(&data.ProductSize{}).Error
}
