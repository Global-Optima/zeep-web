package scheduler

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeInventoryManagersTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/sirupsen/logrus"

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
		storeProvisionList, err := tasks.storeProvisionRepo.GetAllCompletedStoreProvisionList(store.ID)
		if err != nil {
			tasks.logger.Errorf("Failed to fetch stock list for store %d: %v", store.ID, err)
			continue
		}

		expiredStoreProvisionIDs, provisionIDsToRecalculate, err := tasks.formExpirationDetailsAndNotify(storeProvisionList)
		if err != nil {
			tasks.logger.Errorf("failed to send store provision expiration notification: %v", err)
		}

		if len(expiredStoreProvisionIDs) > 0 {
			if err := tasks.storeProvisionRepo.ExpireStoreProvisions(expiredStoreProvisionIDs); err != nil {
				tasks.logger.Errorf("failed to expire store provisions: %v", err)
			}
		}

		if len(provisionIDsToRecalculate) > 0 {
			err = tasks.storeInventoryManagerRepo.RecalculateStoreInventory(store.ID, &storeInventoryManagersTypes.RecalculateInput{
				ProvisionIDs: provisionIDsToRecalculate,
			})
			if err != nil {
				tasks.logger.Errorf("failed to recalculate inventory for expiredProvisions provision ID count: %v", err)
			}
		}
		logrus.Infof("notification done, added %d records", len(expiredStoreProvisionIDs))
	}

	tasks.logger.Info("Check Store Provision Notifications completed.")
}

func (tasks *StoreProvisionCronTasks) formExpirationDetailsAndNotify(
	storeProvisionList []data.StoreProvision,
) (expiredProvisionIDs []uint, provisionIDs []uint, err error) {

	provisionIDSet := make(map[uint]struct{})

	for _, storeProvision := range storeProvisionList {
		expiredProvisionIDs = append(expiredProvisionIDs, storeProvision.ID)
		if _, exists := provisionIDSet[storeProvision.ProvisionID]; !exists {
			provisionIDSet[storeProvision.ProvisionID] = struct{}{}
			provisionIDs = append(provisionIDs, storeProvision.ProvisionID)
		}

		spDetails := &details.StoreProvisionExpirationDetails{
			BaseNotificationDetails: details.BaseNotificationDetails{
				ID:           storeProvision.StoreID,
				FacilityName: storeProvision.Store.Name,
			},
			ItemName:       storeProvision.Provision.Name,
			CompletionDate: storeProvision.CompletedAt.Format("2006-01-02 15:04"),
		}

		err = tasks.notificationService.NotifyStoreProvisionExpiration(spDetails)
	}

	return
}
