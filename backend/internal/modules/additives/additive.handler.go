package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"strconv"

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
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	additives, err := h.service.GetAdditiveCategories(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *AdditiveHandler) CreateAdditiveCategory(c *gin.Context) {
	var dto types.CreateAdditiveCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}
	id, err := h.service.CreateAdditiveCategory(&dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create additive category")
		return
	}

	action := types.CreateAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive category created successfully"})
}

func (h *AdditiveHandler) UpdateAdditiveCategory(c *gin.Context) {
	var dto types.UpdateAdditiveCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update additive category: category not found")
		return
	}

	if err := h.service.UpdateAdditiveCategory(uint(categoryID), &dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update additive category")
		return
	}

	action := types.UpdateAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(categoryID),
			Name: category.Name,
		},
		&dto,
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive category updated successfully"})
}

func (h *AdditiveHandler) DeleteAdditiveCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update additive category: category not found")
		return
	}

	if err := h.service.DeleteAdditiveCategory(uint(categoryID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete additive category")
		return
	}

	action := types.DeleteAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(categoryID),
			Name: category.Name,
		},
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive category deleted successfully"})
}

func (h *AdditiveHandler) GetAdditiveCategoryByID(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additive category")
		return
	}

	utils.SendSuccessResponse(c, category)
}

func (h *AdditiveHandler) GetAdditives(c *gin.Context) {
	var filter types.AdditiveFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Additive{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	filter.Pagination = utils.ParsePagination(c)

	additives, err := h.service.GetAdditives(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *AdditiveHandler) CreateAdditive(c *gin.Context) {
	var dto types.CreateAdditiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	id, err := h.service.CreateAdditive(&dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create additive")
		return
	}

	action := types.CreateAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive created successfully"})
}

func (h *AdditiveHandler) UpdateAdditive(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid additive ID")
		return
	}

	var dto types.UpdateAdditiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update additive: additive not found")
		return
	}

	if err := h.service.UpdateAdditive(uint(additiveID), &dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update additive")
		return
	}

	action := types.UpdateAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(additiveID),
			Name: additive.Name,
		},
		&dto,
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive updated successfully"})
}

func (h *AdditiveHandler) DeleteAdditive(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid additive ID")
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update additive: additive not found")
		return
	}

	if err := h.service.DeleteAdditive(uint(additiveID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete additive")
		return
	}

	action := types.DeleteAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(additiveID),
			Name: additive.Name,
		},
	)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive deleted successfully"})
}

func (h *AdditiveHandler) GetAdditiveByID(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid additive ID")
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additive")
		return
	}

	utils.SendSuccessResponse(c, additive)
}
