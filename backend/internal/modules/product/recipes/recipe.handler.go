package recipes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RecipeHandler struct {
	service RecipeService
}

func NewRecipeHandler(service RecipeService) *RecipeHandler {
	return &RecipeHandler{
		service: service,
	}
}

func (h *RecipeHandler) GetRecipeSteps(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product-id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(uint(productID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve recipe steps")
		return
	}

	utils.SendSuccessResponse(c, recipeSteps)
}

func (h *RecipeHandler) GetRecipeStepDetails(c *gin.Context) {
	recipeStepIDParam := c.Param("id")

	recipeStepID, err := strconv.ParseUint(recipeStepIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid recipe step ID")
		return
	}

	recipeStepDetails, err := h.service.GetRecipeStepByID(uint(recipeStepID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve recipe step details")
		return
	}

	if recipeStepDetails == nil {
		utils.SendNotFoundError(c, "Recipe step not found")
		return
	}

	utils.SendSuccessResponse(c, recipeStepDetails)
}

func (h *RecipeHandler) CreateRecipeStep(c *gin.Context) {
	var input types.CreateRecipeStepDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateRecipeStep(&input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create recipe step")
		return
	}

	utils.SendMessageWithStatus(c, "Recipe step created successfully", http.StatusCreated)
}

func (h *RecipeHandler) UpdateRecipeStep(c *gin.Context) {
	recipeStepIDParam := c.Param("id")
	recipeStepID, err := strconv.ParseUint(recipeStepIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid recipe step ID")
		return
	}

	var input types.UpdateRecipeStepDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateRecipeStep(uint(recipeStepID), &input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update recipe step")
		return
	}

	utils.SendMessageWithStatus(c, "Recipe step updated successfully", http.StatusOK)
}

func (h *RecipeHandler) DeleteRecipeStep(c *gin.Context) {
	recipeStepIDParam := c.Param("id")
	recipeStepID, err := strconv.ParseUint(recipeStepIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid recipe step ID")
		return
	}

	err = h.service.DeleteRecipeStep(uint(recipeStepID))
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete recipe step")
		return
	}

	utils.SendMessageWithStatus(c, "Recipe step deleted successfully", http.StatusOK)
}
