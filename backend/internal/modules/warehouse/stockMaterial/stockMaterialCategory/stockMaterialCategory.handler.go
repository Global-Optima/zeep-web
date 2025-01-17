package stockMaterialCategory

import (
	"strconv"

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
		utils.SendBadRequestError(c, err.Error())
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessCreatedResponse(c, gin.H{"id": id})
}

func (h *StockMaterialCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	response, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, response)
}

func (h *StockMaterialCategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, categories)
}

func (h *StockMaterialCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	var dto types.UpdateStockMaterialCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		utils.SendInternalServerError(c, err.Error())
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
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(204) // No Content
}
