package storeInventoryManagers

import (
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreInventoryManagerHandler struct {
	service StoreInventoryManagerService
}

func NewStoreInventoryManagerHandler(service StoreInventoryManagerService) *StoreInventoryManagerHandler {
	return &StoreInventoryManagerHandler{
		service: service,
	}
}

func (h *StoreInventoryManagerHandler) GetFrozenInventory(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	frozenInventory, err := h.service.CalculateFrozenInventory(storeID, nil)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreInventoryManagerGetFrozenInventory)
		return
	}

	utils.SendSuccessResponse(c, frozenInventory)
}
