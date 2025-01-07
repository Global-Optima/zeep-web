package product

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"net/http"
	"strconv"

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
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	products, err := h.service.GetProducts(&filter)

	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve products")
		return
	}

	utils.SendSuccessResponseWithPagination(c, products, filter.Pagination)
}

func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	productDetails, err := h.service.GetProductByID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "Product not found")
		return
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

func (h *ProductHandler) GetProductSizesByProductID(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	productSizes, err := h.service.GetProductSizesByProductID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product sizes")
		return
	}

	utils.SendSuccessResponse(c, productSizes)
}

func (h *ProductHandler) GetProductSizeByID(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	productSize, err := h.service.GetProductSizeDetailsByID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product sizes")
		return
	}

	utils.SendSuccessResponse(c, productSize)
}

func (h *ProductHandler) CreateProductSize(c *gin.Context) {
	var input types.CreateProductSizeDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON+err.Error())
		return
	}

	_, err := h.service.CreateProductSize(&input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendMessageWithStatus(c, "product created successfully", http.StatusCreated)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	var input *types.UpdateProductDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateProduct(uint(productID), input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve product details")
		return
	}

	utils.SendMessageWithStatus(c, "product updated successfully", http.StatusOK)
}

func (h *ProductHandler) UpdateProductSize(c *gin.Context) {
	productSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	var input *types.UpdateProductSizeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateProductSize(uint(productSizeID), input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update product")
		return
	}

	utils.SendMessageWithStatus(c, "product updated successfully", http.StatusOK)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete product")
		return
	}

	utils.SendMessageWithStatus(c, "product deleted successfully", http.StatusOK)
}
