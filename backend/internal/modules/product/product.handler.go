package product

import (
	"strconv"

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

func (h *ProductHandler) GetStoreProducts(c *gin.Context) {
	storeIDParam := c.Query("storeId")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid storeId")
		return
	}

	categoryIDParam := c.Query("categoryId")
	var categoryID *uint
	if categoryIDParam != "" {
		catID, err := strconv.ParseUint(categoryIDParam, 10, 64)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid categoryId")
			return
		}
		temp := uint(catID)
		categoryID = &temp
	}

	searchQuery := c.Query("search")

	limit, offset := utils.ParsePaginationParams(c)

	filter := types.ProductFilterDao{
		StoreID:     uint(storeID),
		CategoryID:  categoryID,
		SearchQuery: searchQuery,
		Limit:       limit,
		Offset:      offset,
	}

	products, err := h.service.GetStoreProducts(filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products")
		return
	}

	utils.SuccessResponse(c, products)
}

func (h *ProductHandler) GetStoreProductDetails(c *gin.Context) {
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

	productDetails, err := h.service.GetStoreProductDetails(uint(storeID), uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "Product not found")
		return
	}

	utils.SuccessResponse(c, productDetails)
}
