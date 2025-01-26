package notifications

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(service NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	notification, err := h.service.GetNotificationByID(id)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve notification")
		return
	}

	utils.SendSuccessResponse(c, notification)
}

func (h *NotificationHandler) GetNotificationsByEmployee(c *gin.Context) {
	employeeID, err := contexts.GetEmployeeIDFromCtx(c)
	if err != nil {
		utils.SendMessageWithStatus(c, "Employee ID not found in context", 401)
		return
	}

	var filter types.GetNotificationsFilter
	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.EmployeeNotification{}); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	notifications, err := h.service.GetNotificationsByEmployee(employeeID, filter)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to retrieve notifications")
		return
	}

	utils.SendSuccessResponseWithPagination(c, notifications, filter.Pagination)
}

func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	employeeID, err := contexts.GetEmployeeIDFromCtx(c)
	if err != nil {
		utils.SendMessageWithStatus(c, "Employee ID not found in context", 401)
		return
	}

	var dto types.MarkNotificationAsReadDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.MarkNotificationAsRead(dto.NotificationID, employeeID)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to mark notification as read")
		return
	}

	utils.SendSuccessResponse(c, "Notification marked as read")
}

func (h *NotificationHandler) MarkMultipleNotificationsAsRead(c *gin.Context) {
	employeeID, err := contexts.GetEmployeeIDFromCtx(c)
	if err != nil {
		utils.SendMessageWithStatus(c, "Employee ID not found in context", 401)
		return
	}

	var dto types.MarkNotificationsAsReadDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.MarkMultipleNotificationsAsRead(employeeID, dto.NotificationIDs)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to mark notifications as read")
		return
	}

	utils.SendSuccessResponse(c, "Notifications marked as read")
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := utils.ParseParam(c, "id")
	if err != nil {
		utils.SendBadRequestError(c, "Invalid ID")
		return
	}

	err = h.service.DeleteNotification(id)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to delete notification")
		return
	}

	utils.SendSuccessResponse(c, "Notification deleted successfully")
}
