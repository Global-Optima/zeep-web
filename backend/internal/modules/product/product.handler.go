package product

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"net/http"
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
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter types.ProductsFilterDto
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Product{}); err != nil {
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
		"storeId":    c.DefaultQuery("storeId", ""),
		"categoryId": c.DefaultQuery("categoryId", ""),
		"search":     c.DefaultQuery("search", ""),
		"page":       strconv.Itoa(filter.Pagination.Page),
		"pageSize":   strconv.Itoa(filter.Pagination.PageSize),
		"totalCount": strconv.Itoa(filter.Pagination.TotalCount),
		"totalPages": strconv.Itoa(filter.Pagination.TotalPages),
	})

	cacheUtil := utils.GetCacheInstance()

	var cachedProducts []types.StoreProductDTO
	if err := cacheUtil.Get(cacheKey, &cachedProducts); err == nil {
		if !utils.IsEmpty(cachedProducts) {
			utils.SendSuccessResponseWithPagination(c, cachedProducts, filter.Pagination)
			return
		}
	}

	products, err := h.service.GetProducts(&filter)

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
	productIDParam := c.Param("id")

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

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input types.CreateProductDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateProduct(&input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendMessageWithStatus(c, "product created successfully", http.StatusCreated)
}

func (h *ProductHandler) CreateProductSize(c *gin.Context) {
	var input types.CreateProductSizeDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateProductSize(&input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendMessageWithStatus(c, "product created successfully", http.StatusCreated)
}

func (h *ProductHandler) GetStoreProductDetails(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	storeID, errH := getStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
	}

	productDetails, err := h.service.GetStoreProductDetails(storeID, uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendSuccessResponse(c, productDetails)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	var input *types.UpdateProductDTO

	err = h.service.UpdateProduct(uint(productID), input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendMessageWithStatus(c, "product updated successfully", http.StatusOK)
}

func (h *ProductHandler) UpdateProductSize(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	var input *types.UpdateProductSizeDTO

	err = h.service.UpdateProductSize(uint(productID), input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendMessageWithStatus(c, "product updated successfully", http.StatusOK)
}

// getStoreId returns the retrieved id and HandlerError
func getStoreId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, types.ErrUnauthorizedAccess
	}

	var storeID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleDirector {
		storeID = claims.WorkplaceID
	} else {
		id, err := strconv.ParseUint(c.Query("storeId"), 10, 64)
		if err != nil {
			return 0, types.ErrInvalidStoreID
		}
		storeID = uint(id)
	}

	return storeID, nil
}
