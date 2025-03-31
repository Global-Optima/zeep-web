<template>
	<Dialog
		:open="isOpen"
		@update:open="() => emits('close')"
	>
		<DialogContent
			:include-close-button="false"
			class="p-0 border-none !rounded-[52px] max-w-[80vw] h-[95vh] overflow-clip"
		>
			<div class="relative bg-[#F3F4F9] pb-32 overflow-y-scroll text-black no-scrollbar">
				<!-- Loading State -->
				<PageLoader v-if="isFetching" />
				<!-- Error State -->
				<div
					v-else-if="isError || errorMessage"
					class="flex flex-col justify-center items-center p-6 w-full h-screen text-center"
				>
					<h1 class="mt-8 font-bold text-red-600 text-4xl">Ошибка</h1>
					<p class="mt-6 max-w-md text-gray-600 text-2xl">
						{{ errorMessage || 'К сожалению, данный товар временно недоступен, попробуйте позже' }}
					</p>
					<button
						@click="emits('close')"
						class="flex justify-center items-center bg-slate-200 mt-8 px-8 py-5 rounded-3xl h-14 text-slate- text-2xl"
					>
						Вернуться назад
					</button>
				</div>
				<!-- Product Content -->
				<div
					v-else
					class="pb-44"
				>
					<div class="bg-white shadow-gray-200 shadow-xl px-8 pb-6 rounded-b-[48px] w-full">
						<header class="flex justify-between items-center gap-6 pt-6">
							<Button
								size="icon"
								variant="ghost"
								class="size-14"
								@click="emits('close')"
							>
								<ChevronLeft
									class="size-14 text-gray-400"
									stroke-width="1.6"
								/>
							</Button>
							<template v-if="selectedSize">
								<KioskProductRecipeDialog :nutrition="selectedSize.totalNutrition" />
							</template>
						</header>
						<div class="flex flex-col justify-center items-center">
							<LazyImage
								:src="productDetails.imageUrl"
								alt="Изображение товара"
								class="rounded-3xl w-38 h-64 object-contain"
							/>
							<p class="mt-7 font-semibold text-4xl">{{ productDetails.name }}</p>
							<p class="mt-2 text-slate-600 text-xl">{{ productDetails.description }}</p>
						</div>
						<!-- Sticky Section -->
						<div class="top-0 z-10 sticky bg-white mt-8 pb-6">
							<div class="flex justify-between items-center gap-4">
								<!-- Size Selection -->
								<div class="flex items-center gap-2 overflow-x-auto no-scrollbar">
									<KioskDetailsSizes
										v-for="size in sortedSizes"
										:key="size.id"
										:size="size"
										:is-selected="selectedSize?.id === size.id"
										@click:size="onSizeSelect"
									/>
								</div>
								<!-- Update Button with Price -->
								<div class="flex items-center gap-6">
									<p class="font-medium text-4xl">{{ formatPrice(totalPrice) }}</p>
									<button
										@click="handleUpdate"
										class="flex items-center gap-3 bg-primary p-5 rounded-full text-primary-foreground"
									>
										<Pencil class="size-8 sm:size-12" />
									</button>
								</div>
							</div>
						</div>
					</div>
					<!-- Additives Selection -->
					<div class="mt-10">
						<KioskDetailsAdditivesSection
							:categories="additiveCategories ?? []"
							:isAdditiveSelected="isAdditiveSelected"
							@toggle-additive="onAdditiveToggle"
						/>
					</div>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent } from '@/core/components/ui/dialog'
import { formatPrice } from '@/core/utils/price.utils'
import type {
  StoreAdditiveCategoryDTO,
  StoreAdditiveCategoryItemDTO
} from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { type CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import KioskProductRecipeDialog from '@/modules/kiosk/products/components/details/kiosk-product-recipe-dialog.vue'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft, Pencil } from 'lucide-vue-next'
import { computed, ref, watchEffect } from 'vue'

// ----- Props & Emits -----
const props = defineProps<{
  isOpen: boolean;
  cartItem: CartItem;
}>()
const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'update', updatedSize: StoreProductSizeDetailsDTO, updatedAdditives: StoreAdditiveCategoryItemDTO[]): void;
}>()


// ----- Local State -----
// Use the product data from the passed cartItem.
const productDetails = ref(props.cartItem.product)
// Initialize selected size and additives from the passed cartItem.
const selectedSize = ref<StoreProductSizeDetailsDTO>(props.cartItem.size)
const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>(
  props.cartItem.additives.reduce((acc, additive) => {
    acc[additive.category.id] = acc[additive.category.id] || []
    acc[additive.category.id].push(additive)
    return acc
  }, {} as Record<number, StoreAdditiveCategoryItemDTO[]>)
)

// This will hold the additive categories fetched for the selected size.
const additiveCategories = ref<StoreAdditiveCategoryDTO[]>([])
const errorMessage = ref<string | null>(null)

// ----- Fetch Additive Categories -----
// Use Vue Query to fetch additive categories for the current size.
const { data: fetchedAdditives, isPending: isFetching, isError } = useQuery({
  queryKey: computed(() => ['kiosk-additive-categories', selectedSize.value?.id]),
  queryFn: () => {
    if (selectedSize.value) {
      return storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id)
    }
  },
  initialData: [],
  enabled: computed(() => !!selectedSize.value)
})

// When new additive data is fetched, merge the current selections.
watchEffect(() => {
  if (fetchedAdditives.value) {
    additiveCategories.value = fetchedAdditives.value.map(category => ({
      ...category,
      additives: category.additives.map(additive => ({
        ...additive,
        isSelected: isAdditiveSelected(category, additive.additiveId)
      }))
    }))
  }
})

// ----- Computed Values -----
const sortedSizes = computed(() => {
  return productDetails.value?.sizes ? [...productDetails.value.sizes].sort((a, b) => a.size - b.size) : []
})

watchEffect(() => {
  // Ensure a size is selected; if not, select the first one.
  if (!selectedSize.value && sortedSizes.value.length > 0) {
    selectedSize.value = sortedSizes.value[0]
  }
})

const totalPrice = computed(() => {
  if (!selectedSize.value) return 0
  const basePrice = selectedSize.value.storePrice
  const additivesPrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, additive) => sum + additive.storePrice, 0)
  return basePrice + additivesPrice
})

// ----- Helper Functions -----
function isAdditiveSelected(category: StoreAdditiveCategoryDTO, additiveId: number): boolean {
  return selectedAdditives.value[category.id]?.some(a => a.additiveId === additiveId) || false
}

const onSizeSelect = (size: StoreProductSizeDetailsDTO) => {
  if (selectedSize.value?.id === size.id) return
  selectedSize.value = size
  selectedAdditives.value = {}
}

const onAdditiveToggle = (category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO) => {
  const current = selectedAdditives.value[category.id] || []
  const alreadySelected = current.some(a => a.additiveId === additive.additiveId)
  if (category.isMultipleSelect) {
    selectedAdditives.value[category.id] = alreadySelected
      ? current.filter(a => a.additiveId !== additive.additiveId)
      : [...current, additive]
  } else {
    selectedAdditives.value[category.id] = alreadySelected ? [] : [additive]
  }
}

// ----- Action Handler -----
// When the user clicks the update button, emit the updated size and additives.
const handleUpdate = () => {
  const updatedAdditives = Object.values(selectedAdditives.value).flat()
  emits('update', selectedSize.value, updatedAdditives)
  emits('close')
}
</script>

<style scoped>
/* You can add additional styles if needed */
</style>
