package storeAdditives

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreAdditiveHandler struct {
	service           StoreAdditiveService
	additiveService   additives.AdditiveService
	franchiseeService franchisees.FranchiseeService
	auditService      audit.AuditService
	logger            *zap.SugaredLogger
}

func NewStoreAdditiveHandler(
	service StoreAdditiveService,
	additiveService additives.AdditiveService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	logger *zap.SugaredLogger,
) *StoreAdditiveHandler {
	return &StoreAdditiveHandler{
		service:           service,
		additiveService:   additiveService,
		franchiseeService: franchiseeService,
		auditService:      auditService,
		logger:            logger,
	}
}

func (h *StoreAdditiveHandler) GetStoreAdditiveCategories(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}
	productSizeID, err := strconv.ParseUint(c.Param("productSizeId"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid productSizeID")
		return
	}

	var filter types.StoreAdditiveCategoriesFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON+err.Error())
		return
	}

	additivesList, err := h.service.GetStoreAdditiveCategoriesByProductSize(storeID, uint(productSizeID), &filter)
	if err != nil {
		if errors.Is(err, types.ErrStoreAdditiveCategoriesNotFound) {
			utils.SendSuccessResponse(c, []types.StoreAdditiveCategoryDTO{})
			return
		}
		utils.SendInternalServerError(c, "Failed to retrieve store additives")
		return
	}

	utils.SendSuccessResponse(c, additivesList)
}

func (h *StoreAdditiveHandler) GetAdditivesListToAdd(c *gin.Context) {
	var filter additiveTypes.AdditiveFilterQuery

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Additive{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	storeAdditives, err := h.service.GetAdditivesListToAdd(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, storeAdditives, filter.Pagination)
}

func (h *StoreAdditiveHandler) GetStoreAdditives(c *gin.Context) {
	var filter additiveTypes.AdditiveFilterQuery

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreAdditive{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	storeAdditives, err := h.service.GetStoreAdditives(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch additives")
		return
	}

	utils.SendSuccessResponseWithPagination(c, storeAdditives, filter.Pagination)
}

func (h *StoreAdditiveHandler) CreateStoreAdditives(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var dtos []types.CreateStoreAdditiveDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		utils.SendBadRequestError(c, "Invalid input data")
		return
	}

	additiveIDs := make([]uint, len(dtos))
	for i, dto := range dtos {
		additiveIDs[i] = dto.AdditiveID
	}

	ids, err := h.service.CreateStoreAdditives(storeID, dtos)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create additive")
		return
	}

	stockList, err := h.service.GetStoreAdditivesByIDs(storeID, ids)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create additives: additive not found")
		return
	}

	dtoMap := make(map[uint]*types.CreateStoreAdditiveDTO)
	for _, dto := range dtos {
		dtoCopy := dto
		dtoMap[dto.AdditiveID] = &dtoCopy
	}

	var actions []shared.AuditAction

	for _, stock := range stockList {
		matchedDTO, exists := dtoMap[stock.AdditiveID]
		if !exists {
			h.logger.Errorf("Failed to match stock with DTO for stock ID: %d, AdditiveID: %d", stock.ID, stock.AdditiveID)
			continue
		}

		action := types.CreateStoreAdditiveAuditFactory(
			&data.BaseDetails{
				ID:   stock.ID,
				Name: stock.Name,
			},
			matchedDTO, storeID,
		)
		actions = append(actions, &action)
	}

	if len(actions) > 0 {
		_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Additive created successfully"})
}

func (h *StoreAdditiveHandler) UpdateStoreAdditive(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeAdditiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store id")
		return
	}

	storeAdditive, err := h.service.GetStoreAdditiveByID(storeID, uint(storeAdditiveID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update additive: additive not found")
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

	action := types.UpdateStoreAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeAdditiveID),
			Name: storeAdditive.Name,
		},
		&dto, storeID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive updated successfully"})
}

func (h *StoreAdditiveHandler) DeleteStoreAdditive(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeAdditiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store additive ID")
		return
	}

	storeAdditive, err := h.service.GetStoreAdditiveByID(storeID, uint(storeAdditiveID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update additive: additive not found")
		return
	}

	if err := h.service.DeleteStoreAdditive(storeID, uint(storeAdditiveID)); err != nil {
		utils.SendInternalServerError(c, "Failed to delete additive")
		return
	}

	action := types.DeleteStoreAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeAdditiveID),
			Name: storeAdditive.Name,
		},
		struct{}{}, storeID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "Additive deleted successfully"})
}

func (h *StoreAdditiveHandler) GetStoreAdditiveByID(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
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
