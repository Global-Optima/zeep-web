package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500AdminEmployeeGet    = localization.NewResponseKey(500, data.AdminEmployeeComponent, data.GetOperation.ToString())
	Response500AdminEmployeeCreate = localization.NewResponseKey(500, data.AdminEmployeeComponent, data.CreateOperation.ToString())

	Response400AdminEmployee = localization.NewResponseKey(400, data.AdminEmployeeComponent)

	Response201AdminEmployee = localization.NewResponseKey(201, data.AdminEmployeeComponent)
)
