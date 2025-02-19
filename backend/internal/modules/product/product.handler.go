package product

import (
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	products, err := h.service.GetProducts(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	utils.SendSuccessResponseWithPagination(c, products, filter.Pagination)
}

func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	productDetails, err := h.service.GetProductByID(uint(productID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	if productDetails == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, productDetails)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var dto types.CreateProductDTO

	err := c.Request.ParseMultipartForm(30 << 20)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		logrus.Info("Multipart form parsing failed: ", err)
		return
	}

	err = c.ShouldBind(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		logrus.Info(err)
		return
	}

	dto.Image, err = media.GetImageWithFormFile(c)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	dto.Video, err = media.GetVideoWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingVideo)
		return
	}

	id, err := h.service.CreateProduct(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	action := types.CreateProductAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		})

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Product)
}

func (h *ProductHandler) GetProductSizesByProductID(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	productSizes, err := h.service.GetProductSizesByProductID(uint(productID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	utils.SendSuccessResponse(c, productSizes)
}

func (h *ProductHandler) GetProductSizeByID(c *gin.Context) {
	productIDParam := c.Param("id")

	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	productSize, err := h.service.GetProductSizeDetailsByID(uint(productID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	utils.SendSuccessResponse(c, productSize)
}

func (h *ProductHandler) CreateProductSize(c *gin.Context) {
	var input types.CreateProductSizeDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateProductSize(&input)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
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

	localization.SendLocalizedResponseWithKey(c, types.Response201Product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	var dto *types.UpdateProductDTO
	if err := c.ShouldBind(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	dto.Image, err = media.GetImageWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	/*dto.Video, err = media.GetVideoWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingVideo)
		return
	}*/

	existingProduct, err := h.service.UpdateProduct(uint(productID), dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	action := types.UpdateProductAuditFactory(
		&data.BaseDetails{
			ID:   uint(productID),
			Name: existingProduct.Name,
		},
		dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductUpdate)
}

func (h *ProductHandler) UpdateProductSize(c *gin.Context) {
	productSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	var input *types.UpdateProductSizeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	existingProductSize, err := h.service.GetProductSizeDetailsByID(uint(productSizeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	err = h.service.UpdateProductSize(uint(productSizeID), input)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductUpdate)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	existingProduct, err := h.service.DeleteProduct(uint(productID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductDelete)
}

func (h *ProductHandler) DeleteProductSize(c *gin.Context) {
	productSizeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Product)
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(productSizeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
		return
	}

	err = h.service.DeleteProductSize(uint(productSizeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500Product)
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

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductDelete)
}
