export interface Warehouse {
	id: string
	name: string
	facilityAddress: {
		id: number
		address: string
		longitude: number
		latitude: number
	}
}
