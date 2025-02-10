package orders

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/export"
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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	filter.StoreID = &storeID

	orders, err := h.service.GetOrders(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch orders")
		return
	}

	utils.SendSuccessResponseWithPagination(c, orders, filter.Pagination)
}

func (h *OrderHandler) GetAllBaristaOrders(c *gin.Context) {
	var filter types.GetBaristaOrdersFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}
	filter.StoreID = &storeID

	orders, err := h.service.GetAllBaristaOrders(filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch orders")
		return
	}

	utils.SendSuccessResponse(c, orders)
}

func (h *OrderHandler) GetSubOrders(c *gin.Context) {
	orderIDStr := c.Param("orderId")

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	subOrders, err := h.service.GetSubOrders(uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to fetch suborders")
		return
	}

	utils.SendSuccessResponse(c, subOrders)
}

func (h *OrderHandler) CheckCustomerName(c *gin.Context) {
	var dto types.ValidateCustomerNameDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if err := censor.GetCensorValidator().ValidateText(dto.CustomerName); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400OrderCustomerName)
		return
	}

	localization.SendLocalizedResponseWithStatus(c, http.StatusOK)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderDTO types.CreateOrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	createdOrder, err := h.service.CreateOrder(storeID, &orderDTO)
	if err != nil || createdOrder == nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500OrderCreate)
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

func (h *OrderHandler) CompleteSubOrder(c *gin.Context) {
	orderID, errH := utils.ParseParam(c, "orderId")
	if errH != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	subOrderID, errH := utils.ParseParam(c, "subOrderId")
	if errH != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	err := h.service.CompleteSubOrder(orderID, uint(subOrderID))
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

	localization.SendLocalizedResponseWithKey(c, types.Response200OrderUpdate)
}

func (h *OrderHandler) GetSuborderBarcode(c *gin.Context) {
	suborderIDParam := c.Param("subOrderId")
	suborderID, err := strconv.ParseUint(suborderIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	barcodeImage, err := h.service.GenerateSuborderBarcodePDF(uint(suborderID))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	filename := fmt.Sprintf("suborder-barcode-%d.pdf", suborderID)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(barcodeImage)))
	c.Data(http.StatusOK, "application/pdf", barcodeImage)
}

func (h *OrderHandler) CompleteSubOrderByBarcode(c *gin.Context) {
	subOrderID, errH := utils.ParseParam(c, "subOrderId")
	if errH != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	subOrder, err := h.service.CompleteSubOrderByBarcode(subOrderID)
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

	utils.SendSuccessResponse(c, subOrder)
}

func (h *OrderHandler) GeneratePDFReceipt(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
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
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
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
	var filter types.GetBaristaOrdersFilter
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}
	filter.StoreID = &storeID

	conn, err := UpgradeConnection(c)
	if err != nil {
		utils.SendInternalServerError(c, "failed to upgrade WebSocket connection")
		return
	}

	initialOrders, err := h.service.GetAllBaristaOrders(filter)

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
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	orderDetails, err := h.service.GetOrderDetails(uint(orderID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch order details")
		return
	}

	if orderDetails == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, orderDetails)
}

func (h *OrderHandler) ExportOrders(c *gin.Context) {
	var filter types.OrdersExportFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	filter.StoreID = &storeID

	orders, err := h.service.ExportOrders(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to export orders")
		return
	}

	headers := export.RusHeaders
	switch filter.Language {
	case "kk":
		headers = export.KazHeaders
	case "en":
		headers = export.EngHeaders
	}

	excelData, err := export.GenerateSalesExcelV2(orders, headers)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to generate Excel file")
		return
	}

	filename := fmt.Sprintf("orders_export_%s.xlsx", time.Now().Format("02_01_2006"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", fmt.Sprintf("%d", len(excelData)))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

func (h *OrderHandler) AcceptSubOrder(c *gin.Context) {
	subOrderID, err := utils.ParseParam(c, "subOrderId")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	if err := h.service.AcceptSubOrder(uint(subOrderID)); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("failed to accept suborder: %v", err))
		return
	}

	// Optionally fetch the updated order to broadcast the new status.
	order, err := h.service.GetOrderBySubOrder(uint(subOrderID))
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("failed to fetch updated order: %v", err))
		return
	}

	BroadcastOrderUpdated(order.StoreID, types.ConvertOrderToDTO(order))
	localization.SendLocalizedResponseWithKey(c, types.Response200OrderUpdate)
}

func (h *OrderHandler) ChangeSubOrderStatus(c *gin.Context) {
	subOrderID, err := utils.ParseParam(c, "subOrderId")
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Order)
		return
	}

	updatedSuborderDTO, err := h.service.AdvanceSubOrderStatus(uint(subOrderID))
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("failed to update suborder status: %v", err))
		return
	}

	order, err := h.service.GetOrderBySubOrder(uint(subOrderID))
	if err == nil {
		BroadcastOrderUpdated(order.StoreID, types.ConvertOrderToDTO(order))
	}

	utils.SendSuccessResponse(c, updatedSuborderDTO)
}
