package analytics

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/analytics/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service AnalyticsService
}

func NewAnalyticsHandler(service AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	var filter types.AnalyticsFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid filter parameters")
		return
	}

	summary, err := h.service.GetSummary(&filter.StartDate, &filter.EndDate, &filter.StoreID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch summary")
		return
	}

	utils.SendSuccessResponse(c, summary)
}

func (h *AnalyticsHandler) GetSalesByMonth(c *gin.Context) {
	var filter types.AnalyticsFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid filter parameters")
		return
	}

	data, err := h.service.GetSalesByMonth(&filter.StartDate, &filter.EndDate, &filter.StoreID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch sales by month")
		return
	}

	utils.SendSuccessResponse(c, data)
}

func (h *AnalyticsHandler) GetPopularProducts(c *gin.Context) {
	var filter types.AnalyticsFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendBadRequestError(c, "Invalid filter parameters")
		logger.GetZapSugaredLogger().Errorln(err.Error())
		return
	}

	data, err := h.service.GetPopularProducts(&filter.StartDate, &filter.EndDate, &filter.StoreID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch popular products")
		return
	}

	utils.SendSuccessResponse(c, data)
}
