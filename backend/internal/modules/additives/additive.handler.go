package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdditiveHandler struct {
	service AdditiveService
}

func NewAdditiveHandler(service AdditiveService) *AdditiveHandler {
	return &AdditiveHandler{service: service}
}

func (h *AdditiveHandler) GetAdditiveCategories(c *gin.Context) {
	var filter types.AdditiveCategoriesFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	additives, err := h.service.GetAdditiveCategories(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve additives")
		return
	}

	utils.SendSuccessResponse(c, additives)
}

func (h *AdditiveHandler) GetAdditives(c *gin.Context) {
	var filter types.AdditiveFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	filter.Pagination = utils.ParsePagination(c)

	additives, err := h.service.GetAdditives(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}
