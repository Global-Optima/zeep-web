import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'

import { type PaginatedResponse } from '@/core/utils/pagination.utils'

import type { CreateEmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type {
	RegionEmployeeDetailsDTO,
	RegionEmployeeDTO,
	RegionEmployeeFilter,
	UpdateRegionEmployeeDTO,
} from '@/modules/admin/employees/regions/models/region-employees.model'

class RegionEmployeeService {
	private readonly baseUrl = '/employees/region'

	async getRegionEmployees(filter?: RegionEmployeeFilter) {
		const response = await apiClient.get<PaginatedResponse<RegionEmployeeDTO[]>>(this.baseUrl, {
			params: buildRequestFilter(filter),
		})
		return response.data
	}

	async createRegionEmployee(dto: CreateEmployeeDTO, regionId: number) {
		const response = await apiClient.post<void>(this.baseUrl, dto, {
			params: { regionId: regionId },
		})
		return response.data
	}

	async getRegionEmployeeById(id: number) {
		const response = await apiClient.get<RegionEmployeeDetailsDTO>(`${this.baseUrl}/${id}`)
		return response.data
	}

	async updateRegionEmployee(id: number, dto: UpdateRegionEmployeeDTO) {
		const response = await apiClient.put<void>(`${this.baseUrl}/${id}`, dto)
		return response.data
	}

	async deleteRegionEmployee(employeeId: number) {
		const response = await apiClient.delete<void>(`${this.baseUrl}/${employeeId}`)
		return response.data
	}
}

export const regionEmployeeService = new RegionEmployeeService()
