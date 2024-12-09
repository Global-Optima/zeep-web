export enum EmployeeRole {
	RoleAdmin = 'ADMIN',
	RoleDirector = 'DIRECTOR',
	RoleManager = 'MANAGER',
	RoleEmployee = 'EMPLOYEE',
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
}

export interface EmployeeLoginDTO {
	email: string
	password: string
}
