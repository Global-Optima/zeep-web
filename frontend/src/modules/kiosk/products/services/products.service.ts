import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import { buildFormData } from '@/core/utils/request-form-data-builder.utils'
import type {
	CreateProductCategoryDTO,
	CreateProductDTO,
	CreateProductSizeDTO,
	ProductCategoriesFilterDTO,
	ProductCategoryDTO,
	ProductDetailsDTO,
	ProductSizeDTO,
	ProductSizeDetailsDTO,
	ProductsFilterDTO,
	UpdateProductCategoryDTO,
	UpdateProductDTO,
	UpdateProductSizeDTO,
} from '../models/product.model'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'

class ProductsService {
	async getProducts(filter?: ProductsFilterDTO) {
		try {
			const response = await apiClient.get<PaginatedResponse<ProductDetailsDTO[]>>(`/products`, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch products: ', error)
			throw error
		}
	}

	async getProductDetails(id: number) {
		try {
			const response = await apiClient.get<ProductDetailsDTO>(`/products/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch product details for ID ${id}: `, error)
			throw error
		}
	}

	async createProduct(product: CreateProductDTO) {
		const formData = buildFormData<CreateProductDTO>(product)

		try {
			const response = await apiClient.post('/products', formData, {
				headers: {
					'Content-Type': 'multipart/form-data',
				},
			})

			return response.data
		} catch (error) {
			console.error('Error creating product:', error)
			throw error
		}
	}

	async updateProduct(id: number, data: UpdateProductDTO) {
		const formData = buildFormData<UpdateProductDTO>(data)

		try {
			const response = await apiClient.put<void>(`/products/${id}`, formData, {
				headers: {
					'Content-Type': 'multipart/form-data',
				},
			})
			return response.data
		} catch (error) {
			console.error(`Failed to update product with ID ${id}: `, error)
			throw error
		}
	}

	async deleteProduct(id: number) {
		try {
			await apiClient.delete<void>(`/products/${id}`)
		} catch (error) {
			console.error(`Failed to delete product with ID ${id}: `, error)
			throw error
		}
	}

	async getProductSizesByProductID(id: number) {
		try {
			const response = await apiClient.get<ProductSizeDTO[]>(`/products/${id}/sizes`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch sizes for product ID ${id}: `, error)
			throw error
		}
	}

	async getProductSizeById(id: number) {
		try {
			const response = await apiClient.get<ProductSizeDetailsDTO>(`/products/sizes/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch size by product ${id}: `, error)
			throw error
		}
	}

	async createProductSize(data: CreateProductSizeDTO) {
		try {
			const response = await apiClient.post<void>(`/products/sizes`, data)
			return response.data
		} catch (error) {
			console.error('Failed to create product size: ', error)
			throw error
		}
	}

	async updateProductSize(id: number, data: UpdateProductSizeDTO) {
		try {
			const response = await apiClient.put<void>(`/products/sizes/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update product size with ID ${id}: `, error)
			throw error
		}
	}

	async getAllProductCategories(filter?: ProductCategoriesFilterDTO) {
		const response = await apiClient.get<PaginatedResponse<ProductCategoryDTO[]>>(
			'/product-categories',
			{ params: buildRequestFilter(filter) },
		)
		return response.data
	}

	async getProductCategoryByID(id: number) {
		const response = await apiClient.get<ProductCategoryDTO>(`/product-categories/${id}`)
		return response.data
	}

	async updateProductCategory(id: number, dto: UpdateProductCategoryDTO) {
		const response = await apiClient.put(`/product-categories/${id}`, dto)
		return response.data
	}

	async createProductCategory(dto: CreateProductCategoryDTO) {
		const response = await apiClient.post<void>('/product-categories', dto)
		return response.data
	}

	async deleteProductCategory(id: number) {
		return (await apiClient.delete<void>(`/product-categories/${id}`)).data
	}
}

export const productsService = new ProductsService()
