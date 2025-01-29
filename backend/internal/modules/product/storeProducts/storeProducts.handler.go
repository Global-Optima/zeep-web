package storeProducts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
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
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StoreProductID")
		return
	}

	productDetails, err := h.service.GetStoreProductById(storeID, uint(storeProductID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve store product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "store product not found")
		return
	}

	utils.SendSuccessResponse(c, productDetails)
}

func (h *StoreProductHandler) GetProductsListToAdd(c *gin.Context) {
	var filter productTypes.ProductsFilterDto

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Product{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	productDetails, err := h.service.GetProductsListToAdd(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products list to add for store")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "products not found")
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
		utils.SendInternalServerError(c, "Failed to retrieve store product categories")
		return
	}

	if categories == nil {
		utils.SendNotFoundError(c, "store product categories not found")
		return
	}

	utils.SendSuccessResponse(c, categories)
}

func (h *StoreProductHandler) GetStoreProducts(c *gin.Context) {
	var filter types.StoreProductsFilterDTO

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreProduct{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	productDetails, err := h.service.GetStoreProducts(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve store product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "store product not found")
		return
	}

	utils.SendSuccessResponseWithPagination(c, productDetails, filter.Pagination)
}

func (h *StoreProductHandler) GetStoreProductSizeByID(c *gin.Context) {
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProductSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StoreProductSizeID")
		return
	}

	storeProductSize, err := h.service.GetStoreProductSizeByID(storeID, uint(storeProductSizeID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve store product size details")
		return
	}

	if storeProductSize == nil {
		utils.SendNotFoundError(c, "store product size not found")
		return
	}

	utils.SendSuccessResponse(c, storeProductSize)
}

func (h *StoreProductHandler) CreateStoreProduct(c *gin.Context) {
	var dto types.CreateStoreProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	existingProduct, err := h.productService.GetProductByID(dto.ProductID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store product: product not found")
		return
	}

	id, err := h.service.CreateStoreProduct(storeID, &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store product")
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

	utils.SendMessageWithStatus(c, "store product created successfully", http.StatusCreated)
}

func (h *StoreProductHandler) CreateMultipleStoreProducts(c *gin.Context) {
	var dtos []types.CreateStoreProductDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	dtoLength := len(dtos)
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	ids, err := h.service.CreateMultipleStoreProducts(storeID, dtos)
	if err != nil {
		msg := fmt.Sprintf("failed to create %d store products", dtoLength)
		utils.SendInternalServerError(c, msg)
		return
	}

	storeProducts, err := h.service.GetStoreProductsByStoreProductIDs(storeID, ids)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve created store products: product not found")
		return
	}

	dtoMap := make(map[uint]*types.CreateStoreProductDTO)
	for _, dto := range dtos {
		dtoCopy := dto
		dtoMap[dto.ProductID] = &dtoCopy
	}

	actions := make([]shared.AuditAction, 0, len(storeProducts))
	for _, product := range storeProducts {
		matchedDTO, exists := dtoMap[product.ProductID]
		if !exists {
			h.logger.Errorf("Failed to match product with DTO for product ID: %d", product.ProductID)
			continue
		}

		action := types.CreateStoreProductAuditFactory(
			&data.BaseDetails{
				ID:   product.ID,
				Name: product.Name,
			},
			matchedDTO, storeID,
		)
		actions = append(actions, &action)
	}

	go func() {
		_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
	}()

	msg := fmt.Sprintf("%d store product(s) created successfully", dtoLength)
	utils.SendMessageWithStatus(c, msg, http.StatusCreated)
}

func (h *StoreProductHandler) UpdateStoreProduct(c *gin.Context) {
	var dto types.UpdateStoreProductDTO

	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StoreProductID")
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	existingProduct, err := h.service.GetStoreProductById(storeID, uint(storeProductID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store product: product not found")
		return
	}

	err = h.service.UpdateStoreProduct(storeID, uint(storeProductID), &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update store product")
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

	utils.SendMessageWithStatus(c, "store product updated successfully", http.StatusCreated)
}

func (h *StoreProductHandler) DeleteStoreProduct(c *gin.Context) {
	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, types.ErrInvalidStoreProductID.Error())
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	existingProduct, err := h.service.GetStoreProductById(storeID, uint(storeProductID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store product: product not found")
		return
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
		utils.SendInternalServerError(c, "failed to delete store product")
		return
	}

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "store product deleted successfully", http.StatusCreated)
}
