package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetStoreProducts(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Param("store_id"))
	if err != nil || storeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store_id parameter"})
		return
	}

	category := c.Query("category")
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	products, err := h.service.GetStoreProducts(uint(storeID), category, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) SearchStoreProducts(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Param("store_id"))
	if err != nil || storeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store_id parameter"})
		return
	}

	searchQuery := c.Query("q")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	category := c.Query("category")
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	products, err := h.service.SearchStoreProducts(uint(storeID), searchQuery, category, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetStoreProductDetails(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Param("store_id"))
	if err != nil || storeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store_id parameter"})
		return
	}

	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id parameter"})
		return
	}

	product, err := h.service.GetStoreProductDetails(uint(storeID), uint(productID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product details", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}
