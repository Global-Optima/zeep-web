package product

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter types.ProductsFilterDto
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid query parameters")
		return
	}

	if filter.StoreID != nil {
		storeID, err := strconv.ParseUint(c.Query("storeId"), 10, 64)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid storeId")
			return
		}
		temp := uint(storeID)
		filter.StoreID = &temp
	}

	if filter.CategoryID != nil {
		categoryID, err := strconv.ParseUint(c.Query("categoryId"), 10, 64)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid categoryId")
			return
		}
		temp := uint(categoryID)
		filter.CategoryID = &temp
	}

	cacheKey := utils.BuildCacheKey("storeProducts", map[string]string{
		"storeId":    c.Query("storeId"),
		"categoryId": c.Query("categoryId"),
		"search":     c.Query("search"),
		"page":       strconv.Itoa(filter.Pagination.Page),
		"pageSize":   strconv.Itoa(filter.Pagination.PageSize),
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedProducts []types.StoreProductDTO
	if err := cacheUtil.Get(cacheKey, &cachedProducts); err == nil {
		if !utils.IsEmpty(cachedProducts) {
			utils.SendSuccessResponseWithPagination(c, cachedProducts, filter.Pagination)
			return
		}
	}

	products, err := h.service.GetProducts(filter)

	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products")
		return
	}

	if err := cacheUtil.Set(cacheKey, products, 5*time.Minute); err != nil {
		fmt.Printf("Failed to cache products: %v\n", err)
	}

	utils.SendSuccessResponseWithPagination(c, products, filter.Pagination)
}

func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	storeIDParam := c.Query("storeId")
	productIDParam := c.Param("productId")

	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store ID")
		return
	}

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	cacheKey := utils.BuildCacheKey("productDetails", map[string]string{
		"storeId":   storeIDParam,
		"productId": productIDParam,
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedProductDetails *types.StoreProductDetailsDTO
	if err := cacheUtil.Get(cacheKey, &cachedProductDetails); err == nil {
		if !utils.IsEmpty(cachedProductDetails) {
			utils.SendSuccessResponse(c, cachedProductDetails)
			return
		}
	}

	productDetails, err := h.service.GetStoreProductDetails(uint(storeID), uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
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
