package storeInventoryManagers

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"go.uber.org/zap"
)

type StoreInventoryManagerService interface {
	/*RecalculateStoreInventory(storeID uint, input *types.RecalculateInput) error
	RecalculateStoreProduct(storeID uint, storeProductIDs []uint) error
	RecalculateStoreAdditive(storeID uint, storeAdditiveIDs []uint) error*/
	CalculateFrozenInventory(storeID uint, filter *types.FrozenInventoryFilter) (*types.FrozenInventoryDTO, error)
}

type storeInventoryManagerService struct {
	repo                StoreInventoryManagerRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewStoreInventoryManagerService(
	repo StoreInventoryManagerRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) StoreInventoryManagerService {
	return &storeInventoryManagerService{
		repo:                repo,
		notificationService: notificationService,
		logger:              logger,
	}
}

func (s *storeInventoryManagerService) CalculateFrozenInventory(storeID uint, filter *types.FrozenInventoryFilter) (*types.FrozenInventoryDTO, error) {
	frozenInventory, err := s.repo.CalculateFrozenInventory(storeID, filter)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to calculate frozen inventory: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	inventory := frozenInventory.GetIDs()

	ingredientsList, err := s.repo.GetIngredientsByIDs(inventory.IngredientIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get ingredients by ids: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	provisionsList, err := s.repo.GetProvisionsByIDs(inventory.ProvisionIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get provisions by ids: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.ConvertToFrozenInventoryDTO(frozenInventory, ingredientsList, provisionsList), nil
}
