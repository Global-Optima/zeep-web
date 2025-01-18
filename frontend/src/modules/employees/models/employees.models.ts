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

export interface EmployeeAccount {
	firstName: string
	lastName: string
	email: string
}

export interface Employee {
	id: number
	firstName: string
	lastName: string
	phone: string
	email: string
	role: EmployeeRole
	isActive: boolean
	type: EmployeeType
}

export interface StoreEmployee extends Employee {
	storeId: number
	isFranchise: boolean
}

export interface WarehouseEmployee extends Employee {
	warehouseId: number
}

export interface EmployeeLoginDTO {
	email: string
	password: string
}

export interface StoreEmployeesFilter {
	search?: string
	role?: EmployeeRole
	storeId?: number
	isActive?: boolean
}

export interface WarehouseEmployeesFilter {
	search?: string
	role?: EmployeeRole
	warehouseId?: number
	isActive?: boolean
}

export const EMPLOYEE_ROLES_FORMATTED: Record<EmployeeRole, string> = {
	[EmployeeRole.ADMIN]: 'Админ',
	[EmployeeRole.DIRECTOR]: 'Директор',
	[EmployeeRole.MANAGER]: 'Менеджер',
	[EmployeeRole.BARISTA]: 'Бариста',
	[EmployeeRole.WAREHOUSE]: 'Склад',
}

export interface CreateEmployeeDto {
	name: string
	phone: string
	email: string
	role: EmployeeRole
	password: string
	type: EmployeeType
	storeDetails?: {
		storeId: number
	}
	warehouseDetails?: {
		warehouseId: number
	}
}

export interface UpdateEmployeeDto {
	name: string
	phone: string
	email: string
	role: EmployeeRole
}
