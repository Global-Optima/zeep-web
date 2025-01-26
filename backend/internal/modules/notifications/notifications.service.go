package notifications

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/roles"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/types"
	"gorm.io/gorm"
)

type NotificationService interface {
	NotifyStockRequestStatusUpdated(details details.NotificationDetails) error
	NotifyNewOrder(details details.NotificationDetails) error
	NotifyStoreWarehouseRunOut(details details.NotificationDetails) error
	NotifyCentralCatalogUpdate(details details.NotificationDetails) error
	NotifyStockExpiration(details details.NotificationDetails) error
	NotifyOutOfStock(details details.NotificationDetails) error
	NotifyNewStockRequest(details details.NotificationDetails) error
	NotifyPriceChange(details details.NotificationDetails) error

	GetNotificationByID(notificationID uint) (*data.EmployeeNotificationRecipient, error)
	GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]data.EmployeeNotificationRecipient, error)
	MarkNotificationAsRead(notificationID uint, employeeID uint) error
	MarkMultipleNotificationsAsRead(employeeID uint, notificationIDs []uint) error
	DeleteNotification(notificationID uint) error
}

type notificationService struct {
	db          *gorm.DB
	repo        NotificationRepository
	roleManager *roles.NotificationRoleMappingManager
}

func NewNotificationService(db *gorm.DB, repo NotificationRepository, roleManager *roles.NotificationRoleMappingManager) NotificationService {
	return &notificationService{
		db:          db,
		repo:        repo,
		roleManager: roleManager,
	}
}

func (s *notificationService) createNotification(eventType data.NotificationEventType, priority data.NotificationPriority, message string, detailsJSON []byte) error {
	employees, err := s.repo.GetRecipientsForEvent(eventType)
	if err != nil {
		return fmt.Errorf("failed to fetch recipients: %w", err)
	}

	notification := &data.EmployeeNotification{
		EventType: eventType,
		Priority:  priority,
		Message:   message,
		Details:   detailsJSON,
	}
	if err := s.repo.CreateNotification(notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
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

func (s *notificationService) NotifyStockRequestStatusUpdated(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string

	return s.createNotification(data.STOCK_REQUEST_STATUS_UPDATED, data.MEDIUM, message, notificationDetails)
}

func (s *notificationService) NotifyNewOrder(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.NEW_ORDER, data.LOW, message, notificationDetails)
}

func (s *notificationService) NotifyStoreWarehouseRunOut(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.STORE_WAREHOUSE_RUN_OUT, data.HIGH, message, notificationDetails)
}

func (s *notificationService) NotifyCentralCatalogUpdate(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.CENTRAL_CATALOG_UPDATE, data.MEDIUM, message, notificationDetails)
}

func (s *notificationService) NotifyStockExpiration(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.STOCK_EXPIRATION, data.MEDIUM, message, notificationDetails)
}

func (s *notificationService) NotifyOutOfStock(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.OUT_OF_STOCK, data.HIGH, message, notificationDetails)
}

func (s *notificationService) NotifyNewStockRequest(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.NEW_STOCK_REQUEST, data.MEDIUM, message, notificationDetails)
}

func (s *notificationService) NotifyPriceChange(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	var message string
	return s.createNotification(data.PRICE_CHANGE, data.MEDIUM, message, notificationDetails)
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
