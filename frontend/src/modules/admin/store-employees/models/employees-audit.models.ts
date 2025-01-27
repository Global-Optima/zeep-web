import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { EmployeeDTO } from './employees.models'

export enum EmployeeAuditComponentName {
	PRODUCT = 'PRODUCT',
	PRODUCT_CATEGORY = 'PRODUCT_CATEGORY',
	STORE_PRODUCT = 'STORE_PRODUCT',
	EMPLOYEE = 'EMPLOYEE',
	ADDITIVE = 'ADDITIVE',
	ADDITIVE_CATEGORY = 'ADDITIVE_CATEGORY',
	STORE_ADDITIVE = 'STORE_ADDITIVE',
	PRODUCT_SIZE = 'PRODUCT_SIZE',
	RECIPE_STEPS = 'RECIPE_STEPS',
	STORE = 'STORE',
	WAREHOUSE = 'WAREHOUSE',
	STORE_WAREHOUSE_STOCK = 'STORE_WAREHOUSE_STOCK',
	INGREDIENT = 'INGREDIENT',
	INGREDIENT_CATEGORY = 'INGREDIENT_CATEGORY',
}

export enum EmployeeAuditOperationType {
	CREATE = 'CREATE',
	UPDATE = 'UPDATE',
	DELETE = 'DELETE',
}

export const FORMATTED_AUDIT_COMPONENTS: Record<EmployeeAuditComponentName, string> = {
	[EmployeeAuditComponentName.PRODUCT]: 'Продукт',
	[EmployeeAuditComponentName.PRODUCT_CATEGORY]: 'Категория продукта',
	[EmployeeAuditComponentName.STORE_PRODUCT]: 'Продукт в магазине',
	[EmployeeAuditComponentName.EMPLOYEE]: 'Сотрудник',
	[EmployeeAuditComponentName.ADDITIVE]: 'Добавка',
	[EmployeeAuditComponentName.ADDITIVE_CATEGORY]: 'Категория добавок',
	[EmployeeAuditComponentName.STORE_ADDITIVE]: 'Добавка в магазине',
	[EmployeeAuditComponentName.PRODUCT_SIZE]: 'Размер продукта',
	[EmployeeAuditComponentName.RECIPE_STEPS]: 'Шаги рецепта',
	[EmployeeAuditComponentName.STORE]: 'Магазин',
	[EmployeeAuditComponentName.WAREHOUSE]: 'Склад',
	[EmployeeAuditComponentName.STORE_WAREHOUSE_STOCK]: 'Запасы склада магазина',
	[EmployeeAuditComponentName.INGREDIENT]: 'Ингредиент',
	[EmployeeAuditComponentName.INGREDIENT_CATEGORY]: 'Категория ингредиентов',
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
