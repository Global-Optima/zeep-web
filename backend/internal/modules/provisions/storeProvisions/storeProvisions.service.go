package storeProvisions

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"go.uber.org/zap"
	"time"
)

type StoreProvisionService interface {
	GetStoreProvisionByID(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error)
	GetStoreProvisions(storeID uint, filter *types.StoreProvisionFilterDTO) ([]types.StoreProvisionDTO, error)
	CreateStoreProvision(storeID uint, dto *types.CreateStoreProvisionDTO) (*types.StoreProvisionDTO, error)
	UpdateStoreProvision(storeID, storeProvisionID uint, dto *types.UpdateStoreProvisionDTO) (*types.StoreProvisionDTO, error)
	CompleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error)
	DeleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error)
}

type storeProvisionService struct {
	repo                StoreProvisionRepository
	provisionRepo       provisions.ProvisionRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewStoreProvisionService(repo StoreProvisionRepository, notificationService notifications.NotificationService, storageRepo storage.StorageRepository, logger *zap.SugaredLogger) StoreProvisionService {
	return &storeProvisionService{
		repo:                repo,
		logger:              logger,
		notificationService: notificationService,
	}
}

func (s *storeProvisionService) GetStoreProvisionByID(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error) {
	provision, err := s.repo.GetStoreProvisionWithDetailsByID(storeID, storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to get store provision by ID %d: %w", storeProvisionID, err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	return types.MapToStoreProvisionDTO(provision), nil
}

func (s *storeProvisionService) GetStoreProvisions(storeID uint, filter *types.StoreProvisionFilterDTO) ([]types.StoreProvisionDTO, error) {
	storeProvisions, err := s.repo.GetStoreProvisions(storeID, filter)
	if err != nil {
		wrapped := fmt.Errorf("failed to list store provisions: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	dtos := make([]types.StoreProvisionDTO, len(storeProvisions))
	for i, sp := range storeProvisions {
		dtos[i] = *types.MapToStoreProvisionDTO(&sp)
	}

	return dtos, nil
}

func (s *storeProvisionService) CreateStoreProvision(storeID uint, dto *types.CreateStoreProvisionDTO) (*types.StoreProvisionDTO, error) {
	existingProvision, err := s.provisionRepo.GetProvisionWithDetailsByID(dto.ProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to get store provision by ID %d: %w", dto.ProvisionID, err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	provisionIngredientIDs := make([]uint, len(existingProvision.ProvisionIngredients))
	for i, ingredient := range existingProvision.ProvisionIngredients {
		provisionIngredientIDs[i] = ingredient.IngredientID
	}

	if err := validateStoreProvisionIngredients(dto.Ingredients, provisionIngredientIDs); err != nil {
		return nil, fmt.Errorf("ingredient validation failed: %w", err)
	}

	count, err := s.repo.CountStoreProvisionsToday(storeID, existingProvision.ID)
	if err != nil {
		wrapped := fmt.Errorf("failed to count store provisions today: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	if count >= existingProvision.LimitPerDay {
		return nil, types.ErrStoreProvisionDailyLimitReached
	}

	storeProvision := types.CreateToStoreProvisionModel(storeID, dto)

	_, err = s.repo.CreateStoreProvision(storeProvision)
	if err != nil {
		wrapped := fmt.Errorf("failed to create store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}
	return types.MapToStoreProvisionDTO(storeProvision), nil
}

func (s *storeProvisionService) UpdateStoreProvision(storeID, storeProvisionID uint, dto *types.UpdateStoreProvisionDTO) (*types.StoreProvisionDTO, error) {
	storeProvision, err := s.repo.GetStoreProvisionByID(storeID, storeProvisionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find store provision: %w", err)
	}

	if storeProvision.Status != data.PROVISION_STATUS_PREPARING {
		return nil, fmt.Errorf("failed to update store provision: %w", types.ErrProvisionAlreadyCompleted)
	}

	provisionIngredientIDs, err := s.provisionRepo.GetProvisionIngredientIDs(storeProvision.ProvisionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find store provision ingredient IDs: %w", err)
	}

	if err := validateStoreProvisionIngredients(dto.Ingredients, provisionIngredientIDs); err != nil {
		return nil, fmt.Errorf("ingredient validation failed: %w", err)
	}

	if err := types.UpdateToStoreProvisionModel(storeProvision, dto); err != nil {
		return nil, fmt.Errorf("failed to apply update: %w", err)
	}

	err = s.repo.SaveStoreProvision(storeProvision)
	if err != nil {
		wrapped := fmt.Errorf("failed to save store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	return types.MapToStoreProvisionDTO(storeProvision), nil
}

func validateStoreProvisionIngredients(
	ingredients []ingredientTypes.SelectedIngredientDTO,
	provisionIngredientIDs []uint,
) error {
	expectedSet := make(map[uint]struct{}, len(provisionIngredientIDs))
	for _, id := range provisionIngredientIDs {
		expectedSet[id] = struct{}{}
	}

	inputSet := make(map[uint]struct{}, len(ingredients))
	for _, ing := range ingredients {
		inputSet[ing.IngredientID] = struct{}{}
	}

	if len(expectedSet) != len(inputSet) {
		return types.ErrStoreProvisionIngredientMismatch
	}

	for id := range expectedSet {
		if _, found := inputSet[id]; !found {
			return types.ErrStoreProvisionIngredientMismatch
		}
	}

	return nil
}

func (s *storeProvisionService) CompleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error) {
	provision, err := s.repo.GetStoreProvisionByID(storeID, storeProvisionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find store provision: %w", err)
	}

	provision.Status = data.PROVISION_STATUS_COMPLETED

	currentTime := time.Now().UTC()
	provision.CompletedAt = &currentTime

	expirationTime := currentTime.Add(time.Duration(provision.ExpirationInHours) * time.Hour)
	provision.ExpiresAt = &expirationTime

	err = s.repo.SaveStoreProvision(provision)
	if err != nil {
		return nil, fmt.Errorf("failed to complete store provision: %w", err)
	}

	return types.MapToStoreProvisionDTO(provision), nil
}

func (s *storeProvisionService) DeleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error) {
	deleted, err := s.repo.DeleteStoreProvision(storeID, storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to delete store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}
	return types.MapToStoreProvisionDTO(deleted), nil
}
