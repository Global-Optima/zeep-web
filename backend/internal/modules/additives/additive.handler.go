package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
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

	additives, err := h.service.GetAdditiveCategories(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *AdditiveHandler) CreateAdditiveCategory(c *gin.Context) {
	var dto types.CreateAdditiveCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateAdditiveCategory(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
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
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
		return
	}

	if err := h.service.UpdateAdditiveCategory(uint(categoryID), &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
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
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
		return
	}

	if err := h.service.DeleteAdditiveCategory(uint(categoryID)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
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
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategory)
		return
	}

	utils.SendSuccessResponse(c, category)
}

func (h *AdditiveHandler) GetAdditives(c *gin.Context) {
	var filter types.AdditiveFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Additive{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	additives, err := h.service.GetAdditives(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *AdditiveHandler) CreateAdditive(c *gin.Context) {
	var dto types.CreateAdditiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateAdditive(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
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
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	var dto types.UpdateAdditiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
		return
	}

	if err := h.service.UpdateAdditive(uint(additiveID), &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
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
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
		return
	}

	if err := h.service.DeleteAdditive(uint(additiveID)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
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
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Additive)
		return
	}

	utils.SendSuccessResponse(c, additive)
}
