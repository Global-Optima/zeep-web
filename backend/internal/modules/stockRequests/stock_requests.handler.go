package stockRequests

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID, err := h.service.CreateStockRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"requestId": requestID})
}

func (h *StockRequestHandler) GetStockRequests(c *gin.Context) {
	var filter types.StockRequestFilter

	if storeID := c.Query("storeId"); storeID != "" {
		id, err := strconv.ParseUint(storeID, 10, 64)
		if err == nil {
			storeIDUint := uint(id)
			filter.StoreID = &storeIDUint
		}
	}

	if warehouseID := c.Query("warehouseId"); warehouseID != "" {
		id, err := strconv.ParseUint(warehouseID, 10, 64)
		if err == nil {
			warehouseIDUint := uint(id)
			filter.WarehouseID = &warehouseIDUint
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if startDate := c.Query("startDate"); startDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, startDate)
		if err == nil {
			filter.StartDate = &parsedDate
		}
	}

	if endDate := c.Query("endDate"); endDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, endDate)
		if err == nil {
			filter.EndDate = &parsedDate
		}
	}

	requests, err := h.service.GetStockRequests(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *StockRequestHandler) UpdateStockRequestStatus(c *gin.Context) {
	requestID, err := strconv.Atoi(c.Param("requestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var status types.UpdateStockRequestStatusDTO
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStockRequestStatus(uint(requestID), status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) GetLowStockIngredients(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("storeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	ingredients, err := h.service.GetLowStockIngredients(uint(storeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingredients)
}

func (h *StockRequestHandler) GetMarketplaceProducts(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("storeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	var filter types.MarketplaceFilter
	if category := c.Query("category"); category != "" {
		filter.Category = &category
	}
	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	products, err := h.service.GetMarketplaceProducts(uint(storeID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *StockRequestHandler) AddStockRequestIngredient(c *gin.Context) {
	requestID, err := strconv.Atoi(c.Param("requestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var item types.StockRequestItemDTO
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddStockRequestIngredient(uint(requestID), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *StockRequestHandler) DeleteStockRequestIngredient(c *gin.Context) {
	ingredientID, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	if err := h.service.DeleteStockRequestIngredient(uint(ingredientID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
