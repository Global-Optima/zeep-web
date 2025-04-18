package ingredientCategories

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
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
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryCreate)
		return
	}

	action := types.CreateIngredientCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201IngredientCategory)
}

func (h *IngredientCategoryHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400IngredientCategory)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	category, err := h.service.GetTranslatedByID(locale, id)
	if err != nil {
		if errors.Is(err, types.ErrIngredientCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404IngredientCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryGet)
		return
	}

	utils.SendSuccessResponse(c, category)
}

func (h *IngredientCategoryHandler) Update(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400IngredientCategory)
		return
	}

	var dto types.UpdateIngredientCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400IngredientCategory)
		return
	}

	existingCategory, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, types.ErrIngredientCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404IngredientCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryUpdate)
		return
	}

	if err := h.service.Update(id, dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryUpdate)
		return
	}

	action := types.UpdateIngredientCategoryAuditFactory(
		&data.BaseDetails{
			ID:   existingCategory.ID,
			Name: existingCategory.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200IngredientCategoryUpdate)
}

func (h *IngredientCategoryHandler) Delete(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400IngredientCategory)
		return
	}

	existingCategory, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, types.ErrIngredientCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404IngredientCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryDelete)
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, types.ErrIngredientCategoryIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409IngredientCategoryDeleteInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryDelete)
		return
	}

	action := types.DeleteIngredientCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: existingCategory.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200IngredientCategoryDelete)
}

func (h *IngredientCategoryHandler) GetAll(c *gin.Context) {
	var filter types.IngredientCategoryFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.IngredientCategory{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400IngredientCategory)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	categories, err := h.service.GetAll(locale, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.GetPagination())
}

func (h *IngredientCategoryHandler) CreateOrUpdateIngredientCategoryTranslation(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400IngredientCategory)
		return
	}
	var dto types.IngredientCategoryTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}
	if err := h.service.UpsertIngredientCategoryTranslations(id, &dto); err != nil {
		if errors.Is(err, types.ErrIngredientCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404IngredientCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500IngredientCategoryTranslationsUpdate)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200IngredientCategoryTranslationsUpdate)
}
