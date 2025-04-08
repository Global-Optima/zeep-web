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
						{{ errorMessage || 'К сожалению, данный продукт временно недоступен, попробуйте позже' }}
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
								alt="Изображение продукта"
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
										class="flex items-center gap-3 bg-primary p-6 rounded-full text-primary-foreground"
									>
										<Pencil class="size-6 sm:size-10" />
									</button>
								</div>
							</div>
						</div>
					</div>
					<!-- Additives Selection -->
					<div class="mt-10">
						<KioskDetailsAdditivesSection
							:categories="additiveCategories ?? []"
							:selected-additives="selectedAdditives"
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
import type { StoreAdditiveCategoryDTO, StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import type { CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import KioskProductRecipeDialog from '@/modules/kiosk/products/components/details/kiosk-product-recipe-dialog.vue'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft, Pencil } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

// ----- Props & Emits -----
const props = defineProps<{
  isOpen: boolean;
  cartItem: CartItem
}>()
const emits = defineEmits<{
  (e: 'close'): void;
  (e: 'update', updatedSize: StoreProductSizeDetailsDTO, updatedAdditives: StoreAdditiveCategoryItemDTO[]): void;
}>()

// ----- Local State -----
const productDetails = ref(props.cartItem.product)
const selectedSize = ref<StoreProductSizeDetailsDTO>(props.cartItem.size)
const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>(
  props.cartItem.additives.reduce((acc, additive) => {
    acc[additive.category.id] = acc[additive.category.id] || []
    acc[additive.category.id].push(additive)
    return acc
  }, {} as Record<number, StoreAdditiveCategoryItemDTO[]>)
)

const additiveCategories = ref<StoreAdditiveCategoryDTO[]>([])
const errorMessage = ref<string | null>(null)

// ----- Fetch Additive Categories -----
// Use Vue Query to fetch additive categories based on the selected size.
const { data: fetchedAdditives, isPending: isFetching, isError } = useQuery({
  queryKey: computed(() => ['kiosk-additive-categories', selectedSize.value?.id]),
  queryFn: () => {
    if (selectedSize.value) {
      return storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id)
    }
  },
  initialData: [] as StoreAdditiveCategoryDTO[],
  enabled: computed(() => !!selectedSize.value)
})

// Watch fetchedAdditives and enforce required auto‑selection.
watch(
  fetchedAdditives,
  (newAdditives) => {
    if (newAdditives) {
      additiveCategories.value = newAdditives
      newAdditives.forEach(category => {
        if (category.isRequired) {
          const current = selectedAdditives.value[category.id] || []
          if (current.length === 0 && category.additives.length > 0) {
            const defaultAdditive = category.additives.find(a => a.isDefault)
            selectedAdditives.value[category.id] = defaultAdditive ? [defaultAdditive] : [category.additives[0]]
          }
        }
      })
    }
  },
  { immediate: true }
)

// ----- Computed Values -----
const sortedSizes = computed(() => {
  return productDetails.value?.sizes
    ? [...productDetails.value.sizes].sort((a, b) => a.size - b.size)
    : []
})

// Ensure a size is selected when sortedSizes updates.
watch(
  sortedSizes,
  (newSizes) => {
    if (!selectedSize.value && newSizes.length > 0) {
      selectedSize.value = newSizes[0]
    }
  },
  { immediate: true }
)

const totalPrice = computed(() => {
  if (!selectedSize.value) return 0
  const basePrice = selectedSize.value.storePrice
  const additivesPrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, additive) => sum + additive.storePrice, 0)
  return basePrice + additivesPrice
})

// ----- Helper Functions -----
const onSizeSelect = (size: StoreProductSizeDetailsDTO) => {
  if (selectedSize.value?.id === size.id) return
  selectedSize.value = size
  selectedAdditives.value = {}
}

const onAdditiveToggle = (category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO) => {
  const current = selectedAdditives.value[category.id] || []
  const alreadySelected = current.some(a => a.additiveId === additive.additiveId)
  if (alreadySelected) {
    // Prevent unselecting if required category has only one selected or if additive is default.
    if (category.isRequired && current.length === 1) return
    if (additive.isDefault) return
    selectedAdditives.value[category.id] = current.filter(a => a.additiveId !== additive.additiveId)
  } else {
    if (!category.isMultipleSelect) {
      selectedAdditives.value[category.id] = [additive]
    } else {
      selectedAdditives.value[category.id] = [...current, additive]
    }
  }
}

// ----- Action Handler -----
const handleUpdate = () => {
  const updatedAdditives = Object.values(selectedAdditives.value).flat()
  emits('update', selectedSize.value, updatedAdditives)
  emits('close')
}
</script>

<style scoped>
/* Additional styles if needed */
</style>
