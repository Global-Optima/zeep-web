package supplier

import (
	"errors"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	service      SupplierService
	auditService audit.AuditService
}

func NewSupplierHandler(service SupplierService, auditService audit.AuditService) *SupplierHandler {
	return &SupplierHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var createDTO types.CreateSupplierDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := types.ValidateCreateSupplierDTO(createDTO); err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateSupplier(createDTO)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierCreate)
		return
	}

	action := types.CreateSupplierAuditFactory(
		&data.BaseDetails{
			ID:   response.ID,
			Name: response.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Supplier)
}

func (h *SupplierHandler) GetSupplierByID(c *gin.Context) {
	supplierID, err := utils.ParseParam(c, "id")
	if err != nil || supplierID <= 0 {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	response, err := h.service.GetSupplierByID(uint(supplierID))
	if err != nil {
		if errors.Is(err, types.ErrSupplierNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Supplier)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierGet)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	supplierID, err := utils.ParseParam(c, "id")
	if err != nil || supplierID <= 0 {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var updateDTO types.UpdateSupplierDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := types.ValidateUpdateSupplierDTO(updateDTO); err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetSupplierByID(uint(supplierID))
	if err != nil {
		if errors.Is(err, types.ErrSupplierNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Supplier)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierGet)
		return
	}

	err = h.service.UpdateSupplier(uint(supplierID), updateDTO)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierUpdate)
		return
	}

	action := types.UpdateSupplierAuditFactory(
		&data.BaseDetails{
			ID:   uint(supplierID),
			Name: supplier.Name,
		}, &updateDTO)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200SupplierUpdate)
}

func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	supplierID, err := utils.ParseParam(c, "id")
	if err != nil || supplierID <= 0 {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	supplier, err := h.service.GetSupplierByID(uint(supplierID))
	if err != nil {
		if errors.Is(err, types.ErrSupplierNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Supplier)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierGet)
		return
	}

	err = h.service.DeleteSupplier(uint(supplierID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierDelete)
		return
	}

	action := types.DeleteSupplierAuditFactory(
		&data.BaseDetails{
			ID:   uint(supplierID),
			Name: supplier.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200SupplierDelete)
}

func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	var filter types.SuppliersFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Supplier{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	suppliers, err := h.service.GetSuppliers(filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, suppliers, filter.Pagination)
}

func (h *SupplierHandler) UpsertMaterialsForSupplier(c *gin.Context) {
	supplierID, err := utils.ParseParam(c, "id")
	if err != nil || supplierID <= 0 {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.UpsertSupplierMaterialsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	err = h.service.UpsertMaterialsForSupplier(uint(supplierID), dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierUpsertMaterials)
		return
	}

	localization.SendLocalizedResponseWithKey(c, types.Response200SupplierUpsertMaterial)
}

func (h *SupplierHandler) GetMaterialsBySupplier(c *gin.Context) {
	supplierID, err := utils.ParseParam(c, "id")
	if err != nil || supplierID <= 0 {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	materials, err := h.service.GetMaterialsBySupplier(uint(supplierID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500SupplierGetMaterials)
		return
	}

	utils.SendSuccessResponse(c, materials)
}
