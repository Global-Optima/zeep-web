package ingredients

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	service      IngredientService
	auditService audit.AuditService
}

func NewIngredientHandler(service IngredientService, auditService audit.AuditService) *IngredientHandler {
	return &IngredientHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var dto types.CreateIngredientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	id, err := h.service.CreateIngredient(&dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create ingredient")
		return
	}

	action := types.CreateIngredientAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

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

	existingIngredient, err := h.service.GetIngredientByID(uint(ingredientID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to get update ingredient: ingredient not found")
		return
	}

	if err := h.service.UpdateIngredient(uint(ingredientID), &dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update ingredient")
		return
	}

	action := types.UpdateIngredientAuditFactory(
		&data.BaseDetails{
			ID:   uint(ingredientID),
			Name: existingIngredient.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "Ingredient updated successfully"})
}

func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ingredient ID")
		return
	}

	existing, err := h.service.GetIngredientByID(uint(ingredientID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete ingredient: ingredient not found")
		return
	}

	if err := h.service.DeleteIngredient(uint(ingredientID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete ingredient")
		return
	}

	action := types.DeleteIngredientAuditFactory(
		&data.BaseDetails{
			ID:   uint(ingredientID),
			Name: existing.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

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
