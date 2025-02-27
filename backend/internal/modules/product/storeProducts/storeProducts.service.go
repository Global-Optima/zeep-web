package storeProducts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	storeAdditivesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	categoriesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	storeWarehousesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

const DEFAULT_LOW_STOCK_THRESHOLD = 50

type StoreProductService interface {
	GetStoreProductCategories(storeID uint) ([]categoriesTypes.ProductCategoryDTO, error)
	GetStoreProductById(storeID, storeProductID uint) (*types.StoreProductDetailsDTO, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]types.StoreProductDetailsDTO, error)
	GetStoreProductsByStoreProductIDs(storeID uint, storeProductIDs []uint) ([]types.StoreProductDetailsDTO, error)
	GetAvailableProductsToAdd(storeID uint, filter *productTypes.ProductsFilterDto) ([]productTypes.ProductDetailsDTO, error)
	GetRecommendedStoreProducts(storeID uint, excludedStoreProductIDs []uint) ([]types.StoreProductDetailsDTO, error)
	GetStoreProductSizeByID(storeID, storeProductSizeID uint) (*types.StoreProductSizeDetailsDTO, error)
	CreateStoreProduct(storeID uint, dto *types.CreateStoreProductDTO) (uint, error)
	CreateMultipleStoreProducts(storeID uint, dtos []types.CreateStoreProductDTO) ([]uint, error)
	UpdateStoreProduct(storeID, storeProductID uint, dto *types.UpdateStoreProductDTO) error
	DeleteStoreProduct(storeID, storeProductID uint) error
}

type storeProductService struct {
	repo               StoreProductRepository
	productRepo        product.ProductRepository
	ingredientsRepo    ingredients.IngredientRepository
	storeAdditiveRepo  storeAdditives.StoreAdditiveRepository
	storageRepo        storage.StorageRepository
	transactionManager TransactionManager
	logger             *zap.SugaredLogger
}

