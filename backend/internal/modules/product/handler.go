package product

import (
	"net/http"
	"strconv"

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
	categoryIDParam := c.Query("categoryId")
	searchQuery := c.Query("search")

	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid store Id")
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category Id")
		return
	}

	limit, offset := utils.ParsePaginationParams(c)

	products, err := h.service.GetStoreProducts(c, uint(storeID), uint(categoryID), searchQuery, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// meta := gin.H{"limit": limit, "offset": offset, "total": len(products)}

	utils.SuccessResponse(c, products)
}

func (h *ProductHandler) GetStoreProductDetails(c *gin.Context) {
	storeIDParam := c.Param("storeId")
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

	productDetails, err := h.service.GetStoreProductDetails(c, uint(storeID), uint(productID))
	if err != nil {
		utils.SendInternalError(c, "Failed to retrieve product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "Product not found")
		return
	}

	utils.SuccessResponse(c, productDetails)
}
