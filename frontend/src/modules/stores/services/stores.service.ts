import { apiClient } from '@/core/config/axios-instance.config'
import type { Store } from '../models/stores.models'

class StoreService {
	async getStores(): Promise<Store[]> {
		try {
			const response = await apiClient.get<Store[]>('/stores')
			return response.data
		} catch (error) {
			console.error('Failed to fetch stores:', error)
			throw error
		}
	}
}

export const storesService = new StoreService()
