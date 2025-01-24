package regions

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RegionHandler struct {
	service RegionService
}

func NewRegionHandler(service RegionService) *RegionHandler {
	return &RegionHandler{service: service}
}

func (h *RegionHandler) Create(c *gin.Context) {
	var input types.CreateRegionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}
	region, err := h.service.Create(&input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create region")
		return
	}
	utils.SendSuccessResponse(c, region)
}

func (h *RegionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	var input types.UpdateRegionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	if err := h.service.Update(uint(id), &input); err != nil {
		utils.SendInternalServerError(c, "failed to update region")
		return
	}
	utils.SendSuccessResponse(c, gin.H{"message": "region updated successfully"})
}

func (h *RegionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, "failed to delete region")
		return
	}
	utils.SendSuccessResponse(c, gin.H{"message": "region deleted successfully"})
}

func (h *RegionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	region, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve region")
		return
	}
	utils.SendSuccessResponse(c, region)
}

func (h *RegionHandler) GetAll(c *gin.Context) {
	var filter types.RegionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "invalid filter parameters")
		return
	}

	regions, err := h.service.GetAll(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve regions")
		return
	}
	utils.SendSuccessResponse(c, regions)
}
