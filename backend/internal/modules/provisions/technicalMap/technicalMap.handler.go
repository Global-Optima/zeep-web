package technicalMap

import (
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/technicalMap/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TechnicalMapHandler struct {
	technicalMapService TechnicalMapService
}

func NewTechnicalMapHandler(technicalMapService TechnicalMapService) *TechnicalMapHandler {
	return &TechnicalMapHandler{technicalMapService: technicalMapService}
}

func (h *TechnicalMapHandler) GetProvisionTechnicalMapByID(c *gin.Context) {
	ProvisionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	technicalMap, err := h.technicalMapService.GetProvisionTechnicalMapByID(ProvisionID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500TechnicalMapGet)
		return
	}

	utils.SendSuccessResponse(c, technicalMap)
}
