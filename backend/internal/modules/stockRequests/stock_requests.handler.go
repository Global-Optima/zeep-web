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

func (h *StockRequestHandler) GetStockRequests(c *gin.Context) {
	var filter types.GetStockRequestsFilter
	var requests = make([]types.StockRequestResponse, 0)

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

func (h *StockRequestHandler) GetStockRequestByID(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	request, err := h.service.GetStockRequestByID(stockRequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		} else {
			localization.SendLocalizedResponseWithKey(c, types.Response500StockRequest)
		}
		return
	}

	utils.SendSuccessResponse(c, request)
}

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

func (h *StockRequestHandler) RejectWarehouseStatus(c *gin.Context) {
	stockRequestID, err := utils.ParseParam(c, "requestId")
	if err != nil {
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
