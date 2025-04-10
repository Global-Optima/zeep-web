import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { IngredientsDTO } from '../../ingredients/models/ingredients.model'
import type { UnitDTO } from '../../units/models/units.model'
import type { TechnicalMapDTO } from '@/modules/kiosk/products/models/product.model'

export interface ProvisionFilter extends PaginationParams {
	search?: string
}

export interface BaseProvisionDTO {
	name: string
	absoluteVolume: number
	unit: UnitDTO
	preparationInMinutes: number
	netCost: number
	limitPerDay: number
}

export interface ProvisionDTO extends BaseProvisionDTO {
	id: number
}

export interface ProvisionDetailsIngredients {
	ingredient: IngredientsDTO
	quantity: number
}

export interface ProvisionDetailsDTO extends ProvisionDTO {
	ingredients: ProvisionDetailsIngredients[]
}

export interface SelectedProvisionsIngredients {
	ingredientId: number
	quantity: number
}

export interface CreateProvisionDTO {
	name: string
	absoluteVolume: number
	unitId: number
	preparationInMinutes: number
	netCost: number
	limitPerDay: number
	ingredients: SelectedProvisionsIngredients[]
}

export interface UpdateProvisionDTO {
	name: string
	absoluteVolume: number
	unitId: number
	preparationInMinutes: number
	netCost: number
	limitPerDay: number
	ingredients: SelectedProvisionsIngredients[]
}

export interface ProvisionTechnicalMap {
	ingredients: TechnicalMapDTO[]
}
