<template>
	<div class="relative bg-[#F3F4F9] pt-safe text-black no-scrollbar">
		<!-- Loading State -->
		<PageLoader v-if="isFetching" />
		<!-- Error State -->
		<div
			v-else-if="isError || !isAddButtonEnabled"
			class="flex flex-col justify-center items-center p-6 w-full h-screen text-center"
		>
			<h1 class="mt-8 font-bold text-red-600 text-4xl">Ошибка</h1>
			<p class="mt-6 max-w-md text-gray-600 text-2xl">
				К сожалению, данный продукт временно недоступен, попробуйте позже
			</p>

			<button
				@click="onBackClick"
				class="flex justify-center items-center bg-slate-200 mt-8 px-8 py-5 rounded-3xl h-14 text-slate- text-2xl"
			>
				Вернуться назад
			</button>
		</div>

		<!-- Product Content -->
		<div
			v-else-if="productDetails"
			class="pb-44"
		>
			<!-- Non-sticky content -->
			<div class="bg-white shadow-gray-200 shadow-xl px-8 pb-6 w-full">
				<header class="flex justify-between items-center gap-6 pt-6">
					<Button
						size="icon"
						variant="ghost"
						class="size-14"
						@click="onBackClick"
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
				<div class="flex flex-col justify-center items-center mt-4">
					<LazyImage
						:src="productDetails.imageUrl"
						alt="Изображение продукта"
						class="rounded-3xl w-38 h-64 object-contain"
					/>
					<p class="mt-7 font-semibold text-5xl">{{ productDetails.name }}</p>
					<p class="mt-4 text-slate-600 text-2xl">
						{{ productDetails.description }}
					</p>
				</div>
			</div>

			<!-- Sticky Section -->
			<div class="top-0 z-10 sticky bg-white shadow-lg px-8 pt-4 pb-10 rounded-b-[52px]">
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

					<!-- Add to Cart Button + Price -->
					<div class="flex items-center gap-8">
						<p class="font-medium text-5xl truncate">
							{{ formatPrice(totalPrice) }}
						</p>
						<button
							@click="handleAddToCart"
							:disabled="!isAddButtonEnabled"
							class="flex items-center gap-3 bg-primary disabled:bg-muted p-8 rounded-full text-primary-foreground disabled:text-muted-foreground"
						>
							<Plus class="size-8 sm:size-12" />
						</button>
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
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryDTO, StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsAdditivesSection from '@/modules/kiosk/products/components/details/kiosk-details-additives-section.vue'
import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue'
import KioskProductRecipeDialog from '@/modules/kiosk/products/components/details/kiosk-product-recipe-dialog.vue'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft, Plus } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()

const productId = computed(() => Number(route.params.id))
function onBackClick() {
  router.push({ name: getRouteName('KIOSK_HOME') })
}

const selectedSize = ref<StoreProductSizeDetailsDTO | null>(null)
const selectedAdditives = ref<Record<number, StoreAdditiveCategoryItemDTO[]>>({})

const {
  data: productDetails,
  isPending: isFetching,
  isError = true,
} = useQuery({
  queryKey: computed(() => ['kiosk-product-details', productId]),
  queryFn: () => storeProductsService.getStoreProduct(productId.value),
  enabled: computed(() => productId.value > 0),
  refetchInterval: 20_000,
})

const sortedSizes = computed(() => {
  return productDetails.value?.sizes
    ? productDetails.value.sizes.slice().sort((a, b) => a.size - b.size)
    : []
})

watch(
  sortedSizes,
  (newSizes) => {
    if (!selectedSize.value && newSizes.length > 0) {
      selectedSize.value = newSizes[0]
    }
  },
  { immediate: true }
)

const { data: additiveCategories } = useQuery({
  queryKey: computed(() => ['kiosk-additive-categories', selectedSize.value]),
  queryFn: () => {
    if (selectedSize.value) {
      return storeAdditivesService.getStoreAdditiveCategories(selectedSize.value.id)
    }
  },
  enabled: computed(() => !!selectedSize.value),
  refetchInterval: 20_000,
})

watch(
  additiveCategories,
  (newCategories) => {
    newCategories?.forEach((category) => {
      const current = selectedAdditives.value[category.id] || []

      if (category.isRequired && current.length === 0) {
        const nonDefaultAdditive = category.additives.find(a => !a.isDefault)
        if (nonDefaultAdditive) {
          selectedAdditives.value[category.id] = [nonDefaultAdditive]
        } else {
          selectedAdditives.value[category.id] = []
        }
      }
    })
  },
  { immediate: true }
)

const isAddButtonEnabled = computed(() => {
  if (productDetails.value?.isOutOfStock) return false

  const hasDefaultOutOfStock = additiveCategories?.value?.some((category) =>
    category.additives.some((additive) => additive.isDefault && additive.isOutOfStock)
  )

  return !hasDefaultOutOfStock
})

const totalPrice = computed(() => {
  if (!selectedSize.value) return 0

  const storePrice = selectedSize.value.storePrice

  const additivePrice = Object.values(selectedAdditives.value)
    .flat()
    .reduce((sum, add) => {
      if (add.isDefault) {
        return sum
      }
      return sum + add.storePrice
    }, 0)

  return storePrice + additivePrice
})

function onSizeSelect(size: StoreProductSizeDetailsDTO) {
  if (selectedSize.value?.id === size.id) return
  selectedSize.value = size
  selectedAdditives.value = {}
}

function onAdditiveToggle(category: StoreAdditiveCategoryDTO, additive: StoreAdditiveCategoryItemDTO) {
  if (additive.isDefault) {
    return
  }

  const current = selectedAdditives.value[category.id] || []
  const isSelected = current.some((a) => a.additiveId === additive.additiveId)

  if (isSelected) {
    if (category.isRequired && current.length === 1) {
      return
    }

    selectedAdditives.value[category.id] = current.filter(
      (a) => a.additiveId !== additive.additiveId
    )
  } else {
    if (!category.isMultipleSelect) {
      selectedAdditives.value[category.id] = [additive]
    } else {
      selectedAdditives.value[category.id] = [...current, additive]
    }
  }
}

function handleAddToCart() {
  if (!productDetails.value || !selectedSize.value) return
  const allAdditives = Object.values(selectedAdditives.value).flat()
  cartStore.addToCart(productDetails.value, selectedSize.value, allAdditives, 1)
}
</script>

<style scoped></style>
