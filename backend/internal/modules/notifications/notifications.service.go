package notifications

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/types"
	"gorm.io/gorm"
)

type NotificationService interface {
	NotifyStockRequestStatusUpdated(details data.NotificationDetails) error
	NotifyNewOrder(details data.NotificationDetails) error
	NotifyStoreWarehouseRunOut(details data.NotificationDetails) error
	NotifyCentralCatalogUpdate(details data.NotificationDetails) error
	NotifyStockExpiration(details data.NotificationDetails) error
	NotifyOutOfStock(details data.NotificationDetails) error
	NotifyNewStockRequest(details data.NotificationDetails) error
	NotifyPriceChange(details data.NotificationDetails) error

	GetNotificationByID(notificationID uint) (*data.EmployeeNotificationRecipient, error)
	GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]data.EmployeeNotificationRecipient, error)
	MarkNotificationAsRead(notificationID uint, employeeID uint) error
	MarkMultipleNotificationsAsRead(employeeID uint, notificationIDs []uint) error
	DeleteNotification(notificationID uint) error
}

type notificationService struct {
	db          *gorm.DB
	repo        NotificationRepository
	roleManager *shared.NotificationRoleMappingManager
}

func NewNotificationService(db *gorm.DB, repo NotificationRepository, roleManager *shared.NotificationRoleMappingManager) NotificationService {
	return &notificationService{
		db:          db,
		repo:        repo,
		roleManager: roleManager,
	}
}

func (s *notificationService) createNotificationWithDetails(eventType data.NotificationEventType, priority data.NotificationPriority, details data.NotificationDetails) error {
	detailsJSON, err := details.ToDetails()
	if err != nil {
		return err
	}

	employees, err := s.roleManager.GetRecipients(s.db, eventType)
	if err != nil {
		return err
	}

	message, err := s.buildLocalizedMessage(eventType, details)
	if err != nil {
		return fmt.Errorf("failed to build localized message: %w", err)
	}

	notification := &data.EmployeeNotification{
		EventType: eventType,
		Priority:  priority,
		Message:   message.Ru, // Default to Russian
		Details:   detailsJSON,
	}
	if err := s.repo.CreateNotification(notification); err != nil {
		return err
	}

	recipients := make([]data.EmployeeNotificationRecipient, len(employees))
	for i, employee := range employees {
		recipients[i] = data.EmployeeNotificationRecipient{
			NotificationID: notification.ID,
			EmployeeID:     employee.ID,
		}
	}
	return s.repo.CreateNotificationRecipients(recipients)
}

func (s *notificationService) buildLocalizedMessage(eventType data.NotificationEventType, details data.NotificationDetails) (*localization.LocalizedMessages, error) {
	baseDetails := details.GetBaseDetails()
	messageID := localization.FormTranslationKey("notification", eventType.ToString())

	return localization.Translate(messageID, map[string]interface{}{
		"Name": baseDetails.Name,
		"ID":   baseDetails.ID,
	})
}

func (s *notificationService) NotifyStockRequestStatusUpdated(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.STOCK_REQUEST_STATUS_UPDATED, data.MEDIUM, details)
}

func (s *notificationService) NotifyNewOrder(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.NEW_ORDER, data.LOW, details)
}

func (s *notificationService) NotifyStoreWarehouseRunOut(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.STORE_WAREHOUSE_RUN_OUT, data.HIGH, details)
}

func (s *notificationService) NotifyCentralCatalogUpdate(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.CENTRAL_CATALOG_UPDATE, data.MEDIUM, details)
}

func (s *notificationService) NotifyStockExpiration(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.STOCK_EXPIRATION, data.MEDIUM, details)
}

func (s *notificationService) NotifyOutOfStock(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.OUT_OF_STOCK, data.HIGH, details)
}

func (s *notificationService) NotifyNewStockRequest(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.NEW_STOCK_REQUEST, data.MEDIUM, details)
}

func (s *notificationService) NotifyPriceChange(details data.NotificationDetails) error {
	return s.createNotificationWithDetails(data.PRICE_CHANGE, data.MEDIUM, details)
}

// Base notification module methods
func (s *notificationService) GetNotificationByID(notificationID uint) (*data.EmployeeNotificationRecipient, error) {
	return s.repo.GetNotificationByID(notificationID)
}

func (s *notificationService) GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]data.EmployeeNotificationRecipient, error) {
	return s.repo.GetNotificationsByEmployee(employeeID, filter)
}

func (s *notificationService) MarkNotificationAsRead(notificationID uint, employeeID uint) error {
	return s.repo.MarkNotificationAsRead(notificationID, employeeID)
}

func (s *notificationService) MarkMultipleNotificationsAsRead(employeeID uint, notificationIDs []uint) error {
	if len(notificationIDs) == 0 {
		return nil
	}
	return s.repo.MarkMultipleNotificationsAsRead(employeeID, notificationIDs)
}

func (s *notificationService) DeleteNotification(notificationID uint) error {
	return s.repo.DeleteNotification(notificationID)
}
