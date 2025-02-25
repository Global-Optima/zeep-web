import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface EmployeeLoginDTO {
	email: string
	password: string
}

export interface EmployeeAccount {
	firstName: string
	lastName: string
	email: string
}

export enum EmployeeRole {
	ADMIN = 'ADMIN',
	OWNER = 'OWNER',
	STORE_MANAGER = 'STORE_MANAGER',
	BARISTA = 'BARISTA',
	WAREHOUSE_MANAGER = 'WAREHOUSE_MANAGER',
	WAREHOUSE_EMPLOYEE = 'WAREHOUSE_EMPLOYEE',
	FRANCHISEE_MANAGER = 'FRANCHISEE_MANAGER',
	FRANCHISEE_OWNER = 'FRANCHISEE_OWNER',
	REGION_WAREHOUSE_MANAGER = 'REGION_WAREHOUSE_MANAGER',
}

export enum EmployeeType {
	STORE = 'STORE',
	WAREHOUSE = 'WAREHOUSE',
	FRANCHISEE = 'FRANCHISEE',
	REGION = 'REGION',
	ADMIN = 'ADMIN',
}

export const EMPLOYEE_ROLES_FORMATTED: Record<EmployeeRole, string> = {
	[EmployeeRole.ADMIN]: 'Администратор',
	[EmployeeRole.OWNER]: 'Владелец',
	[EmployeeRole.STORE_MANAGER]: 'Менеджер кафе',
	[EmployeeRole.BARISTA]: 'Бариста',
	[EmployeeRole.WAREHOUSE_MANAGER]: 'Менеджер склада',
	[EmployeeRole.WAREHOUSE_EMPLOYEE]: 'Сотрудник склада',
	[EmployeeRole.FRANCHISEE_MANAGER]: 'Менеджер франшизы',
	[EmployeeRole.FRANCHISEE_OWNER]: 'Владелец франшизы',
	[EmployeeRole.REGION_WAREHOUSE_MANAGER]: 'Региональный менеджер',
}

export const EMPLOYEE_TYPES_FORMATTED: Record<EmployeeType, string> = {
	[EmployeeType.STORE]: 'Кафе',
	[EmployeeType.WAREHOUSE]: 'Склад',
	[EmployeeType.FRANCHISEE]: 'Франшиза',
	[EmployeeType.REGION]: 'Регион',
	[EmployeeType.ADMIN]: 'Администрация',
}

export interface CreateOrReplaceWorkdayDTO {
	day: string
	startAt: string
	endAt: string
}

export interface CreateEmployeeDTO {
	firstName: string
	lastName: string
	phone: string
	email: string
	role: EmployeeRole
	password: string
	isActive: boolean
	workdays: CreateOrReplaceWorkdayDTO[]
}

export interface UpdateEmployeeDTO {
	firstName?: string
	lastName?: string
	isActive?: boolean
	workdays?: CreateOrReplaceWorkdayDTO[]
	role?: EmployeeRole
}

export interface ReassignEmployeeTypeDTO {
	employeeType: EmployeeType
	role: EmployeeRole
	workplaceId: number
}

export interface BaseEmployeeDTO {
	firstName: string
	lastName: string
	phone: string
	email: string
	type: EmployeeType
	role: EmployeeRole
	isActive: boolean
}

export interface EmployeeDTO extends BaseEmployeeDTO {
	id: number
}

export interface EmployeeWorkdayDTO {
	id: number
	day: string
	startAt: string
	endAt: string
	employeeId: number
}

export interface BaseEmployeeDetailsDTO extends BaseEmployeeDTO {
	workdays: EmployeeWorkdayDTO[]
	employeeId: number
}

export interface EmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
}

export interface UpdatePasswordDTO {
	oldPassword: string
	newPassword: string
}

export interface EmployeesFilter extends PaginationParams {
	role?: EmployeeRole
	isActive?: boolean
	search?: string
}

export interface ReassignEmployeeTypeDTO {
	employeeType: EmployeeType
	role: EmployeeRole
	workplaceId: number
}
