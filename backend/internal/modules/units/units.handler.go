package units

import (
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	service      UnitService
	auditService audit.AuditService
}

func NewUnitHandler(service UnitService, auditService audit.AuditService) *UnitHandler {
	return &UnitHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var dto types.CreateUnitDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitCreate)
		return
	}

	action := types.CreateUnitAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Unit)
}

func (h *UnitHandler) GetAllUnits(c *gin.Context) {
	var filter types.UnitFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Unit{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Unit)
		return
	}

	units, err := h.service.GetAll(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, units, filter.GetPagination())
}

func (h *UnitHandler) GetUnitByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Unit)
		return
	}

	unit, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrUnitNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Unit)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitGet)
		return
	}

	utils.SendSuccessResponse(c, unit)
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Unit)
		return
	}

	var dto types.UpdateUnitDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Unit)
		return
	}

	unit, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrUnitNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Unit)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitUpdate)
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitUpdate)
		return
	}

	action := types.UpdateUnitAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: unit.Name,
		}, &dto)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200UnitUpdate)
}

func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Unit)
		return
	}

	unit, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrUnitNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Unit)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitDelete)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, types.ErrUnitIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409UnitDeleteInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500UnitDelete)
		return
	}

	action := types.DeleteUnitAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: unit.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200UnitDelete)
}
