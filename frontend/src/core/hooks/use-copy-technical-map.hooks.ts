import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { provisionsService } from '@/modules/admin/provisions/services/provisions.service'
// useTechnicalMapCopy.ts
import type { TechnicalMapDTO } from '@/modules/kiosk/products/models/product.model'
import { productsService } from '@/modules/kiosk/products/services/products.service'
import { ref } from 'vue'

export enum TechnicalMapEntity {
	PRODUCT_SIZE = 'productSize',
	ADDITIVE = 'additive',
	PROVISION = 'provision',
}

// Internal buffer for the technical map reference.
const technicalMapReference = ref<string | null>(null)

/**
 * Sets the technical map reference using an entity type and an ID.
 * The reference is built in the format "entity-{id}" (e.g. "productSize-6") and saved to the internal buffer.
 * It also copies the reference to the user's clipboard.
 *
 * @param entity - The type of entity (using the TechnicalMapEntity enum).
 * @param id - The ID of the entity.
 */
async function setTechnicalMapReference(entity: TechnicalMapEntity, id: number): Promise<void> {
	const refStr = `${entity}-${id}`
	technicalMapReference.value = refStr

	try {
		if (navigator.clipboard && typeof navigator.clipboard.writeText === 'function') {
			await navigator.clipboard.writeText(refStr)
			console.log('Technical map reference copied to clipboard:', refStr)
		}
	} catch (error) {
		console.error('Failed to copy technical map reference to clipboard:', error)
	}
}

/**
 * Returns the current technical map reference.
 */
function getTechnicalMapReference(): string | null {
	return technicalMapReference.value
}

/**
 * Clears the stored technical map reference.
 */
function clearTechnicalMapReference(): void {
	technicalMapReference.value = null
}

async function getTechnicalMapReferenceFromBuffer(): Promise<string | null> {
	if (technicalMapReference.value) {
		return technicalMapReference.value
	}

	try {
		if (navigator.clipboard && typeof navigator.clipboard.readText === 'function') {
			const clipText = await navigator.clipboard.readText()
			// Basic validation: must contain a dash.
			if (clipText && clipText.includes('-')) {
				return clipText
			}
		}
	} catch (error) {
		console.error('Error reading technical map reference from clipboard:', error)
	}

	return null
}

/**
 * Fetches the technical map using the stored reference.
 * Expects the reference in the format "entity-{id}" and calls the appropriate service.
 * Returns a promise that resolves with an array of TechnicalMapDTO, or null if no reference is stored.
 */
async function fetchTechnicalMap(): Promise<TechnicalMapDTO[] | null> {
	const refStr = await getTechnicalMapReferenceFromBuffer()
	if (!refStr) return null

	// Reference should be in the format "entity-{id}"
	const parts = refStr.split('-')
	if (parts.length < 2) {
		throw new Error('Invalid technical map reference format')
	}

	const entity = parts[0] // e.g. "productSize", "additive", "provision"
	const id = Number(parts[1])
	if (isNaN(id)) {
		throw new Error('Invalid id in technical map reference')
	}

	switch (entity) {
		case TechnicalMapEntity.PRODUCT_SIZE: {
			const productMap = await productsService.getProductSizeTechMap(id)
			return productMap.ingredients
		}
		case TechnicalMapEntity.ADDITIVE: {
			const additiveMap = await additivesService.getAdditiveTechMap(id)
			return additiveMap.ingredients
		}
		case TechnicalMapEntity.PROVISION: {
			const provisionMap = await provisionsService.getProvisionTechMap(id)
			return provisionMap.ingredients
		}
		default:
			throw new Error(`Unknown technical map entity: ${entity}`)
	}
}

export function useCopyTechnicalMap() {
	return {
		setTechnicalMapReference,
		getTechnicalMapReference,
		clearTechnicalMapReference,
		fetchTechnicalMap,
		getTechnicalMapReferenceFromBuffer,
	}
}
