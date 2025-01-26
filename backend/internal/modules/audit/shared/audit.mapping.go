package shared

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
)

var auditActions = make(map[AuditActionCore]func() data.AuditDetails)
var defaultFactory = func() data.AuditDetails {
	return &data.BaseDetails{}
}

func GetAuditActionDetailsFactory(core AuditActionCore) func() data.AuditDetails {
	zapLogger := logger.GetZapSugaredLogger()
	factory, ok := auditActions[core]

	if !ok {
		zapLogger.Warnf("audit action not found for '%v', using default factory", core)
		return defaultFactory
	}

	return factory
}

func NewAuditActionBaseFactory(
	operationType data.OperationType,
	componentName data.ComponentName,
) func(details *data.BaseDetails) AuditActionBase {

	core := AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		panic(fmt.Errorf("duplicate audit action core found: %v", core))
	}

	auditActions[core] = func() data.AuditDetails {
		return defaultFactory()
	}

	return func(baseDetails *data.BaseDetails) AuditActionBase {
		return AuditActionBase{
			Core:    core,
			Details: baseDetails,
		}
	}
}

func NewAuditActionExtendedFactory[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) func(baseDetails *data.BaseDetails, dto T) AuditActionExtended {

	core := AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		panic(fmt.Errorf("duplicate audit action core found: %v", core))
	}

	auditActions[core] = func() data.AuditDetails {
		return &data.ExtendedDetails{
			DTO: dto,
		}
	}

	return func(baseDetails *data.BaseDetails, dto T) AuditActionExtended {
		return AuditActionExtended{
			Core: core,
			Details: &data.ExtendedDetails{
				BaseDetails: *baseDetails,
				DTO:         dto,
			},
		}
	}
}

func NewAuditStoreActionExtendedFactory[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) func(baseDetails *data.BaseDetails, dto T, storeID uint) AuditStoreActionExtended {

	core := AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		panic(fmt.Errorf("duplicate audit action core found: %v", core))
	}

	auditActions[core] = func() data.AuditDetails {
		return &data.ExtendedDetailsStore{
			ExtendedDetails: data.ExtendedDetails{
				DTO: dto,
			},
		}
	}

	return func(baseDetails *data.BaseDetails, dto T, storeID uint) AuditStoreActionExtended {
		return AuditStoreActionExtended{
			Core: core,
			Details: &data.ExtendedDetailsStore{
				ExtendedDetails: data.ExtendedDetails{
					BaseDetails: *baseDetails,
					DTO:         dto,
				},
				StoreInfo: data.StoreInfo{
					StoreID: storeID,
				},
			},
		}
	}
}

func NewAuditWarehouseActionExtendedFactory[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) func(baseDetails *data.BaseDetails, dto T, warehouseID uint) AuditWarehouseActionExtended {

	core := AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		panic(fmt.Errorf("duplicate audit action core found: %v", core))
	}

	auditActions[core] = func() data.AuditDetails {
		return &data.ExtendedDetailsWarehouse{
			ExtendedDetails: data.ExtendedDetails{
				DTO: dto,
			},
		}
	}

	return func(baseDetails *data.BaseDetails, dto T, warehouseID uint) AuditWarehouseActionExtended {
		return AuditWarehouseActionExtended{
			Core: core,
			Details: &data.ExtendedDetailsWarehouse{
				ExtendedDetails: data.ExtendedDetails{
					BaseDetails: *baseDetails,
					DTO:         dto,
				},
				WarehouseInfo: data.WarehouseInfo{
					WarehouseID: warehouseID,
				},
			},
		}
	}
}

func NewAuditFranchiseeActionExtendedFactory[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) func(baseDetails *data.BaseDetails, dto T, storeID uint) AuditFranchiseeActionExtended {

	core := AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		panic(fmt.Errorf("duplicate audit action core found: %v", core))
	}

	auditActions[core] = func() data.AuditDetails {
		return &data.ExtendedDetailsStore{
			ExtendedDetails: data.ExtendedDetails{
				DTO: dto,
			},
		}
	}

	return func(baseDetails *data.BaseDetails, dto T, franchiseeID uint) AuditFranchiseeActionExtended {
		return AuditFranchiseeActionExtended{
			Core: core,
			Details: &data.ExtendedDetailsFranchisee{
				ExtendedDetails: data.ExtendedDetails{
					BaseDetails: *baseDetails,
					DTO:         dto,
				},
				FranchiseeInfo: data.FranchiseeInfo{
					FranchiseeID: franchiseeID,
				},
			},
		}
	}
}

func NewAuditRegionActionExtendedFactory[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) func(baseDetails *data.BaseDetails, dto T, storeID uint) AuditRegionActionExtended {

	core := AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		panic(fmt.Errorf("duplicate audit action core found: %v", core))
	}

	auditActions[core] = func() data.AuditDetails {
		return &data.ExtendedDetailsStore{
			ExtendedDetails: data.ExtendedDetails{
				DTO: dto,
			},
		}
	}

	return func(baseDetails *data.BaseDetails, dto T, regionID uint) AuditRegionActionExtended {
		return AuditRegionActionExtended{
			Core: core,
			Details: &data.ExtendedDetailsRegion{
				ExtendedDetails: data.ExtendedDetails{
					BaseDetails: *baseDetails,
					DTO:         dto,
				},
				RegionInfo: data.RegionInfo{
					RegionID: regionID,
				},
			},
		}
	}
}
