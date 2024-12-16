import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type { AdditiveCategories } from '@/modules/admin/additives/models/additives.model'
import type {
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
			const response = await apiClient.get<ProductCategory[]>(`/product-categories`)
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

	async getAdditiveCategories(productSizeId: number): Promise<AdditiveCategories[]> {
		try {
			const response = await apiClient.get<AdditiveCategories[]>(`/additives/categories`, {
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
