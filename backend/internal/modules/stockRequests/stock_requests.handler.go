package stockRequests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockRequestHandler struct {
	service StockRequestService
}

func NewStockRequestHandler(service StockRequestService) *StockRequestHandler {
	return &StockRequestHandler{service: service}
}

func (h *StockRequestHandler) CreateStockRequest(c *gin.Context) {
	var req types.CreateStockRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if len(req.Items) == 0 {
		utils.SendBadRequestError(c, "The cart cannot be empty")
		return
	}

	requestID, err := h.service.CreateStockRequest(req)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessCreatedResponse(c, gin.H{"requestId": requestID})
}

func (h *StockRequestHandler) GetStockRequests(c *gin.Context) {
	var filter types.GetStockRequestsFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	filter.Pagination = utils.ParsePagination(c)

	requests, err := h.service.GetStockRequests(filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, requests, filter.Pagination)
}

func (h *StockRequestHandler) GetStockRequestByID(c *gin.Context) {
	idParam := c.Param("id")
	stockRequestID, err := strconv.Atoi(idParam)
	if err != nil || stockRequestID <= 0 {
		utils.SendBadRequestError(c, "Invalid stock request ID")
		return
	}

	request, err := h.service.GetStockRequestByID(uint(stockRequestID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendNotFoundError(c, "Stock request not found")
		} else {
			utils.SendInternalServerError(c, err.Error())
		}
		return
	}

	utils.SendSuccessResponse(c, request)
}

func (h *StockRequestHandler) UpdateStockRequestStatus(c *gin.Context) {
	requestID, err := strconv.Atoi(c.Param("requestId"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request ID")
		return
	}

	var status types.UpdateStockRequestStatusDTO
	if err := c.ShouldBindJSON(&status); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if err := h.service.UpdateStockRequestStatus(uint(requestID), status); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) UpdateStockRequestIngredients(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("requestId"), 10, 32)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request ID")
		return
	}

	var req []types.CreateStockRequestItemDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input: "+err.Error())
		return
	}

	if len(req) == 0 {
		utils.SendBadRequestError(c, "The cart cannot be empty")
		return
	}

	if err := h.service.UpdateStockRequestIngredients(uint(requestID), req); err != nil {
		utils.SendInternalServerError(c, "Failed to update stock request ingredients: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Stock request ingredients updated successfully"})
}

func (h *StockRequestHandler) GetLowStockIngredients(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("storeId"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	ingredients, err := h.service.GetLowStockIngredients(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, ingredients, nil)
}

func (h *StockRequestHandler) GetAllStockMaterials(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("storeId"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	var filter types.StockMaterialFilter
	if err = c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	products, err := h.service.GetAllStockMaterials(uint(storeID), filter)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponseWithPagination(c, products, nil)
}

func (h *StockRequestHandler) GetStockMaterialsByIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("ingredientId"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ingredient ID")
		return
	}

	var warehouseID *uint
	if warehouseParam := c.Query("warehouseId"); warehouseParam != "" {
		id, err := strconv.ParseUint(warehouseParam, 10, 64)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid warehouse ID")
			return
		}
		wID := uint(id)
		warehouseID = &wID
	}

	availability, err := h.service.GetAvailableStockMaterialsByIngredient(uint(ingredientID), warehouseID)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, availability)
}
