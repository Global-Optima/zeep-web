import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
	EmployeeRole,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'

export interface UpdateFranchiseeEmployeeDTO extends UpdateEmployeeDTO {
	role?: EmployeeRole
	franchiseeId?: number
}

export interface FranchiseeEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface FranchiseeEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
	employeeId: number
	franchisee: FranchiseeDTO
}
