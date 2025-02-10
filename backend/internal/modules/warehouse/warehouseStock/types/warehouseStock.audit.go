package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

type WarehouseStockPayloads struct {
	UpdateWarehouseStockDTO   *UpdateWarehouseStockDTO    `json:"updateWarehouseStockDTO,omitempty"`
	AddWarehouseStockMaterial []AddWarehouseStockMaterial `json:"addWarehouseStockMaterial,omitempty"`
	ReceiveWarehouseDelivery  *ReceiveWarehouseDelivery   `json:"receiveWarehouseDelivery,omitempty"`
}

var (
	UpdateWarehouseStockAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.UpdateOperation, data.WarehouseStockComponent, &WarehouseStockPayloads{})
)
