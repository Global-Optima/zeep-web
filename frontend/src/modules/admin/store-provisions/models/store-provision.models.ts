import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { ProvisionDTO } from '../../provisions/models/provision.models'

export enum StoreProvisionStatus {
	PREPARING = 'PREPARING',
	COMPLETED = 'COMPLETED',
	EXPIRED = 'EXPIRED',
  EMPTY = 'EMPTY',
}

export interface StoreProvisionFilter extends PaginationParams {
	search?: string
  statuses?: StoreProvisionStatus[]
}

export interface StoreProvisionDTO {
	id: number
	provision: ProvisionDTO
	volume: number
  initialVolume: number
	expirationInMinutes: number
	status: StoreProvisionStatus
	createdAt: Date
	completedAt?: Date
	expiresAt?: Date
}

export interface StoreProvisionDetailsIngredients {
	ingredient: IngredientsDTO
	quantity: number
  initialQuantity: number
}

export interface StoreProvisionDetailsDTO extends StoreProvisionDTO {
	ingredients: StoreProvisionDetailsIngredients[]
}

export interface CreateStoreProvisionDTO {
	provisionId: number
	volume: number
	expirationInMinutes: number
}

export interface UpdateStoreProvisionDTO {
	volume: number
	expirationInMinutes: number
}
