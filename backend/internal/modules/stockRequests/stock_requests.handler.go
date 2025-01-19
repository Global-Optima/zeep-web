package stockRequests

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
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
		utils.SendBadRequestError(c, "Invalid body")
		return
	}

	if len(req.StockMaterials) == 0 {
		utils.SendBadRequestError(c, "The cart cannot be empty")
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	requestID, err := h.service.CreateStockRequest(storeID, req)
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to create stock requests: %s", err.Error()))
		return
	}

	utils.SuccessCreatedResponse(c, gin.H{"requestId": requestID})
}

func (h *StockRequestHandler) GetStockRequests(c *gin.Context) {
	var filter types.GetStockRequestsFilter

	storeID, storeErr := contexts.GetStoreId(c)
	warehouseID, warehouseErr := contexts.GetWarehouseId(c)

	if storeErr != nil && warehouseErr != nil {
		utils.SendErrorWithStatus(c, "Unauthorized access", http.StatusForbidden)
		return
	}

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockRequest{}); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	filter.Pagination = utils.ParsePagination(c)

	var requests []types.StockRequestResponse
	var err error

	if storeErr == nil {
		filter.StoreID = &storeID
		requests, err = h.service.GetStockRequests(filter)
	}

	if warehouseErr == nil {
		filter.WarehouseID = &warehouseID
		err = types.ValidateWarehouseStatuses(filter.Statuses)
		if err != nil {
			utils.SendBadRequestError(c, err.Error())
		}
		requests, err = h.service.GetStockRequests(filter)
	}

	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch stock requests")
		return
	}

	utils.SendSuccessResponseWithPagination(c, requests, filter.Pagination)
}

func (h *StockRequestHandler) GetStockRequestByID(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	request, err := h.service.GetStockRequestByID(uint(stockRequestID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendNotFoundError(c, "Stock request not found")
		} else {
			utils.SendInternalServerError(c, "Failed to fetch stock request")
		}
		return
	}

	utils.SendSuccessResponse(c, request)
}

func (h *StockRequestHandler) AcceptWithChangeStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	var dto types.AcceptWithChangeRequestStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid payload for acceptance with change")
		return
	}

	if len(dto.Items) == 0 {
		utils.SendBadRequestError(c, "The cart cannot be empty")
		return
	}

	if err := h.service.AcceptStockRequestWithChange(uint(stockRequestID), dto); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) RejectStoreStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	var dto types.RejectStockRequestStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid payload for rejection by store")
		return
	}

	if err := h.service.RejectStockRequestByStore(uint(stockRequestID), dto); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) RejectWarehouseStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	var dto types.RejectStockRequestStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid payload for rejection by warehouse")
		return
	}

	if err := h.service.RejectStockRequestByWarehouse(uint(stockRequestID), dto); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) SetProcessedStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	if err := h.service.SetProcessedStatus(uint(stockRequestID)); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) SetInDeliveryStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	if err := h.service.SetInDeliveryStatus(uint(stockRequestID)); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) SetCompletedStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	if err := h.service.SetCompletedStatus(uint(stockRequestID)); err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) UpdateStockRequest(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	var req []types.StockRequestStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid input: "+err.Error())
		return
	}

	if len(req) == 0 {
		utils.SendBadRequestError(c, "The cart cannot be empty")
		return
	}

	if err := h.service.UpdateStockRequest(uint(stockRequestID), req); err != nil {
		utils.SendInternalServerError(c, "Failed to update stock request ingredients: "+err.Error())
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "Stock request ingredients updated successfully"})
}

func (h *StockRequestHandler) DeleteStockRequest(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	err = h.service.DeleteStockRequest(uint(stockRequestID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete stock request")
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) GetLastCreatedStockRequest(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	request, err := h.service.GetLastCreatedStockRequest(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch the last cart")
		return
	}

	utils.SendSuccessResponse(c, request)
}

func (h *StockRequestHandler) AddStockMaterialToCart(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var dto types.StockRequestStockMaterialDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid payload")
		return
	}

	err := h.service.AddStockMaterialToCart(uint(storeID), dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed too add to cart")
		return
	}

	c.Status(http.StatusOK)
}
