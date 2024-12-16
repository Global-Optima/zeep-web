import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	AdditiveCategoryDTO,
	ProductCategory,
	Products,
	ProductsFilter,
	StoreProductDetailsDTO,
} from '../models/product.model'

class ProductService {
	async getProducts(filter?: ProductsFilter): Promise<Products[]> {
		try {
			const response = await apiClient.get<Products[]>(`/products`, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch products: `, error)
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

	async getStoreProductDetails(
		productId: number,
		storeId: number,
	): Promise<StoreProductDetailsDTO> {
		try {
			const response = await apiClient.get<StoreProductDetailsDTO>(`/products/${productId}`, {
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

	async getAdditiveCategoriesByProductSize(productSizeId: number): Promise<AdditiveCategoryDTO[]> {
		try {
			const response = await apiClient.get<AdditiveCategoryDTO[]>(`/additives`, {
				params: {
					productSizeId,
				},
			})
			return response.data
		} catch (error) {
			console.error(`Failed to fetch additives for product size ID ${productSizeId}:`, error)
			throw error
		}
	}
}

export const productsService = new ProductService()
