package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500FranchiseeEmployeeGet    = localization.NewResponseKey(500, data.FranchiseeEmployeeComponent, data.GetOperation.ToString())
	Response500FranchiseeEmployeeCreate = localization.NewResponseKey(500, data.FranchiseeEmployeeComponent, data.CreateOperation.ToString())
	Response500FranchiseeEmployeeUpdate = localization.NewResponseKey(500, data.FranchiseeEmployeeComponent, data.UpdateOperation.ToString())
	Response500FranchiseeEmployeeDelete = localization.NewResponseKey(500, data.FranchiseeEmployeeComponent, data.DeleteOperation.ToString())

	Response400FranchiseeEmployee = localization.NewResponseKey(400, data.FranchiseeEmployeeComponent)

	Response200FranchiseeEmployeeUpdate = localization.NewResponseKey(200, data.FranchiseeEmployeeComponent, data.UpdateOperation.ToString())
	Response200FranchiseeEmployeeDelete = localization.NewResponseKey(200, data.FranchiseeEmployeeComponent, data.DeleteOperation.ToString())
	Response201FranchiseeEmployee       = localization.NewResponseKey(201, data.FranchiseeEmployeeComponent)
)
