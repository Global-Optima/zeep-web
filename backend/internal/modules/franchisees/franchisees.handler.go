package franchisees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FranchiseeHandler struct {
	service FranchiseeService
}

func NewFranchiseeHandler(service FranchiseeService) *FranchiseeHandler {
	return &FranchiseeHandler{service: service}
}

func (h *FranchiseeHandler) Create(c *gin.Context) {
	var input types.CreateFranchiseeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}
	_, err := h.service.Create(&input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create franchisee")
		return
	}
	utils.SuccessCreatedResponse(c, gin.H{
		"message": "franchisee was created successfully",
	})
}

func (h *FranchiseeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	var input types.UpdateFranchiseeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	if err := h.service.Update(uint(id), &input); err != nil {
		utils.SendInternalServerError(c, "failed to update franchisee")
		return
	}
	utils.SendSuccessResponse(c, gin.H{"message": "franchisee updated successfully"})
}

func (h *FranchiseeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, "failed to delete franchisee")
		return
	}
	utils.SendSuccessResponse(c, gin.H{"message": "franchisee deleted successfully"})
}

func (h *FranchiseeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	franchisee, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee")
		return
	}
	utils.SendSuccessResponse(c, franchisee)
}

func (h *FranchiseeHandler) GetAll(c *gin.Context) {
	var filter types.FranchiseeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "invalid filter parameters")
		return
	}

	franchisees, err := h.service.GetAll(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisees")
		return
	}
	utils.SendSuccessResponse(c, franchisees)
}
