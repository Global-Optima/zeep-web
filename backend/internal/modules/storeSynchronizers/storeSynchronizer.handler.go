package storeSynchronizers

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreSynchronizerHandler struct {
	service StoreSynchronizeService
}

func NewStoreSynchronizeHandler(manager StoreSynchronizeService) *StoreSynchronizerHandler {
	return &StoreSynchronizerHandler{
		service: manager,
	}
}

func (h *StoreSynchronizerHandler) SynchronizeStore(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := h.service.SynchronizeStoreInventory(storeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreSynchronization)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreSynchronization)
}

func (h *StoreSynchronizerHandler) IsSynchronizedStore(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	syncStatus, err := h.service.GetSynchronizationStatus(storeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreSynchronizationCheck)
		return
	}

	utils.SendResponseWithStatus(c, syncStatus, http.StatusOK)
}
