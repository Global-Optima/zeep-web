package sku

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SKUHandler struct {
	service SKUService
}

func NewSKUHandler(service SKUService) *SKUHandler {
	return &SKUHandler{service: service}
}

func (h *SKUHandler) GetAllSKUs(c *gin.Context) {
	var filter types.SKUFilter

	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}
	if category := c.Query("category"); category != "" {
		filter.Category = &category
	}
	if lowStock := c.Query("lowStock"); lowStock != "" {
		value, err := strconv.ParseBool(lowStock)
		if err == nil {
			filter.LowStock = &value
		}
	}
	if expirationFlag := c.Query("expirationFlag"); expirationFlag != "" {
		value, err := strconv.ParseBool(expirationFlag)
		if err == nil {
			filter.ExpirationFlag = &value
		}
	}
	if isActive := c.Query("isActive"); isActive != "" {
		value, err := strconv.ParseBool(isActive)
		if err == nil {
			filter.IsActive = &value
		}
	}

	skuResponses, err := h.service.GetAllSKUs(&filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, skuResponses)
}

func (h *SKUHandler) GetSKUByID(c *gin.Context) {
	idParam := c.Param("id")
	skuID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid SKU ID")
		return
	}

	skuResponse, err := h.service.GetSKUByID(uint(skuID))
	if err != nil {
		if err.Error() == "SKU not found" {
			utils.SendNotFoundError(c, "SKU not found")
		} else {
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, skuResponse)
}

func (h *SKUHandler) CreateSKU(c *gin.Context) {
	var req types.CreateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	skuResponse, err := h.service.CreateSKU(&req)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, skuResponse)
}

func (h *SKUHandler) UpdateSKU(c *gin.Context) {
	idParam := c.Param("id")
	skuID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid SKU ID")
		return
	}

	var req types.UpdateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	skuResponse, err := h.service.UpdateSKU(uint(skuID), &req)
	if err != nil {
		if err.Error() == "SKU not found" {
			utils.SendNotFoundError(c, "SKU not found")
		} else {
			utils.SendBadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, skuResponse)
}

func (h *SKUHandler) DeleteSKU(c *gin.Context) {
	idParam := c.Param("id")
	skuID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid SKU ID")
		return
	}

	err = h.service.DeleteSKU(uint(skuID))
	if err != nil {
		if err.Error() == "SKU not found" {
			utils.SendNotFoundError(c, "SKU not found")
		} else {
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SKUHandler) DeactivateSKU(c *gin.Context) {
	idParam := c.Param("id")
	skuID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid SKU ID")
		return
	}

	err = h.service.DeactivateSKU(uint(skuID))
	if err != nil {
		if err.Error() == "SKU not found" {
			utils.SendNotFoundError(c, "SKU not found")
		} else {
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
