package stockRequests

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockRequestHandler struct {
	service           StockRequestService
	franchiseeService franchisees.FranchiseeService
	regionService     regions.RegionService
	auditService      audit.AuditService
}

func NewStockRequestHandler(service StockRequestService,
	franchiseeService franchisees.FranchiseeService,
	regionService regions.RegionService,
	auditService audit.AuditService,
) *StockRequestHandler {
	return &StockRequestHandler{
		service:           service,
		franchiseeService: franchiseeService,
		regionService:     regionService,
		auditService:      auditService,
	}
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

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	id, storeName, err := h.service.CreateStockRequest(storeID, req)
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to create stock requests: %s", err.Error()))
		return
	}

	action := types.CreateStockRequestAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: storeName,
		},
		&types.AuditPayloads{
			CreateStockRequestDTO: &req,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "stock request created successfully")
}

func (h *StockRequestHandler) GetStockRequests(c *gin.Context) {
	var filter types.GetStockRequestsFilter
	var requests = make([]types.StockRequestResponse, 0)

	// Check user access roles
	storeID, storeErr := h.franchiseeService.CheckFranchiseeStore(c)
	warehouseID, warehouseErr := h.regionService.CheckRegionWarehouse(c)

	// If user has neither store nor warehouse access, return unauthorized
	if storeErr != nil && warehouseErr != nil {
		utils.SendErrorWithStatus(c, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Parse filters
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockRequest{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	// Fetch requests for Store users
	if storeErr == nil {
		filter.StoreID = &storeID
		storeRequests, err := h.service.GetStockRequests(filter)
		if err != nil {
			utils.SendInternalServerError(c, "Failed to fetch stock requests")
			return
		}
		requests = append(requests, storeRequests...)
	}

	// Fetch requests for Warehouse users
	if warehouseErr == nil {
		filter.WarehouseID = &warehouseID

		// Validate warehouse statuses
		if err := types.ValidateWarehouseStatuses(filter.Statuses); err != nil {
			utils.SendBadRequestError(c, err.Error())
			return
		}

		// Apply default warehouse statuses if none are provided
		if len(filter.Statuses) == 0 {
			filter.Statuses = types.DefaultWarehouseStatuses()
		}

		warehouseRequests, err := h.service.GetStockRequests(filter)
		if err != nil {
			utils.SendInternalServerError(c, "Failed to fetch stock requests")
			return
		}
		requests = append(requests, warehouseRequests...)
	}

	// Always return an array, even if empty
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

	request, err := h.service.AcceptStockRequestWithChange(uint(stockRequestID), dto)
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Store.Name,
		},
		&types.AuditPayloads{
			AcceptWithChangeRequestStatusDTO: &dto,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

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

	request, err := h.service.RejectStockRequestByStore(uint(stockRequestID), dto)
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Store.Name,
		},
		&types.AuditPayloads{
			RejectStockRequestStatusDTO: &dto,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

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

	request, err := h.service.RejectStockRequestByWarehouse(uint(stockRequestID), dto)
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Warehouse.Name,
		},
		&types.AuditPayloads{
			RejectStockRequestStatusDTO: &dto,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) SetProcessedStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	request, err := h.service.SetProcessedStatus(uint(stockRequestID))
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Store.Name,
		},
		&types.AuditPayloads{},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) SetInDeliveryStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	request, err := h.service.SetInDeliveryStatus(uint(stockRequestID))
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Warehouse.Name,
		},
		&types.AuditPayloads{},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) SetCompletedStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	request, err := h.service.SetCompletedStatus(uint(stockRequestID))
	if err != nil {
		utils.SendInternalServerError(c, fmt.Sprintf("Failed to update status: %s", err.Error()))
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Store.Name,
		},
		&types.AuditPayloads{},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

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

	request, err := h.service.UpdateStockRequest(uint(stockRequestID), req)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update stock request ingredients: "+err.Error())
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Store.Name,
		},
		&types.AuditPayloads{
			StockMaterials: req,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "Stock request ingredients updated successfully"})
}

func (h *StockRequestHandler) DeleteStockRequest(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if utils.SendBadRequestInvalidParam(c, "requestId", err) {
		return
	}

	request, err := h.service.DeleteStockRequest(uint(stockRequestID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete stock request")
		return
	}

	action := types.DeleteStockRequestAuditFactory(
		&data.BaseDetails{
			ID:   stockRequestID,
			Name: request.Store.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	c.Status(http.StatusOK)
}

func (h *StockRequestHandler) GetLastCreatedStockRequest(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
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
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var dto types.StockRequestStockMaterialDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, "Invalid payload")
		return
	}

	request, err := h.service.AddStockMaterialToCart(uint(storeID), dto)
	if err != nil {
		utils.SendInternalServerError(c, "Failed too add to cart")
		return
	}

	action := types.UpdateStockRequestStatusAuditFactory(
		&data.BaseDetails{
			ID:   request.ID,
			Name: request.Store.Name,
		},
		&types.AuditPayloads{
			StockMaterials: []types.StockRequestStockMaterialDTO{dto},
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	c.Status(http.StatusOK)
}
