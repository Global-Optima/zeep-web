package additives

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AdditiveHandler struct {
	service      AdditiveService
	auditService audit.AuditService
}

func NewAdditiveHandler(service AdditiveService, auditService audit.AuditService) *AdditiveHandler {
	return &AdditiveHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *AdditiveHandler) GetAdditiveCategories(c *gin.Context) {
	var filter types.AdditiveCategoriesFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.AdditiveCategory{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	additives, err := h.service.GetAdditiveCategories(locale, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *AdditiveHandler) CreateAdditiveCategory(c *gin.Context) {
	var dto types.CreateAdditiveCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateAdditiveCategory(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryCreate)
		return
	}

	action := types.CreateAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201AdditiveCategory)
}

func (h *AdditiveHandler) UpdateAdditiveCategory(c *gin.Context) {
	var dto types.UpdateAdditiveCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	categoryID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(categoryID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryUpdate)
		return
	}

	if err := h.service.UpdateAdditiveCategory(categoryID, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryUpdate)
		return
	}

	action := types.UpdateAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(categoryID),
			Name: category.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveCategoryUpdate)
}

func (h *AdditiveHandler) DeleteAdditiveCategory(c *gin.Context) {
	categoryID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(categoryID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryDelete)
		return
	}

	if err := h.service.DeleteAdditiveCategory(uint(categoryID)); err != nil {
		if errors.Is(err, types.ErrAdditiveCategoryIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409AdditiveCategoryDeleteInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryDelete)
		return
	}

	action := types.DeleteAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(categoryID),
			Name: category.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveCategoryDelete)
}

func (h *AdditiveHandler) GetAdditiveCategoryByID(c *gin.Context) {
	categoryID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	category, err := h.service.GetTranlsatedAdditiveCategoryByID(locale, categoryID)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrAdditiveCategoryNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404AdditiveCategory)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryGet)
			return
		}
	}

	utils.SendSuccessResponse(c, category)
}

func (h *AdditiveHandler) GetAdditives(c *gin.Context) {
	var filter types.AdditiveFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Additive{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	additives, err := h.service.GetAdditives(locale, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *AdditiveHandler) CreateAdditive(c *gin.Context) {
	var dto types.CreateAdditiveDTO
	var err error

	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := ingredientTypes.ParseJSONIngredientsFromString(c.PostForm(types.INGREDIENTS_FORM_DATA_KEY), &dto.Ingredients); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := provisionsTypes.ParseJSONProvisionsFromString(c.PostForm(types.PROVISIONS_FORM_DATA_KEY), &dto.Provisions); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	dto.Image, err = media.GetImageWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	id, err := h.service.CreateAdditive(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCreate)
		return
	}

	action := types.CreateAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Additive)
}

func (h *AdditiveHandler) UpdateAdditive(c *gin.Context) {
	additiveID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	var dto types.UpdateAdditiveDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := ingredientTypes.ParseJSONIngredientsFromString(c.PostForm(types.INGREDIENTS_FORM_DATA_KEY), &dto.Ingredients); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := provisionsTypes.ParseJSONProvisionsFromString(c.PostForm(types.PROVISIONS_FORM_DATA_KEY), &dto.Provisions); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	dto.Image, err = media.GetImageWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	additive, err := h.service.UpdateAdditive(additiveID, &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveUpdate)
		return
	}

	action := types.UpdateAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(additiveID),
			Name: additive.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveUpdate)
}

func (h *AdditiveHandler) DeleteAdditive(c *gin.Context) {
	additiveID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	additive, err := h.service.GetAdditiveByID(additiveID)
	if err != nil {
		if errors.Is(err, types.ErrAdditiveNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Additive)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveDelete)
		return
	}

	if err := h.service.DeleteAdditive(uint(additiveID)); err != nil {
		if errors.Is(err, types.ErrAdditiveInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409AdditiveDeleteInUse)
			return
		}

		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveDelete)
		return
	}

	action := types.DeleteAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(additiveID),
			Name: additive.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveDelete)
}

func (h *AdditiveHandler) GetAdditiveByID(c *gin.Context) {
	additiveID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	locale := contexts.GetLocaleFromCtx(c)

	additive, err := h.service.GetTranslatedAdditiveByID(locale, additiveID)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrAdditiveNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404Additive)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveGet)
			return
		}
	}

	utils.SendSuccessResponse(c, additive)
}

func (h *AdditiveHandler) CreateOrUpdateAdditiveTranslation(c *gin.Context) {
	additiveID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	var dto types.AdditiveTranslationsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.UpsertAdditiveTranslations(additiveID, &dto); err != nil {
		if errors.Is(err, types.ErrAdditiveNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Additive)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveTranslationsCreate)
		return
	}

	utils.SendSuccessCreatedResponse(c, "successfully created translations")
}

func (h *AdditiveHandler) CreateOrUpdateAdditiveCategoryTranslation(c *gin.Context) {
	categoryID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	var dto types.AdditiveCategoryTranslationsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := h.service.UpsertAdditiveCategoryTranslations(categoryID, &dto); err != nil {
		if errors.Is(err, types.ErrAdditiveCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404AdditiveCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveTranslationsCreate)
		return
	}

	utils.SendSuccessCreatedResponse(c, "successfully created translations")
}
