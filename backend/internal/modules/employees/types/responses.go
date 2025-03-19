package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500EmployeeGet            = localization.NewResponseKey(500, data.EmployeeComponent, "get")
	Response500EmployeeUpdatePassword = localization.NewResponseKey(500, data.EmployeeComponent, "updatePassword")
	Response500EmployeeReassignType   = localization.NewResponseKey(500, data.EmployeeComponent, "reassignType")
	Response500EmployeeGetWorkday     = localization.NewResponseKey(500, data.EmployeeComponent, "getWorkday")
	Response500EmployeeGetWorkdays    = localization.NewResponseKey(500, data.EmployeeComponent, "getWorkdays")
	Response400Employee               = localization.NewResponseKey(400, data.EmployeeComponent)
	Response401Employee               = localization.NewResponseKey(401, data.EmployeeComponent)
)
