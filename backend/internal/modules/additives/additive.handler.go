package additives

import (
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdditiveHandler struct {
	service AdditiveService
}

func NewAdditiveHandler(service AdditiveService) *AdditiveHandler {
	return &AdditiveHandler{service: service}
}

func (h *AdditiveHandler) GetAdditivesByStoreAndProduct(c *gin.Context) {
	productSizeIDParam := c.Query("productSizeId")

	productSizeID, err := strconv.ParseUint(productSizeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product size ID")
		return
	}

	additives, err := h.service.GetAdditivesByStoreAndProductSize(uint(productSizeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve additives")
		return
	}

	utils.SuccessResponse(c, additives)

}
