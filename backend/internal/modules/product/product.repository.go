package product

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	CheckProductExists(productName string) (bool, error)
	CreateProduct(product *data.Product) (uint, error)
	GetProducts(locale data.LanguageCode, filter *types.ProductsFilterDto) ([]data.Product, error)
	GetRawProductByID(productID uint) (*data.Product, error)
	GetProductByID(productID uint) (*data.Product, error)
	GetTranslatedProductByID(locale data.LanguageCode, productID uint) (*data.Product, error)

	SaveProduct(product *data.Product) error
	DeleteProduct(productID uint) (*data.Product, error)

	CreateProductSize(createModels *data.ProductSize) (uint, error)
	GetProductSizesByProductID(productID uint) ([]data.ProductSize, error)
	GetProductSizeById(productSizeID uint) (*data.ProductSize, error)
	GetProductSizeDetailsByID(productSizeID uint) (*data.ProductSize, error)
	UpdateProductSizeWithAssociations(id uint, updateModels *types.ProductSizeModels) error
	DeleteProductSize(productID uint) error
	FindRawProductByID(productID uint, product *data.Product) error

	CloneWithTransaction(tx *gorm.DB) ProductRepository
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CloneWithTransaction(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}

func (r *productRepository) CheckProductExists(productName string) (bool, error) {
	var product data.Product
	err := r.db.Where(&data.Product{Name: productName}).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
		Preload("Additives.Additive.Ingredients.Ingredient").
		Preload("ProductSizeIngredients.Ingredient.IngredientCategory").
		Preload("ProductSizeIngredients.Ingredient.Unit").
		Preload("ProductSizeProvisions.Provision.Unit").
		First(&productSize, productSizeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrProductSizeNotFound
		}
		return nil, fmt.Errorf("failed to fetch ProductSize ID %d: %w", productSizeID, err)
	}

	return &productSize, nil
}

