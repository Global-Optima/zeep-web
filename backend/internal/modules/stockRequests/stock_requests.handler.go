package stockRequests

import (
	"errors"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// StockRequestHandler handles stock request endpoints.
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

// CreateStockRequest godoc
// @Summary Create a new stock request
// @Description Creates a new stock request using the provided stock materials.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param input body types.CreateStockRequestDTO true "Stock Request Data"
// @Success 201 {object} map[string]interface{} "Created stock request"
// @Failure 400 {object} map[string]interface{} "Bad Request – invalid input or missing stock materials"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests [post]
func (h *StockRequestHandler) CreateStockRequest(c *gin.Context) {
	var req types.CreateStockRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if len(req.StockMaterials) == 0 {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockRequest)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	id, storeName, err := h.service.CreateStockRequest(storeID, req)
	if err != nil {
		if errors.Is(err, types.ErrExistingRequest) {
			localization.SendLocalizedResponseWithKey(c, types.Response400StockRequestExistingRequest)
			return
		}
		if errors.Is(err, types.ErrOneRequestPerDay) {
			localization.SendLocalizedResponseWithKey(c, types.Response400StockRequestOnlyOneRequestPerDay)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201StockRequest)
}

// GetStockRequests godoc
// @Summary Get stock requests
// @Description Retrieves stock requests filtered by store, warehouse, or date range. Supports pagination.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param startDate query string false "Start date (ISO8601 format)"
// @Param endDate query string false "End date (ISO8601 format)"
// @Param search query string false "Search term"
// @Param statuses query []string false "Stock request statuses" collectionFormat(multi)
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "Stock requests with pagination"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests [get]
func (h *StockRequestHandler) GetStockRequests(c *gin.Context) {
	var filter types.GetStockRequestsFilter
	requests := make([]types.StockRequestResponse, 0)

	storeID, storeErr := h.franchiseeService.CheckFranchiseeStore(c)
	warehouseID, warehouseErr := h.regionService.CheckRegionWarehouse(c)

	if storeErr != nil && warehouseErr != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response403StockRequest)
		return
	}

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StockRequest{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	if storeErr == nil {
		filter.StoreID = &storeID
		storeRequests, err := h.service.GetStockRequests(filter)
		if err != nil {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
			return
		}
		requests = append(requests, storeRequests...)
	}

	if warehouseErr == nil {
		filter.WarehouseID = &warehouseID
		if err := types.ValidateWarehouseStatuses(filter.Statuses); err != nil {
			localization.SendLocalizedResponseWithKey(c, types.Response400StockRequest)
			return
		}
		if len(filter.Statuses) == 0 {
			filter.Statuses = types.DefaultWarehouseStatuses()
		}
		warehouseRequests, err := h.service.GetStockRequests(filter)
		if err != nil {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
			return
		}
		requests = append(requests, warehouseRequests...)
	}

	utils.SendSuccessResponseWithPagination(c, requests, filter.Pagination)
}

// GetStockRequestByID godoc
// @Summary Get a stock request by ID
// @Description Retrieves the stock request details by its ID.
// @Tags stock-requests
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Success 200 {object} types.StockRequestResponse
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /api/v1/stock-requests/{requestId} [get]
func (h *StockRequestHandler) GetStockRequestByID(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	request, err := h.service.GetStockRequestByID(stockRequestID)
	if err != nil {
		if errors.Is(err, types.ErrStockRequestNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404StockRequest)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
		return
	}

	utils.SendSuccessResponse(c, request)
}

// AcceptWithChangeStatus godoc
// @Summary Accept stock request with changes
// @Description Accepts a stock request, applying changes as specified.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Param input body types.AcceptWithChangeRequestStatusDTO true "Acceptance details with changes"
// @Success 200 {object} map[string]interface{} "Stock request status updated"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/status/{requestId}/accept-with-change [put]
func (h *StockRequestHandler) AcceptWithChangeStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.AcceptWithChangeRequestStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if len(dto.Items) == 0 {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockRequest)
		return
	}

	request, err := h.service.AcceptStockRequestWithChange(stockRequestID, dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// RejectStoreStatus godoc
// @Summary Reject stock request by store
// @Description Rejects a stock request at the store level with a comment.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Param input body types.RejectStockRequestStatusDTO true "Rejection details"
// @Success 200 {object} map[string]interface{} "Stock request status updated"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/status/{requestId}/reject-store [put]
func (h *StockRequestHandler) RejectStoreStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var dto types.RejectStockRequestStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	request, err := h.service.RejectStockRequestByStore(stockRequestID, dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// RejectWarehouseStatus godoc
// @Summary Reject stock request by warehouse
// @Description Rejects a stock request at the warehouse level with a comment.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Param input body types.RejectStockRequestStatusDTO true "Rejection details"
// @Success 200 {object} map[string]interface{} "Stock request status updated"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/status/{requestId}/reject-warehouse [put]
func (h *StockRequestHandler) RejectWarehouseStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		// In this handler the error is silently returned.
		return
	}

	var dto types.RejectStockRequestStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	request, err := h.service.RejectStockRequestByWarehouse(stockRequestID, dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// SetProcessedStatus godoc
// @Summary Set stock request to processed
// @Description Updates the stock request status to "processed".
// @Tags stock-requests
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Success 200 {object} map[string]interface{} "Stock request updated to processed"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/status/{requestId}/processed [put]
func (h *StockRequestHandler) SetProcessedStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	request, err := h.service.SetProcessedStatus(stockRequestID)
	if err != nil {
		if errors.Is(err, types.ErrOneRequestPerDay) {
			localization.SendLocalizedResponseWithKey(c, types.Response400StockRequestOnlyOneRequestPerDay)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// SetInDeliveryStatus godoc
// @Summary Set stock request to in-delivery
// @Description Updates the stock request status to "in delivery".
// @Tags stock-requests
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Success 200 {object} map[string]interface{} "Stock request updated to in delivery"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/status/{requestId}/in-delivery [put]
func (h *StockRequestHandler) SetInDeliveryStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	request, err := h.service.SetInDeliveryStatus(stockRequestID)
	if err != nil {
		if errors.Is(err, types.ErrInsufficientStock) {
			localization.SendLocalizedResponseWithKey(c, types.Response400StockRequestInsufficientStock)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// SetCompletedStatus godoc
// @Summary Set stock request to completed
// @Description Updates the stock request status to "completed".
// @Tags stock-requests
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Success 200 {object} map[string]interface{} "Stock request updated to completed"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/status/{requestId}/completed [put]
func (h *StockRequestHandler) SetCompletedStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	request, err := h.service.SetCompletedStatus(uint(stockRequestID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// UpdateStockRequest godoc
// @Summary Update stock request
// @Description Updates the stock materials for an existing stock request.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Param input body []types.StockRequestStockMaterialDTO true "Stock materials update data"
// @Success 200 {object} map[string]interface{} "Stock request updated"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/{requestId} [put]
func (h *StockRequestHandler) UpdateStockRequest(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	var req []types.StockRequestStockMaterialDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if len(req) == 0 {
		localization.SendLocalizedResponseWithKey(c, types.Response400StockRequest)
		return
	}

	request, err := h.service.UpdateStockRequest(stockRequestID, req)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestUpdate)
}

// DeleteStockRequest godoc
// @Summary Delete stock request
// @Description Deletes a stock request based on its ID.
// @Tags stock-requests
// @Produce json
// @Param requestId path int true "Stock Request ID"
// @Success 200 {object} map[string]interface{} "Stock request deleted"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/{requestId} [delete]
func (h *StockRequestHandler) DeleteStockRequest(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	request, err := h.service.DeleteStockRequest(stockRequestID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200StockRequestDelete)
}

// GetLastCreatedStockRequest godoc
// @Summary Get last created stock request
// @Description Retrieves the most recently created stock request for the franchisee store.
// @Tags stock-requests
// @Produce json
// @Success 200 {object} types.StockRequestResponse
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/current [get]
func (h *StockRequestHandler) GetLastCreatedStockRequest(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	request, err := h.service.GetLastCreatedStockRequest(storeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
		return
	}

	utils.SendSuccessResponse(c, request)
}

// AddStockMaterialToCart godoc
// @Summary Add stock material to cart
// @Description Adds a stock material entry to the current cart for the store.
// @Tags stock-requests
// @Accept json
// @Produce json
// @Param input body types.StockRequestStockMaterialDTO true "Stock Material Data"
// @Success 201 {object} map[string]interface{} "Stock material added to cart"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/stock-requests/add-material-to-latest-cart [post]
func (h *StockRequestHandler) AddStockMaterialToCart(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var dto types.StockRequestStockMaterialDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	request, err := h.service.AddStockMaterialToCart(storeID, dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201StockRequest)
}
