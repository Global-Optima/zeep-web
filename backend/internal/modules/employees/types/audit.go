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

	CreateRegionManagerAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.CreateOperation, data.RegionManagerComponent, &CreateEmployeeDTO{})

	UpdateStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreEmployeeComponent, &UpdateStoreEmployeeDTO{})

	UpdateWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.UpdateOperation, data.WarehouseEmployeeComponent, &UpdateWarehouseEmployeeDTO{})

	UpdateFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.UpdateOperation, data.FranchiseeEmployeeComponent, &UpdateFranchiseeEmployeeDTO{})

	UpdateRegionManagerEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.UpdateOperation, data.RegionManagerComponent, &UpdateRegionManagerEmployeeDTO{})

	DeleteFranchiseeEmployeeAuditFactory = shared.NewAuditFranchiseeActionExtendedFactory(
		data.DeleteOperation, data.FranchiseeEmployeeComponent, struct{}{})

	DeleteStoreEmployeeAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreEmployeeComponent, struct{}{})

	DeleteWarehouseEmployeeAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.DeleteOperation, data.WarehouseEmployeeComponent, struct{}{})

	DeleteRegionManagerEmployeeAuditFactory = shared.NewAuditRegionActionExtendedFactory(
		data.DeleteOperation, data.RegionManagerComponent, struct{}{})
)
