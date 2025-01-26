package ingredients

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	service IngredientService
}

func NewIngredientHandler(service IngredientService) *IngredientHandler {
	return &IngredientHandler{service: service}
}

func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var dto types.CreateIngredientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if err := h.service.CreateIngredient(&dto); err != nil {
		utils.SendInternalServerError(c, "Failed to create ingredient")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Ingredient created successfully"})
}

func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ingredient ID")
		return
	}

	var dto types.UpdateIngredientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if err := h.service.UpdateIngredient(uint(ingredientID), &dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update ingredient")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Ingredient updated successfully"})
}

func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ingredient ID")
		return
	}

	if err := h.service.DeleteIngredient(uint(ingredientID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete ingredient")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Ingredient deleted successfully"})
}

func (h *IngredientHandler) GetIngredientByID(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ingredient ID")
		return
	}

	ingredient, err := h.service.GetIngredientByID(uint(ingredientID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch ingredient")
		return
	}

	utils.SendSuccessResponse(c, ingredient)
}

func (h *IngredientHandler) GetIngredients(c *gin.Context) {
	var filter types.IngredientFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Ingredient{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	ingredients, err := h.service.GetIngredients(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch ingredients")
		return
	}

	utils.SendSuccessResponseWithPagination(c, ingredients, filter.Pagination)
}
