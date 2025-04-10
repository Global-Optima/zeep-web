package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

type UpdateStoreProvisionFields struct {
	*UpdateStoreProvisionDTO
	Status *data.StoreProvisionStatus `json:"status,omitempty"`
}

var (
	CreateStoreProvisionAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreProvisionComponent, &CreateStoreProvisionDTO{})

	UpdateStoreProvisionAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreProvisionComponent, &UpdateStoreProvisionFields{})

	DeleteStoreProvisionAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreProvisionComponent, &struct{}{})
)
