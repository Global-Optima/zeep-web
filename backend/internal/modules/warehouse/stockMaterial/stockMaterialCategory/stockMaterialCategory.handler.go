package stockMaterialCategory

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StockMaterialCategoryHandler struct {
	service      StockMaterialCategoryService
	auditService audit.AuditService
}

func NewStockMaterialCategoryHandler(service StockMaterialCategoryService, auditService audit.AuditService) *StockMaterialCategoryHandler {
	return &StockMaterialCategoryHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *StockMaterialCategoryHandler) Create(c *gin.Context) {
	var dto types.CreateStockMaterialCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryCreate)
		return
	}

	action := types.CreateStockMaterialAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StockMaterialCategory)
}

func (h *StockMaterialCategoryHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	response, err := h.service.GetByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrNotFound):
			localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryGet)
			return
		}
	}

	utils.SendSuccessResponse(c, response)
}

func (h *StockMaterialCategoryHandler) GetAll(c *gin.Context) {
	var filter types.StockMaterialCategoryFilter

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockMaterialCategory{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	categories, err := h.service.GetAll(filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.Pagination)
}

func (h *StockMaterialCategoryHandler) Update(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.UpdateStockMaterialCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	response, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockMaterialCategory)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryGet)
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryUpdate)
		return
	}

	action := types.UpdateStockMaterialAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: response.Name,
		}, &dto)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StockMaterialCategoryUpdate)
}

func (h *StockMaterialCategoryHandler) Delete(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	response, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrStockMaterialCategoryNotFound) {
			utils.SendNotFoundError(c, types.ErrStockMaterialCategoryNotFound.Error())
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryGet)
		}
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockMaterialCategoryDelete)
		return
	}

	action := types.DeleteStockMaterialAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: response.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StockMaterialCategoryDelete)
}
