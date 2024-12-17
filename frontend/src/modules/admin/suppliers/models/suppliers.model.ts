export interface SuppliersFilter {
	searchTerm: string
}

export interface Suppliers {
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
