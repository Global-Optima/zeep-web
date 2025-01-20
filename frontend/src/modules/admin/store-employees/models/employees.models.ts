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
	DIRECTOR = 'DIRECTOR',
	MANAGER = 'MANAGER',
	BARISTA = 'BARISTA',
	WAREHOUSE = 'WAREHOUSE_EMPLOYEE',
}

export enum EmployeeType {
	STORE = 'STORE',
	WAREHOUSE = 'WAREHOUSE',
}

export const EMPLOYEE_ROLES_FORMATTED: Record<EmployeeRole, string> = {
	[EmployeeRole.ADMIN]: 'Админ',
	[EmployeeRole.DIRECTOR]: 'Директор',
	[EmployeeRole.MANAGER]: 'Менеджер',
	[EmployeeRole.BARISTA]: 'Бариста',
	[EmployeeRole.WAREHOUSE]: 'Работник Склада',
}

export interface CreateWorkdayDTO {
	day: string
	startAt: string
	endAt: string
}

export interface UpdateWorkdayDTO {
	startAt?: string
	endAt?: string
}

export interface EmployeeWorkdayDTO {
	id: number
	day: string
	startAt: string
	endAt: string
	employeeId: number
}

export interface CreateEmployeeDTO {
	firstName: string
	lastName: string
	phone?: string
	email: string
	role: EmployeeRole
	password: string
	isActive: boolean
	workdays: CreateWorkdayDTO[]
}

export interface CreateStoreEmployeeDTO extends CreateEmployeeDTO {
	storeId: number
	isFranchise?: boolean
}

export interface CreateWarehouseEmployeeDTO extends CreateEmployeeDTO {
	warehouseId: number
}

export interface UpdateEmployeeDTO {
	firstName?: string
	lastName?: string
	phone?: string
	email?: string
	role?: EmployeeRole
	isActive?: boolean
}

export interface UpdateStoreEmployeeDTO extends UpdateEmployeeDTO {
	storeId?: number
	isFranchise?: boolean
}

export interface UpdateWarehouseEmployeeDTO extends UpdateEmployeeDTO {
	warehouseId?: number
}

export interface EmployeeDTO {
	id: number
	firstName: string
	lastName: string
	phone: string
	email: string
	type: EmployeeType
	role: EmployeeRole
	isActive: boolean
}

export interface StoreEmployeeDTO extends EmployeeDTO {
	storeId: number
	isFranchise: boolean
}

export interface WarehouseEmployeeDTO extends EmployeeDTO {
	warehouseId: number
}

export interface RoleDTO {
	name: string
}

export interface UpdatePasswordDTO {
	oldPassword: string
	newPassword: string
}

export interface GetStoreEmployeesFilter extends PaginationParams {
	role?: string
	isActive?: boolean
	search?: string
}

export interface GetWarehouseEmployeesFilter extends PaginationParams {
	role?: EmployeeRole
	isActive?: boolean
	search?: string
}
