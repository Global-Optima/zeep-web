export interface Employee {
	id: number
	name: string
	phone: string
	email: string
	storeId: boolean
	role: EmployeeRoles
}

export interface EmployeeLoginDTO {
	employeeId: number
	password: string
}

export enum EmployeeRoles {
	ADMIN = 'ADMIN',
	DIRECTOR = 'DIRECTOR',
	MANAGER = 'MANAGER',
	EMPLOYEE = 'EMPLOYEE',
}
