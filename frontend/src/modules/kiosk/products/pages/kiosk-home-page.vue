<template>
	<div
		class="z-0 pt-32 h-screen overflow-y-auto no-scrollbar"
		ref="scrollContainer"
	>
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
		<section class="px-8 pb-4">
			<!-- Search Mode -->
			<div v-if="searchTerm">
				<div v-if="isSearchProductsPending">
					<div class="gap-4 grid grid-cols-2 sm:grid-cols-3">
						<Skeleton
							v-for="n in 6"
							:key="n"
							class="bg-slate-200 bg-opacity-80 rounded-[38px] w-full h-[404px]"
						/>
					</div>
				</div>
				<div
					v-else-if="searchProducts?.data?.length === 0"
					class="flex justify-center items-center h-20 text-gray-500"
				>
					<p class="text-lg">Ничего не найдено</p>
				</div>
				<div v-else>
					<div class="gap-4 grid grid-cols-2 sm:grid-cols-3">
						<KioskHomeProductCard
							v-for="product in sortedSearchProducts"
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
					class="mb-16"
				>
					<h2
						:ref="el => setCategoryRef(category.id, el)"
						class="my-8 px-4 text-4xl"
					>
						{{ category.name }}
					</h2>
					<!-- Show skeleton grid if products are still loading -->
					<div v-if="categoryProductsQueries[index].isPending">
						<div class="gap-4 grid grid-cols-2 sm:grid-cols-3">
							<Skeleton
								v-for="n in 6"
								:key="n"
								class="bg-slate-200 bg-opacity-80 rounded-[38px] w-full h-[404px]"
							/>
						</div>
					</div>
					<!-- Show products if loaded -->
					<div v-else>
						<div class="gap-4 grid grid-cols-2 sm:grid-cols-3">
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
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeToolbarMobile from '@/modules/kiosk/products/components/home/toolbar/kiosk-home-toolbar-mobile.vue'
import { useQueries, useQuery } from '@tanstack/vue-query'
import { useDebounceFn, useScroll } from '@vueuse/core'
import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  shallowRef,
  watch,
  watchEffect,
  type ComponentPublicInstance
} from 'vue'

/* --- Reactive State --- */
// Use shallowRef for refs that don't need deep reactivity
const scrollContainer = shallowRef<HTMLElement | null>(null)
const toolbar = shallowRef<InstanceType<typeof KioskHomeToolbarMobile> | null>(null)

// Use ref for primitive values
const searchTerm = ref('')
const selectedCategoryId = ref<number | null>(null)
const categoryRefs = ref<Record<number, HTMLElement>>({})
const isManualScroll = ref(false)
const scrollEndTimeout = ref<number | null>(null)

/* --- Fetch Categories --- */
const { data: categories, isPending: categoriesLoading } = useQuery({
  queryKey: ['store-product-categories'],
  queryFn: () => storeProductsService.getStoreProductCategories(),
  refetchInterval: 20_000,
  initialData: []
})

// Initialize selected category when categories are loaded
watchEffect(() => {
  if (categories.value.length > 0 && selectedCategoryId.value === null) {
    selectedCategoryId.value = categories.value[0].id
  }
})

