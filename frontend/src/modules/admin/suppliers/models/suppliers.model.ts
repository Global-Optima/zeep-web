import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { StockMaterialsDTO } from '../../stock-materials/models/stock-materials.model'

export interface CreateSupplierDTO {
	name: string
	contactEmail: string
	contactPhone: string
	city: string
	address: string
}

export interface UpdateSupplierDTO {
	name?: string
	contactEmail?: string
	contactPhone?: string
	city?: string
	address?: string
}

export interface SupplierDTO {
	id: number
	name: string
	contactEmail: string
	contactPhone: string
	city: string
	address: string
	createdAt: string
	updatedAt: string
}

export interface SuppliersFilterDTO extends PaginationParams {
	search?: string
}

export interface GetMaterialsBySupplierFilterDTO extends PaginationParams {
	search?: string
}

export interface SupplierMaterialResponse {
	stockMaterial: StockMaterialsDTO
	basePrice: number
}

export interface UpdateSupplierMaterialDTO {
	stockMaterialId: number
	basePrice: number
}

export interface UpsertSupplierMaterialsDTO {
	materials: UpdateSupplierMaterialDTO[]
}
