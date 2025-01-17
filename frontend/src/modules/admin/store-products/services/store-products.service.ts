import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateStoreProductDTO,
	StoreProductDetailsDTO,
	StoreProductDTO,
	StoreProductsFilterDTO,
	UpdateStoreProductDTO,
} from '../models/store-products.model'

class StoreProductsService {
	/**
	 * Fetch all store products with optional filters
	 */
	async getStoreProducts(filter?: StoreProductsFilterDTO) {
		try {
			const response = await apiClient.get<PaginatedResponse<StoreProductDTO[]>>(
				`/store-products`,
				{
					params: buildRequestFilter(filter),
				},
			)
			return response.data
		} catch (error) {
			console.error('Failed to fetch store products: ', error)
			throw error
		}
	}

	/**
	 * Fetch details of a specific store product by ID
	 */
	async getStoreProduct(id: number) {
		try {
			const response = await apiClient.get<StoreProductDetailsDTO>(`/store-products/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch store product details for ID ${id}: `, error)
			throw error
		}
	}

	/**
	 * Create a new store product
	 */
	async createStoreProduct(data: CreateStoreProductDTO) {
		try {
			const response = await apiClient.post<void>(`/store-products`, data)
			return response.data
		} catch (error) {
			console.error('Failed to create store product: ', error)
			throw error
		}
	}

	/**
	 * Create multiple store products at once
	 */
	async createMultipleStoreProducts(data: CreateStoreProductDTO[]) {
		try {
			const response = await apiClient.post<void>(`/store-products/multiple`, data)
			return response.data
		} catch (error) {
			console.error('Failed to create multiple store products: ', error)
			throw error
		}
	}

	/**
	 * Update an existing store product by ID
	 */
	async updateStoreProduct(id: number, data: UpdateStoreProductDTO) {
		try {
			const response = await apiClient.put<void>(`/store-products/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update store product with ID ${id}: `, error)
			throw error
		}
	}

	/**
	 * Delete a store product by ID
	 */
	async deleteStoreProduct(id: number) {
		try {
			await apiClient.delete<void>(`/store-products/${id}`)
		} catch (error) {
			console.error(`Failed to delete store product with ID ${id}: `, error)
			throw error
		}
	}
}

export const storeProductsService = new StoreProductsService()
