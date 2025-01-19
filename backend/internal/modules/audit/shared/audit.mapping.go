package shared

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
)

var auditActions = make(map[types.AuditActionCore]func() data.AuditDetails)
var defaultFactory = func() data.AuditDetails {
	return &data.BaseDetails{}
}

func GetAuditActionDetailsFactory(core types.AuditActionCore) func() data.AuditDetails {
	zapLogger := logger.GetZapSugaredLogger()
	factory, ok := auditActions[core]

	if !ok {
		zapLogger.Warnf("audit action not found for '%v', using default factory", core)
		return defaultFactory
	}

	return factory
}

func NewAuditActionExtendedFactory[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) func(baseDetails *data.BaseDetails, dto T) types.AuditActionExtended {
	zapLogger := logger.GetZapSugaredLogger()

	core := types.AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		zapLogger.Warn("Duplicate audit action core found, first copy was overwritten!")
	}

	auditActions[core] = func() data.AuditDetails {
		return &data.ExtendedDetails{
			DTO: dto,
		}
	}

	return func(baseDetails *data.BaseDetails, dto T) types.AuditActionExtended {
		return types.AuditActionExtended{
			Core: core,
			Details: &data.ExtendedDetails{
				BaseDetails: *baseDetails,
				DTO:         dto,
			},
		}
	}
}

func NewAuditActionBaseFactory(
	operationType data.OperationType,
	componentName data.ComponentName,
) func(details *data.BaseDetails) types.AuditActionBase {
	zapLogger := logger.GetZapSugaredLogger()

	core := types.AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}

	if _, ok := auditActions[core]; ok {
		zapLogger.Warn("Duplicate audit action core found, first copy was overwritten!")
	}

	auditActions[core] = func() data.AuditDetails {
		return defaultFactory()
	}

	return func(baseDetails *data.BaseDetails) types.AuditActionBase {
		return types.AuditActionBase{
			Core:    core,
			Details: baseDetails,
		}
	}
}

/*func NewAuditActionMultiple[T any](
	operationType data.OperationType,
	componentName data.ComponentName,
	dto T,
) types.AuditAction[*data.MultipleItemDetails[T]] {
	zapLogger := logger.GetZapSugaredLogger()

	core := types.AuditActionCore{
		OperationType: operationType,
		ComponentName: componentName,
	}
	details := defaultFactory()

	if _, ok := auditActions[core]; ok {
		zapLogger.Warn("Duplicate audit action core found, first copy was overwritten!")
	}

	auditActions[core] = func() data.AuditDetails {
		return details
	}

	return types.AuditAction[*data.MultipleItemDetails[T]]{
		OperationType: operationType,
		ComponentName: componentName,
		Details: &data.MultipleItemDetails[T]{
			DTO: dto,
		},
	}
}*/
