package scheduler

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"go.uber.org/zap"
)

type StoreProvisionCronTasks struct {
	storeProvisionService     storeProvisions.StoreProvisionService
	storeProvisionRepo        storeProvisions.StoreProvisionRepository
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository
	storeService              stores.StoreService
	notificationService       notifications.NotificationService
	logger                    *zap.SugaredLogger
}

func NewStoreProvisionCronTasks(
	storeProvisionService storeProvisions.StoreProvisionService,
	storeProvisionRepo storeProvisions.StoreProvisionRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	storeService stores.StoreService,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) *StoreProvisionCronTasks {
	return &StoreProvisionCronTasks{
		storeProvisionService:     storeProvisionService,
		storeProvisionRepo:        storeProvisionRepo,
		storeInventoryManagerRepo: storeInventoryManagerRepo,
		storeService:              storeService,
		notificationService:       notificationService,
		logger:                    logger,
	}
}

func (tasks *StoreProvisionCronTasks) CheckStoreProvisionNotifications() {
	tasks.logger.Info("Running CheckStoreProvisionNotifications...")

	storesList, err := tasks.storeService.GetAllStoresForNotifications()
	if err != nil {
		tasks.logger.Errorf("Failed to fetch stores: %v", err)
		return
	}

	for _, store := range storesList {
		processedProvisions := make(map[uint]bool)

		storeProvisionList, err := tasks.storeProvisionRepo.GetAllStoreProvisionList(store.ID)
		if err != nil {
			tasks.logger.Errorf("Failed to fetch stock list for store %d: %v", store.ID, err)
			continue
		}

		provisionIDsToRecalculate := make(map[uint]struct{})
		for _, storeProvision := range storeProvisionList {
			if storeProvision.ExpiresAt != nil && storeProvision.ExpiresAt.UTC().Before(time.Now().UTC()) {
				// write unique provisionIDs for recalculation
				if _, exists := provisionIDsToRecalculate[storeProvision.ProvisionID]; exists {
					provisionIDsToRecalculate[storeProvision.ProvisionID] = struct{}{}
				}

				spDetails := &details.StoreProvisionExpirationDetails{
					BaseNotificationDetails: details.BaseNotificationDetails{
						ID:           storeProvision.StoreID,
						FacilityName: storeProvision.Store.Name,
					},
					ItemName:       storeProvision.Provision.Name,
					ExpirationDate: storeProvision.ExpiresAt.Format("2006-01-02 15:04"),
				}

				if err := tasks.notificationService.NotifyStoreProvisionExpiration(spDetails); err != nil {
					tasks.logger.Errorf("failed to send store provision expiration notification: %v", err)
				}
			}

			processedProvisions[storeProvision.ID] = true
		}

		provisionIDs := make([]uint, 0, len(provisionIDsToRecalculate))
		for id := range provisionIDsToRecalculate {
			provisionIDs = append(provisionIDs, id)
		}

		err = tasks.storeInventoryManagerRepo.RecalculateStoreInventory(store.ID, &storeInventoryManagersTypes.RecalculateInput{
			ProvisionIDs: provisionIDs,
		})
	}

	tasks.logger.Info("Check Store Provision Notifications completed.")
}
