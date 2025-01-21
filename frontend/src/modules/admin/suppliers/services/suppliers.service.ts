import { apiClient } from '@/core/config/axios-instance.config'
import type { PaginatedResponse } from '@/core/utils/pagination.utils'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateSupplierDTO,
	GetMaterialsBySupplierFilterDTO,
	SupplierDTO,
	SupplierMaterialResponse,
	SuppliersFilterDTO,
	UpdateSupplierDTO,
  UpsertSupplierMaterialsDTO,
} from '../models/suppliers.model'

class SuppliersService {
	async getSuppliers(filter?: SuppliersFilterDTO) {
		try {
			const response = await apiClient.get<PaginatedResponse<SupplierDTO[]>>(`/suppliers`, {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch suppliers:', error)
			throw error
		}
	}

	async getSupplierByID(id: number) {
		try {
			const response = await apiClient.get<SupplierDTO>(`/suppliers/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch supplier details for ID ${id}:`, error)
			throw error
		}
	}

	async createSupplier(data: CreateSupplierDTO) {
		try {
			const response = await apiClient.post<void>(`/suppliers`, data)
			return response.data
		} catch (error) {
			console.error('Failed to create supplier:', error)
			throw error
		}
	}

	async updateSupplier(id: number, data: UpdateSupplierDTO) {
		try {
			const response = await apiClient.put<void>(`/suppliers/${id}`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update supplier with ID ${id}:`, error)
			throw error
		}
	}

	async updateSupplierMaterials(supplierId: number, data: UpsertSupplierMaterialsDTO) {
		try {
			const response = await apiClient.put<void>(`/suppliers/${supplierId}/materials`, data)
			return response.data
		} catch (error) {
			console.error(`Failed to update supplier materials with ID ${supplierId}:`, error)
			throw error
		}
	}

	async deleteSupplier(id: number) {
		try {
			await apiClient.delete<void>(`/suppliers/${id}`)
		} catch (error) {
			console.error(`Failed to delete supplier with ID ${id}:`, error)
			throw error
		}
	}

	async getMaterialsBySupplier(id: number, filter?: GetMaterialsBySupplierFilterDTO) {
		try {
			const response = await apiClient.get<SupplierMaterialResponse[]>(
				`/suppliers/${id}/materials`,
				{ params: buildRequestFilter(filter) },
			)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch materials for supplier ID ${id}:`, error)
			throw error
		}
	}
}

export const suppliersService = new SuppliersService()