func (r *productRepository) GetProducts(locale data.LanguageCode, filter *types.ProductsFilterDto) ([]data.Product, error) {
	var products []data.Product

	base := r.db.Model(&data.Product{})

	if filter.CategoryID != nil {
		base = base.Where("products.category_id = ?", *filter.CategoryID)
	}

	if filter.Search != nil && *filter.Search != "" {
		term := "%" + *filter.Search + "%"
		base = base.Where(
			"products.name ILIKE ? OR products.description ILIKE ?",
			term, term,
		)
	}

	pg, err := utils.ApplySortedPaginationForModel(
		base, filter.Pagination, filter.Sort, &data.Product{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	pg = utils.ApplyLocalizedPreloads(
		pg,
		locale,
		types.ProductPreloadMap,
	)

	if err := pg.Find(&products).Error; err != nil {
		return nil, err
	}
	if products == nil {
		products = []data.Product{}
	}
	return products, nil
}

func (r *productRepository) GetRawProductByID(productID uint) (*data.Product, error) {
	var product data.Product

	err := r.db.
		Where("id = ?", productID).
		First(&product).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
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
			return nil, types.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetTranslatedProductByID(locale data.LanguageCode, productID uint) (*data.Product, error) {
	var prod data.Product

	q := r.db.Model(&data.Product{}).Where("id = ?", productID)

	q = utils.ApplyLocalizedPreloads(q, locale, types.ProductPreloadMap)

	if err := q.First(&prod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrProductNotFound
		}
		return nil, err
	}
	return &prod, nil
}

func (r *productRepository) CreateProduct(product *data.Product) (uint, error) {
	if err := r.db.Create(product).Error; err != nil {
		return 0, err
	}

	return product.ID, nil
}

func (r *productRepository) CreateProductSize(productSize *data.ProductSize) (uint, error) {
	if err := r.db.Create(productSize).Error; err != nil {
		if strings.Contains(err.Error(), "23505") {
			return 0, types.ErrProductSizeUniqueName
		}
		return 0, err
	}

	return productSize.ID, nil
}

func (r *productRepository) UpdateProductSizeWithAssociations(id uint, updateModels *types.ProductSizeModels) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels == nil {
			return fmt.Errorf("nothing to update")
		}

		additivesUpdated := false
		provisionsUpdated := false

		if updateModels.Additives != nil {
			if err := r.updateAdditives(tx, id, updateModels.Additives); err != nil {
				return err
			}
			additivesUpdated = true
		}

		if updateModels.Provisions != nil {
			if err := r.updateProvisions(tx, id, updateModels.Provisions); err != nil {
				return err
			}
			provisionsUpdated = true
		}

		if updateModels.ProductSize != nil {
			if err := r.updateProductSize(tx, id, updateModels.ProductSize, additivesUpdated, provisionsUpdated); err != nil {
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

func (r *productRepository) updateProductSize(tx *gorm.DB, id uint, productSize *data.ProductSize, additivesUpdated, provisionsUpdated bool) error {
	if additivesUpdated {
		productSize.AdditivesUpdatedAt = time.Now().UTC()
	}

	if provisionsUpdated {
		productSize.ProvisionsUpdatedAt = time.Now().UTC()
	}

	if err := tx.Model(&data.ProductSize{}).
		Where("id = ?", id).
		Updates(productSize).Error; err != nil {
		return fmt.Errorf("failed to update product size: %w", err)
	}
	return nil
}

func (r *productRepository) updateAdditives(tx *gorm.DB, productSizeID uint, additives []data.ProductSizeAdditive) error {
	var existingAdditives []data.ProductSizeAdditive
	if err := tx.Where("product_size_id = ?", productSizeID).Find(&existingAdditives).Error; err != nil {
		return fmt.Errorf("failed to fetch existing additives: %w", err)
	}

	existingMap := make(map[uint]data.ProductSizeAdditive)
	for _, additive := range existingAdditives {
		existingMap[additive.AdditiveID] = additive
	}

	var toInsert []data.ProductSizeAdditive
	var toUpdate []data.ProductSizeAdditive
	existingIDs := make(map[uint]struct{})

	for _, additive := range additives {
		existing, exists := existingMap[additive.AdditiveID]

		if exists {
			if existing.IsDefault != additive.IsDefault || existing.IsHidden != additive.IsHidden {
				existing.IsDefault = additive.IsDefault
				existing.IsHidden = additive.IsHidden
				toUpdate = append(toUpdate, existing)
			}

			existingIDs[additive.AdditiveID] = struct{}{}
		} else {
			toInsert = append(toInsert, data.ProductSizeAdditive{
				ProductSizeID: productSizeID,
				AdditiveID:    additive.AdditiveID,
				IsDefault:     additive.IsDefault,
				IsHidden:      additive.IsHidden,
			})
		}
	}

	var toDeleteIDs []uint
	for id := range existingMap {
		if _, found := existingIDs[id]; !found {
			toDeleteIDs = append(toDeleteIDs, id)
		}
	}

	if len(toUpdate) > 0 {
		if err := tx.Save(&toUpdate).Error; err != nil {
			return fmt.Errorf("failed to update additives: %w", err)
		}
	}

	if len(toDeleteIDs) > 0 {
		if err := tx.Unscoped().Where("product_size_id = ? AND additive_id IN (?)", productSizeID, toDeleteIDs).Delete(&data.ProductSizeAdditive{}).Error; err != nil {
			return fmt.Errorf("failed to delete old additives: %w", err)
		}
	}

	if len(toInsert) > 0 {
		if err := tx.Create(&toInsert).Error; err != nil {
			return fmt.Errorf("failed to insert new additives: %w", err)
		}
	}

	return nil
}

func (r *productRepository) updateIngredients(tx *gorm.DB, productSizeID uint, ingredients []data.ProductSizeIngredient) error {
	var existingIngredients []data.ProductSizeIngredient
	if err := tx.Where("product_size_id = ?", productSizeID).Find(&existingIngredients).Error; err != nil {
		return fmt.Errorf("failed to fetch existing ingredients: %w", err)
	}

	existingMap := make(map[uint]data.ProductSizeIngredient)
	for _, ing := range existingIngredients {
		existingMap[ing.IngredientID] = ing
	}

	var toInsert []data.ProductSizeIngredient
	var toUpdate []data.ProductSizeIngredient
	existingIDs := make(map[uint]struct{})

	for _, ingredient := range ingredients {
		existing, exists := existingMap[ingredient.IngredientID]

		if exists {
			if existing.Quantity != ingredient.Quantity {
				existing.Quantity = ingredient.Quantity
				toUpdate = append(toUpdate, existing)
			}
			existingIDs[ingredient.IngredientID] = struct{}{}
		} else {
			toInsert = append(toInsert, data.ProductSizeIngredient{
				ProductSizeID: productSizeID,
				IngredientID:  ingredient.IngredientID,
				Quantity:      ingredient.Quantity,
			})
		}
	}

	var toDeleteIDs []uint
	for id := range existingMap {
		if _, found := existingIDs[id]; !found {
			toDeleteIDs = append(toDeleteIDs, id)
		}
	}

	if len(toUpdate) > 0 {
		if err := tx.Save(&toUpdate).Error; err != nil {
			return fmt.Errorf("failed to update ingredients: %w", err)
		}
	}

	if len(toDeleteIDs) > 0 {
		if err := tx.Unscoped().Where("product_size_id = ? AND ingredient_id IN (?)", productSizeID, toDeleteIDs).Delete(&data.ProductSizeIngredient{}).Error; err != nil {
			return fmt.Errorf("failed to delete old ingredients: %w", err)
		}
	}

	if len(toInsert) > 0 {
		if err := tx.Create(&toInsert).Error; err != nil {
			return fmt.Errorf("failed to insert new ingredients: %w", err)
		}
	}

	return nil
}

func (r *productRepository) updateProvisions(tx *gorm.DB, productSizeID uint, provisions []data.ProductSizeProvision) error {
	var existingProvisions []data.ProductSizeProvision
	if err := tx.Where("product_size_id = ?", productSizeID).Find(&existingProvisions).Error; err != nil {
		return fmt.Errorf("failed to fetch existing provisions: %w", err)
	}

	existingMap := make(map[uint]data.ProductSizeProvision)
	for _, prov := range existingProvisions {
		existingMap[prov.ProvisionID] = prov
	}

	var toInsert []data.ProductSizeProvision
	var toUpdate []data.ProductSizeProvision
	existingIDs := make(map[uint]struct{})

	for _, provision := range provisions {
		existing, exists := existingMap[provision.ProvisionID]

		if exists {
			if existing.Volume != provision.Volume {
				existing.Volume = provision.Volume
				toUpdate = append(toUpdate, existing)
			}
			existingIDs[provision.ProvisionID] = struct{}{}
		} else {
			toInsert = append(toInsert, data.ProductSizeProvision{
				ProductSizeID: productSizeID,
				ProvisionID:   provision.ProvisionID,
				Volume:        provision.Volume,
			})
		}
	}

	var toDeleteIDs []uint
	for id := range existingMap {
		if _, found := existingIDs[id]; !found {
			toDeleteIDs = append(toDeleteIDs, id)
		}
	}

	if len(toUpdate) > 0 {
		if err := tx.Save(&toUpdate).Error; err != nil {
			return fmt.Errorf("failed to update provisions: %w", err)
		}
	}

	if len(toDeleteIDs) > 0 {
		if err := tx.Unscoped().Where("product_size_id = ? AND provision_id IN (?)", productSizeID, toDeleteIDs).Delete(&data.ProductSizeProvision{}).Error; err != nil {
			return fmt.Errorf("failed to delete old provisions: %w", err)
		}
	}

	if len(toInsert) > 0 {
		if err := tx.Create(&toInsert).Error; err != nil {
			return fmt.Errorf("failed to insert new provisions: %w", err)
		}
	}

	return nil
}

func (r *productRepository) SaveProduct(product *data.Product) error {
	err := r.db.Save(product).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) DeleteProduct(productID uint) (*data.Product, error) {
	var product data.Product

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkProductInActiveOrders(tx, productID); err != nil {
			return err
		}

		if err := tx.Clauses(clause.Returning{}).
			Where("id = ?", productID).
			Delete(&product).Error; err != nil {
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
		return nil, err
	}
	return &product, nil
}

func checkProductInActiveOrders(db *gorm.DB, productID uint) error {
	var exists bool

	query := db.Table("suborders").
		Joins("JOIN store_product_sizes ON suborders.store_product_size_id = store_product_sizes.id").
		Joins("JOIN store_products ON store_product_sizes.store_product_id = store_products.id").
		Where("store_products.product_id = ?", productID).
		Where("suborders.status IN ?", []string{
			string(data.SubOrderStatusPending),
			string(data.SubOrderStatusPreparing),
		}).
		Limit(1).
		Select("1").Scan(&exists)

	if query.Error != nil {
		return query.Error
	}

	if exists {
		return types.ErrProductIsInUse // Custom error indicating the product is in use
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := checkProductSizeInActiveOrders(tx, productSizeID); err != nil {
			return err
		}

		result := tx.Where("id = ?", productSizeID).Delete(&data.ProductSize{})
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func checkProductSizeInActiveOrders(db *gorm.DB, productSizeID uint) error {
	var exists bool

	query := db.Table("suborders").
		Joins("JOIN store_product_sizes ON suborders.store_product_size_id = store_product_sizes.id").
		Where("store_product_sizes.product_size_id = ?", productSizeID).
		Where("suborders.status IN ?", []string{
			string(data.SubOrderStatusPending),
			string(data.SubOrderStatusPreparing),
		}).
		Limit(1).
		Select("1").Scan(&exists)

	if query.Error != nil {
		return query.Error
	}

	if exists {
		return types.ErrProductSizeIsInUse // Custom error indicating the ProductSize is in use
	}

	return nil
}

func (r *productRepository) FindRawProductByID(productID uint, product *data.Product) error {
	err := r.db.
		Where("id = ?", productID).
		First(product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ErrProductNotFound
		}
		return err
	}

	return nil
}
