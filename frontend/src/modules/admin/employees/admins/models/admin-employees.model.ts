import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'

export interface AdminEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface AdminEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
}
