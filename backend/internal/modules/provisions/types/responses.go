package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500ProvisionCreate          = localization.NewResponseKey(500, data.ProvisionComponent, data.CreateOperation.ToString())
	Response500ProvisionUpdate          = localization.NewResponseKey(500, data.ProvisionComponent, data.UpdateOperation.ToString())
	Response500ProvisionDelete          = localization.NewResponseKey(500, data.ProvisionComponent, data.DeleteOperation.ToString())
	Response500ProvisionGet             = localization.NewResponseKey(500, data.ProvisionComponent, data.GetOperation.ToString())
	Response409ProvisionCreateDuplicate = localization.NewResponseKey(409, data.ProvisionComponent, data.CreateOperation.ToString(), "DUPLICATE")
	Response409ProvisionDeleteInUse     = localization.NewResponseKey(409, data.ProvisionComponent, data.DeleteOperation.ToString(), "IN_USE")
	Response404Provision                = localization.NewResponseKey(404, data.ProvisionComponent)
	Response400Provision                = localization.NewResponseKey(400, data.ProvisionComponent)
	Response201Provision                = localization.NewResponseKey(201, data.ProvisionComponent)
	Response200ProvisionUpdate          = localization.NewResponseKey(200, data.ProvisionComponent, data.UpdateOperation.ToString())
	Response200ProvisionDelete          = localization.NewResponseKey(200, data.ProvisionComponent, data.DeleteOperation.ToString())
)
