import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
	EmployeesFilter,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'

export interface UpdateFranchiseeEmployeeDTO extends UpdateEmployeeDTO {
	franchiseeId?: number
}

export interface FranchiseeEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface FranchiseeEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
	franchisee: FranchiseeDTO
}

export interface FranchiseEmployeeFilter extends EmployeesFilter {
	franchiseeId?: number
}
