package notifications

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationService interface {
	NotifyStockRequestStatusUpdated(details details.NotificationDetails) error
	NotifyNewOrder(details details.NotificationDetails) error
	NotifyStoreWarehouseRunOut(details details.NotificationDetails) error
	NotifyCentralCatalogUpdate(details details.NotificationDetails) error
	NotifyStoreStockExpiration(details details.NotificationDetails) error
	NotifyStoreProvisionExpiration(details details.NotificationDetails) error
	NotifyWarehouseStockExpiration(details details.NotificationDetails) error
	NotifyOutOfStock(details details.NotificationDetails) error
	NotifyNewStockRequest(details details.NotificationDetails) error
	NotifyPriceChange(details details.NotificationDetails) error
	NotifyNewProductAdded(details details.NotificationDetails) error
	NotifyNewProductSizeAdded(details details.NotificationDetails) error
	NotifyNewAdditiveAdded(details details.NotificationDetails) error

	GetNotificationByID(notificationID, employeeID uint) (*types.NotificationDTO, error)
	GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]types.NotificationDTO, error)
	MarkNotificationAsRead(notificationID uint, employeeID uint) error
	MarkMultipleNotificationsAsRead(employeeID uint, notificationIDs []uint) error
	DeleteNotification(notificationID uint) error
}

type notificationService struct {
	db          *gorm.DB
	repo        NotificationRepository
	roleManager *shared.NotificationRoleMappingManager
	logger      *zap.SugaredLogger
}

func NewNotificationService(db *gorm.DB, repo NotificationRepository, roleManager *shared.NotificationRoleMappingManager, logger *zap.SugaredLogger) NotificationService {
	return &notificationService{
		db:          db,
		repo:        repo,
		roleManager: roleManager,
		logger:      logger,
	}
}

// func (s *notificationService) createNotification(eventType data.NotificationEventType, priority data.NotificationPriority, detailsJSON []byte) error {
// 	employees, err := s.repo.GetRecipientsForEvent(eventType)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch recipients: %w", err)
// 	}

// 	if len(employees) == 0 {
// 		return fmt.Errorf("no recipients found for event type %s", eventType)
// 	}

// 	notification := &data.EmployeeNotification{
// 		EventType: eventType,
// 		Priority:  priority,
// 		Details:   detailsJSON,
// 	}
// 	if err := s.repo.CreateNotification(notification); err != nil {
// 		return fmt.Errorf("failed to create notification: %w", err)
// 	}

// 	recipients := make([]data.EmployeeNotificationRecipient, len(employees))
// 	for i, employee := range employees {
// 		recipients[i] = data.EmployeeNotificationRecipient{
// 			NotificationID: notification.ID,
// 			EmployeeID:     employee.ID,
// 		}
// 	}
// 	return s.repo.CreateNotificationRecipients(recipients)
// }

func (s *notificationService) createNotificationAsync(eventType data.NotificationEventType, priority data.NotificationPriority, detailsJSON []byte, baseDetails *details.BaseNotificationDetails) {
	go func() {
		if baseDetails == nil {
			s.logger.Errorf("Base details are not provided, don't know where to send notifications!")
			return
		}

		employees, err := s.repo.GetRecipientsForEvent(eventType, *baseDetails)
		if err != nil {
			s.logger.Errorf("Failed to fetch recipients for event %s: %v", eventType, err)
			return
		}

		if len(employees) == 0 {
			s.logger.Infof("No recipients found for event type %s", eventType)
			return
		}

		notification := &data.EmployeeNotification{
			EventType: eventType,
			Priority:  priority,
			Details:   detailsJSON,
		}

		if err := s.repo.CreateNotification(notification); err != nil {
			s.logger.Errorf("Failed to create notification for event type %s: %v", eventType, err)
			return
		}

		recipients := make([]data.EmployeeNotificationRecipient, len(employees))
		for i, employee := range employees {
			recipients[i] = data.EmployeeNotificationRecipient{
				NotificationID: notification.ID,
				EmployeeID:     employee.ID,
			}
		}

		if err := s.repo.CreateNotificationRecipients(recipients); err != nil {
			s.logger.Errorf("Failed to create recipients for notification %d: %v", notification.ID, err)
		}
	}()
}

// notifiers
func (s *notificationService) NotifyStockRequestStatusUpdated(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.STOCK_REQUEST_STATUS_UPDATED, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyNewOrder(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}
	s.createNotificationAsync(data.NEW_ORDER, data.LOW, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyStoreWarehouseRunOut(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.STORE_WAREHOUSE_RUN_OUT, data.HIGH, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyCentralCatalogUpdate(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.CENTRAL_CATALOG_UPDATE, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyStoreStockExpiration(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.STORE_STOCK_EXPIRATION, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyStoreProvisionExpiration(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.STORE_PROVISION_EXPIRATION, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyWarehouseStockExpiration(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.WAREHOUSE_STOCK_EXPIRATION, data.HIGH, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyOutOfStock(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.WAREHOUSE_OUT_OF_STOCK, data.HIGH, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyNewStockRequest(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.NEW_STOCK_REQUEST, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyNewProductAdded(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.NEW_PRODUCT, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyNewProductSizeAdded(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.NEW_PRODUCT_SIZE, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyNewAdditiveAdded(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.NEW_ADDITIVE, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyPriceChange(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.PRICE_CHANGE, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

func (s *notificationService) NotifyNewProductSize(details details.NotificationDetails) error {
	notificationDetails, err := details.ToDetails()
	if err != nil {
		return err
	}

	s.createNotificationAsync(data.PRICE_CHANGE, data.MEDIUM, notificationDetails, details.GetBaseDetails())
	return nil
}

// Base notification module methods
func (s *notificationService) GetNotificationByID(notificationID, employeeID uint) (*types.NotificationDTO, error) {
	employeeNotification, err := s.repo.GetNotificationByID(notificationID, employeeID)
	if err != nil {
		return nil, err
	}

	return types.ConvertNotificationToDTO(&employeeNotification.Notification, employeeNotification.IsRead)
}

func (s *notificationService) GetNotificationsByEmployee(employeeID uint, filter types.GetNotificationsFilter) ([]types.NotificationDTO, error) {
	employeeNotifications, err := s.repo.GetNotificationsByEmployee(employeeID, filter)
	if err != nil {
		return nil, err
	}

	return types.ConvertNotificationsToDTOs(employeeNotifications)
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
