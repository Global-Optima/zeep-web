import type { EmployeeDTO } from '@/modules/admin/employees/models/employees.models'

export function getEmployeeFullName(employee: EmployeeDTO) {
	return `${employee.firstName} ${employee.lastName}`
}

export function getEmployeeShortName(employee: EmployeeDTO) {
	return `${employee.lastName} ${employee.firstName.charAt(0)}.`
}

export function getEmployeeInitials(employee: EmployeeDTO) {
	return `${employee.firstName.charAt(0)}${employee.lastName.charAt(0)}`
}
