import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
	EmployeeRole,
	EmployeesFilter,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'

export interface UpdateStoreEmployeeDTO extends UpdateEmployeeDTO {
	role?: EmployeeRole
	storeId?: number
}

export interface StoreEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface StoreEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
	employeeId: number
	store: StoreDTO
}

export interface StoreEmployeeFilter extends EmployeesFilter {
	storeId?: number
}
