import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface FranchiseeDTO {
	id: number
	name: string
	description: string
}

export interface CreateFranchiseeDTO {
	name: string
	description?: string
}

export interface UpdateFranchiseeDTO {
	name?: string
	description?: string
}

export interface FranchiseeFilterDTO extends PaginationParams {
	search?: string
}

export interface AllFranchiseeFilter {
	search: string
}
