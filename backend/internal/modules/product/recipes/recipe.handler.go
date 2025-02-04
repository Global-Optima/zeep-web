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
		utils.SendInternalServerError(c, "Не удалось получить шаги рецепта")
		return
	}

	utils.SendSuccessResponse(c, recipeSteps)
}

func (h *RecipeHandler) GetRecipeStepDetails(c *gin.Context) {
	recipeStepIDParam := c.Param("id")

	recipeStepID, err := strconv.ParseUint(recipeStepIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "Неверный ID шага рецепта")
		return
	}

	recipeStepDetails, err := h.service.GetRecipeStepByID(uint(recipeStepID))
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось получить детали шага рецепта")
		return
	}

	if recipeStepDetails == nil {
		utils.SendNotFoundError(c, "Шаг рецепта не найден")
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
		utils.SendInternalServerError(c, "Не удалось создать шаги рецепта")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось создать шаги рецепта: шаг рецепта не найден")
		return
	}

	dtoMap := make(map[int]*types.CreateOrReplaceRecipeStepDTO)
	for _, dto := range input {
		dtoCopy := dto
		dtoMap[dto.Step] = &dtoCopy
	}

	var actions []shared.AuditAction

	for _, recipeStep := range recipeSteps {
		matchedDTO, exists := dtoMap[recipeStep.Step]
		if !exists {
			utils.SendInternalServerError(c, fmt.Sprintf("Не найдено соответствие для шага %d", recipeStep.Step))
			return
		}

		action := types.CreateRecipeStepsAuditFactory(
			&data.BaseDetails{
				ID:   recipeStep.ID,
				Name: recipeStep.Name,
			},
			matchedDTO,
		)
		actions = append(actions, &action)
	}

	if len(actions) > 0 {
		go func() {
			_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
		}()
	}

	utils.SendMessageWithStatus(c, "Шаги рецепта успешно созданы", http.StatusCreated)
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
		utils.SendInternalServerError(c, "Не удалось обновить шаг рецепта")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось создать шаги рецепта: шаг рецепта не найден")
		return
	}

	dtoMap := make(map[int]*types.CreateOrReplaceRecipeStepDTO)
	for _, dto := range input {
		dtoCopy := dto
		dtoMap[dto.Step] = &dtoCopy
	}

	actions := make([]shared.AuditAction, len(recipeSteps))

	for i, recipeStep := range recipeSteps {
		matchedDTO, exists := dtoMap[recipeStep.Step]
		if !exists {
			utils.SendInternalServerError(c, fmt.Sprintf("Не найдено соответствие для шага %d", recipeStep.Step))
			return
		}

		action := types.UpdateRecipeStepsAuditFactory(
			&data.BaseDetails{
				ID:   recipeStep.ID,
				Name: recipeStep.Name,
			},
			matchedDTO,
		)
		actions[i] = &action
	}

	go func() {
		_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
	}()

	utils.SendMessageWithStatus(c, "Шаг рецепта успешно обновлен", http.StatusOK)
}

func (h *RecipeHandler) DeleteRecipeSteps(c *gin.Context) {
	productID, errH := getProductID(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := h.service.DeleteRecipeStep(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось удалить шаг рецепта")
		return
	}

	recipeSteps, err := h.service.GetRecipeStepsByProductID(productID)
	if err != nil {
		utils.SendInternalServerError(c, "Не удалось создать шаги рецепта: шаг рецепта не найден")
		return
	}

	actions := make([]shared.AuditAction, len(recipeSteps))

	for i, recipeStep := range recipeSteps {
		action := types.DeleteRecipeStepsAuditFactory(
			&data.BaseDetails{
				ID:   recipeStep.ID,
				Name: recipeStep.Name,
			},
		)
		actions[i] = &action
	}

	go func() {
		_ = h.auditService.RecordMultipleEmployeeActions(c, actions)
	}()

	utils.SendMessageWithStatus(c, "Шаг рецепта успешно удален", http.StatusOK)
}

func getProductID(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	productID, err := strconv.ParseUint(c.Param("product-id"), 10, 64)
	if err != nil || productID == 0 {
		return 0, types.ErrInvalidProductID
	}
	return uint(productID), nil
}
