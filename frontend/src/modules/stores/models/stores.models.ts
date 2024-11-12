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

export interface Employee {
	id: number
	name: string
	phone: string
	email: string
	isActive: boolean
	role: {
		id: number
		name: string
	}
}
