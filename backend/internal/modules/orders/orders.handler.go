package orders

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"log"
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

func (h *OrderHandler) GetOrders(c *gin.Context) {
	var filter types.OrdersFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Order{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	filter.Pagination = utils.ParsePagination(c)

	orders, err := h.service.GetOrders(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch orders")
		return
	}

	utils.SendSuccessResponseWithPagination(c, orders, filter.Pagination)
}

func (h *OrderHandler) GetAllBaristaOrders(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	orders, err := h.service.GetAllBaristaOrders(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch orders")
		return
	}

	utils.SendSuccessResponse(c, orders)
}

// Get all suborders for a specific order
func (h *OrderHandler) GetSubOrders(c *gin.Context) {
	orderIDStr := c.Param("orderId")

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		utils.SendBadRequestError(c, "invalid order ID")
		return
	}

	subOrders, err := h.service.GetSubOrders(uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch suborders")
		return
	}

	utils.SendSuccessResponse(c, subOrders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderDTO types.CreateOrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		utils.SendBadRequestError(c, "invalid input data")
		return
	}

	createdOrder, err := h.service.CreateOrder(&orderDTO)
	if err != nil || createdOrder == nil {
		utils.SendInternalServerError(c, fmt.Sprintf("failed to create order: %s", err.Error()))
		return
	}

	createdOrderWithPreloads, err := h.service.GetOrderById(createdOrder.ID)
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("failed to fetch created order with ID %d: %s", createdOrder.ID, err.Error()))
		return
	}

	BroadcastOrderCreated(orderDTO.StoreID, createdOrderWithPreloads)

	utils.SendSuccessResponse(c, createdOrder)
}

// Complete a suborder
func (h *OrderHandler) CompleteSubOrder(c *gin.Context) {
	subOrderIDStr := c.Param("subOrderId")

	subOrderID, err := strconv.ParseUint(subOrderIDStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid suborder ID")
		return
	}

	err = h.service.CompleteSubOrder(uint(subOrderID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to complete suborder")
		return
	}

	order, err := h.service.GetOrderBySubOrder(uint(subOrderID))
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get order: %v", err)
		utils.SendInternalServerError(c, errorMessage)
		return
	}

	BroadcastOrderUpdated(order.StoreID, types.ConvertOrderToDTO(order))

	utils.SendMessageWithStatus(c, "Sub order completed", http.StatusOK)
}

// Generate a PDF receipt for a specific order
func (h *OrderHandler) GeneratePDFReceipt(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
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

// Get count of orders grouped by statuses
func (h *OrderHandler) GetStatusesCount(c *gin.Context) {
	storeIDStr := c.Query("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	statusesWithCounts, err := h.service.GetStatusesCount(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch statuses count")
		return
	}

	utils.SendSuccessResponse(c, statusesWithCounts)
}

func (h *OrderHandler) ServeWS(c *gin.Context) {
	storeIDStr := c.Param("storeId")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 64)
	if err != nil || storeID == 0 {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	conn, err := UpgradeConnection(c)
	if err != nil {
		utils.SendInternalServerError(c, "failed to upgrade WebSocket connection")
		return
	}

	initialOrders, err := h.service.GetAllBaristaOrders(uint(storeID))

	if err != nil {
		log.Printf("Failed to fetch initial orders for store %d: %v", storeID, err)
		utils.SendInternalServerError(c, "failed to fetch initial orders")
		return
	}

	HandleClient(uint(storeID), conn, initialOrders)
}

func (h *OrderHandler) GetOrderDetails(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid order ID")
		return
	}

	orderDetails, err := h.service.GetOrderDetails(uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch order details")
		return
	}

	if orderDetails == nil {
		utils.SendNotFoundError(c, "Order not found")
		return
	}

	utils.SendSuccessResponse(c, orderDetails)
}
