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
	storeIDParam := c.Query("storeId")
	productIDParam := c.Query("productId")

	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	additives, err := h.service.GetAdditivesByStoreAndProduct(uint(storeID), uint(productID))
	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve additives")
		return
	}

	utils.SuccessResponse(c, additives)

}
