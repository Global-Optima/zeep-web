import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface AdditiveDTO {
	id: number
	name: string
	description: string
	price: number
	imageUrl: string
	size: string
	category: {
		id: number
		name: string
		isMultipleSelect: boolean
	}
}

// DTO for Additives within Categories
export interface AdditiveCategoryItemDTO {
	id: number
	name: string
	description: string
	price: number
	imageUrl: string
	size: string
	categoryId: number
}

// Base DTO for Additive Categories
export interface AdditiveCategoryDTO {
	id: number
	name: string
	description: string
	additives: AdditiveCategoryItemDTO[]
	isMultipleSelect: boolean
}

// Filter Query DTO for Additives
export interface AdditiveFilterQuery extends PaginationParams {
	search?: string
	minPrice?: number
	maxPrice?: number
	categoryId?: number
	productSizeId?: number
}

// Filter Query DTO for Additive Categories
export interface AdditiveCategoriesFilterQuery extends PaginationParams {
	productSizeId?: number
	search?: string
}

// Create DTOs
export interface CreateAdditiveDTO {
	name: string
	description: string
	price: number
	imageUrl?: string
	size: string
	additiveCategoryId: number
}

export interface CreateAdditiveCategoryDTO {
	name: string
	description?: string
	isMultipleSelect: boolean
}

// Update DTOs
export interface UpdateAdditiveDTO {
	name?: string
	description?: string
	price?: number
	imageUrl?: string
	size?: string
	additiveCategoryId?: number
}

export interface UpdateAdditiveCategoryDTO {
	id: number
	name?: string
	description?: string
	isMultipleSelect?: boolean
}
