<script setup lang="ts">
import { ref, type Ref } from 'vue'

// UI Components (shadcn or your custom wrappers)
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/core/components/ui/card'
import { Checkbox } from '@/core/components/ui/checkbox'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/core/components/ui/table'

// Icons
import { ChevronLeft, Trash } from 'lucide-vue-next'

// Dialog for selecting products

// Types
import AdminSelectAvailableToAddProductDialog from '@/modules/admin/store-products/components/admin-select-available-to-add-product-dialog.vue'
import type { CreateStoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { ProductDetailsDTO } from '@/modules/kiosk/products/models/product.model'

// Props & Emits
const emits = defineEmits<{
  onSubmit: [dto: CreateStoreProductDTO[]],
  onCancel: []
}>()

/** Sub-item (size) data structure in your local state */
interface ProductSizeLocal {
  id: number
  name: string
  storePrice: number
  basePrice: number
  selected: boolean
}

/** Main product structure for your local state */
interface CreateProductsList {
  id: number
  name: string
  imageUrl: string
  categoryName: string
  isAvailable: boolean
  sizes: ProductSizeLocal[]
}

// Reactive array storing user-selected products
const productsList: Ref<CreateProductsList[]> = ref([])

// State to control the product selection dialog
const openProductDialog = ref(false)

/**
 * Submit the entire list as CreateStoreProductDTO[]
 */
function onSubmit() {
  // Map local data to final DTO
  const createDTO: CreateStoreProductDTO[] = productsList.value.map(product => {
    return {
      productId: product.id,
      isAvailable: product.isAvailable,
      productSizes: product.sizes
        .filter(s => s.selected) // Only include sizes that are selected
        .map(s => ({
          productSizeID: s.id,
          storePrice: s.storePrice
        })),
    }
  })

  emits('onSubmit', createDTO)
}

/**
 * Cancel - clear the list and emit "onCancel"
 */
function onCancel() {
  productsList.value = []
  emits('onCancel')
}

/**
 * Called when a product is selected in the dialog.
 * If the product already exists in the list, ignore it.
 */
function selectProduct(product: ProductDetailsDTO) {
  const exists = productsList.value.some(p => p.id === product.id)
  if (exists) {
    // You might optionally show a notification or toast here
    // e.g., toast("Этот товар уже добавлен.")
    openProductDialog.value = false
    return
  }

  // Otherwise, add a new product with sizes
  const dto: CreateProductsList = {
    id: product.id,
    name: product.name,
    imageUrl: product.imageUrl,
    categoryName: product.category.name,
    isAvailable: true,
    sizes: product.sizes.map(s => ({
      id: s.id,
      name: s.name,
      basePrice: s.basePrice,
      storePrice: s.basePrice,
      selected: true,
    }))
  }

  productsList.value.push(dto)
  openProductDialog.value = false
}

/** Remove an entire product card from the list */
function removeProduct(index: number) {
  productsList.value.splice(index, 1)
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
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
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Создать товар в магазине
			</h1>

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
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- List of selected products -->
		<div class="flex flex-col gap-2">
			<template
				v-for="(product, productIndex) in productsList"
				:key="product.id"
			>
				<Card>
					<CardHeader>
						<div class="flex justify-between">
							<div>
								<CardTitle>{{ product.name }}</CardTitle>
								<CardDescription class="mt-2">{{ product.categoryName }}</CardDescription>
							</div>
							<!-- Trash icon to remove entire product card -->
							<Button
								variant="ghost"
								size="icon"
								@click="removeProduct(productIndex)"
							>
								<Trash class="w-5 h-5 text-red-500 hover:text-red-700" />
							</Button>
						</div>
					</CardHeader>

					<CardContent>
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
								<template v-if="product.sizes.length > 0">
									<TableRow
										v-for="(size, sizeIndex) in product.sizes"
										:key="sizeIndex"
									>
										<TableCell>{{ size.name }}</TableCell>
										<TableCell>
											<Checkbox
												v-model="size.selected"
												:checked="size.selected"
												@update:checked="v => size.selected = v"
												class="mr-2"
											/>
										</TableCell>
										<TableCell>{{ size.basePrice }}</TableCell>
										<TableCell>
											<Input
												type="number"
												placeholder="Цена"
												v-model.number="size.storePrice"
												class="w-24"
											/>
										</TableCell>
									</TableRow>
								</template>
								<template v-else>
									<TableRow>
										<TableCell
											colspan="3"
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
			</template>

			<!-- Optional: If no products are selected yet -->
			<div
				v-if="productsList.length === 0"
				class="mt-4 text-center text-gray-500"
			>
				Нет добавленных товаров
			</div>
		</div>

		<!-- Button to open product selection dialog -->
		<div class="flex justify-center items-center my-4">
			<Button
				type="button"
				variant="outline"
				@click="openProductDialog = true"
			>
				Добавить товар
			</Button>
		</div>

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

		<!-- Dialog for selecting product -->
		<AdminSelectAvailableToAddProductDialog
			:open="openProductDialog"
			@close="openProductDialog = false"
			@select="selectProduct"
		/>
	</div>
</template>
