import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
	EmployeesFilter,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'

export interface UpdateWarehouseEmployeeDTO extends UpdateEmployeeDTO {
	warehouseId?: number
}

export interface WarehouseEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface WarehouseEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
	warehouse: WarehouseDTO
}

export interface WarehouseEmployeeFilter extends EmployeesFilter {
	warehouseId?: number
}
