import type { BaseEmployeeDTO } from '@/modules/admin/employees/models/employees.models'

export function getEmployeeFullName(employee: BaseEmployeeDTO) {
	return `${employee.firstName} ${employee.lastName}`
}

export function getEmployeeShortName(employee: BaseEmployeeDTO) {
	return `${employee.lastName} ${employee.firstName.charAt(0)}.`
}

export function getEmployeeInitials(employee: BaseEmployeeDTO) {
	return `${employee.firstName.charAt(0)}${employee.lastName.charAt(0)}`
}
