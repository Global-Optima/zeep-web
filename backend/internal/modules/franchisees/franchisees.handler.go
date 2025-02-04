package franchisees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FranchiseeHandler struct {
	service      FranchiseeService
	auditService audit.AuditService
}

func NewFranchiseeHandler(service FranchiseeService, auditService audit.AuditService) *FranchiseeHandler {
	return &FranchiseeHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *FranchiseeHandler) CreateFranchisee(c *gin.Context) {
	var input types.CreateFranchiseeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}
	id, err := h.service.CreateFranchisee(&input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create franchisee")
		return
	}

	action := types.CreateFranchiseeAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: input.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "franchisee was created successfully")
}

func (h *FranchiseeHandler) UpdateFranchisee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	var input types.UpdateFranchiseeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	franchisee, err := h.service.GetFranchiseeByID(uint(id))
	if err != nil {
		utils.SendNotFoundError(c, "failed to update franchisee: franchisee not found")
		return
	}

	if err := h.service.UpdateFranchisee(uint(id), &input); err != nil {
		utils.SendInternalServerError(c, "failed to update franchisee")
		return
	}

	action := types.UpdateFranchiseeAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: franchisee.Name,
		},
		&input,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "franchisee updated successfully"})
}

func (h *FranchiseeHandler) DeleteFranchisee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	franchisee, err := h.service.GetFranchiseeByID(uint(id))
	if err != nil {
		utils.SendNotFoundError(c, "failed to delete franchisee: franchisee not found")
		return
	}

	if err := h.service.DeleteFranchisee(uint(id)); err != nil {
		utils.SendInternalServerError(c, "failed to delete franchisee")
		return
	}

	action := types.DeleteFranchiseeAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: franchisee.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	if err := h.service.DeleteFranchisee(uint(id)); err != nil {
		utils.SendInternalServerError(c, "failed to delete franchisee")
		return
	}
	utils.SendSuccessResponse(c, gin.H{"message": "franchisee deleted successfully"})
}

func (h *FranchiseeHandler) GetFranchiseeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	franchisee, err := h.service.GetFranchiseeByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee")
		return
	}
	utils.SendSuccessResponse(c, franchisee)
}

func (h *FranchiseeHandler) GetMyFranchisee(c *gin.Context) {
	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
	}

	franchisee, err := h.service.GetFranchiseeByID(franchiseeID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee")
		return
	}
	utils.SendSuccessResponse(c, franchisee)
}

func (h *FranchiseeHandler) GetFranchisees(c *gin.Context) {
	var filter types.FranchiseeFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Franchisee{}); err != nil {
		utils.SendBadRequestError(c, "invalid filter parameters")
		return
	}

	franchisees, err := h.service.GetFranchisees(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisees")
		return
	}
	utils.SendSuccessResponseWithPagination(c, franchisees, filter.Pagination)
}

func (h *FranchiseeHandler) GetAllFranchisees(c *gin.Context) {
	var filter types.FranchiseeFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{})

	warehouses, err := h.service.GetAllFranchisees(&filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, warehouses)
}
