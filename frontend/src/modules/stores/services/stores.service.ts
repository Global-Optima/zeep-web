import { apiClient } from '@/core/config/axios-instance.config'
import type { Employee, Store } from '../models/stores.models'

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

	async getStoreEmployees(storeID: number): Promise<Employee[]> {
		try {
			const response = await apiClient.get<Employee[]>(`/stores/${storeID}/employees`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch employees for store ID ${storeID}:`, error)
			throw error
		}
	}
}

export const storesService = new StoreService()
