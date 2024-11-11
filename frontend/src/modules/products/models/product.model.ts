export interface Products {
	id: number
	title: string
	description?: string
	image: string
	category: string
	startPrice: number
}

export interface ProductSizes {
	id: number
	label: string
	volume: number
	price: number
}
