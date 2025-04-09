package storeProvisions

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions/types"
	"go.uber.org/zap"
)

type StoreProvisionService interface {
	GetStoreProvisionByID(storeID, storeProvisionID uint) (*types.StoreProvisionDetailsDTO, error)
	GetStoreProvisions(storeID uint, filter *types.StoreProvisionFilterDTO) ([]types.StoreProvisionDTO, error)
	CreateStoreProvision(storeID uint, dto *types.CreateStoreProvisionDTO) (*types.StoreProvisionDTO, error)
	UpdateStoreProvision(storeID, storeProvisionID uint, dto *types.UpdateStoreProvisionDTO) (*types.StoreProvisionDTO, error)
	CompleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error)
	DeleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error)
}

type storeProvisionService struct {
	repo                StoreProvisionRepository
	ingredientsRepo     ingredients.IngredientRepository
	provisionRepo       provisions.ProvisionRepository
	notificationService notifications.NotificationService
	transactionManager  TransactionManager
	logger              *zap.SugaredLogger
}

func NewStoreProvisionService(
	repo StoreProvisionRepository,
	ingredientsRepo ingredients.IngredientRepository,
	provisionRepo provisions.ProvisionRepository,
	notificationService notifications.NotificationService,
	transactionManager TransactionManager,
	logger *zap.SugaredLogger,
) StoreProvisionService {
	return &storeProvisionService{
		repo:                repo,
		ingredientsRepo:     ingredientsRepo,
		provisionRepo:       provisionRepo,
		notificationService: notificationService,
		transactionManager:  transactionManager,
		logger:              logger,
	}
}

func (s *storeProvisionService) GetStoreProvisionByID(storeID, storeProvisionID uint) (*types.StoreProvisionDetailsDTO, error) {
	provision, err := s.repo.GetStoreProvisionWithDetailsByID(storeID, storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to get store provision by ID %d: %w", storeProvisionID, err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	return types.MapToStoreProvisionDetailsDTO(provision), nil
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
	centralCatalogProvision, err := s.provisionRepo.GetProvisionWithDetailsByID(dto.ProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to get store provision by ID %d: %w", dto.ProvisionID, err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	count, err := s.repo.CountStoreProvisionsToday(storeID, centralCatalogProvision.ID)
	if err != nil {
		wrapped := fmt.Errorf("failed to count store provisions today: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	if count >= centralCatalogProvision.LimitPerDay {
		return nil, types.ErrStoreProvisionDailyLimitReached
	}

	storeProvision := types.CreateToStoreProvisionModel(storeID, dto, centralCatalogProvision)

	ingredientIDs, err := s.formAddStockDTOsFromProvisions([]uint{storeProvision.ProvisionID})
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store provision: ", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	_, err = s.transactionManager.CreateStoreProvisionWithStocks(storeProvision, ingredientIDs)
	if err != nil {
		wrapped := fmt.Errorf("failed to create store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}
	return types.MapToStoreProvisionDTO(storeProvision), nil
}

func (s *storeProvisionService) UpdateStoreProvision(storeID, storeProvisionID uint, dto *types.UpdateStoreProvisionDTO) (*types.StoreProvisionDTO, error) {
	storeProvision, err := s.repo.GetStoreProvisionWithDetailsByID(storeID, storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to get store provision by ID %d: %w", storeProvisionID, err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	if storeProvision.Status != data.STORE_PROVISION_STATUS_PREPARING {
		wrapped := fmt.Errorf("failed to update store provision: %w", types.ErrProvisionCompleted)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	updateModels, err := types.UpdateToStoreProvisionModels(storeProvision, dto)
	if err != nil {
		wrapped := fmt.Errorf("failed to update store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	err = s.repo.SaveStoreProvisionWithAssociations(updateModels)
	if err != nil {
		wrapped := fmt.Errorf("failed to save store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	return types.MapToStoreProvisionDTO(storeProvision), nil
}

func (s *storeProvisionService) CompleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error) {
	provision, err := s.repo.GetStoreProvisionByID(storeID, storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to find store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	if provision.Status != data.STORE_PROVISION_STATUS_PREPARING {
		wrapped := fmt.Errorf("failed to update store provision: %w", types.ErrProvisionCompleted)
		return nil, wrapped
	}

	provision.Status = data.STORE_PROVISION_STATUS_COMPLETED

	currentTime := time.Now().UTC()
	provision.CompletedAt = &currentTime

	expirationTime := currentTime.Add(time.Duration(provision.ExpirationInMinutes) * time.Minute)
	provision.ExpiresAt = &expirationTime

	err = s.repo.SaveStoreProvision(provision)
	if err != nil {
		wrapped := fmt.Errorf("failed to complete store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	return types.MapToStoreProvisionDTO(provision), nil
}

func (s *storeProvisionService) DeleteStoreProvision(storeID, storeProvisionID uint) (*types.StoreProvisionDTO, error) {
	storeProvision, err := s.repo.GetStoreProvisionByID(storeID, storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to get store provision by ID %d: %w", storeProvisionID, err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}

	if storeProvision.Status == data.STORE_PROVISION_STATUS_COMPLETED {
		wrapped := fmt.Errorf("failed to delete store provision by ID %d: %w", storeProvisionID, types.ErrProvisionCompleted)
		return nil, wrapped
	}

	err = s.repo.DeleteStoreProvision(storeProvisionID)
	if err != nil {
		wrapped := fmt.Errorf("failed to delete store provision: %w", err)
		s.logger.Error(wrapped)
		return nil, wrapped
	}
	return types.MapToStoreProvisionDTO(storeProvision), nil
}

func (s *storeProvisionService) formAddStockDTOsFromProvisions(provisionIDs []uint) ([]uint, error) {
	ingredientsList, err := s.ingredientsRepo.GetIngredientsForProvisions(provisionIDs)
	if err != nil {
		return nil, utils.WrapError("could not get ingredients", err)
	}

	ingredientIDs := make([]uint, len(ingredientsList))
	for i, ingredient := range ingredientsList {
		ingredientIDs[i] = ingredient.ID
	}
	return ingredientIDs, nil
}
