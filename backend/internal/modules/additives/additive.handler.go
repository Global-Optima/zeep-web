package additives

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AdditiveHandler struct {
	service      AdditiveService
	auditService audit.AuditService
}

func NewAdditiveHandler(service AdditiveService, auditService audit.AuditService) *AdditiveHandler {
	return &AdditiveHandler{
		service:      service,
		auditService: auditService,
	}
}

// GetAdditiveCategories godoc
// @Summary Get additive categories
// @Description Returns list of additive categories with filters and pagination
// @Tags additive-categories
// @Accept json
// @Produce json
// @Param includeEmpty query bool false "Include empty categories"
// @Param productSizeId query int false "Filter by product size ID"
// @Param isMultipleSelect query bool false "Filter by multiple select flag"
// @Param isRequired query bool false "Filter by required flag"
// @Param search query string false "Search query"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/additives/categories [get]
func (h *AdditiveHandler) GetAdditiveCategories(c *gin.Context) {
	var filter types.AdditiveCategoriesFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.AdditiveCategory{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	additives, err := h.service.GetAdditiveCategories(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

// CreateAdditiveCategory godoc
// @Summary Create additive category
// @Description Adds a new additive category
// @Tags additive-categories
// @Accept json
// @Produce json
// @Param input body types.CreateAdditiveCategoryDTO true "Category data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/additives/categories [post]
func (h *AdditiveHandler) CreateAdditiveCategory(c *gin.Context) {
	var dto types.CreateAdditiveCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateAdditiveCategory(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryCreate)
		return
	}

	action := types.CreateAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201AdditiveCategory)
}

// UpdateAdditiveCategory godoc
// @Summary Update additive category
// @Description Updates existing additive category
// @Tags additive-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param input body types.UpdateAdditiveCategoryDTO true "Updated category data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/additives/categories/{id} [put]
func (h *AdditiveHandler) UpdateAdditiveCategory(c *gin.Context) {
	var dto types.UpdateAdditiveCategoryDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryUpdate)
		return
	}

	if err := h.service.UpdateAdditiveCategory(uint(categoryID), &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryUpdate)
		return
	}

	action := types.UpdateAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(categoryID),
			Name: category.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveCategoryUpdate)
}

// DeleteAdditiveCategory godoc
// @Summary Delete additive category
// @Description Deletes an additive category by ID
// @Tags additive-categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{} "If the category is in use"
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/additives/categories/{id} [delete]
func (h *AdditiveHandler) DeleteAdditiveCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryDelete)
		return
	}

	if err := h.service.DeleteAdditiveCategory(uint(categoryID)); err != nil {
		if errors.Is(err, types.ErrAdditiveCategoryIsInUse) {
			localization.SendLocalizedResponseWithKey(c, types.Response409AdditiveCategoryDeleteInUse)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryDelete)
		return
	}

	action := types.DeleteAdditiveCategoryAuditFactory(
		&data.BaseDetails{
			ID:   uint(categoryID),
			Name: category.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveCategoryDelete)
}

// GetAdditiveCategoryByID godoc
// @Summary Get additive category details
// @Description Returns a single additive category by ID
// @Tags additive-categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} types.AdditiveCategoryDetailsDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/additives/categories/{id} [get]
func (h *AdditiveHandler) GetAdditiveCategoryByID(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdditiveCategory)
		return
	}

	category, err := h.service.GetAdditiveCategoryByID(uint(categoryID))
	if err != nil {
		switch {
		case errors.Is(err, types.ErrAdditiveCategoryNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404AdditiveCategory)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryGet)
			return
		}
	}

	utils.SendSuccessResponse(c, category)
}

