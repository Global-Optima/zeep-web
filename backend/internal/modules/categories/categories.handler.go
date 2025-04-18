package categories

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service      CategoryService
	auditService audit.AuditService
}

func NewCategoryHandler(service CategoryService, auditService audit.AuditService) *CategoryHandler {
	return &CategoryHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	var filter types.ProductCategoriesFilterDTO

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.ProductCategory{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	categories, err := h.service.GetCategories(locale, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.Pagination)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	category, err := h.service.GetTranslatedCategoryByID(locale, id)
	if err != nil {
		if errors.Is(err, types.ErrCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404ProductCategory)
			return
		}

		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	utils.SendSuccessResponse(c, category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var dto types.CreateProductCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	id, err := h.service.CreateCategory(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	action := types.CreateProductCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201ProductCategory)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	var dto types.UpdateProductCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	existingCategory, err := h.service.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, types.ErrCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404ProductCategory)
			return
		}

		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	err = h.service.UpdateCategory(id, &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	action := types.UpdateProductCategoryAuditFactory(
		&data.BaseDetails{
			ID:   existingCategory.ID,
			Name: existingCategory.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductCategoryUpdate)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	existingCategory, err := h.service.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, types.ErrCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404ProductCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		if errors.Is(err, types.ErrCategoryIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409ProductCategoryDeleteInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	action := types.DeleteProductCategoryAuditFactory(
		&data.BaseDetails{
			ID:   existingCategory.ID,
			Name: existingCategory.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductCategoryDelete)
}

func (h *CategoryHandler) CreateOrUpdateCategoryTranslation(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}
	var dto types.ProductCategoryTranslationsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}
	if err := h.service.UpsertProductCategoryTranslations(id, &dto); err != nil {
		if errors.Is(err, types.ErrCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404ProductCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategoryTranslationsUpdate)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductCategoryTranslationsUpdate)
}
