<template>
	<div class="flex h-screen sm:flex-row flex-col pt-safe">
		<!-- Sidebar for tablet and larger screens -->
		<div class="pl-4 py-4">
			<KioskHomeSidebarTablet
				@update:category="onUpdateCategory"
				:categories="categories"
				:selected-category-id="selectedCategoryId"
			/>
		</div>

		<!-- Main Content -->
		<div class="flex-1 flex flex-col">
			<!-- Toolbar for mobile view -->
			<KioskHomeToolbarMobile
				v-if="!categoriesLoading"
				class="block sm:hidden"
				:categories="categories"
				:selected-category-id="selectedCategoryId"
				:search-term="searchTerm"
				@update:category="onUpdateCategory"
				@update:search-term="onUpdateSearchTerm"
			/>
			<div
				v-else
				class="w-full py-4 sm:py-6 px-4 flex items-center gap-2 overflow-x-auto no-scrollbar sticky top-0 z-10 sm:hidden"
			>
				<Skeleton
					v-for="n in 4"
					:key="n"
					class="h-16 w-32 rounded-full bg-gray-200"
				/>
			</div>

			<!-- Search Bar for tablet and larger screens -->
			<div class="hidden sm:block px-4 pt-4">
				<KioskHomeToolbarTablet
					:search-term="searchTerm"
					@update:search-term="onUpdateSearchTerm"
				/>
			</div>

			<!-- Products Grid -->
			<section class="flex-1 p-4 overflow-y-auto">
				<div
					v-if="productsLoading"
					class="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-4 gap-2 sm:gap-4"
				>
					<Skeleton
						v-for="n in 8"
						:key="n"
						class="rounded-lg w-full h-48 bg-gray-200"
					/>
				</div>

				<div
					v-else-if="products.length === 0"
					class="flex items-center justify-center h-20 text-gray-500"
				>
					<p class="text-lg">Ничего не найдено</p>
				</div>

				<div
					v-else
					class="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-3 gap-2 sm:gap-4"
				>
					<KioskHomeProductCard
						v-for="product in products"
						:key="product.id"
						:product="product"
					/>
				</div>
			</section>
		</div>

		<!-- Cart Button for mobile -->
		<div class="fixed bottom-6 right-6 flex justify-center">
			<KioskHomeCart />
		</div>
	</div>
</template>

<script setup lang="ts">
import { Skeleton } from '@/core/components/ui/skeleton'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeSidebarTablet from '@/modules/kiosk/products/components/home/kiosk-home-sidebar-tablet.vue'
import KioskHomeToolbarTablet from '@/modules/kiosk/products/components/home/kiosk-home-toolbar-tablet.vue'
import type { ProductCategory, StoreProducts } from '@/modules/kiosk/products/models/product.model'
import { productService } from '@/modules/kiosk/products/services/products.service'
import { useQuery } from '@tanstack/vue-query'
import { useDebounceFn } from '@vueuse/core'
import { computed, defineAsyncComponent, ref, watch } from 'vue'

// Asynchronous Components
const KioskHomeCart = defineAsyncComponent(() =>
  import('@/modules/kiosk/products/components/home/kiosk-home-cart.vue')
)

const KioskHomeToolbarMobile = defineAsyncComponent(() =>
  import('@/modules/kiosk/products/components/home/kiosk-home-toolbar-mobile.vue')
)

// Reactive References
const selectedCategoryId = ref<number | null>(null)
const searchTerm = ref('')
const previousCategoryId = ref<number | null>(null)

// Query Keys
const productsQueryKey = computed(() => [
  'products',
  { categoryId: selectedCategoryId.value, searchTerm: searchTerm.value },
])

const categoriesQueryKey = ['categories']

// Fetch Categories
const { data: categories, isLoading: categoriesLoading } = useQuery<ProductCategory[]>({
  queryKey: categoriesQueryKey,
  queryFn: () => productService.getStoreCategories(),
  initialData: []
})

// Watch Categories to Set Default Selection
watch(
  categories,
  (newCategories) => {
    if (newCategories && newCategories.length > 0 && selectedCategoryId.value === null && searchTerm.value === '') {
      selectedCategoryId.value = newCategories[0].id
    }
  },
  { immediate: true }
)

// Fetch Products
const { data: products, isLoading: productsLoading } = useQuery<StoreProducts[]>({
  queryKey: productsQueryKey,
  queryFn: () =>
    productService.getStoreProducts(
      selectedCategoryId.value,
      searchTerm.value
    ),
  enabled: computed(() => Boolean(selectedCategoryId.value) || searchTerm.value.trim() !== ''),
  initialData: []
})

// Handle Category Update
function onUpdateCategory(categoryId: number) {
  // If a search is active, ignore category updates
  if (searchTerm.value.trim() !== '') return

  selectedCategoryId.value = categoryId
}

// Debounced Search Term Update
const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  if (newTerm.trim() !== '') {
    if (searchTerm.value.trim() === '') {
      // Store the current category before searching
      previousCategoryId.value = selectedCategoryId.value
      // Unselect the category
      selectedCategoryId.value = null
    }
    searchTerm.value = newTerm
  } else {
    // Restore the previous category when search is cleared
    searchTerm.value = ''
    selectedCategoryId.value = previousCategoryId.value
    previousCategoryId.value = null
  }
}, 500)

// Handle Search Term Update
function onUpdateSearchTerm(newSearchTerm: string) {
  debouncedEmitSearchTerm(newSearchTerm)
}
</script>

<style scoped lang="scss">
/* Add any necessary scoped styles here */
</style>
