import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface SuppliersFilter extends PaginationParams {
	search?: string
}

export interface SuppliersDTO {
	id: number
	name: string
	contactEmail: string
	contactPhone: string
	address: string
	createdAt: Date
	updatedAt: Date
}

export interface CreateSupplierDTO {
	name: string
	contactEmail: string
	contactPhone: string
	address: string
}

export interface UpdateSupplierDTO {
	name?: string
	contactEmail?: string
	contactPhone?: string
	address?: string
}
