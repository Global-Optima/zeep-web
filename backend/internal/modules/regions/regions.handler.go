package regions

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}
	id, err := h.service.CreateRegion(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionCreate)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201Region)
}

func (h *RegionHandler) UpdateRegion(c *gin.Context) {
	regionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.UpdateRegionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	region, err := h.service.GetRegionByID(uint(regionID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response404Region)
		return
	}

	if err := h.service.UpdateRegion(uint(regionID), &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionUpdate)
		return
	}

	action := types.UpdateRegionAuditFactory(
		&data.BaseDetails{
			ID:   uint(regionID),
			Name: region.Name,
		},
		&dto,
	)
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200RegionUpdate)
}

func (h *RegionHandler) DeleteRegion(c *gin.Context) {
	regionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	region, err := h.service.GetRegionByID(uint(regionID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response404Region)
		return
	}

	if err := h.service.DeleteRegion(uint(regionID)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionDelete)
		return
	}

	action := types.DeleteRegionAuditFactory(
		&data.BaseDetails{
			ID:   uint(regionID),
			Name: region.Name,
		},
	)
	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200RegionDelete)
}

func (h *RegionHandler) GetRegionByID(c *gin.Context) {
	regionID, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	region, err := h.service.GetRegionByID(uint(regionID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionGet)
		return
	}
	utils.SendSuccessResponse(c, region)
}

func (h *RegionHandler) GetRegions(c *gin.Context) {
	var filter types.RegionFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Region{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	regions, err := h.service.GetRegions(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionGet)
		return
	}
	utils.SendSuccessResponseWithPagination(c, regions, filter.Pagination)
}

func (h *RegionHandler) GetAllRegions(c *gin.Context) {
	var filter types.RegionFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Region{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	allRegions, err := h.service.GetAllRegions(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionGet)
		return
	}
	utils.SendSuccessResponse(c, allRegions)
}
