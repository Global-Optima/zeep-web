<template>
	<!-- Main Content -->
	<div class="relative flex flex-col flex-1 pt-safe">
		<!-- Toolbar for mobile view -->
		<KioskHomeToolbarMobile
			v-if="!categoriesLoading"
			:categories="categories"
			:selected-category-id="selectedCategoryId"
			:search-term="searchTerm"
			@update:category="onUpdateCategory"
			@update:search-term="onUpdateSearchTerm"
		/>

		<!-- Products Grid Loading Skeleton -->
		<div
			v-if="isProductsPending"
			class="gap-3 grid grid-cols-2 sm:grid-cols-3 px-4"
		>
			<Skeleton
				v-for="(_, index) in skeletonsArray"
				:key="index"
				class="bg-slate-200 bg-opacity-80 rounded-[38px] w-full h-96"
			/>
		</div>

		<!-- Scrollable Products List -->
		<section
			v-else
			class="flex-1 px-4 pb-4 overflow-y-auto no-scrollbar"
			ref="scrollContainer"
		>
			<div
				v-if="!products || products.data.length === 0"
				class="flex justify-center items-center h-20 text-gray-500"
			>
				<p class="text-lg">Ничего не найдено</p>
			</div>

			<div
				v-else
				class="gap-3 grid grid-cols-2 sm:grid-cols-3"
			>
				<KioskHomeProductCard
					v-for="product in products.data"
					:key="product.id"
					:product="product"
				/>
			</div>
		</section>
	</div>
</template>

<script setup lang="ts">
import { Skeleton } from '@/core/components/ui/skeleton'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeToolbarMobile from '@/modules/kiosk/products/components/home/toolbar/kiosk-home-toolbar-mobile.vue'
import { useQuery } from '@tanstack/vue-query'
import { useDebounceFn } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

/* --------------------------
   Reactive References
--------------------------*/
const selectedCategoryId = ref<number | null>(null)
const searchTerm = ref('')
const previousCategoryId = ref<number | null>(null)
const skeletonsArray: number[] = new Array(9).fill(5)

/* --------------------------
   Fetch Categories
--------------------------*/
const { data: categories, isPending: categoriesLoading } = useQuery({
  queryKey: ['store-product-categories'],
  queryFn: () => storeProductsService.getStoreProductCategories(),
  initialData: []
})

watch(
  categories,
  (newCategories) => {
    if (newCategories && newCategories.length > 0 && selectedCategoryId.value === null && searchTerm.value === '') {
      selectedCategoryId.value = newCategories[0].id
    }
  },
  { immediate: true }
)

/* --------------------------
   Infinite Products Query
--------------------------*/
const {
  data: products,
  isPending: isProductsPending
} = useQuery({
  queryKey: computed(() => [
    'kiosk-products',
    { categoryId: selectedCategoryId.value, searchTerm: searchTerm.value },
  ]),
  queryFn: () =>
    storeProductsService.getStoreProducts({
      categoryId: selectedCategoryId.value!,
      search: searchTerm.value,
      isAvailable: true,
      pageSize: 100
    }),
  enabled: computed(() => Boolean(selectedCategoryId.value))
})


/* --------------------------
   Toolbar Interactions
--------------------------*/
function onUpdateCategory(categoryId: number) {
  searchTerm.value = ""
  selectedCategoryId.value = categoryId
}

const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  if (newTerm.trim() !== '') {
    if (searchTerm.value.trim() === '') {
      previousCategoryId.value = selectedCategoryId.value
      selectedCategoryId.value = null
    }
    searchTerm.value = newTerm
  } else {
    searchTerm.value = ''
    selectedCategoryId.value = previousCategoryId.value
    previousCategoryId.value = null
  }
}, 500)

function onUpdateSearchTerm(newSearchTerm: string) {
  debouncedEmitSearchTerm(newSearchTerm)
}
</script>

<style scoped lang="scss"></style>
