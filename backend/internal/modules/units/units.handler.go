package units

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	service UnitService
}

func NewUnitHandler(service UnitService) *UnitHandler {
	return &UnitHandler{service: service}
}

func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var dto types.CreateUnitDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	_, err := h.service.Create(dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessCreatedResponse(c, "unit created successfully")
}

func (h *UnitHandler) GetAllUnits(c *gin.Context) {
	var filter types.UnitFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Unit{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	units, err := h.service.GetAll(&filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, units, filter.GetPagination())
}

func (h *UnitHandler) GetUnitByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	unit, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, unit)
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	var dto types.UpdateUnitDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
