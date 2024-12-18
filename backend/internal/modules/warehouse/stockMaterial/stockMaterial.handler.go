package stockMaterial

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StockMaterialHandler struct {
	service StockMaterialService
}

func NewStockMaterialHandler(service StockMaterialService) *StockMaterialHandler {
	return &StockMaterialHandler{service: service}
}

func (h *StockMaterialHandler) GetAllStockMaterials(c *gin.Context) {
	var filter types.StockMaterialFilter

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

	stockMaterialResponses, err := h.service.GetAllStockMaterials(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve stock materials")
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponses)
}

func (h *StockMaterialHandler) GetStockMaterialByID(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	stockMaterialResponse, err := h.service.GetStockMaterialByID(uint(stockMaterialID))
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to retrieve stockMaterial")
		}
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) CreateStockMaterial(c *gin.Context) {
	var req types.CreateStockMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	stockMaterialResponse, err := h.service.CreateStockMaterial(&req)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create stockMaterial")
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) UpdateStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	var req types.UpdateStockMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	stockMaterialResponse, err := h.service.UpdateStockMaterial(uint(stockMaterialID), &req)
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to update stockMaterial")
		}
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) DeleteStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	err = h.service.DeleteStockMaterial(uint(stockMaterialID))
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to delete stockMaterial")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *StockMaterialHandler) DeactivateStockMaterial(c *gin.Context) {
	idParam := c.Param("id")
	stockMaterialID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StockMaterial ID")
		return
	}

	err = h.service.DeactivateStockMaterial(uint(stockMaterialID))
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendInternalServerError(c, "failed to deactivate stockMaterial")
		}
		return
	}

	c.Status(http.StatusNoContent)
}
