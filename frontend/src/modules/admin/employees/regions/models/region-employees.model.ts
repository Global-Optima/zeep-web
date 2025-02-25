import type {
	BaseEmployeeDetailsDTO,
	BaseEmployeeDTO,
	EmployeesFilter,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import type { RegionDTO } from '@/modules/admin/regions/models/regions.model'

export interface UpdateRegionEmployeeDTO extends UpdateEmployeeDTO {
	regionId?: number
}

export interface RegionEmployeeDTO extends BaseEmployeeDTO {
	id: number
	employeeId: number
}

export interface RegionEmployeeDetailsDTO extends BaseEmployeeDetailsDTO {
	id: number
	region: RegionDTO
}

export interface RegionEmployeeFilter extends EmployeesFilter {
	regionId?: number
}
