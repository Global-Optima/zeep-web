package scheduler

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
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

		var expiredStoreProvisionIDs []uint

		provisionIDsToRecalculate := make(map[uint]struct{})
		for _, storeProvision := range storeProvisionList {
			if storeProvision.Status == data.STORE_PROVISION_STATUS_COMPLETED &&
				storeProvision.ExpiresAt != nil &&
				storeProvision.ExpiresAt.UTC().Before(time.Now().UTC()) {

				// write unique provisionIDs for recalculation
				if _, exists := provisionIDsToRecalculate[storeProvision.ProvisionID]; exists {
					provisionIDsToRecalculate[storeProvision.ProvisionID] = struct{}{}
				}

				//write expired storeProvisionIDs
				expiredStoreProvisionIDs = append(expiredStoreProvisionIDs, storeProvision.ID)

				spDetails := &details.StoreProvisionExpirationDetails{
					BaseNotificationDetails: details.BaseNotificationDetails{
						ID:           storeProvision.StoreID,
						FacilityName: storeProvision.Store.Name,
					},
					ItemName:       storeProvision.Provision.Name,
					CompletionDate: storeProvision.CompletedAt.Format("2006-01-02 15:04"),
				}

				if err := tasks.notificationService.NotifyStoreProvisionExpiration(spDetails); err != nil {
					tasks.logger.Errorf("failed to send store provision expiration notification: %v", err)
				}
			}

			processedProvisions[storeProvision.ID] = true
		}

		if len(expiredStoreProvisionIDs) > 0 {
			if err := tasks.storeProvisionRepo.ExpireStoreProvisions(expiredStoreProvisionIDs); err != nil {
				tasks.logger.Errorf("failed to expire store provisions: %v", err)
			}
		}

		if len(provisionIDsToRecalculate) > 0 {
			provisionIDs := make([]uint, 0, len(provisionIDsToRecalculate))
			for id := range provisionIDsToRecalculate {
				provisionIDs = append(provisionIDs, id)
			}

			err = tasks.storeInventoryManagerRepo.RecalculateStoreInventory(store.ID, &storeInventoryManagersTypes.RecalculateInput{
				ProvisionIDs: provisionIDs,
			})
			if err != nil {
				tasks.logger.Errorf("failed to recalculate inventory for expiredProvisions provision ID count: %v", err)
			}
		}
	}

	tasks.logger.Info("Check Store Provision Notifications completed.")
}
