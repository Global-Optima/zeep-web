package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StoreProvisionCreate              = localization.NewResponseKey(500, data.StoreProvisionComponent, data.CreateOperation.ToString())
	Response500StoreProvisionUpdate              = localization.NewResponseKey(500, data.StoreProvisionComponent, data.UpdateOperation.ToString())
	Response500StoreProvisionDelete              = localization.NewResponseKey(500, data.StoreProvisionComponent, data.DeleteOperation.ToString())
	Response500StoreProvisionGet                 = localization.NewResponseKey(500, data.StoreProvisionComponent, data.GetOperation.ToString())
	Response409StoreProvisionLimit               = localization.NewResponseKey(409, data.StoreProvisionComponent, "LIMIT")
	Response409StoreProvisionCompleted           = localization.NewResponseKey(409, data.StoreProvisionComponent, "COMPLETED")
	Response409StoreProvisionIngredientsMismatch = localization.NewResponseKey(409, data.StoreProvisionComponent, "INGREDIENTS_MISMATCH")
	Response404StoreProvision                    = localization.NewResponseKey(404, data.StoreProvisionComponent)
	Response400StoreProvision                    = localization.NewResponseKey(400, data.StoreProvisionComponent)
	Response201StoreProvision                    = localization.NewResponseKey(201, data.StoreProvisionComponent)
	Response200StoreProvisionUpdate              = localization.NewResponseKey(200, data.StoreProvisionComponent, data.UpdateOperation.ToString())
	Response200StoreProvisionDelete              = localization.NewResponseKey(200, data.StoreProvisionComponent, data.DeleteOperation.ToString())
)
