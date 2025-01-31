package shared

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

var NotificationRoleMappingManagerInstance *NotificationRoleMappingManager

type NotificationRoleMapping struct {
	EventType     data.NotificationEventType
	EmployeeRoles []data.EmployeeRole
}

type NotificationRoleMappingManager struct {
	mappings []NotificationRoleMapping
}

// TODO check role mapping
func NewNotificationRoleMappingManager() *NotificationRoleMappingManager {
	if NotificationRoleMappingManagerInstance == nil {
		NotificationRoleMappingManagerInstance = &NotificationRoleMappingManager{
			mappings: []NotificationRoleMapping{
				{
					EventType:     data.STOCK_REQUEST_STATUS_UPDATED,
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee, data.RoleAdmin, data.RoleStoreManager},
				},
				{
					EventType:     data.NEW_ORDER,
					EmployeeRoles: []data.EmployeeRole{data.RoleStoreManager, data.RoleBarista},
				},
				{
					EventType:     data.STORE_WAREHOUSE_RUN_OUT,
					EmployeeRoles: []data.EmployeeRole{data.RoleStoreManager, data.RoleBarista},
				},
				{
					EventType:     data.CENTRAL_CATALOG_UPDATE,
					EmployeeRoles: []data.EmployeeRole{data.RoleAdmin},
				},
				{
					EventType:     data.WAREHOUSE_STOCK_EXPIRATION,
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee},
				},
				{
					EventType:     data.WAREHOUSE_OUT_OF_STOCK,
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee},
				},
				{
					EventType:     data.NEW_STOCK_REQUEST,
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee},
				},
				{
					EventType:     data.PRICE_CHANGE,
					EmployeeRoles: []data.EmployeeRole{data.RoleStoreManager, data.RoleBarista},
				},
			},
		}
	}
	return NotificationRoleMappingManagerInstance
}

func (m *NotificationRoleMappingManager) GetMappings() []NotificationRoleMapping {
	return m.mappings
}

func (m *NotificationRoleMappingManager) GetMappingByEventType(eventType data.NotificationEventType) *NotificationRoleMapping {
	for _, mapping := range m.mappings {
		if mapping.EventType == eventType {
			return &mapping
		}
	}
	return nil
}

func (m *NotificationRoleMappingManager) AddMapping(mapping NotificationRoleMapping) {
	m.mappings = append(m.mappings, mapping)
}

func (m *NotificationRoleMappingManager) UpdateMapping(eventType data.NotificationEventType, updatedMapping NotificationRoleMapping) bool {
	for i, mapping := range m.mappings {
		if mapping.EventType == eventType {
			m.mappings[i] = updatedMapping
			return true
		}
	}
	return false
}

func (m *NotificationRoleMappingManager) DeleteMapping(eventType data.NotificationEventType) bool {
	for i, mapping := range m.mappings {
		if mapping.EventType == eventType {
			m.mappings = append(m.mappings[:i], m.mappings[i+1:]...)
			return true
		}
	}
	return false
}
