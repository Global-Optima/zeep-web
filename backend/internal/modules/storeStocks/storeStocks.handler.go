package storeStocks

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreStockHandler struct {
	service           StoreStockService
	auditService      audit.AuditService
	ingredientService ingredients.IngredientService
	franchiseeService franchisees.FranchiseeService
	logger            *zap.SugaredLogger
}

func NewStoreStockHandler(
	service StoreStockService,
	ingredientService ingredients.IngredientService,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	logger *zap.SugaredLogger,
) *StoreStockHandler {
	return &StoreStockHandler{
		service:           service,
		ingredientService: ingredientService,
		auditService:      auditService,
		franchiseeService: franchiseeService,
		logger:            logger,
	}
}

func (h *StoreStockHandler) GetAvailableIngredientsToAdd(c *gin.Context) {
	var filter ingredientTypes.IngredientFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Ingredient{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	ingredientsList, err := h.service.GetAvailableIngredientsToAdd(storeID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	utils.SendSuccessResponseWithPagination(c, ingredientsList, filter.Pagination)
}

func (h *StoreStockHandler) AddStoreStock(c *gin.Context) {
	var dto types.AddStoreStockDTO

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	ingredient, err := h.ingredientService.GetIngredientByID(dto.IngredientID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	id, err := h.service.AddStock(storeID, &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	action := types.CreateStoreStockAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: ingredient.Name,
		},
		&dto, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StoreStock)
}

func (h *StoreStockHandler) AddMultipleStoreStock(c *gin.Context) {
	var dto types.AddMultipleStoreStockDTO

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	IDs, err := h.service.AddMultipleStock(storeID, &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	stockList, err := h.service.GetStockListByIDs(storeID, IDs)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	dtoMap := make(map[uint]*types.AddStoreStockDTO)
	for _, stockDTO := range dto.IngredientStocks {
		dtoCopy := stockDTO
		dtoMap[stockDTO.IngredientID] = &dtoCopy
	}

	var actions []shared.AuditAction
	for _, stock := range stockList {
		matchedDTO, exists := dtoMap[stock.Ingredient.ID]
		if !exists {
			h.logger.Errorf("Failed to match stock with DTO for stock ID: %d", stock.ID)
			continue
		}

		action := types.CreateStoreStockAuditFactory(
			&data.BaseDetails{
				ID:   stock.ID,
				Name: stock.Name,
			},
			matchedDTO, storeID,
		)
		actions = append(actions, &action)
	}

	go func() {
		_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StoreStockMultiple)
}

func (h *StoreStockHandler) GetStoreStockList(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockFilter := &types.GetStockFilterQuery{}
	if err := utils.ParseQueryWithBaseFilter(c, stockFilter, &data.StoreStock{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	stockList, err := h.service.GetStockList(storeID, stockFilter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	utils.SendSuccessResponseWithPagination(c, stockList, stockFilter.GetPagination())
}

func (h *StoreStockHandler) GetStoreStockById(c *gin.Context) {
	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreStock)
		return
	}

	filter, errH := contexts.GetStoreContextFilter(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stock, err := h.service.GetStockById(uint(stockId), filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	utils.SendSuccessResponse(c, stock)
}

func (h *StoreStockHandler) UpdateStoreStockById(c *gin.Context) {
	var input types.UpdateStoreStockDTO

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreStock)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	stock, err := h.service.GetStockById(uint(stockId), &contexts.StoreContextFilter{StoreID: &storeID})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	err = h.service.UpdateStockById(storeID, uint(stockId), &input)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	action := types.UpdateStoreStockAuditFactory(
		&data.BaseDetails{
			ID:   uint(stockId),
			Name: stock.Name,
		},
		&input, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreStockUpdate)
}

func (h *StoreStockHandler) DeleteStoreStockById(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreStock)
		return
	}

	stock, err := h.service.GetStockById(uint(stockId), &contexts.StoreContextFilter{StoreID: &storeID})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	err = h.service.DeleteStockById(storeID, uint(stockId))
	if err != nil {
		if errors.Is(err, types.ErrStockIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response500StoreStockIsInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreStock)
		return
	}

	action := types.DeleteStoreStockAuditFactory(
		&data.BaseDetails{
			ID:   uint(stockId),
			Name: stock.Name,
		},
		struct{}{}, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreStockDelete)
}
