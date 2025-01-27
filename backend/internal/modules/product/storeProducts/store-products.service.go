package storeProducts

import (
	"fmt"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	storeWarehousesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

const DEFAULT_LOW_STOCK_THRESHOLD = 50

type StoreProductService interface {
	GetStoreProductById(storeID, storeProductID uint) (*types.StoreProductDetailsDTO, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]types.StoreProductDetailsDTO, error)
	GetProductsListToAdd(storeID uint, filter *productTypes.ProductsFilterDto) ([]productTypes.ProductDetailsDTO, error)
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
	transactionManager TransactionManager
	logger             *zap.SugaredLogger
}

func NewStoreProductService(
	repo StoreProductRepository,
	productRepo product.ProductRepository,
	ingredientRepo ingredients.IngredientRepository,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) StoreProductService {
	return &storeProductService{
		repo:               repo,
		productRepo:        productRepo,
		ingredientsRepo:    ingredientRepo,
		transactionManager: transactionManager,
		logger:             logger,
	}
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

func (s *storeProductService) GetProductsListToAdd(storeID uint, filter *productTypes.ProductsFilterDto) ([]productTypes.ProductDetailsDTO, error) {
	products, err := s.repo.GetProductsListToAdd(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to get products list to add", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	productDTOs := make([]productTypes.ProductDetailsDTO, len(products))
	for i, product := range products {
		productDTOs[i] = *productTypes.MapToProductDetailsDTO(&product)
	}

	return productDTOs, nil
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

	ingredientsList, err := s.ingredientsRepo.GetIngredientsForProductSizes(productSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store product: could not get ingredients", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	addStockDTOs := make([]storeWarehousesTypes.AddStoreStockDTO, len(ingredientsList))
	for i, ingredient := range ingredientsList {
		addStockDTOs[i] = storeWarehousesTypes.AddStoreStockDTO{
			IngredientID:      ingredient.ID,
			Quantity:          0,
			LowStockThreshold: DEFAULT_LOW_STOCK_THRESHOLD,
		}
	}

	storeProduct := types.CreateToStoreProduct(dto)
	id, err := s.transactionManager.CreateStoreProductWithStocks(storeID, storeProduct, addStockDTOs)
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

	addStockDTOs, err := s.formAddStockDTOsFromIngredients(inputSizeIDs)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store products", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	ids, err := s.transactionManager.CreateMultipleStoreProductsWithStocks(storeID, storeProducts, addStockDTOs)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create %d store product: %w", len(dtos), err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	return ids, nil
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

		if err := s.validateProductSizesByProductID(inputSizeIDs, storeProductDetails.ID); err != nil {
			wrappedErr := utils.WrapError("failed to update store product: ", err)
			return wrappedErr
		}
	}

	addStockDTOs, err := s.formAddStockDTOsFromIngredients(inputSizeIDs)
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

func (s *storeProductService) formAddStockDTOsFromIngredients(productSizeIDs []uint) ([]storeWarehousesTypes.AddStoreStockDTO, error) {
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
		wrappedErr := fmt.Errorf("%w: %w", types.ErrValidationFailed, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if err := types.ValidateStoreProductSizes(productSizeIDs, productSizes); err != nil {
		wrappedErr := fmt.Errorf("%w: %w", types.ErrValidationFailed, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}
