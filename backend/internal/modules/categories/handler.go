package categories

import (
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
	categories, err := h.service.GetCategories()

	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve categories")
		return
	}

	utils.SuccessResponse(c, categories)
}
