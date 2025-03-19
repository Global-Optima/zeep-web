import type { FranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import type { WarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'

export interface StoreDTO {
	id: number
	name: string
	franchisee?: FranchiseeDTO
	warehouse: WarehouseDTO
	facilityAddress: {
		id: number
		address: string
	}
	isActive: boolean
	contactPhone: string
	contactEmail: string
	storeHours: string
}
