import type { LocalizedMessage } from './localized.model'

export interface LocalizedError {
	message: LocalizedMessage
	status: number
	timestamp: Date
	path: string
}
