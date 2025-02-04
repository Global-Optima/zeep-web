import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { EmployeeAuditDTO, EmployeeAuditFilter } from '../models/employees-audit.models'

export class EmployeeAuditService {
	private readonly baseUrl: string = '/audits'

	async getAudits(filter?: EmployeeAuditFilter) {
		try {
			const response = await apiClient.get<PaginatedResponse<EmployeeAuditDTO[]>>(this.baseUrl, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch employee audits:', error)
			throw error
		}
	}
}

export const employeeAuditService = new EmployeeAuditService()
