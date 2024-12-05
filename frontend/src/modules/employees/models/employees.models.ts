export enum EmployeeRole {
	RoleAdmin = 'Admin',
	RoleDirector = 'Director',
	RoleManager = 'Manager',
	RoleEmployee = 'Employee',
}

export interface Employee {
	id: number
	name: string
	phone: string
	email: string
	storeId: boolean
	role: EmployeeRole
}

export interface EmployeeLoginDTO {
	employeeId: number
	password: string
}
