package product

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
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
			// If is_default is being set to true, first set all other sizes to false
			if updateModels.ProductSize.IsDefault {
				var productSize data.ProductSize
				if err := tx.Select("product_id").First(&productSize, id).Error; err != nil {
					return fmt.Errorf("failed to get product size: %w", err)
				}

				// Set all other sizes for this product to non-default
				if err := tx.Model(&data.ProductSize{}).
					Where("product_id = ? AND id != ?", productSize.ProductID, id).
					Update("is_default", false).Error; err != nil {
					return fmt.Errorf("failed to update other product sizes: %w", err)
				}
			}

			if err := tx.Model(&data.ProductSize{}).
				Where("id = ?", id).
				Updates(updateModels.ProductSize).Error; err != nil {
				return fmt.Errorf("failed to update product size: %w", err)
			}
		}

		if updateModels.Additives != nil {
			// Remove existing ingredients
			if err := tx.Where("product_size_id = ?", id).Delete(&data.ProductSizeAdditive{}).Error; err != nil {
				return fmt.Errorf("failed to delete additives: %w", err)
			}
			// Add new additives
			additives := make([]data.ProductSizeAdditive, len(updateModels.Additives))
			for i, additive := range updateModels.Additives {
				additives[i] = data.ProductSizeAdditive{
					ProductSizeID: id,
					AdditiveID:    additive.AdditiveID,
					IsDefault:     additive.IsDefault,
				}
			}

			if err := tx.Create(additives).Error; err != nil {
				return fmt.Errorf("failed to create additive: %w", err)
			}
		}

		if updateModels.Ingredients != nil {
			// Remove existing ingredients
			if err := tx.Where("product_size_id = ?", id).Delete(&data.ProductSizeIngredient{}).Error; err != nil {
				return fmt.Errorf("failed to delete ingredients: %w", err)
			}

			// Add new ingredients
			ingredients := make([]data.ProductSizeIngredient, len(updateModels.Ingredients))
			for i, ingredient := range updateModels.Ingredients {
				ingredients[i] = data.ProductSizeIngredient{
					ProductSizeID: id,
					IngredientID:  ingredient.IngredientID,
					Quantity:      ingredient.Quantity,
				}
			}
			if err := tx.Create(ingredients).Error; err != nil {
				return fmt.Errorf("failed to create ingredients: %w", err)
			}
		}

		return nil
	})
}

func (r *productRepository) setNewDefaultProductSize(productSizeID, productID uint) error {
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

	err := r.db.
		Preload("Unit").
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
