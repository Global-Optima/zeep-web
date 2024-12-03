export interface Store {
	id: number
	name: string
	facilityAddress: {
		id: number
		address: string
		longitude: string
		latitude: string
	}
	isFranchise: boolean
}
