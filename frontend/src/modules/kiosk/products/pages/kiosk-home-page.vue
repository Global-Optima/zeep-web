<template>
	<div class="flex sm:flex-row flex-col pt-safe h-screen">
		<!-- Sidebar for tablet and larger screens -->
		<div class="sm:block hidden py-4 pl-4">
			<KioskHomeSidebarTablet
				@update:category="onUpdateCategory"
				:categories="categories?.data ?? []"
				:selected-category-id="selectedCategoryId"
			/>
		</div>

		<!-- Main Content -->
		<div class="flex flex-col flex-1">
			<!-- Toolbar for mobile view -->
			<KioskHomeToolbarMobile
				v-if="!categoriesLoading"
				class="block sm:hidden"
				:categories="categories?.data ?? []"
				:selected-category-id="selectedCategoryId"
				:search-term="searchTerm"
				@update:category="onUpdateCategory"
				@update:search-term="onUpdateSearchTerm"
			/>

			<!-- Search Bar for tablet and larger screens -->
			<div class="sm:block hidden pt-4 pr-4 pl-3">
				<KioskHomeToolbarTablet
					:search-term="searchTerm"
					@update:search-term="onUpdateSearchTerm"
				/>
			</div>

			<!-- Products Grid -->
			<section class="flex-1 pt-3 pr-4 pb-4 pl-3 overflow-y-auto no-scrollbar">
				<div
					v-if="products?.data.length === 0"
					class="flex justify-center items-center h-20 text-gray-500"
				>
					<p class="text-lg">Ничего не найдено</p>
				</div>

				<div
					v-else
					class="gap-2 sm:gap-2 grid grid-cols-2 sm:grid-cols-3"
				>
					<KioskHomeProductCard
						v-for="product in products?.data ?? []"
						:key="product.id"
						:product="product"
					/>
				</div>
			</section>
		</div>

		<!-- Cart Button for mobile -->
		<div
			v-if="!cartStore.isEmpty"
			class="right-8 bottom-8 fixed flex justify-center"
		>
			<KioskHomeCart />
		</div>
	</div>
</template>

<script setup lang="ts">
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeSidebarTablet from '@/modules/kiosk/products/components/home/kiosk-home-sidebar-tablet.vue'
import KioskHomeToolbarTablet from '@/modules/kiosk/products/components/home/kiosk-home-toolbar-tablet.vue'
import { productsService } from '@/modules/kiosk/products/services/products.service'
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

const cartStore = useCartStore()

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
// TODO: Make it more button or unpaged
const { data: categories, isLoading: categoriesLoading } = useQuery({
  queryKey: categoriesQueryKey,
  queryFn: () => productsService.getAllProductCategories({pageSize: 1000}),
})

// Watch Categories to Set Default Selection
watch(
  categories,
  (newCategories) => {
    if (newCategories && newCategories.data.length > 0 && selectedCategoryId.value === null && searchTerm.value === '') {
      selectedCategoryId.value = newCategories.data[0].id
    }
  },
  { immediate: true }
)


// Fetch Products
const { data: products } = useQuery({
  queryKey: productsQueryKey,
  queryFn: () =>
    storeProductsService.getStoreProducts({
       categoryId: selectedCategoryId.value!, search: searchTerm.value,

    }),
  enabled: computed(() => Boolean(selectedCategoryId.value) || searchTerm.value.trim() !== ''),
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
