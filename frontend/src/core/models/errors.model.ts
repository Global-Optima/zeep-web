export interface LocalizedError {
	message: {
		en: string
		ru: string
		kk: string
	}
	status: number
	timestamp: Date
	path: string
}
