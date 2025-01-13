package categories

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
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

	var filter types.CategoriesFilterDTO

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.ProductCategory{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	categories, err := h.service.GetCategories(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve categories")
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.Pagination)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve category")
		return
	}

	utils.SendSuccessResponse(c, category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var dto types.CreateCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateCategory(&dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create category")
		return
	}

	utils.SendMessageWithStatus(c, "category created successfully", http.StatusCreated)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	var dto types.UpdateCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateCategory(uint(id), &dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update category")
		return
	}

	utils.SendMessageWithStatus(c, "category updated successfully", http.StatusOK)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete category")
		return
	}

	utils.SendMessageWithStatus(c, "Category deleted successfully", http.StatusOK)
}
