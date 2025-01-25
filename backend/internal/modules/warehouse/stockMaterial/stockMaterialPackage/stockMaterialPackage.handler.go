package stockMaterialPackage

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StockMaterialPackageHandler struct {
	service StockMaterialPackageService
}

func NewStockMaterialPackageHandler(service StockMaterialPackageService) *StockMaterialPackageHandler {
	return &StockMaterialPackageHandler{service: service}
}

func (h *StockMaterialPackageHandler) Create(c *gin.Context) {
	var dto types.CreateStockMaterialPackageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SuccessCreatedResponse(c, gin.H{"id": id})
}

func (h *StockMaterialPackageHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	packageEntity, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, packageEntity)
}

func (h *StockMaterialPackageHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	var dto types.UpdateStockMaterialPackageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	if err := h.service.Update(uint(id), dto); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockMaterialPackageHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *StockMaterialPackageHandler) GetAll(c *gin.Context) {
	packages, err := h.service.GetAll()
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, packages)
}
