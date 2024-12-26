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
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	filter.Pagination = utils.ParsePagination(c)

	stockMaterialResponses, err := h.service.GetAllStockMaterials(&filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, stockMaterialResponses, filter.Pagination)
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
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	utils.SendSuccessResponse(c, stockMaterialResponse)
}

func (h *StockMaterialHandler) CreateStockMaterial(c *gin.Context) {
	var req types.CreateStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	stockMaterialResponse, err := h.service.CreateStockMaterial(&req)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
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

	var req types.UpdateStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	stockMaterialResponse, err := h.service.UpdateStockMaterial(uint(stockMaterialID), &req)
	if err != nil {
		if err.Error() == "StockMaterial not found" {
			utils.SendNotFoundError(c, "StockMaterial not found")
		} else {
			utils.SendBadRequestError(c, err.Error())
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
			utils.SendInternalServerError(c, err.Error())
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
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
