package shared

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

var NotificationRoleMappingManagerInstance *NotificationRoleMappingManager

type NotificationRoleMapping struct {
	EventType     data.NotificationEventType
	EmployeeTypes []data.EmployeeType
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
					EmployeeTypes: []data.EmployeeType{data.WarehouseEmployeeType, data.StoreEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee, data.RoleAdmin, data.RoleStoreManager},
				},
				{
					EventType:     data.NEW_ORDER,
					EmployeeTypes: []data.EmployeeType{data.StoreEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleStoreManager, data.RoleBarista},
				},
				{
					EventType:     data.STORE_WAREHOUSE_RUN_OUT,
					EmployeeTypes: []data.EmployeeType{data.StoreEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager},
				},
				{
					EventType:     data.CENTRAL_CATALOG_UPDATE,
					EmployeeTypes: []data.EmployeeType{data.StoreEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleAdmin},
				},
				{
					EventType:     data.WAREHOUSE_STOCK_EXPIRATION,
					EmployeeTypes: []data.EmployeeType{data.WarehouseEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee},
				},
				{
					EventType:     data.WAREHOUSE_OUT_OF_STOCK,
					EmployeeTypes: []data.EmployeeType{data.WarehouseEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleAdmin, data.RoleWarehouseManager, data.RoleWarehouseEmployee},
				},
				{
					EventType:     data.NEW_STOCK_REQUEST,
					EmployeeTypes: []data.EmployeeType{data.WarehouseEmployeeType, data.StoreEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleWarehouseManager, data.RoleWarehouseEmployee},
				},
				{
					EventType:     data.PRICE_CHANGE,
					EmployeeTypes: []data.EmployeeType{data.StoreEmployeeType},
					EmployeeRoles: []data.EmployeeRole{data.RoleStoreManager, data.RoleAdmin},
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
