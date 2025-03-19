import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { EmployeeDTO } from './employees.models'

export enum EmployeeAuditOperationType {
	CREATE = 'CREATE',
	UPDATE = 'UPDATE',
	DELETE = 'DELETE',
}

export enum EmployeeAuditComponentName {
	FRANCHISEE = 'FRANCHISEE',
	REGION = 'REGION',
	PRODUCT = 'PRODUCT',
	PRODUCT_CATEGORY = 'PRODUCT_CATEGORY',
	STORE_PRODUCT = 'STORE_PRODUCT',
	EMPLOYEE = 'EMPLOYEE',
	STORE_EMPLOYEE = 'STORE_EMPLOYEE',
	WAREHOUSE_EMPLOYEE = 'WAREHOUSE_EMPLOYEE',
	FRANCHISEE_EMPLOYEE = 'FRANCHISEE_EMPLOYEE',
	REGION_EMPLOYEE = 'REGION_EMPLOYEE',
	ADMIN_EMPLOYEE = 'ADMIN_EMPLOYEE',
	ADDITIVE = 'ADDITIVE',
	ADDITIVE_CATEGORY = 'ADDITIVE_CATEGORY',
	STORE_ADDITIVE = 'STORE_ADDITIVE',
	PRODUCT_SIZE = 'PRODUCT_SIZE',
	RECIPE_STEPS = 'RECIPE_STEPS',
	STORE = 'STORE',
	WAREHOUSE = 'WAREHOUSE',
	STORE_STOCK = 'STORE_STOCK',
	INGREDIENT = 'INGREDIENT',
	INGREDIENT_CATEGORY = 'INGREDIENT_CATEGORY',
	STOCK_REQUEST = 'STOCK_REQUEST',
	STOCK_MATERIAL = 'STOCK_MATERIAL',
	STOCK_MATERIAL_CATEGORY = 'STOCK_MATERIAL_CATEGORY',
	WAREHOUSE_STOCK = 'WAREHOUSE_STOCK',
	SUPPLIER = 'SUPPLIER',
	UNIT = 'UNIT',
	ORDER = 'ORDER',
}

export const FORMATTED_AUDIT_COMPONENTS: Record<EmployeeAuditComponentName, string> = {
	[EmployeeAuditComponentName.FRANCHISEE]: 'Франчайзи',
	[EmployeeAuditComponentName.REGION]: 'Регион',
	[EmployeeAuditComponentName.PRODUCT]: 'Продукт',
	[EmployeeAuditComponentName.PRODUCT_CATEGORY]: 'Категория продукта',
	[EmployeeAuditComponentName.STORE_PRODUCT]: 'Продукт в кафе',
	[EmployeeAuditComponentName.EMPLOYEE]: 'Сотрудник',
	[EmployeeAuditComponentName.STORE_EMPLOYEE]: 'Сотрудник кафе',
	[EmployeeAuditComponentName.WAREHOUSE_EMPLOYEE]: 'Сотрудник склада',
	[EmployeeAuditComponentName.FRANCHISEE_EMPLOYEE]: 'Сотрудник франчайзи',
	[EmployeeAuditComponentName.REGION_EMPLOYEE]: 'Сотрудник региона',
	[EmployeeAuditComponentName.ADMIN_EMPLOYEE]: 'Администратор',
	[EmployeeAuditComponentName.ADDITIVE]: 'Добавка',
	[EmployeeAuditComponentName.ADDITIVE_CATEGORY]: 'Категория добавок',
	[EmployeeAuditComponentName.STORE_ADDITIVE]: 'Добавка в кафе',
	[EmployeeAuditComponentName.PRODUCT_SIZE]: 'Размер продукта',
	[EmployeeAuditComponentName.RECIPE_STEPS]: 'Шаги рецепта',
	[EmployeeAuditComponentName.STORE]: 'Кафе',
	[EmployeeAuditComponentName.WAREHOUSE]: 'Склад',
	[EmployeeAuditComponentName.STORE_STOCK]: 'Запасы в кафе',
	[EmployeeAuditComponentName.INGREDIENT]: 'Ингредиент',
	[EmployeeAuditComponentName.INGREDIENT_CATEGORY]: 'Категория ингредиентов',
	[EmployeeAuditComponentName.STOCK_REQUEST]: 'Запрос на запасы',
	[EmployeeAuditComponentName.STOCK_MATERIAL]: 'Складской материал',
	[EmployeeAuditComponentName.STOCK_MATERIAL_CATEGORY]: 'Категория складских материалов',
	[EmployeeAuditComponentName.WAREHOUSE_STOCK]: 'Запасы склада',
	[EmployeeAuditComponentName.SUPPLIER]: 'Поставщик',
	[EmployeeAuditComponentName.UNIT]: 'Единица измерения',
	[EmployeeAuditComponentName.ORDER]: 'Заказ',
}

export const FORMATTED_AUDIT_OPERATION: Record<EmployeeAuditOperationType, string> = {
	[EmployeeAuditOperationType.CREATE]: 'Создание',
	[EmployeeAuditOperationType.UPDATE]: 'Обновление',
	[EmployeeAuditOperationType.DELETE]: 'Удаление',
}

export interface EmployeeAuditDTO {
	id: number
	timestamp: Date
	operationType: EmployeeAuditOperationType
	componentName: EmployeeAuditComponentName
	localizedMessages: {
		en: string
		ru: string
		kk: string
	}
	ipAddress: string
	resourceUrl: string
	method: string
	details: Record<string, unknown>
	employee: EmployeeDTO
}

export interface EmployeeAuditFilter extends PaginationParams {
	minTimestamp?: Date
	maxTimestamp?: Date
	operationType?: string
	componentName?: string
	employeeId?: number
	method?: string
	search?: string
}
