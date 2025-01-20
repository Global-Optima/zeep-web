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
