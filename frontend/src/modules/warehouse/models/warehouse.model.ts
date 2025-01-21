export interface WarehouseDTO {
	id: string
	name: string
	facilityAddress: {
		id: number
		address: string
		longitude: number
		latitude: number
	}
}
