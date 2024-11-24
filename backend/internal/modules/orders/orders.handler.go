package orders

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	status := c.Query("status")

	orders, err := h.service.GetAllOrders(&status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetSubOrders(c *gin.Context) {
	orderIDStr := c.Query("order_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing order_id"})
		return
	}

	subOrders, err := h.service.GetSubOrders(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch suborders"})
		return
	}

	c.JSON(http.StatusOK, subOrders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderDTO types.CreateOrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	err := h.service.CreateOrder(&orderDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *OrderHandler) CompleteSubOrder(c *gin.Context) {
	subOrderIDStr := c.Param("subOrderId")
	subOrderID, err := strconv.ParseUint(subOrderIDStr, 10, 64)
	if err != nil || subOrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing subOrderId"})
		return
	}

	err = h.service.CompleteSubOrder(uint(subOrderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete suborder"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderHandler) GeneratePDFReceipt(c *gin.Context) {
	orderIDParam := c.Param("order_id")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 64)
	if err != nil || orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	pdfData, err := h.service.GeneratePDFReceipt(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PDF receipt"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=order_%d_receipt.pdf", orderID))
	c.Writer.Write(pdfData)
}

func (h *OrderHandler) GetStatusesCount(c *gin.Context) {
	counts, err := h.service.GetStatusesCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch statuses count"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

func (h *OrderHandler) GetSubOrderCount(c *gin.Context) {
	orderIDStr := c.Query("order_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing order_id"})
		return
	}

	count, err := h.service.GetSubOrderCount(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch suborder count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
