package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StoreProvisionDTO struct {
	ID                  uint                      `json:"id"`
	Provision           types.ProvisionDTO        `json:"provision"`
	ExpirationInMinutes uint                      `json:"expirationInMinutes"`
	Volume              float64                   `json:"volume"`
	InitialVolume       float64                   `json:"initialVolume"`
	Status              data.StoreProvisionStatus `json:"status"`
	CompletedAt         *time.Time                `json:"completedAt,omitempty"`
	ExpiresAt           *time.Time                `json:"expiresAt,omitempty"`
	CreatedAt           time.Time                 `json:"createdAt"`
}

type StoreProvisionDetailsDTO struct {
	StoreProvisionDTO
	Ingredients []StoreProvisionIngredientDTO `json:"ingredients"`
}

type StoreProvisionIngredientDTO struct {
	Ingredient      ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity        float64                       `json:"quantity"`
	InitialQuantity float64                       `json:"initialQuantity"`
}

type CreateStoreProvisionDTO struct {
	ProvisionID         uint    `json:"provisionId" binding:"required"`
	Volume              float64 `json:"volume" binding:"required,gt=0"`
	ExpirationInMinutes uint    `json:"expirationInMinutes" binding:"required,gt=0"`
}

type UpdateStoreProvisionDTO struct {
	Volume              *float64 `json:"volume" binding:"omitempty,gt=0"`
	ExpirationInMinutes *uint    `json:"expirationInMinutes" binding:"omitempty,gt=0"`
}

type StoreProvisionFilterDTO struct {
	utils.BaseFilter
	Statuses       []data.StoreProvisionStatus `form:"statuses[]"`
	Search         *string                     `form:"search"`
	MinCompletedAt *time.Time                  `form:"minCompletedAt"`
	MaxCompletedAt *time.Time                  `form:"maxCompletedAt"`
	IsExpired      *bool                       `form:"isExpired"`
}
