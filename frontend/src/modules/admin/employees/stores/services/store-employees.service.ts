import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import { type PaginatedResponse } from '@/core/utils/pagination.utils'
import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type {
	StoreEmployeeDetailsDTO,
	StoreEmployeeDTO,
	StoreEmployeeFilter,
	UpdateStoreEmployeeDTO,
} from '@/modules/admin/employees/stores/models/store-employees.model'

class StoreEmployeeService {
	private readonly baseUrl = '/employees/store'

	async getStoreEmployees(filter?: StoreEmployeeFilter) {
		const response = await apiClient.get<PaginatedResponse<StoreEmployeeDTO[]>>(this.baseUrl, {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async createStoreEmployee(dto: CreateEmployeeDTO, storeId: number) {
		const response = await apiClient.post<void>(this.baseUrl, dto, {
			params: { storeId: storeId },
		})
		return response.data
	}

	async getStoreEmployeeById(id: number) {
		const response = await apiClient.get<StoreEmployeeDetailsDTO>(`${this.baseUrl}/${id}`)
		return response.data
	}

	async updateStoreEmployee(id: number, dto: UpdateStoreEmployeeDTO) {
		const response = await apiClient.put<void>(`${this.baseUrl}/${id}`, dto)
		return response.data
	}

	async deleteStoreEmployee(employeeId: number) {
		const response = await apiClient.delete<void>(`${this.baseUrl}/${employeeId}`)
		return response.data
	}
}

export const storeEmployeeService = new StoreEmployeeService()