/* --- Fetch Products for Each Category --- */
const categoryProductsQueries = useQueries({
  queries: computed(() =>
    categories.value.map((category) => ({
      queryKey: ['kiosk-products', category.id],
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

// Use computed for derived state
const categoryProducts = computed(() => {
  return categories.value.reduce((acc, category, index) => {
    const products = categoryProductsQueries.value[index].data?.data || []
    acc[category.id] = products.slice().sort((a, b) => {
      if (a.isOutOfStock === b.isOutOfStock) return 0
      return a.isOutOfStock ? 1 : -1
    })
    return acc
  }, {} as Record<number, StoreProductDTO[]>)
})

/* --- Fetch Search Products --- */
const { data: searchProducts, isPending: isSearchProductsPending } = useQuery({
  queryKey: computed(() => ['kiosk-products', searchTerm.value]),
  queryFn: () =>
    storeProductsService.getStoreProducts({
      search: searchTerm.value,
      isAvailable: true,
      pageSize: 100
    }),
  enabled: computed(() => !!searchTerm.value)
})

const sortedSearchProducts = computed(() => {
  if (!searchProducts.value?.data) {
    return []
  }

  return searchProducts.value.data.slice().sort((a, b) => {
    if (a.isOutOfStock === b.isOutOfStock) return 0
    return a.isOutOfStock ? 1 : -1
  })
})

const { y: scrollY } = useScroll(
  () => scrollContainer.value,
)

// Function to update active category based on scroll position
const updateActiveCategory = () => {
  // Skip if in search mode or manual scrolling
  if (searchTerm.value || isManualScroll.value) return
  const toolbarHeight = toolbar.value?.$el?.offsetHeight || 0
  let closestId: number | null = null
  let minDistance = Infinity

  // Find the category closest to the top of the viewport
  Object.entries(categoryRefs.value).forEach(([id, el]) => {
    const rect = el.getBoundingClientRect()
    const distance = Math.abs(rect.top - toolbarHeight)

    if (distance < minDistance) {
      minDistance = distance
      closestId = Number(id)
    }
  })

  // Update selected category if it changed
  if (closestId !== null && closestId !== selectedCategoryId.value) {
    selectedCategoryId.value = closestId
  }
}

// Watch scroll position to update active category
watch(scrollY, () => {
  if (!isManualScroll.value) {
    updateActiveCategory()
  }
})

/* --- Category Reference Management --- */
const setCategoryRef = (id: number, el: Element | ComponentPublicInstance | null) => {
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

/* --- Toolbar Interaction: Scroll to a Specific Category Section --- */
const onUpdateCategory = (categoryId: number) => {
  // Clear search mode
  searchTerm.value = ''

  // Set flag to ignore scroll events during programmatic scrolling
  isManualScroll.value = true

  // Update category ID immediately
  selectedCategoryId.value = categoryId

  // Clear any existing timeout
  if (scrollEndTimeout.value !== null) {
    window.clearTimeout(scrollEndTimeout.value)
  }

  // Scroll to the selected category
  const el = categoryRefs.value[categoryId]
  if (el && scrollContainer.value) {
    const toolbarHeight = toolbar.value?.$el?.offsetHeight || 0

    scrollContainer.value.scrollTo({
      top: el.offsetTop - toolbarHeight - 20, // Add some padding
      behavior: 'smooth'
    })

    // Reset the manual scroll flag after animation completes
    scrollEndTimeout.value = window.setTimeout(() => {
      isManualScroll.value = false
      scrollEndTimeout.value = null
    }, 1000) // Typical smooth scroll duration
  } else {
    // If no scroll happened, reset the flag immediately
    isManualScroll.value = false
  }
}

/* --- Search Term Interaction --- */
const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  searchTerm.value = newTerm.trim() !== '' ? newTerm : ''
}, 300) // Reduced debounce time for more responsive search

const onUpdateSearchTerm = (newSearchTerm: string) => {
  debouncedEmitSearchTerm(newSearchTerm)
}

/* --- Touch Event Handling for Horizontal Swipe Gestures --- */
// Use shallowRefs for touch data
const touchStartX = shallowRef(0)
const touchStartY = shallowRef(0)

const handleTouchStart = (e: TouchEvent) => {
  if (e.touches.length > 0) {
    const touch = e.touches[0]
    touchStartX.value = touch.clientX
    touchStartY.value = touch.clientY
  }
}

const handleTouchEnd = (e: TouchEvent) => {
  if (e.changedTouches.length > 0) {
    const touch = e.changedTouches[0]
    const deltaX = touch.clientX - touchStartX.value
    const deltaY = touch.clientY - touchStartY.value
    const threshold = 50 // pixels

    // If horizontal swipe is dominant
    if (Math.abs(deltaX) > threshold && Math.abs(deltaX) > Math.abs(deltaY)) {
      if (deltaX < 0) {
        goToNextSection()
      } else {
        goToPreviousSection()
      }
    }
  }
}

const goToNextSection = () => {
  const currentIndex = categories.value.findIndex(cat => cat.id === selectedCategoryId.value)
  if (currentIndex >= 0 && currentIndex < categories.value.length - 1) {
    const nextCat = categories.value[currentIndex + 1]
    onUpdateCategory(nextCat.id)
  }
}

const goToPreviousSection = () => {
  const currentIndex = categories.value.findIndex(cat => cat.id === selectedCategoryId.value)
  if (currentIndex > 0) {
    const prevCat = categories.value[currentIndex - 1]
    onUpdateCategory(prevCat.id)
  }
}

/* --- Lifecycle Hooks --- */
onMounted(() => {
  // Set up event listeners
  if (scrollContainer.value) {
    scrollContainer.value.addEventListener('touchstart', handleTouchStart)
    scrollContainer.value.addEventListener('touchend', handleTouchEnd)
  }

  // Initial update after DOM is ready
  nextTick(() => {
    updateActiveCategory()
  })
})

onBeforeUnmount(() => {
  // Clean up event listeners
  if (scrollContainer.value) {
    scrollContainer.value.removeEventListener('touchstart', handleTouchStart)
    scrollContainer.value.removeEventListener('touchend', handleTouchEnd)
  }

  // Clear any pending timeouts
  if (scrollEndTimeout.value !== null) {
    window.clearTimeout(scrollEndTimeout.value)
  }
})
</script>

<style scoped lang="scss"></style>
