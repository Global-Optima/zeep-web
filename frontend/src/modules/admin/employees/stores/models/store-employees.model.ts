import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
	EmployeesFilter,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import type { StoreDTO } from '@/modules/admin/stores/models/stores.models'

export interface UpdateStoreEmployeeDTO extends UpdateEmployeeDTO {
	storeId?: number
}

export interface StoreEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface StoreEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
	store: StoreDTO
}

export interface StoreEmployeeFilter extends EmployeesFilter {
	storeId?: number
}
