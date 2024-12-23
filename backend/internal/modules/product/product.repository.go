package product

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *data.Product) (uint, error)
	GetProducts(filter *types.ProductsFilterDto) ([]data.Product, error)
	//GetProductDetails(productID uint) (*data.Product, error)
	GetProductDetails(productID uint) (*types.ProductDetailsDTO, error)
	UpdateProduct(id uint, product *data.Product, defaultAdditiveIDs []uint) error
	DeleteProduct(productID uint) error

	//GetProductSizes(productID uint) ([]data.ProductSize, error)
	CreateProductSize(productSize *data.ProductSize) (uint, error)
	GetProductSizeById(productSizeID uint) (*data.ProductSize, error)
	UpdateProductSizeWithAssociations(id uint, productSize *data.ProductSize, additiveIDs, ingredientIDs []uint) error
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

func (r *productRepository) GetProducts(filter *types.ProductsFilterDto) ([]data.Product, error) {
	var products []data.Product

	query := r.db.
		Model(&data.Product{}).
		Joins("JOIN store_products ON store_products.product_id = products.id").
		Preload("ProductSizes", "is_default = TRUE").
		Preload("ProductSizes.ProductIngredients").
		Preload("ProductSizes.ProductIngredients.Ingredient")

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

func (r *productRepository) GetProductDetails(productID uint) (*types.ProductDetailsDTO, error) {
	var product data.Product

	err := r.db.
		Preload("ProductSizes").
		Preload("DefaultAdditives.Additive").
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

	productDTO := types.MapToProductDetailsDTO(&product, defaultAdditives)

	return productDTO, nil
}

func (r *productRepository) CreateProduct(product *data.Product) (uint, error) {
	if err := r.db.Create(product).Error; err != nil {
		return 0, err
	}

	return product.ID, nil
}

/*func (r *productRepository) CreateProduct(product *data.Product, defaultAdditiveIDs []uint) (uint, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		defaultAdditives := make([]data.DefaultProductAdditive, 0)
		for _, id := range defaultAdditiveIDs {
			defaultAdditives = append(defaultAdditives, data.DefaultProductAdditive{
				ProductID:  product.ID,
				AdditiveID: id,
			})
		}

		if len(defaultAdditives) > 0 {
			if err := tx.Create(&defaultAdditives).Error; err != nil {
				return err
			}
		} else {
			return errors.New("no default additives attached")
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return product.ID, nil
}*/

func (r *productRepository) CreateProductSize(productSize *data.ProductSize) (uint, error) {
	if err := r.db.Create(productSize).Error; err != nil {
		return 0, err
	}
	return productSize.ID, nil
}

func (r *productRepository) UpdateProductSizeWithAssociations(id uint,
	productSize *data.ProductSize,
	additiveIDs, ingredientIDs []uint,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update the ProductSize model
		if productSize != nil {
			if err := tx.Model(&data.ProductSize{}).Where("id = ?", id).Updates(productSize).Error; err != nil {
				return fmt.Errorf("failed to update product size: %w", err)
			}
		}

		// Update Additives
		if additiveIDs != nil {
			// Remove existing additives
			if err := tx.Where("product_size_id = ?", id).Delete(&data.ProductAdditive{}).Error; err != nil {
				return fmt.Errorf("failed to delete additives: %w", err)
			}
			// Add new additives
			for _, additiveID := range additiveIDs {
				newAdditive := data.ProductAdditive{
					ProductSizeID: id,
					AdditiveID:    additiveID,
				}
				if err := tx.Create(&newAdditive).Error; err != nil {
					return fmt.Errorf("failed to create additive: %w", err)
				}
			}
		}

		// Update Ingredients
		if ingredientIDs != nil {
			// Remove existing ingredients
			if err := tx.Where("product_size_id = ?", id).Delete(&data.ProductIngredient{}).Error; err != nil {
				return fmt.Errorf("failed to delete ingredients: %w", err)
			}
			// Add new ingredients
			for _, ingredientID := range ingredientIDs {
				newIngredient := data.ProductIngredient{
					ProductSizeID:    id,
					ItemIngredientID: ingredientID,
				}
				if err := tx.Create(&newIngredient).Error; err != nil {
					return fmt.Errorf("failed to create ingredient: %w", err)
				}
			}
		}

		return nil
	})
}

/*func (r *productRepository) CreateProductWithSizes(dto *types.CreateProductWithAttachesDTO) (uint, error) {
	product := types.CreateToProductModel(dto)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		productSizes := types.ToProductSizesModels(dto.ProductSizes, product.ID)
		if err := r.createProductSizes(tx, productSizes, product.ID); err != nil {
			return err
		}

		if err := r.handleAdditives(tx, dto.Additives, product.ID, productSizes); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to finish transaction for product creation: %w", err)
	}

	return product.ID, nil
}*/

func (r *productRepository) createProduct(tx *gorm.DB, product *data.Product) (uint, error) {
	if err := tx.Create(product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (r *productRepository) createProductSizes(tx *gorm.DB, productSizes []data.ProductSize, productID uint) error {
	for i := range productSizes {
		productSizes[i].ProductID = productID
	}
	return tx.Create(&productSizes).Error
}

func (r *productRepository) handleAdditives(tx *gorm.DB, additives []types.SelectedAdditiveTypesDTO, productID uint, productSizes []data.ProductSize) error {
	var defaultAdditives []data.DefaultProductAdditive
	productSizeAdditives := make(map[uint][]data.ProductAdditive)

	for _, additive := range additives {
		if additive.IsDefault {
			defaultAdditives = append(defaultAdditives, data.DefaultProductAdditive{
				ProductID:  productID,
				AdditiveID: additive.AdditiveID,
			})
		} else {
			for _, size := range productSizes {
				productSizeAdditives[size.ID] = append(productSizeAdditives[size.ID], data.ProductAdditive{
					ProductSizeID: size.ID,
					AdditiveID:    additive.AdditiveID,
				})
			}
		}
	}

	if len(defaultAdditives) > 0 {
		if err := tx.Create(&defaultAdditives).Error; err != nil {
			return err
		}
	}

	for productSizeID, productAdditives := range productSizeAdditives {
		if len(additives) > 0 {
			for i := range additives {
				productAdditives[i].ProductSizeID = productSizeID
			}
			if err := tx.Create(&productAdditives).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *productRepository) UpdateProduct(productID uint, product *data.Product, defaultAdditiveIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&data.Product{}).Where("id = ?", productID).Updates(product).Error
		if err != nil {
			return err
		}

		if defaultAdditiveIDs != nil {
			if err := tx.Where("product_id = ?", productID).Delete(&data.DefaultProductAdditive{}).Error; err != nil {
				return err
			}

			for _, additiveID := range defaultAdditiveIDs {
				newAdditive := data.DefaultProductAdditive{
					ProductID:  productID,
					AdditiveID: additiveID,
				}
				if err := tx.Create(&newAdditive).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *productRepository) UpdateProductSize(productSize *data.ProductSize) error {
	return r.db.Updates(productSize).Error
}

func (r *productRepository) DeleteProduct(productID uint) error {
	return r.db.Where("id = ?", productID).Delete(&data.Product{}).Error
}

func (r *productRepository) DeleteProductSize(productSizeID uint) error {
	return r.db.Where("id = ?", productSizeID).Delete(&data.ProductSize{}).Error
}
