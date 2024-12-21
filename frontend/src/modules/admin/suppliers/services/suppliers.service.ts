import { apiClient } from '@/core/config/axios-instance.config'
import { buildRequestFilter } from '@/core/utils/request-filters.utils'
import type {
	CreateSupplierDTO,
	Suppliers,
	SuppliersFilter,
	UpdateSupplierDTO,
} from '@/modules/admin/suppliers/models/suppliers.model'

class SuppliersService {
	async getSuppliers(filter?: SuppliersFilter): Promise<Suppliers[]> {
		try {
			const response = await apiClient.get<Suppliers[]>('/suppliers', {
				params: buildRequestFilter(filter),
			})
			return response.data
		} catch (error) {
			console.error('Failed to fetch suppliers:', error)
			throw error
		}
	}

	async getSupplier(id: number): Promise<Suppliers> {
		try {
			const response = await apiClient.get<Suppliers>(`/suppliers/${id}`)
			return response.data
		} catch (error) {
			console.error(`Failed to fetch supplier by id ${id}:`, error)
			throw error
		}
	}

	async updateSupplier(id: number, dto: UpdateSupplierDTO) {
		try {
			const response = await apiClient.put<void>(`/suppliers/${id}`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to update supplier by id ${id}:`, error)
			throw error
		}
	}

	async createSupplier(dto: CreateSupplierDTO) {
		try {
			const response = await apiClient.post<void>(`/suppliers`, dto)
			return response.data
		} catch (error) {
			console.error(`Failed to create supplier:`, error)
			throw error
		}
	}
}

export const suppliersService = new SuppliersService()
