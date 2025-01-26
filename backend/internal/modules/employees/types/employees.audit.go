package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.CreateOperation, data.FranchiseeEmployeeComponent, &CreateEmployeeDTO{})

	CreateStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreEmployeeComponent, &CreateEmployeeDTO{})

	CreateWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.CreateOperation, data.WarehouseEmployeeComponent, &CreateEmployeeDTO{})

	CreateRegionEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.CreateOperation, data.RegionEmployeeComponent, &CreateEmployeeDTO{})

	UpdateStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreEmployeeComponent, &UpdateStoreEmployeeDTO{})

	UpdateWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.UpdateOperation, data.WarehouseEmployeeComponent, &UpdateWarehouseEmployeeDTO{})

	UpdateFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.UpdateOperation, data.FranchiseeEmployeeComponent, &UpdateFranchiseeEmployeeDTO{})

	UpdateRegionEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.UpdateOperation, data.RegionEmployeeComponent, &UpdateRegionEmployeeDTO{})

	DeleteFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.DeleteOperation, data.FranchiseeEmployeeComponent, struct{}{})

	DeleteStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreEmployeeComponent, struct{}{})

	DeleteWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.DeleteOperation, data.WarehouseEmployeeComponent, struct{}{})

	DeleteRegionEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.DeleteOperation, data.RegionEmployeeComponent, struct{}{})
)
