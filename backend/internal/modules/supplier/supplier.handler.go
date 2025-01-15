package supplier

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	service SupplierService
}

func NewSupplierHandler(service SupplierService) *SupplierHandler {
	return &SupplierHandler{service}
}

func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var createDTO types.CreateSupplierDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if err := types.ValidateCreateSupplierDTO(createDTO); err != nil {
		utils.SendBadRequestError(c, "invalid request")
		return
	}

	response, err := h.service.CreateSupplier(createDTO)
	if err != nil {
		utils.SendInternalServerError(c, "fail to create supplier")
		return
	}

	utils.SendResponseWithStatus(c, response, http.StatusCreated)
}

func (h *SupplierHandler) GetSupplierByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.SendBadRequestError(c, "invalid ID")
		return
	}

	response, err := h.service.GetSupplierByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve supplier")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.SendBadRequestError(c, "invalid ID")
		return
	}

	var updateDTO types.UpdateSupplierDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	if err := types.ValidateUpdateSupplierDTO(updateDTO); err != nil {
		utils.SendBadRequestError(c, "invalid request")
		return
	}

	err = h.service.UpdateSupplier(uint(id), updateDTO)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update supplier")
		return
	}

	utils.SendMessageWithStatus(c, "supplier updated successfully", http.StatusOK)
}

func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.SendBadRequestError(c, "invalid ID")
		return
	}

	err = h.service.DeleteSupplier(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete supplier")
		return
	}

	utils.SendMessageWithStatus(c, "supplier deleted successfully", http.StatusOK)
}

func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	suppliers, err := h.service.GetSuppliers()
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve suppliers")
		return
	}

	utils.SendSuccessResponse(c, suppliers)
}

func (h *SupplierHandler) AssociateMaterialToSupplier(c *gin.Context) {
	supplierID, err := strconv.Atoi(c.Param("id"))
	if err != nil || supplierID <= 0 {
		utils.SendBadRequestError(c, "Invalid supplier ID")
		return
	}

	var dto types.CreateSupplierMaterialDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid request body")
		return
	}

	err = h.service.AddMaterialToSupplier(uint(supplierID), dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to associate material to supplier")
		return
	}

	utils.SendMessageWithStatus(c, "Material associated with supplier successfully", http.StatusCreated)
}

func (h *SupplierHandler) GetMaterialsBySupplier(c *gin.Context) {
	supplierID, err := strconv.Atoi(c.Param("id"))
	if err != nil || supplierID <= 0 {
		utils.SendBadRequestError(c, "Invalid supplier ID")
		return
	}

	materials, err := h.service.GetMaterialsBySupplier(uint(supplierID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve materials for the supplier")
		return
	}

	utils.SendSuccessResponse(c, materials)
}
