import type { PaginationParams } from '@/core/utils/pagination.utils'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import type {
	ProductSizeIngredientDTO,
	ProductSizeProvisionDTO,
} from '@/modules/admin/store-products/models/store-products.model'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'

export enum ProductSizeNames {
	S = 'S',
	M = 'M',
	L = 'L',
}

export enum MachineCategory {
	TEA = 'TEA',
	COFFEE = 'COFFEE',
	ICE_CREAM = 'ICE_CREAM',
	OTHERS = 'OTHERS',
}

export const MACHINE_CATEGORY_OPTIONS = [
	{ label: 'Чай', value: MachineCategory.TEA },
	{ label: 'Кофе', value: MachineCategory.COFFEE },
	{ label: 'Мороженое', value: MachineCategory.ICE_CREAM },
	{ label: 'Другое', value: MachineCategory.OTHERS },
]

export const MACHINE_CATEGORY_FORMATTED: Record<MachineCategory, string> = {
	[MachineCategory.TEA]: 'Чай',
	[MachineCategory.COFFEE]: 'Кофе',
	[MachineCategory.ICE_CREAM]: 'Мороженое',
	[MachineCategory.OTHERS]: 'Другое',
}

export interface BaseProductDTO {
	name: string
	description: string
	imageUrl: string
	videoUrl: string
	category: ProductCategoryDTO
}

export interface ProductDTO extends BaseProductDTO {
	id: number
	productSizeCount: number
	basePrice: number
}

export interface ProductDetailsDTO extends ProductDTO {
	sizes: ProductSizeDTO[]
}

export interface ProductTotalNutrition {
	ingredients: string[]
	allergenIngredients: string[]
	calories: number
	proteins: number
	fats: number
	carbs: number
}

export interface BaseProductSizeDTO {
	name: string
	basePrice: number
	productId: number
	unit: UnitDTO
	size: number
	machineId: string
}

export interface ProductSizeDTO extends BaseProductSizeDTO {
	id: number
}

export interface ProductSizeDetailsDTO extends ProductSizeDTO {
	additives: ProductSizeAdditiveDTO[]
	ingredients: ProductSizeIngredientDTO[]
	provisions: ProductSizeProvisionDTO[]
}

export interface ProductSizeAdditiveDTO extends AdditiveDTO {
	isDefault: boolean
	isHidden: boolean
}

export interface CreateProductDTO {
	name: string
	description?: string
	categoryId: number
	image?: File
	video?: File
}

export interface SelectedAdditiveDTO {
	additiveId: number
	isDefault: boolean
	isHidden: boolean
}

export interface SelectedIngredientDTO {
	ingredientId: number
	quantity: number
}

export interface SelectedProvisionDTO {
	provisionId: number
	volume: number
}

export interface CreateProductSizeDTO {
	productId: number
	name: string
	size: number
	unitId: number
	basePrice: number
	machineId: string
	additives?: SelectedAdditiveDTO[]
	ingredients: SelectedIngredientDTO[]
	provisions: SelectedProvisionDTO[]
}

export interface UpdateProductDTO {
	name?: string
	description?: string
	categoryId?: number
	image?: File
	video?: File
	deleteImage: boolean
	deleteVideo: boolean
}

export interface UpdateProductSizeDTO {
	name?: string
	basePrice?: number
	size?: number
	unitId?: number
	isDefault?: boolean
	machineId?: string
	additives?: SelectedAdditiveDTO[]
	ingredients?: SelectedIngredientDTO[]
	provisions?: SelectedProvisionDTO[]
}

export interface ProductsFilterDTO extends PaginationParams {
	categoryId?: number
	search?: string
}

export interface ProductCategoryDTO {
	id: number
	name: string
	description: string
	machineCategory: MachineCategory
}

export interface ProductCategoriesFilterDTO extends PaginationParams {
	search?: string
}

export interface CreateProductCategoryDTO {
	name: string
	description?: string
	machineCategory?: MachineCategory
}

export interface UpdateProductCategoryDTO {
	name?: string
	description?: string
	machineCategory?: MachineCategory
}
