package ingredientCategories

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type IngredientCategoryHandler struct {
	service      IngredientCategoryService
	auditService audit.AuditService
}

func NewIngredientCategoryHandler(service IngredientCategoryService, auditService audit.AuditService) *IngredientCategoryHandler {
	return &IngredientCategoryHandler{
		service:      service,
		auditService: auditService,
	}
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

	action := types.CreateIngredientCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessCreatedResponse(c, fmt.Sprintf("id: %d", id))
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
		utils.SendBadRequestError(c, "failed to update ingredient category: category not found")
		return
	}

	existingCategory, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to update ingredient category")
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	action := types.UpdateIngredientCategoryAuditFactory(
		&data.BaseDetails{
			ID:   existingCategory.ID,
			Name: existingCategory.Name,
		},
		&dto,
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Ingredient category updated successfully"})
}

func (h *IngredientCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	existingCategory, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to update ingredient category")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	action := types.DeleteIngredientCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: existingCategory.Name,
		},
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Ingredient category deleted successfully"})
}

func (h *IngredientCategoryHandler) GetAll(c *gin.Context) {
	var filter types.IngredientCategoryFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.IngredientCategory{}); err != nil {
		utils.SendBadRequestError(c, "Failed to parse pagination or filtering parameters")
		return
	}

	categories, err := h.service.GetAll(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch ingredient categories")
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.GetPagination())
}
