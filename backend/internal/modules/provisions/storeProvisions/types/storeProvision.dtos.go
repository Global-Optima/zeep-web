package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"time"
)

type StoreProvisionDTO struct {
	ID uint `json:"id"`
	types.BaseProvisionDTO
	ProvisionID       uint                      `json:"provisionId"`
	ExpirationInHours int                       `json:"expirationInHours"`
	Volume            float64                   `json:"volume"`
	Status            data.StoreProvisionStatus `json:"status"`
	CompletedAt       *time.Time                `json:"completedAt,omitempty"`
	ExpiresAt         *time.Time                `json:"expiresAt"`
	CreatedAt         time.Time                 `json:"createdAt"`
}

type StoreProvisionDetailsDTO struct {
	StoreProvisionDTO
	Ingredients []StoreProvisionIngredientDTO `json:"ingredients"`
}

type StoreProvisionIngredientDTO struct {
	Ingredient ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity   float64                       `json:"quantity"`
}

type CreateStoreProvisionDTO struct {
	ProvisionID       uint    `json:"provisionId" binding:"required"`
	Volume            float64 `json:"volume" binding:"required,gt=0"`
	ExpirationInHours int     `json:"expirationInHours" binding:"required,gt=0"`
}

type UpdateStoreProvisionDTO struct {
	Volume            *float64 `json:"volume" binding:"omitempty,gt=0"`
	ExpirationInHours *int     `json:"expirationInHours" binding:"omitempty,gt=0"`
}

type StoreProvisionFilterDTO struct {
	utils.BaseFilter
	Search         *string    `form:"search"`
	MinCompletedAt *time.Time `form:"minCompletedAt"`
	MaxCompletedAt *time.Time `form:"maxCompletedAt"`
	IsExpired      *bool      `form:"isExpired"`
	Status         *string    `form:"status"`
}