// GetAdditives godoc
// @Summary Get list of additives
// @Description Retrieves all additives with optional filters and pagination
// @Tags additives
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param minPrice query number false "Minimum price"
// @Param maxPrice query number false "Maximum price"
// @Param categoryId query int false "Filter by category ID"
// @Param productSizeId query int false "Filter by product size ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/additives [get]
func (h *AdditiveHandler) GetAdditives(c *gin.Context) {
	var filter types.AdditiveFilterQuery
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Additive{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	additives, err := h.service.GetAdditives(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCategoryGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, additives, filter.Pagination)
}

// CreateAdditive godoc
// @Summary Create a new additive
// @Description Creates an additive, including ingredients and optional image
// @Tags additives
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Additive name"
// @Param description formData string false "Description"
// @Param basePrice formData number true "Base price"
// @Param size formData number true "Size"
// @Param unitId formData int true "Unit ID"
// @Param additiveCategoryId formData int true "Category ID"
// @Param machineId formData string true "Machine ID"
// @Param image formData file false "Additive image"
// @Param ingredients formData string false "JSON-encoded ingredients array"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/additives [post]
func (h *AdditiveHandler) CreateAdditive(c *gin.Context) {
	var dto types.CreateAdditiveDTO
	var err error

	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	ingredientsJSON := c.PostForm("ingredients")
	if ingredientsJSON != "" {
		err := json.Unmarshal([]byte(ingredientsJSON), &dto.Ingredients)
		if err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			return
		}

		for _, ingredient := range dto.Ingredients {
			if ingredient.IngredientID == 0 || ingredient.Quantity <= 0 {
				localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
				return
			}
		}
	}

	dto.Image, err = media.GetImageWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	id, err := h.service.CreateAdditive(&dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveCreate)
		return
	}

	action := types.CreateAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   id,
			Name: dto.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201Additive)
}

// UpdateAdditive godoc
// @Summary Update an existing additive
// @Description Updates additive fields and image
// @Tags additives
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Additive ID"
// @Param name formData string false "Name"
// @Param description formData string false "Description"
// @Param basePrice formData number false "Base price"
// @Param size formData number false "Size"
// @Param unitId formData int false "Unit ID"
// @Param additiveCategoryId formData int false "Category ID"
// @Param machineId formData string false "Machine ID"
// @Param image formData file false "Image file"
// @Param deleteImage formData bool false "Delete existing image"
// @Param ingredients formData string false "JSON-encoded ingredients array"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/additives/{id} [put]
func (h *AdditiveHandler) UpdateAdditive(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	var dto types.UpdateAdditiveDTO
	if err := utils.ParseRequestBody(c, &dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	ingredientsJSON := c.PostForm("ingredients")
	if ingredientsJSON != "" {
		err = json.Unmarshal([]byte(ingredientsJSON), &dto.Ingredients)
		if err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			return
		}
	}

	dto.Image, err = media.GetImageWithFormFile(c)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	additive, err := h.service.UpdateAdditive(uint(additiveID), &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveUpdate)
		return
	}

	action := types.UpdateAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(additiveID),
			Name: additive.Name,
		},
		&dto,
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveUpdate)
}

// DeleteAdditive godoc
// @Summary Delete an additive
// @Description Deletes an additive by ID
// @Tags additives
// @Produce json
// @Param id path int true "Additive ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/additives/{id} [delete]
func (h *AdditiveHandler) DeleteAdditive(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		if errors.Is(err, types.ErrAdditiveNotFound) {
			localization.SendLocalizedResponseWithKey(c, types.Response404Additive)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveDelete)
		return
	}

	if err := h.service.DeleteAdditive(uint(additiveID)); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveDelete)
		return
	}

	action := types.DeleteAdditiveAuditFactory(
		&data.BaseDetails{
			ID:   uint(additiveID),
			Name: additive.Name,
		},
	)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200AdditiveDelete)
}

// GetAdditiveByID godoc
// @Summary Get additive details
// @Description Returns a single additive by ID
// @Tags additives
// @Produce json
// @Param id path int true "Additive ID"
// @Success 200 {object} types.AdditiveDetailsDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/additives/{id} [get]

func (h *AdditiveHandler) GetAdditiveByID(c *gin.Context) {
	additiveID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Additive)
		return
	}

	additive, err := h.service.GetAdditiveByID(uint(additiveID))
	if err != nil {
		switch {
		case errors.Is(err, types.ErrAdditiveNotFound):
			localization.SendLocalizedResponseWithKey(c, types.Response404Additive)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500AdditiveGet)
			return
		}
	}

	utils.SendSuccessResponse(c, additive)
}
