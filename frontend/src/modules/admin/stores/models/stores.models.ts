export interface StoreDTO {
	id: number
	name: string
	isFranchise: boolean
	facilityAddress: {
		id: number
		address: string
		longitude: number
		latitude: number
	}
	contactPhone: string
	contactEmail: string
	storeHours: string
}
