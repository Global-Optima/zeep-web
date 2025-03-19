import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import { type PaginatedResponse } from '@/core/utils/pagination.utils'

import type {
	AdminEmployeeDetailsDTO,
	AdminEmployeeDTO,
} from '@/modules/admin/employees/admins/models/admin-employees.model'
import type {
	CreateEmployeeDTO,
	EmployeesFilter,
} from '@/modules/admin/employees/models/employees.models'

class AdminEmployeeService {
	private readonly baseUrl = '/employees/admin'

	async getAdminEmployees(filter?: EmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<AdminEmployeeDTO[]>>(this.baseUrl, {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async createAdminEmployee(dto: CreateEmployeeDTO) {
		const response = await apiClient.post<void>(this.baseUrl, dto)
		return response.data
	}

	async getAdminEmployeeById(id: number) {
		const response = await apiClient.get<AdminEmployeeDetailsDTO>(`${this.baseUrl}/${id}`)
		return response.data
	}
}

export const adminEmployeeService = new AdminEmployeeService()
