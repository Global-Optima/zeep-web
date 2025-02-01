import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { AllRegionsFilter, RegionDTO } from '@/modules/admin/regions/models/regions.model'

class RegionsService {
	async getAllRegions(filter?: AllRegionsFilter) {
		try {
			const response = await apiClient.get<RegionDTO[]>('/regions/all', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}
}

export const regionsService = new RegionsService()
