package storeProducts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type StoreProductHandler struct {
	service StoreProductService
}

func NewStoreProductHandler(service StoreProductService) *StoreProductHandler {
	return &StoreProductHandler{
		service: service,
	}
}

func (h *StoreProductHandler) GetStoreProduct(c *gin.Context) {
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid StoreProductID")
		return
	}

	cacheKey := utils.BuildCacheKey("productDetails", map[string]string{
		"storeProductId": strconv.FormatUint(storeProductID, 10),
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedStoreProductDetails *types.StoreProductDetailsDTO
	if err := cacheUtil.Get(cacheKey, &cachedStoreProductDetails); err == nil {
		if !utils.IsEmpty(cachedStoreProductDetails) {
			utils.SendSuccessResponse(c, cachedStoreProductDetails)
			return
		}
	}

	productDetails, err := h.service.GetStoreProductById(storeID, uint(storeProductID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve store product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "Product not found")
		return
	}

	if err := cacheUtil.Set(cacheKey, productDetails, 10*time.Minute); err != nil {
		fmt.Printf("Failed to cache product details: %v\n", err)
	}

	utils.SendSuccessResponse(c, productDetails)
}

func (h *StoreProductHandler) GetStoreProducts(c *gin.Context) {
	var filter types.StoreProductsFilterDTO

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreProduct{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	cacheKey := utils.BuildCacheKey("productDetails", map[string]string{
		"storeId":       strconv.FormatUint(uint64(storeID), 10),
		"categoryId":    c.DefaultQuery("categoryId", ""),
		"isAvailable":   c.DefaultQuery("isAvailable", ""),
		"search":        c.DefaultQuery("search", ""),
		"page":          strconv.Itoa(filter.Pagination.Page),
		"pageSize":      strconv.Itoa(filter.Pagination.PageSize),
		"sortField":     filter.Sort.Field,
		"sortDirection": filter.Sort.Direction,
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedData utils.PaginatedData
	if err := cacheUtil.Get(cacheKey, &cachedData); err == nil {
		if !utils.IsEmpty(cachedData.Data) {
			utils.SendSuccessResponse(c, cachedData)
			return
		}
	}

	productDetails, err := h.service.GetStoreProducts(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "Product not found")
		return
	}

	cachedData.Data, cachedData.Pagination = productDetails, *filter.Pagination
	if err := cacheUtil.Set(cacheKey, cachedData, 10*time.Minute); err != nil {
		fmt.Printf("Failed to cache product details: %v\n", err)
	}

	utils.SendSuccessResponseWithPagination(c, productDetails, filter.Pagination)
}

func (h *StoreProductHandler) CreateStoreProduct(c *gin.Context) {
	var dto types.CreateStoreProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	_, err := h.service.CreateStoreProduct(storeID, &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create store product")
		return
	}
	utils.SendMessageWithStatus(c, "store product created successfully", http.StatusCreated)
}

func (h *StoreProductHandler) CreateMultipleStoreProducts(c *gin.Context) {
	var dto []types.CreateStoreProductDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	dtoLength := len(dto)
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	_, err := h.service.CreateMultipleStoreProducts(storeID, dto)
	if err != nil {
		msg := fmt.Sprintf("failed to create %d store products", dtoLength)
		utils.SendInternalServerError(c, msg)
		return
	}
	msg := fmt.Sprintf("%d store product created successfully", dtoLength)
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

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err = h.service.UpdateStoreProduct(storeID, uint(storeProductID), &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update store product")
		return
	}

	cacheKey := utils.BuildCacheKey("productDetails", map[string]string{
		"storeProductId": strconv.FormatUint(storeProductID, 10),
	})

	cacheUtil := utils.GetCacheInstance()
	if err := cacheUtil.Delete(cacheKey); err != nil {
		fmt.Printf("Failed to clear cache product details: %v\n", err)
	}

	utils.SendMessageWithStatus(c, "store product updated successfully", http.StatusCreated)
}

func (h *StoreProductHandler) DeleteStoreProduct(c *gin.Context) {
	storeProductID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, types.ErrInvalidStoreProductID.Error())
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err = h.service.DeleteStoreProduct(storeID, uint(storeProductID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete store product")
		return
	}

	cacheKey := utils.BuildCacheKey("productDetails", map[string]string{
		"storeProductId": strconv.FormatUint(storeProductID, 10),
	})

	cacheUtil := utils.GetCacheInstance()
	if err := cacheUtil.Delete(cacheKey); err != nil {
		fmt.Printf("Failed to clear cache product details: %v\n", err)
	}

	utils.SendMessageWithStatus(c, "store product deleted successfully", http.StatusCreated)
}
