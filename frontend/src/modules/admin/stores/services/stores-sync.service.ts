import { apiClient } from '@/core/config/axios-instance.config'
import type { StoreSynchronizationStatus } from '../models/stores-sync.models'

class StoreSyncService {
	async isStoreSynchronized() {
		try {
			const response = await apiClient.get<StoreSynchronizationStatus>('/sync/store')
			return response.data
		} catch (error) {
			console.error('Failed to fetch store sync status:', error)
			throw error
		}
	}

	async syncStoreStocksAndAdditives() {
		try {
			const response = await apiClient.post<void>('/sync/store')
			return response.data
		} catch (error) {
			console.error('Failed to post sync store:', error)
			throw error
		}
	}
}

export const storeSyncService = new StoreSyncService()
