package product

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

	"github.com/Global-Optima/zeep-web/backend/internal/data"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service      ProductService
	auditService audit.AuditService
}

func NewProductHandler(service ProductService, auditService audit.AuditService) *ProductHandler {
	return &ProductHandler{
		service:      service,
		auditService: auditService,
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
		utils.SendInternalServerError(c, "Не удалось получить список продуктов")
		return
	}

	utils.SendSuccessResponseWithPagination(c, products, filter.Pagination)
}

func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID продукта")
		return
	}

	productDetails, err := h.service.GetProductByID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось получить данные продукта")
		return
	}

	if productDetails == nil {
		utils.SendNotFoundError(c, "Продукт не найден")
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

	id, err := h.service.CreateProduct(&input)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось создать продукт")
		return
	}

	action := types.CreateProductAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: input.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "Продукт успешно создан", http.StatusCreated)
}

func (h *ProductHandler) GetProductSizesByProductID(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID продукта")
		return
	}

	productSizes, err := h.service.GetProductSizesByProductID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось получить размеры продукта")
		return
	}

	utils.SendSuccessResponse(c, productSizes)
}

func (h *ProductHandler) GetProductSizeByID(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID продукта")
		return
	}

	productSize, err := h.service.GetProductSizeDetailsByID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось получить размер продукта")
		return
	}

	utils.SendSuccessResponse(c, productSize)
}

func (h *ProductHandler) CreateProductSize(c *gin.Context) {
	var input types.CreateProductSizeDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	id, err := h.service.CreateProductSize(&input)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось создать размер продукта")
		return
	}

	action := types.CreateProductSizeAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: input.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "Размер продукта успешно создан", http.StatusCreated)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID продукта")
		return
	}

	var input *types.UpdateProductDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось обновить данные продукта: продукт не найден")
		return
	}

	err = h.service.UpdateProduct(uint(productID), input)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось обновить данные продукта")
		return
	}

	action := types.UpdateProductAuditFactory(
		&data.BaseDetails{
			ID:   uint(productID),
			Name: existingProduct.Name,
		},
		input,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "Продукт успешно обновлен", http.StatusOK)
}

func (h *ProductHandler) UpdateProductSize(c *gin.Context) {
	productSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID размера продукта")
		return
	}

	var input *types.UpdateProductSizeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	existingProductSize, err := h.service.GetProductSizeDetailsByID(uint(productSizeID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось обновить размер продукта: размер не найден")
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendBadRequestError(c, "Отсутствует ID магазина")
		return
	}

	err = h.service.UpdateProductSize(storeID, uint(productSizeID), input)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось обновить размер продукта")
		return
	}

	action := types.UpdateProductSizeAuditFactory(
		&data.BaseDetails{
			ID:   uint(productSizeID),
			Name: existingProductSize.Name,
		},
		input,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "Размер продукта успешно обновлен", http.StatusOK)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID продукта")
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось удалить продукт: продукт не найден")
		return
	}

	err = h.service.DeleteProduct(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось удалить продукт")
		return
	}

	action := types.DeleteProductAuditFactory(
		&data.BaseDetails{
			ID:   uint(productID),
			Name: existingProduct.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "Продукт успешно удален", http.StatusOK)
}

func (h *ProductHandler) DeleteProductSize(c *gin.Context) {
	productSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID размера продукта")
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(productSizeID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось удалить размер продукта: размер не найден")
		return
	}

	err = h.service.DeleteProduct(uint(productSizeID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось удалить размер продукта")
		return
	}

	action := types.DeleteProductSizeAuditFactory(
		&data.BaseDetails{
			ID:   uint(productSizeID),
			Name: existingProduct.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendMessageWithStatus(c, "Размер продукта успешно удален", http.StatusOK)
}
