package storeAdditives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
)

type StoreAdditiveService interface {
	CreateStoreAdditives(storeID uint, dtos []types.CreateStoreAdditiveDTO) ([]uint, error)
	UpdateStoreAdditive(storeID, storeAdditiveID uint, dto *types.UpdateStoreAdditiveDTO) error
	GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]types.StoreAdditiveDTO, error)
	GetStoreAdditiveCategories(storeID uint, filter *additiveTypes.AdditiveCategoriesFilterQuery) ([]types.StoreAdditiveCategoryDTO, error)
	GetStoreAdditiveByID(storeID, storeAdditiveID uint) (*types.StoreAdditiveDTO, error)
	DeleteStoreAdditive(storeID, storeAdditiveID uint) error
}

type storeAdditiveService struct {
	repo   StoreAdditiveRepository
	logger *zap.SugaredLogger
}

func NewStoreAdditiveService(repo StoreAdditiveRepository, logger *zap.SugaredLogger) StoreAdditiveService {
	return &storeAdditiveService{
		repo:   repo,
		logger: logger,
	}
}

func (s *storeAdditiveService) CreateStoreAdditives(storeID uint, dtos []types.CreateStoreAdditiveDTO) ([]uint, error) {
	var storeAdditives []data.StoreAdditive
	for _, dto := range dtos {
		storeAdditives = append(storeAdditives, *types.CreateToStoreAdditive(&dto, storeID))
	}

	ids, err := s.repo.CreateStoreAdditives(storeAdditives)
	if err != nil {
		wrappedError := utils.WrapError("failed to create store additives", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	return ids, nil
}

func (s *storeAdditiveService) GetStoreAdditiveCategories(storeID uint, filter *additiveTypes.AdditiveCategoriesFilterQuery) ([]types.StoreAdditiveCategoryDTO, error) {
	categories, err := s.repo.GetStoreAdditiveCategories(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if len(categories) == 0 {
		return []types.StoreAdditiveCategoryDTO{}, nil
	}

	var categoryDTOs []types.StoreAdditiveCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, *types.ConvertToStoreAdditiveCategoryDTO(&category))
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
	for i, additive := range storeAdditives {
		storeAdditiveDTOs[i] = *types.ConvertToStoreAdditiveDTO(&additive)
	}

	return storeAdditiveDTOs, nil
}

func (s *storeAdditiveService) GetStoreAdditiveByID(storeID, storeAdditiveID uint) (*types.StoreAdditiveDTO, error) {
	storeAdditive, err := s.repo.GetStoreAdditiveByID(storeID, storeAdditiveID)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve store additive", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	return types.ConvertToStoreAdditiveDTO(storeAdditive), nil
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
