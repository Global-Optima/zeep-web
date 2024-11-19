package additives

import (
	"fmt"
	"strconv"
	"time"

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

func (h *AdditiveHandler) GetAdditivesByStoreAndProduct(c *gin.Context) {
	productSizeIDParam := c.Query("productSizeId")

	productSizeID, err := strconv.ParseUint(productSizeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product size ID")
		return
	}

	cacheKey := utils.BuildCacheKey("additives", map[string]string{
		"productSizeId": productSizeIDParam,
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedAdditives []types.AdditiveCategoryDTO
	if err := cacheUtil.Get(cacheKey, &cachedAdditives); err == nil {
		utils.SuccessResponse(c, cachedAdditives)
		return
	}

	additives, err := h.service.GetAdditivesByStoreAndProductSize(uint(productSizeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve additives")
		return
	}

	if err := cacheUtil.Set(cacheKey, additives, 10*time.Minute); err != nil {
		fmt.Printf("Failed to cache additives: %v\n", err)
	}

	utils.SuccessResponse(c, additives)
}
