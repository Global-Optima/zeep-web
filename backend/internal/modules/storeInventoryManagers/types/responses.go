package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StoreInventoryManagerGetFrozenInventory   = localization.NewResponseKey(500, data.StoreInventoryManagerComponent, data.GetOperation.ToString(), "FROZEN_INVENTORY")
	Response500StoreInventoryManagerRecalculateInventory = localization.NewResponseKey(500, data.StoreInventoryManagerComponent, "RECALCULATE_INVENTORY")
	Response200StoreInventoryManagerRecalculateInventory = localization.NewResponseKey(200, data.StoreInventoryManagerComponent, "RECALCULATE_INVENTORY")
)
