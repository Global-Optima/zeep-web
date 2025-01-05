package storeAdditives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreAdditiveHandler struct {
	service StoreAdditiveService
}

func NewStoreAdditiveHandler(service StoreAdditiveService) *StoreAdditiveHandler {
	return &StoreAdditiveHandler{service: service}
}

func (h *StoreAdditiveHandler) GetStoreAdditiveCategories(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var filter additiveTypes.AdditiveCategoriesFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.AdditiveCategory{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	additives, err := h.service.GetStoreAdditiveCategories(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve store additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *StoreAdditiveHandler) GetStoreAdditives(c *gin.Context) {
	var filter additiveTypes.AdditiveFilterQuery

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreAdditive{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	additives, err := h.service.GetStoreAdditives(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

func (h *StoreAdditiveHandler) CreateStoreAdditives(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var dtos []types.CreateStoreAdditiveDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	if _, err := h.service.CreateStoreAdditives(storeID, dtos); err != nil {
		utils.SendInternalServerError(c, "Failed to create additive")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive created successfully"})
}

func (h *StoreAdditiveHandler) UpdateStoreAdditive(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeAdditiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store id")
		return
	}

	var dto types.UpdateStoreAdditiveDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	if err := h.service.UpdateStoreAdditive(storeID, uint(storeAdditiveID), &dto); err != nil {
		utils.SendInternalServerError(c, "Failed to update additive")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive updated successfully"})
}

func (h *StoreAdditiveHandler) DeleteStoreAdditive(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeAdditiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store additive ID")
		return
	}

	if err := h.service.DeleteStoreAdditive(storeID, uint(storeAdditiveID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete additive")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive deleted successfully"})
}

func (h *StoreAdditiveHandler) GetStoreAdditiveByID(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store additive ID")
		return
	}

	additive, err := h.service.GetStoreAdditiveByID(storeID, uint(additiveID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch store additive")
		return
	}

	utils.SendSuccessResponse(c, additive)
}
