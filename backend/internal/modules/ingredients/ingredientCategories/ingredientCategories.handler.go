package ingredientCategories

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type IngredientCategoryHandler struct {
	service IngredientCategoryService
}

func NewIngredientCategoryHandler(service IngredientCategoryService) *IngredientCategoryHandler {
	return &IngredientCategoryHandler{service: service}
}

func (h *IngredientCategoryHandler) Create(c *gin.Context) {
	var dto types.CreateIngredientCategoryDTO
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

func (h *IngredientCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	category, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, category)
}

func (h *IngredientCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	var dto types.UpdateIngredientCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *IngredientCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *IngredientCategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, categories)
}
