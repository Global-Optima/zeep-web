package categories

import (
	"errors"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service      CategoryService
	auditService audit.AuditService
}

func NewCategoryHandler(service CategoryService, auditService audit.AuditService) *CategoryHandler {
	return &CategoryHandler{
		service:      service,
		auditService: auditService,
	}
}

// GetAllCategories godoc
// @Summary Get all product categories
// @Description Returns paginated list of product categories
// @Tags product-categories
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/product-categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	var filter types.ProductCategoriesFilterDTO

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.ProductCategory{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	categories, err := h.service.GetCategories(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	utils.SendSuccessResponseWithPagination(c, categories, filter.Pagination)
}

// GetCategoryByID godoc
// @Summary Get product category by ID
// @Description Retrieves product category details by ID
// @Tags product-categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} types.ProductCategoryDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	utils.SendSuccessResponse(c, category)
}

// CreateCategory godoc
// @Summary Create new product category
// @Description Creates a new product category
// @Tags product-categories
// @Accept json
// @Produce json
// @Param input body types.CreateProductCategoryDTO true "New category data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var dto types.CreateProductCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	id, err := h.service.CreateCategory(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	action := types.CreateProductCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201ProductCategory)
}

// UpdateCategory godoc
// @Summary Update a product category
// @Description Updates an existing product category
// @Tags product-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param input body types.UpdateProductCategoryDTO true "Updated category data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	var dto types.UpdateProductCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	existingCategory, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	err = h.service.UpdateCategory(uint(id), &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	action := types.UpdateProductCategoryAuditFactory(
		&data.BaseDetails{
			ID:   existingCategory.ID,
			Name: existingCategory.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductCategoryUpdate)
}

// DeleteCategory godoc
// @Summary Delete product category
// @Description Deletes a product category by ID
// @Tags product-categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400ProductCategory)
		return
	}

	existingCategory, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrCategoryNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404ProductCategory)
			return
		}
		if errors.Is(err, types.ErrCategoryIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409ProductCategoryDeleteInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500ProductCategory)
		return
	}

	action := types.DeleteProductCategoryAuditFactory(
		&data.BaseDetails{
			ID:   existingCategory.ID,
			Name: existingCategory.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200ProductCategoryDelete)
}
