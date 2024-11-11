import { apiClient } from '@/core/config/axios-instance.config'
import { CURRENT_STORE_COOKIES_CONFIG } from '@/modules/stores/constants/store-cookies.constant'
import type {
	AdditiveCategory,
	ProductCategory,
	StoreProductDetails,
	StoreProducts,
} from '../models/product.model'

class ProductService {
	async getStoreProducts(
		categoryId: number,
		searchQuery = '',
		limit = 10,
		offset = 0,
	): Promise<StoreProducts[]> {
		const storeId = localStorage.getItem(CURRENT_STORE_COOKIES_CONFIG.key)

		try {
			const response = await apiClient.get<StoreProducts[]>(`/products`, {
				params: {
					storeId,
					categoryId,
					search: searchQuery,
					limit,
					offset,
				},
			})
			return response.data
		} catch (error) {
			console.error(
				`Failed to fetch products for store ID ${storeId} and category ID ${categoryId}:`,
				error,
			)
			throw error
		}
	}

	async getStoreCategories(): Promise<ProductCategory[]> {
		try {
			const response = await apiClient.get<ProductCategory[]>(`/categories`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch categories:`, error)
			throw error
		}
	}

	async getStoreProductDetails(productId: number): Promise<StoreProductDetails> {
		const storeId = localStorage.getItem(CURRENT_STORE_COOKIES_CONFIG.key)

		try {
			const response = await apiClient.get<StoreProductDetails>(`/products/${productId}`, {
				params: {
					storeId,
				},
			})
			return response.data
		} catch (error) {
			console.error(
				`Failed to fetch product details for store ID ${storeId} and product ID ${productId}:`,
				error,
			)
			throw error
		}
	}

	async getAdditiveCategoriesByStoreAndProduct(productId: number): Promise<AdditiveCategory[]> {
		const storeId = localStorage.getItem(CURRENT_STORE_COOKIES_CONFIG.key)

		try {
			const response = await apiClient.get<AdditiveCategory[]>(`/additives`, {
				params: {
					storeId,
					productId,
				},
			})
			return response.data
		} catch (error) {
			console.error(
				`Failed to fetch additives for store ID ${storeId} and product ID ${productId}:`,
				error,
			)
			throw error
		}
	}
}

export const productService = new ProductService()
