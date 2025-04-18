import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import { buildFormData } from '@/core/utils/request-form-data-builder.utils'
import type {
	CreateProductCategoryDTO,
	CreateProductDTO,
	CreateProductSizeDTO,
	ProductCategoriesFilterDTO,
	ProductCategoryDTO,
	ProductCategoryTranslationsDTO,
	ProductDetailsDTO,
	ProductSizeDTO,
	ProductSizeDetailsDTO,
	ProductSizeTechnicalMap,
	ProductTranslationsDTO,
	ProductsFilterDTO,
	UpdateProductCategoryDTO,
	UpdateProductDTO,
	UpdateProductSizeDTO,
} from '../models/product.model'

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

	async getProductSizeTechMap(sizeId: number) {
		try {
			const response = await apiClient.get<ProductSizeTechnicalMap>(
				`/products/sizes/${sizeId}/technical-map`,
			)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch product size technical map for ID ${sizeId}: `, error)
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

	async upsertProductTranslations(productId: number, data: ProductTranslationsDTO) {
		try {
			const response = await apiClient.post<void>(`/products/${productId}/translations`, data)
			return response.data
		} catch (error) {
			console.error('Failed to upsert product translations: ', error)
			throw error
		}
	}

	async getProductTranslations(productId: number) {
		try {
			const response = await apiClient.get<ProductTranslationsDTO>(
				`/products/${productId}/translations`,
			)
			return response.data
		} catch (error) {
			console.error('Failed to get product translations: ', error)
			throw error
		}
	}

	async upsertProductCategoryTranslations(id: number, data: ProductCategoryTranslationsDTO) {
		try {
			const response = await apiClient.post<void>(`/product-categories/${id}/translations`, data)
			return response.data
		} catch (error) {
			console.error('Failed to upsert product category translations: ', error)
			throw error
		}
	}

	async getProductCategoryTranslations(id: number) {
		try {
			const response = await apiClient.get<ProductCategoryTranslationsDTO>(
				`/product-categories/${id}/translations`,
			)
			return response.data
		} catch (error) {
			console.error('Failed to get product category translations: ', error)
			throw error
		}
	}
}

export const productsService = new ProductsService()
