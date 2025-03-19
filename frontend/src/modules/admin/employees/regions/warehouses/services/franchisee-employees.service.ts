import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import { type PaginatedResponse } from '@/core/utils/pagination.utils'

import type {
	FranchiseeEmployeeDetailsDTO,
	FranchiseeEmployeeDTO,
	UpdateFranchiseeEmployeeDTO,
} from '@/modules/admin/employees/franchisees/models/franchisees-employees.model'
import type {
	CreateEmployeeDTO,
	EmployeesFilter,
} from '@/modules/admin/employees/models/employees.models'

class FranchiseeEmployeeService {
	private readonly baseUrl = '/employees/franchisee'

	async getFranchiseeEmployees(filter?: EmployeesFilter) {
		const response = await apiClient.get<PaginatedResponse<FranchiseeEmployeeDTO[]>>(this.baseUrl, {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async createFranchiseeEmployee(dto: CreateEmployeeDTO) {
		const response = await apiClient.post<void>(this.baseUrl, dto)
		return response.data
	}

	async getFranchiseeEmployeeById(id: number) {
		const response = await apiClient.get<FranchiseeEmployeeDetailsDTO>(`${this.baseUrl}/${id}`)
		return response.data
	}

	async updateFranchiseeEmployee(id: number, dto: UpdateFranchiseeEmployeeDTO) {
		const response = await apiClient.put<void>(`${this.baseUrl}/${id}`, dto)
		return response.data
	}

	async deleteFranchiseeEmployee(employeeId: number) {
		const response = await apiClient.delete<void>(`${this.baseUrl}/${employeeId}`)
		return response.data
	}
}

export const franchiseeEmployeeService = new FranchiseeEmployeeService()
