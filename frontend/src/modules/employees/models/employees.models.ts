export enum EmployeeRole {
	ADMIN = 'ADMIN',
	DIRECTOR = 'DIRECTOR',
	MANAGER = 'MANAGER',
	BARISTA = 'BARISTA',
	WAREHOUSE = 'WAREHOUSE_EMPLOYEE',
}

export enum EmployeeType {
	Store = 'STORE',
	Warehouse = 'WAREHOUSE',
}

export interface Employee {
	id: number
	name: string
	phone: string
	email: string
	role: EmployeeRole
	isActive: boolean
	type: EmployeeType
	imageUrl?: string

	storeDetails?: {
		storeId: number
	}
	warehouseDetails?: {
		warehouseId: number
	}
}

export interface EmployeeLoginDTO {
	email: string
	password: string
}

export interface EmployeesFilter {
	search?: string
	type?: EmployeeType
	role?: EmployeeRole
	storeId?: number
	warehouseId?: string
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
}

export interface UpdateEmployeeDto {
	name: string
	phone: string
	email: string
	role: EmployeeRole
}
