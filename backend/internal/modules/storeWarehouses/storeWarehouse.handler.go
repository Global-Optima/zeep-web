package storeWarehouses

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"go.uber.org/zap"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreWarehouseHandler struct {
	service           StoreWarehouseService
	auditService      audit.AuditService
	ingredientService ingredients.IngredientService
	franchiseeService franchisees.FranchiseeService
	logger            *zap.SugaredLogger
}

func NewStoreWarehouseHandler(
	service StoreWarehouseService,
	ingredientService ingredients.IngredientService,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	logger *zap.SugaredLogger,
) *StoreWarehouseHandler {
	return &StoreWarehouseHandler{
		service:           service,
		ingredientService: ingredientService,
		auditService:      auditService,
		franchiseeService: franchiseeService,
		logger:            logger,
	}
}

func (h *StoreWarehouseHandler) GetAvailableIngredientsToAdd(c *gin.Context) {
	var filter ingredientTypes.IngredientFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Ingredient{}); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	ingredientsList, err := h.service.GetAvailableIngredientsToAdd(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch ingredients")
		return
	}

	utils.SendSuccessResponseWithPagination(c, ingredientsList, filter.Pagination)
}

func (h *StoreWarehouseHandler) AddStoreWarehouseStock(c *gin.Context) {
	var dto types.AddStoreStockDTO

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	ingredient, err := h.ingredientService.GetIngredientByID(dto.IngredientID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to add new stock: ingredient not found")
		return
	}

	id, err := h.service.AddStock(storeID, &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to add new stock")
		return
	}

	action := types.CreateStoreWarehouseStockAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: ingredient.Name,
		},
		&dto, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{
		"message": fmt.Sprintf("store warehouse stock with id %d successfully created", id),
	})
}

func (h *StoreWarehouseHandler) AddMultipleStoreWarehouseStock(c *gin.Context) {
	var dto types.AddMultipleStoreStockDTO

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	IDs, err := h.service.AddMultipleStock(storeID, &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to add new multiple stocks")
		return
	}

	stockList, err := h.service.GetStockListByIDs(storeID, IDs)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve added stock: ingredient not found")
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

		action := types.CreateStoreWarehouseStockAuditFactory(
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

	utils.SendSuccessResponse(c, gin.H{
		"message": "success",
	})
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockList(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockFilter := &types.GetStockFilterQuery{}
	if err := utils.ParseQueryWithBaseFilter(c, stockFilter, &data.StoreWarehouseStock{}); err != nil {
		utils.SendBadRequestError(c, "failed to parse pagination parameters")
		return
	}

	stockList, err := h.service.GetStockList(storeID, stockFilter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to to retrieve stock list")
		return
	}

	utils.SendSuccessResponseWithPagination(c, stockList, stockFilter.GetPagination())
}

func (h *StoreWarehouseHandler) GetStoreWarehouseStockById(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	stock, err := h.service.GetStockById(storeID, uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve stock")
		return
	}

	utils.SendSuccessResponse(c, stock)
}

func (h *StoreWarehouseHandler) UpdateStoreWarehouseStockById(c *gin.Context) {
	var input types.UpdateStoreStockDTO

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "failed to bind json body")
		return
	}

	stock, err := h.service.GetStockById(storeID, uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, "failed to add update stock: stock not found")
		return
	}

	err = h.service.UpdateStockById(storeID, uint(stockId), &input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update stock")
		return
	}

	action := types.UpdateStoreWarehouseStockAuditFactory(
		&data.BaseDetails{
			ID:   uint(stockId),
			Name: stock.Name,
		},
		&input, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "stock updated successfully"})
}

func (h *StoreWarehouseHandler) DeleteStoreWarehouseStockById(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	stockId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store warehouse stock id")
		return
	}

	stock, err := h.service.GetStockById(storeID, uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, "failed to add update stock: stock not found")
		return
	}

	err = h.service.DeleteStockById(storeID, uint(stockId))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete stock")
		return
	}

	action := types.DeleteStoreWarehouseStockAuditFactory(
		&data.BaseDetails{
			ID:   uint(stockId),
			Name: stock.Name,
		},
		struct{}{}, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "stock deleted successfully"})
}
