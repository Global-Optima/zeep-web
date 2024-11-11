export interface StoreProducts {
	id: number
	name: string
	description: string
	imageUrl: string
	category: string
	basePrice: number
}

export interface ProductCategory {
	id: number
	name: string
	description: string
}

export interface StoreProductDetails {
	id: number
	name: string
	description: string
	imageUrl: string
	basePrice: number
	sizes: ProductSize[]
	defaultAdditives: Additive[]
	recipeSteps: RecipeStep[]
}

export interface ProductSize {
	id: number
	name: string
	basePrice: number
	measure: string
}

export interface Additive {
	id: number
	name: string
	price: number
	imageUrl: string
}

export interface AdditiveCategory {
	id: number
	name: string
	additives: Additive[]
}

export interface RecipeStep {
	id: number
	description: string
	imageUrl: string
	step: number
}
