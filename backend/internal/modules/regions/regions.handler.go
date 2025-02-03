package regions

import (
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	service      RegionService
	auditService audit.AuditService
}

func NewRegionHandler(service RegionService, auditService audit.AuditService) *RegionHandler {
	return &RegionHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *RegionHandler) CreateRegion(c *gin.Context) {
	var dto types.CreateRegionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}
	id, err := h.service.CreateRegion(&dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create region")
		return
	}

	action := types.CreateRegionAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "region created successfully")
}

func (h *RegionHandler) UpdateRegion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	var dto types.UpdateRegionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	region, err := h.service.GetRegionByID(uint(id))
	if err != nil {
		utils.SendNotFoundError(c, "failed to update region: region not found")
		return
	}

	if err := h.service.UpdateRegion(uint(id), &dto); err != nil {
		utils.SendInternalServerError(c, "failed to update region")
		return
	}

	action := types.UpdateRegionAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: region.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "region updated successfully"})
}

func (h *RegionHandler) DeleteRegion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	region, err := h.service.GetRegionByID(uint(id))
	if err != nil {
		utils.SendNotFoundError(c, "failed to delete region: region not found")
		return
	}

	if err := h.service.DeleteRegion(uint(id)); err != nil {
		utils.SendInternalServerError(c, "failed to delete region")
		return
	}

	action := types.DeleteRegionAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: region.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "region deleted successfully"})
}

func (h *RegionHandler) GetRegionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	region, err := h.service.GetRegionByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve region")
		return
	}
	utils.SendSuccessResponse(c, region)
}

func (h *RegionHandler) GetRegions(c *gin.Context) {
	var filter types.RegionFilter
	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{})

	if err != nil {
		utils.SendBadRequestError(c, "invalid filter parameters")
		return
	}

	regions, err := h.service.GetRegions(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve regions")
		return
	}
	utils.SendSuccessResponseWithPagination(c, regions, filter.Pagination)
}

func (h *RegionHandler) GetAllRegions(c *gin.Context) {
	var filter types.RegionFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Warehouse{})
	if err != nil {
		utils.SendBadRequestError(c, "invalid filter parameters")
		return
	}

	warehouses, err := h.service.GetAllRegions(&filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, warehouses)
}
