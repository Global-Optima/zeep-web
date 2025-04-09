package storeAdditives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type StoreAdditiveService interface {
	CreateStoreAdditives(storeID uint, dtos []types.CreateStoreAdditiveDTO) ([]uint, error)
	UpdateStoreAdditive(storeID, storeAdditiveID uint, dto *types.UpdateStoreAdditiveDTO) error
	GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]types.StoreAdditiveDTO, error)
	GetAdditivesListToAdd(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]additiveTypes.AdditiveDTO, error)
	GetStoreAdditivesByIDs(storeID uint, IDs []uint) ([]types.StoreAdditiveDTO, error)
	GetStoreAdditiveCategoriesByProductSize(storeID, productSizeID uint, filter *types.StoreAdditiveCategoriesFilter) ([]types.StoreAdditiveCategoryDTO, error)
	GetStoreAdditiveByID(storeAdditiveID uint, filter *contexts.StoreContextFilter) (*types.StoreAdditiveDetailsDTO, error)
	DeleteStoreAdditive(storeID, storeAdditiveID uint) error
}

type storeAdditiveService struct {
	repo               StoreAdditiveRepository
	ingredientsRepo    ingredients.IngredientRepository
	storageRepo        storage.StorageRepository
	transactionManager TransactionManager
	logger             *zap.SugaredLogger
}

func NewStoreAdditiveService(
	repo StoreAdditiveRepository,
	ingredientsRepo ingredients.IngredientRepository,
	storageRepo storage.StorageRepository,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) StoreAdditiveService {
	return &storeAdditiveService{
		repo:               repo,
		ingredientsRepo:    ingredientsRepo,
		storageRepo:        storageRepo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (s *storeAdditiveService) CreateStoreAdditives(storeID uint, dtos []types.CreateStoreAdditiveDTO) ([]uint, error) {
	inputAdditiveIDs := make([]uint, len(dtos))
	storeAdditives := make([]data.StoreAdditive, len(dtos))

	for i, dto := range dtos {
		storeAdditives[i] = *types.CreateToStoreAdditive(&dto, storeID)
		storeAdditives[i].StoreID = storeID
		inputAdditiveIDs[i] = storeAdditives[i].AdditiveID
	}

	ingredientIDs, err := s.formAddStockDTOsFromAdditives(inputAdditiveIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("error forming additional stock DTOs: %w", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	ids, err := s.transactionManager.CreateStoreAdditivesWithStocks(storeID, storeAdditives, ingredientIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create store additives: %w", err)
		s.logger.Error(wrappedErr)
		return nil, err
	}

	return ids, nil
}

func (s *storeAdditiveService) GetStoreAdditiveCategoriesByProductSize(storeID, storeProductSizeID uint, filter *types.StoreAdditiveCategoriesFilter) ([]types.StoreAdditiveCategoryDTO, error) {
	categories, err := s.repo.GetStoreAdditiveCategories(storeID, storeProductSizeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	categoryDTOs := make([]types.StoreAdditiveCategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = *types.ConvertToStoreAdditiveCategoryDTO(&category)
	}

	return categoryDTOs, nil
}

func (s *storeAdditiveService) GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]types.StoreAdditiveDTO, error) {
	storeAdditives, err := s.repo.GetStoreAdditives(storeID, filter)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve store additives", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	storeAdditiveDTOs := make([]types.StoreAdditiveDTO, len(storeAdditives))
	for i, storeAdditive := range storeAdditives {
		storeAdditiveDTOs[i] = *types.ConvertToStoreAdditiveDTO(&storeAdditive)
	}

	return storeAdditiveDTOs, nil
}

func (s *storeAdditiveService) GetAdditivesListToAdd(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]additiveTypes.AdditiveDTO, error) {
	additives, err := s.repo.GetAvailableAdditivesToAdd(storeID, filter)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve list additives to add for store", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	additiveDTOs := make([]additiveTypes.AdditiveDTO, len(additives))
	for i, additive := range additives {
		additiveDTOs[i] = *additiveTypes.ConvertToAdditiveDTO(&additive)
	}

	return additiveDTOs, nil
}

func (s *storeAdditiveService) GetStoreAdditivesByIDs(storeID uint, IDs []uint) ([]types.StoreAdditiveDTO, error) {
	storeAdditives, err := s.repo.GetStoreAdditivesByIDs(storeID, IDs)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve store additives by id list", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	storeAdditiveDTOs := make([]types.StoreAdditiveDTO, len(storeAdditives))
	for i, storeAdditive := range storeAdditives {
		storeAdditiveDTOs[i] = *types.ConvertToStoreAdditiveDTO(&storeAdditive)
	}

	return storeAdditiveDTOs, nil
}

func (s *storeAdditiveService) GetStoreAdditiveByID(storeAdditiveID uint, filter *contexts.StoreContextFilter) (*types.StoreAdditiveDetailsDTO, error) {
	storeAdditive, err := s.repo.GetStoreAdditiveWithDetailsByID(storeAdditiveID, filter)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve store additive", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	return types.ConvertToStoreAdditiveDetailsDTO(storeAdditive), nil
}

func (s *storeAdditiveService) UpdateStoreAdditive(storeID, storeAdditiveID uint, dto *types.UpdateStoreAdditiveDTO) error {
	storeAdditive := types.UpdateToStoreAdditive(dto)

	err := s.repo.UpdateStoreAdditive(storeID, storeAdditiveID, storeAdditive)
	if err != nil {
		wrappedError := utils.WrapError("failed to update store additive", err)
		s.logger.Error(wrappedError)
		return wrappedError
	}

	return nil
}

func (s *storeAdditiveService) DeleteStoreAdditive(storeID, storeAdditiveID uint) error {
	err := s.repo.DeleteStoreAdditive(storeID, storeAdditiveID)
	if err != nil {
		wrappedError := utils.WrapError("failed to delete store additive", err)
		s.logger.Error(wrappedError)
		return wrappedError
	}

	return nil
}

func (s *storeAdditiveService) formAddStockDTOsFromAdditives(additiveIDs []uint) ([]uint, error) {
	ingredientsList, err := s.ingredientsRepo.GetIngredientsForAdditives(additiveIDs)
	if err != nil {
		return nil, utils.WrapError("could not get ingredients", err)
	}

	ingredientIDs := make([]uint, len(ingredientsList))
	for i, ingredient := range ingredientsList {
		ingredientIDs[i] = ingredient.ID
	}
	return ingredientIDs, nil
}
