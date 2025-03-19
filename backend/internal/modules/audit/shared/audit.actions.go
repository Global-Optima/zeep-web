package shared

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

/*type AuditAction[T data.AuditDetails] struct {
	OperationType data.OperationType
	ComponentName data.ComponentName
	Details       T
}*/

type AuditAction interface {
	GetActionCore() AuditActionCore
	GetActionDetails() data.AuditDetails
}

type AuditActionCore struct {
	OperationType data.OperationType
	ComponentName data.ComponentName
}

func (a AuditActionCore) GetActionCore() AuditActionCore {
	return a
}

func (a AuditActionCore) ToString() string {
	return fmt.Sprintf("%s %s", a.OperationType.ToString(), a.ComponentName.ToString())
}

type AuditActionExtended struct {
	Core    AuditActionCore
	Details *data.ExtendedDetails
}

func (a *AuditActionExtended) GetActionCore() AuditActionCore {
	return a.Core
}

func (a *AuditActionExtended) GetActionDetails() data.AuditDetails {
	return a.Details
}

type AuditStoreActionExtended struct {
	Core    AuditActionCore
	Details *data.ExtendedDetailsStore
}

func (a *AuditStoreActionExtended) GetActionCore() AuditActionCore {
	return a.Core
}

func (a *AuditStoreActionExtended) GetActionDetails() data.AuditDetails {
	return a.Details
}

func (a *AuditStoreActionExtended) GetStoreActionDetails() data.StoreInfo {
	return a.Details.StoreInfo
}

type AuditWarehouseActionExtended struct {
	Core    AuditActionCore
	Details *data.ExtendedDetailsWarehouse
}

func (a *AuditWarehouseActionExtended) GetActionCore() AuditActionCore {
	return a.Core
}

func (a *AuditWarehouseActionExtended) GetActionDetails() data.AuditDetails {
	return a.Details
}

func (a *AuditWarehouseActionExtended) GetWarehouseActionDetails() data.WarehouseInfo {
	return a.Details.WarehouseInfo
}

type AuditActionBase struct {
	Core    AuditActionCore
	Details *data.BaseDetails
}

func (a *AuditActionBase) GetActionCore() AuditActionCore {
	return a.Core
}

func (a *AuditActionBase) GetActionDetails() data.AuditDetails {
	return a.Details
}

type AuditFranchiseeActionExtended struct {
	Core    AuditActionCore
	Details *data.ExtendedDetailsFranchisee
}

func (a *AuditFranchiseeActionExtended) GetActionCore() AuditActionCore {
	return a.Core
}

func (a *AuditFranchiseeActionExtended) GetActionDetails() data.AuditDetails {
	return a.Details
}

func (a *AuditFranchiseeActionExtended) GetFranchiseeActionDetails() data.FranchiseeInfo {
	return a.Details.FranchiseeInfo
}

type AuditRegionActionExtended struct {
	Core    AuditActionCore
	Details *data.ExtendedDetailsRegion
}

func (a *AuditRegionActionExtended) GetActionCore() AuditActionCore {
	return a.Core
}

func (a *AuditRegionActionExtended) GetActionDetails() data.AuditDetails {
	return a.Details
}

func (a *AuditRegionActionExtended) GetRegionActionDetails() data.RegionInfo {
	return a.Details.RegionInfo
}
