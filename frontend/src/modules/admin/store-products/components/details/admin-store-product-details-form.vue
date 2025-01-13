<script setup lang="ts">
import { computed, reactive } from 'vue'

// UI Components (shadcn or custom)
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { Checkbox } from '@/core/components/ui/checkbox'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'

// Icons
import type { StoreProductDetailsDTO, UpdateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { ProductDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronLeft } from 'lucide-vue-next'




/**
 * Props define:
 * 1. `initialStoreProduct` - existing store product data to update
 * 2. `product` - the full product definition (with sizes, basePrices, etc.)
 */
const props = defineProps<{
  initialStoreProduct: StoreProductDetailsDTO
  product: ProductDetailsDTO
}>()


const emits = defineEmits<{
  onSubmit: [dto: UpdateStoreProductDTO]
  onCancel: []
}>()

// ----- Local State -----------------------------------
/**
 * Local reactive state to manage store-product updates.
 * We'll merge `initialStoreProduct` with product's full info.
 */
interface ProductSizeLocal {
  id: number
  name: string
  basePrice: number
  storePrice: number
  selected: boolean
}

interface StoreProductLocal {
  isAvailable: boolean
  sizes: ProductSizeLocal[]
}

// Initialize local reactive data
const storeProductLocal = reactive<StoreProductLocal>({
  isAvailable: props.initialStoreProduct.isAvailable,
  sizes: [],
})

/**
 * Set up the local `sizes` array by merging:
 * - The product’s full size list
 * - The store’s existing `storePrice` (if present)
 */
function initializeSizes() {
  // Map all possible sizes from the product
  // If store has a matching size, use its storePrice; otherwise default to product basePrice
  storeProductLocal.sizes = props.product.sizes.map((sz) => {
    const existingSize = props.initialStoreProduct.sizes.find(
      (ps) => ps.id === sz.id
    )
    return {
      id: sz.id,
      name: sz.name,
      basePrice: sz.basePrice,
      storePrice: existingSize?.storePrice ?? sz.basePrice,
      selected: !!existingSize, // If it's in the store product, mark as selected
    }
  })
}

// Call once on mount
initializeSizes()

// ----- Computed ---------------------------------------

// Check if there are any sizes at all
const hasSizes = computed(() => storeProductLocal.sizes.length > 0)

// Submit updated data
function onSubmit() {
  const updateDTO: UpdateStoreProductDTO = {
    isAvailable: storeProductLocal.isAvailable,
    productSizes: storeProductLocal.sizes
      .filter((sz) => sz.selected)
      .map((sz) => ({
        productSizeID: sz.id,
        storePrice: sz.storePrice,
      })),
  }
  emits('onSubmit', updateDTO)
}

function onResetSizes() {
  initializeSizes()
}

// Cancel
function onCancel() {
  emits('onCancel')
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				type="button"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>

			<h1 class="font-semibold text-xl tracking-tight shrink-0">Обновить товар</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					@click="onSubmit"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Product Card -->
		<Card>
			<CardHeader>
				<div class="flex justify-between">
					<div>
						<CardTitle>{{ product.name }}</CardTitle>
						<CardDescription class="mt-1">{{ product.category.name }}</CardDescription>
					</div>

					<Button
						variant="outline"
						type="button"
						@click="onResetSizes"
					>
						Сбросить размеры к исходным
					</Button>
				</div>
			</CardHeader>

			<CardContent class="space-y-4">
				<!-- Availability checkbox -->
				<div class="flex items-center gap-2">
					<Checkbox v-model="storeProductLocal.isAvailable" />
					<span>Товар доступен к продаже</span>
				</div>

				<!-- Table of sizes -->
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Размер</TableHead>
							<TableHead>Включён</TableHead>
							<TableHead>Базовая цена</TableHead>
							<TableHead>Цена в магазине</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<template v-if="hasSizes">
							<TableRow
								v-for="(size) in storeProductLocal.sizes"
								:key="size.id"
							>
								<TableCell>{{ size.name }}</TableCell>
								<TableCell>
									<Checkbox v-model="size.selected" />
								</TableCell>
								<TableCell>{{ size.basePrice }}</TableCell>
								<TableCell>
									<Input
										type="number"
										v-model.number="size.storePrice"
										class="w-24"
									/>
								</TableCell>
							</TableRow>
						</template>
						<template v-else>
							<TableRow>
								<TableCell
									colspan="4"
									class="text-center text-gray-500"
								>
									Нет размеров
								</TableCell>
							</TableRow>
						</template>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<!-- Footer (mobile only) -->
		<div class="flex justify-center items-center gap-2 md:hidden mt-4">
			<Button
				variant="outline"
				type="button"
				@click="onCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>
	</div>
</template>
