package notifications

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(notification *data.EmployeeNotification) error
	CreateNotificationRecipients(recipients []data.EmployeeNotificationRecipient) error

	GetNotificationByID(notificationID, employeeID uint) (*data.EmployeeNotificationRecipient, error)
	GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]data.EmployeeNotificationRecipient, error)

	MarkNotificationAsRead(notificationID uint, employeeID uint) error
	MarkMultipleNotificationsAsRead(employeeID uint, notificationIDs []uint) error

	DeleteNotification(notificationID uint) error
	GetRecipientsForEvent(eventType data.NotificationEventType) ([]data.Employee, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(notification *data.EmployeeNotification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) CreateNotificationRecipients(recipients []data.EmployeeNotificationRecipient) error {
	return r.db.Create(&recipients).Error
}

func (r *notificationRepository) GetNotificationByID(notificationID, employeeID uint) (*data.EmployeeNotificationRecipient, error) {
	var notification data.EmployeeNotificationRecipient
	err := r.db.Preload("Notification").Preload("Employee").
		First(&notification, "notification_id = ? AND employee_id = ?", notificationID, employeeID).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]data.EmployeeNotificationRecipient, error) {
	var notifications []data.EmployeeNotificationRecipient

	query := r.db.Model(&data.EmployeeNotificationRecipient{}).
		Preload("Notification").Preload("Employee").
		Joins("JOIN employee_notifications n ON n.id = employee_notification_recipients.notification_id").
		Where("employee_notification_recipients.employee_id = ?", employeeID)

	if filter.Priority != nil {
		query = query.Where("n.priority = ?", *filter.Priority)
	}

	if filter.EventType != nil {
		query = query.Where("n.event_type = ?", *filter.EventType)
	}

	if filter.IsRead != nil {
		query = query.Where("employee_notification_recipients.is_read = ?", *filter.IsRead)
	}

	if filter.StartDate != nil {
		query = query.Where("n.created_at >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("n.created_at <= ?", *filter.EndDate)
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("n.details::text ILIKE ?", searchTerm)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.BaseFilter.Pagination, filter.BaseFilter.Sort, &data.EmployeeNotificationRecipient{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) MarkNotificationAsRead(notificationID uint, employeeID uint) error {
	return r.db.Model(&data.EmployeeNotificationRecipient{}).
		Where("notification_id = ? AND employee_id = ?", notificationID, employeeID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkMultipleNotificationsAsRead(employeeID uint, notificationIDs []uint) error {
	return r.db.Model(&data.EmployeeNotificationRecipient{}).
		Where("employee_id = ? AND notification_id IN ?", employeeID, notificationIDs).
		Update("is_read", true).Error
}

func (r *notificationRepository) DeleteNotification(notificationID uint) error {
	return r.db.Delete(&data.EmployeeNotificationRecipient{}, "notification_id = ?", notificationID).Error
}

func (r *notificationRepository) GetRecipientsForEvent(eventType data.NotificationEventType) ([]data.Employee, error) {
	mapping := shared.NotificationRoleMappingManagerInstance.GetMappingByEventType(eventType)
	if mapping == nil {
		return nil, nil
	}

	var employees []data.Employee
	err := r.db.Where("type IN ? AND role IN ?", mapping.EmployeeTypes, mapping.EmployeeRoles).Find(&employees).Error
	return employees, err
}
