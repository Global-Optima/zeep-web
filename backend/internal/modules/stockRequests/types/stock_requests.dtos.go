package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	storeTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/stores/types"
	stockMaterialTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	warehouseTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type CreateStockRequestDTO struct {
	StockMaterials []StockRequestStockMaterialDTO `json:"items" binding:"required"`
}

type StockRequestStockMaterialDTO struct {
	StockMaterialID uint    `json:"stockMaterialId" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
}

type RejectStockRequestStatusDTO struct {
	Comment *string `json:"comment" binding:"required"`
}

type AcceptWithChangeRequestStatusDTO struct {
	Comment *string                        `json:"comment" binding:"required"`
	Items   []StockRequestStockMaterialDTO `json:"items" binding:"required"`
}

type UpdateIngredientDates struct {
	DeliveredDate  time.Time
	ExpirationDate time.Time
}

type StockRequestResponse struct {
	RequestID        uint                             `json:"requestId"`
	Status           data.StockRequestStatus          `json:"status"`
	StoreComment     *string                          `json:"storeComment,omitempty"`
	WarehouseComment *string                          `json:"warehouseComment,omitempty"`
	Store            storeTypes.StoreDTO              `json:"store"`
	Warehouse        warehouseTypes.WarehouseResponse `json:"warehouse"`
	StockMaterials   []StockRequestMaterial           `json:"stockMaterials"`
	CreatedAt        time.Time                        `json:"createdAt"`
	UpdatedAt        time.Time                        `json:"updatedAt"`
}

type StockRequestMaterial struct {
	StockMaterial stockMaterialTypes.StockMaterialsDTO `json:"stockMaterial"`
	Quantity      float64                              `json:"quantity"`
}

type GetStockRequestsFilter struct {
	utils.BaseFilter
	StoreID     *uint      `form:"storeId"`
	WarehouseID *uint      `form:"warehouseId"`
	StartDate   *time.Time `form:"startDate"`
	EndDate     *time.Time `form:"endDate"`

	Statuses []data.StockRequestStatus `form:"statuses[]"`
}
