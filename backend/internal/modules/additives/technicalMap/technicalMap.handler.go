package technicalMap

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/technicalMap/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TechnicalMapHandler struct {
	technicalMapService TechnicalMapService
}

func NewTechnicalMapHandler(technicalMapService TechnicalMapService) *TechnicalMapHandler {
	return &TechnicalMapHandler{technicalMapService: technicalMapService}
}

func (h *TechnicalMapHandler) GetAdditiveTechnicalMapByID(c *gin.Context) {
	additiveID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	technicalMap, err := h.technicalMapService.GetAdditiveTechnicalMapByID(additiveID)
	if err != nil {
		if errors.Is(err, types.ErrorTechnicalMapNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404TechnicalMap)
			return
		}

		localization.SendLocalizedResponseWithKey(c, types.Response500TechnicalMapGet)
		return
	}

	utils.SendSuccessResponse(c, technicalMap)
}
