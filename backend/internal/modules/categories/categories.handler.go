package categories

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	cacheKey := "categories:all"

	cacheUtil := utils.GetCacheInstance()

	var cachedCategories []types.CategoryDTO
	if err := cacheUtil.Get(cacheKey, &cachedCategories); err == nil {
		utils.SuccessResponse(c, cachedCategories)
		return
	}

	categories, err := h.service.GetCategories(c)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve categories")
		return
	}

	if err := cacheUtil.Set(cacheKey, categories, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache categories: %v\n", err)
	}

	utils.SuccessResponse(c, categories)
}
