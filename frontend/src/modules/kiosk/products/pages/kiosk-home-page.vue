<template>
	<div class="relative flex flex-col flex-1 pt-safe">
		<!-- Toolbar for mobile view -->
		<KioskHomeToolbarMobile
			ref="toolbar"
			v-if="!categoriesLoading"
			:categories="categories"
			:selected-category-id="selectedCategoryId"
			:search-term="searchTerm"
			@update:category="onUpdateCategory"
			@update:search-term="onUpdateSearchTerm"
		/>

		<!-- Scrollable Products List -->
		<section
			class="relative flex-1 px-4 pb-4 overflow-y-auto no-scrollbar"
			ref="scrollContainer"
		>
			<!-- Search Mode -->
			<div v-if="searchTerm">
				<div v-if="isSearchProductsPending">
					<Skeleton class="h-20" />
				</div>
				<div
					v-else-if="searchProducts?.data?.length === 0"
					class="flex justify-center items-center h-20 text-gray-500"
				>
					<p class="text-lg">Ничего не найдено</p>
				</div>
				<div v-else>
					<div class="gap-3 grid grid-cols-2 sm:grid-cols-3">
						<KioskHomeProductCard
							v-for="product in searchProducts?.data"
							:key="product.id"
							:product="product"
						/>
					</div>
				</div>
			</div>
			<!-- Category Mode -->
			<div v-else>
				<div
					v-for="(category, index) in categories"
					:key="category.id"
				>
					<h2
						:ref="el => setCategoryRef(category.id, el)"
						:class="cn('mb-6 text-4xl', index !== 0 && 'mt-12')"
					>
						{{ category.name }}
					</h2>
					<div v-if="categoryProductsQueries[index].isPending">
						<Skeleton class="h-20" />
					</div>
					<div
						v-else-if="categoryProductsQueries[index].isError"
						class="text-red-500"
					>
						<p>Error loading products</p>
					</div>
					<div v-else>
						<div class="gap-3 grid grid-cols-2 sm:grid-cols-3">
							<KioskHomeProductCard
								v-for="product in categoryProducts[category.id]"
								:key="product.id"
								:product="product"
							/>
						</div>
					</div>
				</div>
			</div>
		</section>
	</div>
</template>

<script setup lang="ts">
import { Skeleton } from '@/core/components/ui/skeleton'
import { cn } from '@/core/utils/tailwind.utils'
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeToolbarMobile from '@/modules/kiosk/products/components/home/toolbar/kiosk-home-toolbar-mobile.vue'
import { useQueries, useQuery } from '@tanstack/vue-query'
import { useDebounceFn, useScroll } from '@vueuse/core'
import { computed, ref, useTemplateRef, watch, type ComponentPublicInstance } from 'vue'

/* Reactive References */
const selectedCategoryId = ref<number | null>(null)
const searchTerm = ref('')
const previousCategoryId = ref<number | null>(null)
const scrollContainer = useTemplateRef('scrollContainer')
const categoryRefs = ref<Record<number, HTMLElement>>({})
const toolbar = useTemplateRef<InstanceType<typeof KioskHomeToolbarMobile>>("toolbar")

/* Fetch Categories */
const { data: categories, isPending: categoriesLoading } = useQuery({
  queryKey: ['store-product-categories'],
  queryFn: () => storeProductsService.getStoreProductCategories(),
  refetchInterval: 20_000,
  initialData: []
})

/* Set Initial Selected Category */
watch(
  categories,
  (newCategories) => {
    if (newCategories.length > 0 && selectedCategoryId.value === null) {
      selectedCategoryId.value = newCategories[0].id
    }
  },
  { immediate: true }
)

/* Fetch Products for Each Category */
const categoryProductsQueries = useQueries({
  queries: computed(() =>
    categories.value.map((category) => ({
      queryKey: ['store-products', category.id],
      queryFn: () =>
        storeProductsService.getStoreProducts({
          categoryId: category.id,
          isAvailable: true,
          pageSize: 100
        }),
      enabled: !!category.id
    }))
  )
})

const categoryProducts = computed(() => {
  return categories.value.reduce((acc, category, index) => {
    acc[category.id] = categoryProductsQueries.value[index].data?.data || []
    return acc
  }, {} as Record<number, StoreProductDTO[]>)
})

/* Fetch Search Products */
const { data: searchProducts, isPending: isSearchProductsPending } = useQuery({
  queryKey: computed(() => ['search-products', searchTerm.value]),
  queryFn: () =>
    storeProductsService.getStoreProducts({
      search: searchTerm.value,
      isAvailable: true,
      pageSize: 100
    }),
  enabled: computed(() => !!searchTerm.value)
})

/* Category Refs & Scroll Handling */
function setCategoryRef(id: number, el: Element | ComponentPublicInstance | null) {
  let element: HTMLElement | null = null
  if (el) {
    if (el instanceof HTMLElement) {
      element = el
    } else if ('$el' in el && el.$el instanceof HTMLElement) {
      element = el.$el
    }
  }
  if (element) {
    categoryRefs.value[id] = element
  } else {
    delete categoryRefs.value[id]
  }
}

// Track scroll position
const { y: scrollY } = useScroll(scrollContainer)

// Update selected category based on scroll position
const updateSelectedCategory = useDebounceFn(() => {
  if (searchTerm.value) return // Skip if in search mode

  let closestId: number | null = null
  let maxOffsetTop = -Infinity

  for (const [id, el] of Object.entries(categoryRefs.value)) {
    const offsetTop = el.offsetTop
    if (offsetTop <= scrollY.value && offsetTop > maxOffsetTop) {
      maxOffsetTop = offsetTop
      closestId = Number(id)
    }
  }

  // Fallback to first category if none is found (e.g., at top)
  if (closestId === null && categories.value.length > 0) {
    closestId = categories.value[0].id
  }

  if (closestId !== null) {
    selectedCategoryId.value = closestId
  }
}, 100) // Debounce by 100ms for performance

// Watch scroll position changes
watch(scrollY, () => {
  updateSelectedCategory()
})

/* Toolbar Interactions */
function onUpdateCategory(categoryId: number) {
  searchTerm.value = ''
  selectedCategoryId.value = categoryId
  const el = categoryRefs.value[categoryId]
  if (el && scrollContainer.value) {
    const headerOffset = toolbar.value?.$el.offsetHeight || 0
    const targetScrollTop = el.offsetTop - headerOffset
    scrollTo({ top: targetScrollTop, behavior: 'smooth' })
  }
}

const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  if (newTerm.trim() !== '') {
    if (searchTerm.value.trim() === '') {
      previousCategoryId.value = selectedCategoryId.value
    }
    searchTerm.value = newTerm
  } else {
    searchTerm.value = ''
    const prevId = previousCategoryId.value
    previousCategoryId.value = null
    if (prevId && categoryRefs.value[prevId]) {
      categoryRefs.value[prevId].scrollIntoView({ behavior: 'smooth' })
    }
  }
}, 500)

function onUpdateSearchTerm(newSearchTerm: string) {
  debouncedEmitSearchTerm(newSearchTerm)
}
</script>

<style scoped lang="scss">
/* Additional styling if needed */
</style>
