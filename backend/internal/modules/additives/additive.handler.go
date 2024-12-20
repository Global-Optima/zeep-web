package additives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdditiveHandler struct {
	service AdditiveService
}

func NewAdditiveHandler(service AdditiveService) *AdditiveHandler {
	return &AdditiveHandler{service: service}
}

func (h *AdditiveHandler) GetAdditiveCategories(c *gin.Context) {
	var filter types.AdditiveCategoriesFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	additives, err := h.service.GetAdditiveCategories(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve additives")
		return
	}

	utils.SendSuccessResponse(c, additives)
}

func (h *AdditiveHandler) CreateAdditiveCategory(c *gin.Context) {
	var dto types.CreateAdditiveCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	if err := h.service.CreateAdditiveCategory(&dto); err != nil {
		utils.SendInternalServerError(c, "Failed to create additive category")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive category created successfully"})
}

func (h *AdditiveHandler) UpdateAdditiveCategory(c *gin.Context) {
	var dto types.UpdateAdditiveCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	if err := h.service.UpdateAdditiveCategory(&dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update additive category")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive category updated successfully"})
}

func (h *AdditiveHandler) DeleteAdditiveCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID")
		return
	}

	if err := h.service.DeleteAdditiveCategory(uint(categoryID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete additive category")
		return
	}

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

	additives, err := h.service.GetAdditives(filter)
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

	if err := h.service.CreateAdditive(&dto); err != nil {
		utils.SendInternalServerError(c, "Failed to create additive")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive created successfully"})
}

func (h *AdditiveHandler) UpdateAdditive(c *gin.Context) {
	var dto types.UpdateAdditiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	if err := h.service.UpdateAdditive(&dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update additive")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive updated successfully"})
}

func (h *AdditiveHandler) DeleteAdditive(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid additive ID")
		return
	}

	if err := h.service.DeleteAdditive(uint(additiveID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete additive")
		return
	}

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
