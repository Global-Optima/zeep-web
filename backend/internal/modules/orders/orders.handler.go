package orders

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	status := c.Query("status")

	orders, err := h.service.GetAllOrders(uint(storeID), &status)
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch orders")
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetSubOrders(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	orderIDStr := c.Param("orderId")

	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		utils.SendBadRequestError(c, "invalid order ID")
		return
	}

	subOrders, err := h.service.GetSubOrders(uint(storeID), uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch suborders")
		return
	}

	c.JSON(http.StatusOK, subOrders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	var orderDTO types.CreateOrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		utils.SendBadRequestError(c, "invalid input data")
		return
	}

	orderDTO.StoreID = uint(storeID)

	orderId, err := h.service.CreateOrder(&orderDTO)
	if err != nil || orderId == nil {
		utils.SendInternalServerError(c, fmt.Sprintf("failed to create order: %s", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"orderId": orderId})
}

func (h *OrderHandler) CompleteSubOrder(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	orderIDStr := c.Param("orderId")
	subOrderIDStr := c.Param("subOrderId")

	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid order ID")
		return
	}

	subOrderID, err := strconv.ParseUint(subOrderIDStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid suborder ID")
		return
	}

	err = h.service.CompleteSubOrder(uint(subOrderID), uint(orderID), uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to complete suborder")
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderHandler) GeneratePDFReceipt(c *gin.Context) {
	orderIDParam := c.Param("orderId")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 64)
	if err != nil || orderID == 0 {
		utils.SendBadRequestError(c, "invalid order ID")
		return
	}

	pdfData, err := h.service.GeneratePDFReceipt(uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to generate PDF receipt")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=order_%d_receipt.pdf", orderID))
	if _, err := c.Writer.Write(pdfData); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetStatusesCount(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	counts, err := h.service.GetStatusesCount(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch statuses count")
		return
	}

	c.JSON(http.StatusOK, counts)
}

func (h *OrderHandler) GetSubOrderCount(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	orderIDStr := c.Param("orderId")

	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid order ID")
		return
	}

	count, err := h.service.GetSubOrderCount(uint(orderID), uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch suborder count")
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *OrderHandler) GetActiveOrder(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	orderIDStr := c.Param("orderId")

	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		utils.SendBadRequestError(c, "invalid order ID")
		return
	}

	subOrders, err := h.service.GetActiveOrderEvent(uint(storeID), uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch suborders")
		return
	}

	c.JSON(http.StatusOK, subOrders)
}
