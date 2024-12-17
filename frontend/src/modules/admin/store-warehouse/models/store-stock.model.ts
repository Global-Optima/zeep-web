export interface StoreStocks {
	id: number
	name: string
	quantity: number
	unit: string
	lowStockAlert: boolean
	lowStockThreshold: number
}

export interface StoreStocksFilter {
	search?: string
	lowStockOnly?: boolean
	page?: number
	pageSize?: number
}

interface CreateStoreStockItem {
	ingredientId: number
	quantity: number
	lowStockThreshold: number
}

export interface CreateMultipleStoreStock {
	ingredientStocks: CreateStoreStockItem[]
}

export interface UpdateStoreStock {
	quantity: number
	lowStockThreshold: number
}
