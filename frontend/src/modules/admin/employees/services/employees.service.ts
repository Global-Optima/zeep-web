import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import type {
	CreateEmployeeDTO,
	EmployeeDetailsDTO,
	EmployeeDTO,
	EmployeesFilter,
	UpdateEmployeeDTO,
} from '@/modules/admin/employees/models/employees.models'
import { type PaginatedResponse } from '../../../../core/utils/pagination.utils'

class EmployeeService {
	async getEmployees(filter?: EmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<EmployeeDTO[]>>(`/employees`, {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async createEmployee(dto: CreateEmployeeDTO) {
		const response = await apiClient.post<void>(`/employees`, dto)
		return response.data
	}

	async getEmployeeById(id: number) {
		const response = await apiClient.get<EmployeeDetailsDTO>(`/employees/${id}`)
		return response.data
	}

	async updateEmployee(id: number, dto: UpdateEmployeeDTO) {
		const response = await apiClient.put<void>(`/employees/${id}`, dto)
		return response.data
	}

	async getCurrentEmployee(): Promise<EmployeeDTO> {
		const response = await apiClient.get<EmployeeDTO>(`/employees/current`)
		return response.data
	}
}

export const employeesService = new EmployeeService()
