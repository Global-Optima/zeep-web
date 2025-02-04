import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import { type PaginatedResponse } from '@/core/utils/pagination.utils'

import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type {
	UpdateWarehouseEmployeeDTO,
	WarehouseEmployeeDetailsDTO,
	WarehouseEmployeeDTO,
	WarehouseEmployeeFilter,
} from '@/modules/admin/employees/warehouses/models/warehouse-employees.model'

class WarehouseEmployeeService {
	private readonly baseUrl = '/employees/warehouse'

	async getWarehouseEmployees(filter?: WarehouseEmployeeFilter) {
		const response = await apiClient.get<PaginatedResponse<WarehouseEmployeeDTO[]>>(this.baseUrl, {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async createWarehouseEmployee(dto: CreateEmployeeDTO) {
		const response = await apiClient.post<void>(this.baseUrl, dto)
		return response.data
	}

	async getWarehouseEmployeeById(id: number) {
		const response = await apiClient.get<WarehouseEmployeeDetailsDTO>(`${this.baseUrl}/${id}`)
		return response.data
	}

	async updateWarehouseEmployee(id: number, dto: UpdateWarehouseEmployeeDTO) {
		const response = await apiClient.put<void>(`${this.baseUrl}/${id}`, dto)
		return response.data
	}

	async deleteWarehouseEmployee(employeeId: number) {
		const response = await apiClient.delete<void>(`${this.baseUrl}/${employeeId}`)
		return response.data
	}
}

export const warehouseEmployeeService = new WarehouseEmployeeService()
