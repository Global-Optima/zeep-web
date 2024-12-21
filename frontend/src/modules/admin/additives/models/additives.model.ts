export interface AdditiveCategoriesFilterQuery {
	search?: string
	productSizeId?: number
}

export interface AdditiveFilterQuery {
	search?: string
	minPrice?: number
	maxPrice?: number
	categoryId?: number
	productSizeId?: number
	page?: number
	pageSize?: number
}

export interface AdditiveCategoryItem {
	id: number
	name: string
	description: string
	price: number
	size: string
	imageUrl: string
	categoryId: number
}

export interface AdditiveCategories {
	id: number
	name: string
	description: string
	additives: AdditiveCategoryItem[]
	isMultipleSelect: boolean
}

export interface Additives {
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