func NewStoreProductService(
	repo StoreProductRepository,
	productRepo product.ProductRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storageRepo storage.StorageRepository,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) StoreProductService {
	return &storeProductService{
		repo:               repo,
		productRepo:        productRepo,
		ingredientsRepo:    ingredientRepo,
		storeAdditiveRepo:  storeAdditiveRepo,
		storageRepo:        storageRepo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (s *storeProductService) GetStoreProductCategories(storeID uint) ([]categoriesTypes.ProductCategoryDTO, error) {
	categories, err := s.repo.GetStoreProductCategories(storeID)
	if err != nil {
		wrappedErr := fmt.Errorf("error getting store product categories: %w", err)
		s.logger.Errorw(wrappedErr.Error())
		return nil, err
	}

	dtos := make([]categoriesTypes.ProductCategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = *categoriesTypes.MapCategoryToDTO(category)
	}
	return dtos, nil
}

func (s *storeProductService) GetStoreProductById(storeID uint, storeProductID uint) (*types.StoreProductDetailsDTO, error) {
	storeProduct, err := s.repo.GetStoreProductById(storeID, storeProductID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store product by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dto := types.MapToStoreProductDetailsDTO(storeProduct)

	return &dto, nil
}

func (s *storeProductService) GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]types.StoreProductDetailsDTO, error) {
	storeProducts, err := s.repo.GetStoreProducts(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreProductDetailsDTO, len(storeProducts))
	for i, storeProduct := range storeProducts {
		dtos[i] = types.MapToStoreProductDetailsDTO(&storeProduct)
	}

	return dtos, nil
}

func (s *storeProductService) GetStoreProductsByStoreProductIDs(storeID uint, storeProductIDs []uint) ([]types.StoreProductDetailsDTO, error) {
	storeProducts, err := s.repo.GetStoreProductsByStoreProductIDs(storeID, storeProductIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreProductDetailsDTO, len(storeProducts))
	for i, storeProduct := range storeProducts {
		dtos[i] = types.MapToStoreProductDetailsDTO(&storeProduct)
	}

	return dtos, nil
}

func (s *storeProductService) GetAvailableProductsToAdd(storeID uint, filter *productTypes.ProductsFilterDto) ([]productTypes.ProductDetailsDTO, error) {
	productsList, err := s.repo.GetAvailableProductsToAdd(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get available products to add", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	productDTOs := make([]productTypes.ProductDetailsDTO, len(productsList))
	for i, productItem := range productsList {
		productDTOs[i] = *productTypes.MapToProductDetailsDTO(&productItem)
	}

	return productDTOs, nil
}

func (s *storeProductService) GetRecommendedStoreProducts(storeID uint, excludedStoreProductIDs []uint) ([]types.StoreProductDetailsDTO, error) {
	storeProducts, err := s.repo.GetRecommendedStoreProducts(storeID, excludedStoreProductIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get recommended products to add", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreProductDetailsDTO, len(storeProducts))
	for i, storeProduct := range storeProducts {
		dtos[i] = types.MapToStoreProductDetailsDTO(&storeProduct)
	}

	return dtos, nil
}

func (s *storeProductService) GetStoreProductSizeByID(storeID, storeProductSizeID uint) (*types.StoreProductSizeDetailsDTO, error) {
	storeProduct, err := s.repo.GetStoreProductSizeById(storeID, storeProductSizeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get store product size by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dto := types.MapToStoreProductSizeDetailsDTO(*storeProduct)
	return &dto, nil
}

func (s *storeProductService) CreateStoreProduct(storeID uint, dto *types.CreateStoreProductDTO) (uint, error) {

	if len(dto.ProductSizes) > 0 {
		inputSizeIDs := make([]uint, len(dto.ProductSizes))
		for i, productSize := range dto.ProductSizes {
			inputSizeIDs[i] = productSize.ProductSizeID
		}

		if err := s.validateProductSizesByProductID(inputSizeIDs, dto.ProductID); err != nil {
			return 0, utils.WrapError("failed to create store product", err)
		}
	}

	productSizeIDs := make([]uint, len(dto.ProductSizes))
	for i, size := range dto.ProductSizes {
		productSizeIDs[i] = size.ProductSizeID
	}

	addStockDTOs, err := s.formAddStockDTOsFromProductSizes(productSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update store product: ", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	additiveIDs, err := s.storeAdditiveRepo.GetMissingStoreAdditiveIDsForProductSizes(storeID, productSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get additives of store product", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	storeAdditiveList := make([]data.StoreAdditive, len(additiveIDs))
	for _, additiveID := range additiveIDs {
		storeAdditivesTypes.CreateToStoreAdditive(&storeAdditivesTypes.CreateStoreAdditiveDTO{
			AdditiveID: additiveID,
			StorePrice: nil,
		}, storeID)
	}

	storeProduct := types.CreateToStoreProduct(dto)
	id, _, err := s.transactionManager.CreateStoreProductWithStocks(storeID, storeProduct, storeAdditiveList, addStockDTOs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store product", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
}

func (s *storeProductService) CreateMultipleStoreProducts(storeID uint, dtos []types.CreateStoreProductDTO) ([]uint, error) {
	var inputSizeIDs []uint
	storeProducts := make([]data.StoreProduct, len(dtos))

	for i, dto := range dtos {
		storeProducts[i] = *types.CreateToStoreProduct(&dto)
		storeProducts[i].StoreID = storeID

		if len(dto.ProductSizes) > 0 {
			inputSizeIDs = make([]uint, len(dto.ProductSizes))
			for j, productSize := range dto.ProductSizes {
				inputSizeIDs[j] = productSize.ProductSizeID
			}

			if err := s.validateProductSizesByProductID(inputSizeIDs, dto.ProductID); err != nil {
				return nil, utils.WrapError("failed to create store products", err)
			}
		}
	}

	addStockDTOs, err := s.formAddStockDTOsFromProductSizes(inputSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	var productSizeIDs []uint
	for _, dto := range dtos {
		for _, size := range dto.ProductSizes {
			productSizeIDs = append(productSizeIDs, size.ProductSizeID)
		}
	}

	additiveIDs, err := s.storeAdditiveRepo.GetMissingStoreAdditiveIDsForProductSizes(storeID, productSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get additives of store product", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	storeAdditiveList := make([]data.StoreAdditive, len(additiveIDs))
	for i, additiveID := range additiveIDs {
		storeAdditiveList[i] = *storeAdditivesTypes.CreateToStoreAdditive(&storeAdditivesTypes.CreateStoreAdditiveDTO{
			AdditiveID: additiveID,
			StorePrice: nil,
		}, storeID)
	}

	storeProductIDs, _, err := s.transactionManager.CreateMultipleStoreProductsWithStocks(storeID, storeProducts, storeAdditiveList, addStockDTOs)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create %d store product: %w", len(dtos), err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	return storeProductIDs, nil
}

func (s *storeProductService) UpdateStoreProduct(storeID, storeProductID uint, dto *types.UpdateStoreProductDTO) error {
	var inputSizeIDs []uint
	if len(dto.ProductSizes) > 0 {
		inputSizeIDs = make([]uint, len(dto.ProductSizes))
		for i, productSize := range dto.ProductSizes {
			inputSizeIDs[i] = productSize.ProductSizeID
		}

		storeProductDetails, err := s.repo.GetStoreProductById(storeID, storeProductID)
		if err != nil {
			wrappedErr := utils.WrapError("failed to update store product", err)
			s.logger.Error(wrappedErr)
			return wrappedErr
		}

		if err := s.validateProductSizesByProductID(inputSizeIDs, storeProductDetails.ProductID); err != nil {
			wrappedErr := utils.WrapError("failed to update store product: ", err)
			return wrappedErr
		}
	}

	addStockDTOs, err := s.formAddStockDTOsFromProductSizes(inputSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update store product: ", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	updateModels := types.UpdateToStoreProductModels(dto)
	err = s.transactionManager.UpdateStoreProductWithStocks(storeID, storeProductID, updateModels, addStockDTOs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update store product", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *storeProductService) DeleteStoreProduct(storeID, storeProductID uint) error {
	err := s.repo.DeleteStoreProductWithSizes(storeID, storeProductID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete store product", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *storeProductService) formAddStockDTOsFromProductSizes(productSizeIDs []uint) ([]storeWarehousesTypes.AddStoreStockDTO, error) {
	ingredientsList, err := s.ingredientsRepo.GetIngredientsForProductSizes(productSizeIDs)
	if err != nil {
		return nil, utils.WrapError("could not get ingredients", err)
	}

	addStockDTOs := make([]storeWarehousesTypes.AddStoreStockDTO, len(ingredientsList))
	for i, ingredient := range ingredientsList {
		addStockDTOs[i] = storeWarehousesTypes.AddStoreStockDTO{
			IngredientID:      ingredient.ID,
			Quantity:          0,
			LowStockThreshold: DEFAULT_LOW_STOCK_THRESHOLD,
		}
	}
	return addStockDTOs, nil
}

func (s *storeProductService) validateProductSizesByProductID(productSizeIDs []uint, productID uint) error {
	productSizes, err := s.productRepo.GetProductSizesByProductID(productID)
	if err != nil {
		wrappedErr := fmt.Errorf("%w: %w", moduleErrors.ErrValidation, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if err := types.ValidateStoreProductSizes(productSizeIDs, productSizes); err != nil {
		wrappedErr := fmt.Errorf("%w: %w", moduleErrors.ErrValidation, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}
