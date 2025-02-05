package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

type AuditPayloads struct {
	CreateStockRequestDTO            *CreateStockRequestDTO            `json:"createStockRequestDTO,omitempty"`
	RejectStockRequestStatusDTO      *RejectStockRequestStatusDTO      `json:"rejectStockRequestStatusDTO,omitempty"`
	AcceptWithChangeRequestStatusDTO *AcceptWithChangeRequestStatusDTO `json:"acceptWithChangeRequestStatusDTO,omitempty"`
	StockMaterials                   []StockRequestStockMaterialDTO    `json:"stockRequestStockMaterials,omitempty"`
}

var (
	CreateStockRequestAuditFactory = shared.NewAuditActionExtendedFactory(
		data.CreateOperation, data.StockRequestComponent, &AuditPayloads{})

	UpdateStockRequestStatusAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.StockRequestComponent, &AuditPayloads{})

	DeleteStockRequestAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.StockRequestComponent)
)
