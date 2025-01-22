package recipes

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RecipeHandler struct {
	service      RecipeService
	auditService audit.AuditService
}

func NewRecipeHandler(service RecipeService, auditService audit.AuditService) *RecipeHandler {
	return &RecipeHandler{
		service:      service,
		auditService: auditService,
	}
}

func (h *RecipeHandler) GetRecipeSteps(c *gin.Context) {
	productID, errH := getProductID(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
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

func (h *RecipeHandler) CreateRecipeSteps(c *gin.Context) {
	var input []types.CreateOrReplaceRecipeStepDTO
	productID, errH := getProductID(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateOrReplaceRecipeSteps(productID, input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create recipe steps")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create recipe steps: recipe step not found")
		return
	}

	actions := make([]shared.AuditAction, len(recipeSteps))

	for i, recipeStep := range recipeSteps {
		var matchedDTO *types.CreateOrReplaceRecipeStepDTO
		for _, dto := range input {
			if recipeStep.Step == dto.Step {
				matchedDTO = &dto
				break
			}
		}

		if matchedDTO == nil {
			utils.SendInternalServerError(c, fmt.Sprintf("No matching DTO found for step %d", recipeStep.Step))
			return
		}

		action := types.CreateProductAuditFactory(
			&data.BaseDetails{
				ID:   recipeStep.ID,
				Name: recipeStep.Name,
			},
			matchedDTO,
		)
		actions[i] = &action
	}

	_ = h.auditService.RecordMultipleEmployeeActions(c, actions)

	utils.SendMessageWithStatus(c, "Recipe steps created successfully", http.StatusCreated)
}

func (h *RecipeHandler) UpdateRecipeSteps(c *gin.Context) {
	productID, errH := getProductID(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input []types.CreateOrReplaceRecipeStepDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateOrReplaceRecipeSteps(productID, input)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to update recipe step")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create recipe steps: recipe step not found")
		return
	}

	actions := make([]shared.AuditAction, len(recipeSteps))

	for i, recipeStep := range recipeSteps {
		var matchedDTO *types.CreateOrReplaceRecipeStepDTO
		for _, dto := range input {
			if recipeStep.Step == dto.Step {
				matchedDTO = &dto
				break
			}
		}

		if matchedDTO == nil {
			utils.SendInternalServerError(c, fmt.Sprintf("No matching DTO found for step %d", recipeStep.Step))
			return
		}

		action := types.UpdateProductAuditFactory(
			&data.BaseDetails{
				ID:   recipeStep.ID,
				Name: recipeStep.Name,
			},
			matchedDTO,
		)
		actions[i] = &action
	}

	_ = h.auditService.RecordMultipleEmployeeActions(c, actions)

	utils.SendMessageWithStatus(c, "Recipe step updated successfully", http.StatusOK)
}

func (h *RecipeHandler) DeleteRecipeSteps(c *gin.Context) {
	productID, errH := getProductID(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := h.service.DeleteRecipeStep(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete recipe step")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to create recipe steps: recipe step not found")
		return
	}

	actions := make([]shared.AuditAction, len(recipeSteps))

	for i, recipeStep := range recipeSteps {
		action := types.DeleteProductAuditFactory(
			&data.BaseDetails{
				ID:   recipeStep.ID,
				Name: recipeStep.Name,
			},
		)
		actions[i] = &action
	}

	_ = h.auditService.RecordMultipleEmployeeActions(c, actions)

	utils.SendMessageWithStatus(c, "Recipe step deleted successfully", http.StatusOK)
}

func getProductID(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	productID, err := strconv.ParseUint(c.Param("product-id"), 10, 64)
	if err != nil || productID == 0 {
		return 0, types.ErrInvalidProductID
	}
	return uint(productID), nil
}
