package ingredients

import (
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateIngredient(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCreate)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201IngredientCreate)
}

func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Ingredient)
		return
	}

	var dto types.UpdateIngredientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	existingIngredient, err := h.service.GetIngredientByID(uint(ingredientID))
	if err != nil {
		if errors.Is(err, types.ErrIngredientNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Ingredient)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientUpdate)
		return
	}

	if err := h.service.UpdateIngredient(uint(ingredientID), &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientUpdate)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200IngredientUpdate)
}

func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Ingredient)
		return
	}

	existing, err := h.service.GetIngredientByID(uint(ingredientID))
	if err != nil {
		if errors.Is(err, types.ErrIngredientNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Ingredient)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response404Ingredient)
		return
	}

	if err := h.service.DeleteIngredient(uint(ingredientID)); err != nil {
		if errors.Is(err, types.ErrIngredientIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409IngredientInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientDelete)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200IngredientDelete)
}

func (h *IngredientHandler) GetIngredientByID(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Ingredient)
		return
	}

	ingredient, err := h.service.GetIngredientByID(uint(ingredientID))
	if err != nil {
		if errors.Is(err, types.ErrIngredientNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Ingredient)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientDelete)
		return
	}

	utils.SendSuccessResponse(c, ingredient)
}

func (h *IngredientHandler) GetIngredients(c *gin.Context) {
	var filter types.IngredientFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Ingredient{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	ingredients, err := h.service.GetIngredients(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientDelete)
		return
	}

	utils.SendSuccessResponseWithPagination(c, ingredients, filter.Pagination)
}
