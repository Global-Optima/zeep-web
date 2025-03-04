package storeProducts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"go.uber.org/zap"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreProductHandler struct {
	service           StoreProductService
	productService    product.ProductService
	franchiseeService franchisees.FranchiseeService
	auditService      audit.AuditService
	logger            *zap.SugaredLogger
}

func NewStoreProductHandler(
	service StoreProductService,
	productService product.ProductService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	logger *zap.SugaredLogger,
) *StoreProductHandler {
	return &StoreProductHandler{
		service:           service,
		productService:    productService,
		franchiseeService: franchiseeService,
		auditService:      auditService,
		logger:            logger,
	}
}

func (h *StoreProductHandler) GetStoreProduct(c *gin.Context) {
	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	filter, errH := contexts.GetStoreContextFilter(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	productDetails, err := h.service.GetStoreProductById(uint(storeProductID), filter)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrNotFound):
			localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
			return
		}
	}

	if productDetails == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, productDetails)
}

func (h *StoreProductHandler) GetAvailableProductsToAdd(c *gin.Context) {
	var filter productTypes.ProductsFilterDto

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Product{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	productDetails, err := h.service.GetAvailableProductsToAdd(storeID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	if productDetails == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponseWithPagination(c, productDetails, filter.Pagination)
}

func (h *StoreProductHandler) GetStoreProductCategories(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	categories, err := h.service.GetStoreProductCategories(storeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	if categories == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, categories)
}

func (h *StoreProductHandler) GetStoreProducts(c *gin.Context) {
	var filter types.StoreProductsFilterDTO

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreProduct{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	productDetails, err := h.service.GetStoreProducts(storeID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	if productDetails == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponseWithPagination(c, productDetails, filter.Pagination)
}

func (h *StoreProductHandler) GetRecommendedStoreProducts(c *gin.Context) {
	var filter types.ExcludedStoreProductsFilterDTO
	if err := c.ShouldBindJSON(&filter); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	recommendedStoreProducts, err := h.service.GetRecommendedStoreProducts(storeID, filter.StoreProductIDs)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	if recommendedStoreProducts == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, recommendedStoreProducts)
}

func (h *StoreProductHandler) GetStoreProductSizeByID(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProductSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeProductSize, err := h.service.GetStoreProductSizeByID(storeID, uint(storeProductSizeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	if storeProductSize == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, storeProductSize)
}

func (h *StoreProductHandler) CreateStoreProduct(c *gin.Context) {
	var dto types.CreateStoreProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	existingProduct, err := h.productService.GetProductByID(dto.ProductID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	id, err := h.service.CreateStoreProduct(storeID, &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	action := types.CreateStoreProductAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: existingProduct.Name,
		},
		&dto, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StoreProduct)
}

func (h *StoreProductHandler) CreateMultipleStoreProducts(c *gin.Context) {
	var dtos []types.CreateStoreProductDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	ids, err := h.service.CreateMultipleStoreProducts(storeID, dtos)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	storeProducts, err := h.service.GetStoreProductsByStoreProductIDs(storeID, ids)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	dtoMap := make(map[uint]*types.CreateStoreProductDTO)
	for _, dto := range dtos {
		dtoCopy := dto
		dtoMap[dto.ProductID] = &dtoCopy
	}

	actions := make([]shared.AuditAction, 0, len(storeProducts))
	for _, storeProduct := range storeProducts {
		matchedDTO, exists := dtoMap[storeProduct.ProductID]
		if !exists {
			h.logger.Errorf("Не удалось сопоставить продукт с DTO для продукта ID: %d", storeProduct.ProductID)
			continue
		}

		action := types.CreateStoreProductAuditFactory(
			&data.BaseDetails{
				ID:   storeProduct.ID,
				Name: storeProduct.Name,
			},
			matchedDTO, storeID,
		)
		actions = append(actions, &action)
	}

	go func() {
		_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StoreProductMultiple)
}

func (h *StoreProductHandler) UpdateStoreProduct(c *gin.Context) {
	var dto types.UpdateStoreProductDTO

	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	existingProduct, err := h.service.GetStoreProductById(uint(storeProductID), &contexts.StoreContextFilter{StoreID: &storeID})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	err = h.service.UpdateStoreProduct(storeID, uint(storeProductID), &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	action := types.UpdateStoreProductAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeProductID),
			Name: existingProduct.Name,
		},
		&dto, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreProductUpdate)
}

func (h *StoreProductHandler) DeleteStoreProduct(c *gin.Context) {
	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreProduct)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	existingProduct, err := h.service.GetStoreProductById(uint(storeProductID), &contexts.StoreContextFilter{StoreID: &storeID})
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrNotFound):
			localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
			return
		}
	}

	action := types.DeleteStoreProductAuditFactory(
		&data.BaseDetails{
			ID:   uint(storeProductID),
			Name: existingProduct.Name,
		},
		nil, storeID,
	)

	err = h.service.DeleteStoreProduct(storeID, uint(storeProductID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreProduct)
		return
	}

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreProductDelete)
}
