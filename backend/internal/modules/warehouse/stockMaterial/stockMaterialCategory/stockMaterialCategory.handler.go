package stockMaterialCategory

import (
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StockMaterialCategoryHandler struct {
	service StockMaterialCategoryService
}

func NewStockMaterialCategoryHandler(service StockMaterialCategoryService) *StockMaterialCategoryHandler {
	return &StockMaterialCategoryHandler{service: service}
}

func (h *StockMaterialCategoryHandler) Create(c *gin.Context) {
	var dto types.CreateStockMaterialCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid request body")
		return
	}

	_, err := h.service.Create(dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create stock material category")
		return
	}

	utils.SendSuccessCreatedResponse(c, "stock material created successfully")
}

func (h *StockMaterialCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	response, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch stock material category")
		return
	}

	utils.SendSuccessResponse(c, response)
}

func (h *StockMaterialCategoryHandler) GetAll(c *gin.Context) {
	var filter types.StockMaterialCategoryFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockMaterialCategory{})
	if err != nil {
		utils.SendBadRequestError(c, "Failed to parse queries")
		return
	}

	categories, err := h.service.GetAll(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch categories")
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.Pagination)
}

func (h *StockMaterialCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	var dto types.UpdateStockMaterialCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid body")
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update stock material category")
		return
	}

	c.Status(204) // No Content
}

func (h *StockMaterialCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete stock material category")
		return
	}

	c.Status(204) // No Content
}
