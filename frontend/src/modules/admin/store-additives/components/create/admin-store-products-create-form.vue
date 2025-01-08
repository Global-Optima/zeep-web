<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// Shadcn UI / custom wrappers

// We will have a parent-level array of these store products
export interface CreateStoreProductDTO {
  productId: number
  isAvailable: boolean
  productName?: string        // For display only
  basePrice?: number          // For display only
  categoryName?: string       // For display only
  imageUrl?: string           // For display only
  productSizes?: CreateStoreProductSizeDTO[]
}

export interface CreateStoreProductSizeDTO {
  productSizeID: number
  storePrice?: number
  sizeName?: string           // For display only
  measure?: string            // For display only
  basePrice?: number          // For display only
}

// The top-level form is going to contain an array of these store products.
const storeProducts = ref<CreateStoreProductDTO[]>([])

/**
 * Example function to create a new "empty" store product object,
 * including some display fields if you have them from a real API selection.
 */
function createEmptyStoreProduct(): CreateStoreProductDTO {
  return {
    productId: 0,
    productName: '',
    categoryName: '',
    imageUrl: '',
    basePrice: 0,
    isAvailable: true,
    productSizes: []
  }
}

/**
 * Example function to create a new "empty" store product size
 */
function createEmptyStoreProductSize(): CreateStoreProductSizeDTO {
  return {
    productSizeID: 0,
    sizeName: '',
    measure: '',
    basePrice: 0,
    storePrice: 0,
  }
}

/**
 * 1) Zod schema for array of CreateStoreProductDTO.
 *    We only validate the necessary fields, not the display fields
 *    (which are optional).
 */
const StoreProductSizeSchema = z.object({
  productSizeID: z.number().min(1, 'ID размера должен быть > 0'),
  storePrice: z.number().min(0, 'Цена в магазине >= 0').optional(),
})

const StoreProductSchema = z.object({
  productId: z.number().min(1, 'Выберите корректный ID продукта'),
  isAvailable: z.boolean(),
  productSizes: z.array(StoreProductSizeSchema).optional(),
})

const storeProductsSchema = toTypedSchema(z.array(StoreProductSchema))

/**
 * 2) Use vee-validate to manage an array of store products
 */
const { handleSubmit, errors } = useForm<CreateStoreProductDTO[]>({
  validationSchema: storeProductsSchema,
  initialValues: storeProducts.value, // Start empty or pass initial data if needed
})

/**
 * 3) Add / Remove store products from the array
 */
function addStoreProduct() {
  storeProducts.value.push(createEmptyStoreProduct())
}

function removeStoreProduct(index: number) {
  storeProducts.value.splice(index, 1)
}

/**
 * Add / Remove sub items (product sizes) for a specific store product
 */
function addSize(index: number) {
  storeProducts.value[index].productSizes?.push(createEmptyStoreProductSize())
}

function removeSize(productIndex: number, sizeIndex: number) {
  storeProducts.value[productIndex].productSizes?.splice(sizeIndex, 1)
}

/**
 * 4) Submitting the entire array
 */
function onSubmit(values: CreateStoreProductDTO[]) {
  // This is the final array of store products
  console.log('Final Store Products:', values)
  // You might send it to the server or emit upwards, etc.
}

function onCancel() {
  // Example: Clear everything or navigate away
  storeProducts.value = []
}
</script>
