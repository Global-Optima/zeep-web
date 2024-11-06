import type { Additives } from '@/modules/additives/models/additive.model'
import type { Products } from '@/modules/products/models/product.model'

export interface CartItems {
	product: Products
	quantity: number
	additives: Additives[]
}
