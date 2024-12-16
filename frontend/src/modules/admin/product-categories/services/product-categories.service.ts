import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { ProductCategories, ProductCategoriesFilter } from '../models/product-categories.model'

class ProductCategoriesService {
	async getProductCategories(filter?: ProductCategoriesFilter) {
		try {
			const response = await apiClient.get<ProductCategories[]>('/product-categories', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch product categories:', error)
			throw error
		}
	}
}

export const productCategoriesService = new ProductCategoriesService()
