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
			ref="scrollContainer"
			class="relative flex-1 px-8 pb-4 overflow-y-auto no-scrollbar"
		>
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
import { computed, onBeforeUnmount, onMounted, ref, watch, type ComponentPublicInstance } from 'vue'

/* --- Reactive References --- */
const selectedCategoryId = ref<number | null>(null)
const searchTerm = ref('')
const scrollContainer = ref<HTMLElement | null>(null)
const categoryRefs = ref<Record<number, HTMLElement>>({})
const toolbar = ref<InstanceType<typeof KioskHomeToolbarMobile> | null>(null)

/* --- Fetch Categories --- */
const { data: categories, isPending: categoriesLoading } = useQuery({
  queryKey: ['store-product-categories'],
  queryFn: () => storeProductsService.getStoreProductCategories(),
  refetchInterval: 20000,
  initialData: []
})

// Set initial active category from the first category
watch(categories, (newCategories) => {
  if (newCategories.length > 0 && selectedCategoryId.value === null) {
    selectedCategoryId.value = newCategories[0].id
  }
}, { immediate: true })

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
  return (searchProducts.value?.data || []).slice().sort((a, b) => {
    if (a.isOutOfStock === b.isOutOfStock) return 0
    return a.isOutOfStock ? 1 : -1
  })
})

/* --- Vertical Auto-Update & Snap-to-Section --- */
// We use useScroll to track vertical scroll in the container.
const { y: scrollY } = useScroll(scrollContainer)

// This function calculates the nearest section (using offsetTop adjusted by toolbar height)
// and updates the active category.
function updateActiveCategory() {
  if (searchTerm.value || !scrollContainer.value) return
  const currentScroll = scrollContainer.value.scrollTop
  const toolbarHeight = toolbar.value?.$el?.offsetHeight || 0
  let closestId: number | null = null
  let minDistance = Infinity

  Object.entries(categoryRefs.value).forEach(([id, el]) => {
    // Calculate distance from section's top (minus toolbar height) to current scroll.
    const distance = Math.abs((el.offsetTop - toolbarHeight) - currentScroll)
    if (distance < minDistance) {
      minDistance = distance
      closestId = Number(id)
    }
  })

  if (closestId !== null && closestId !== selectedCategoryId.value) {
    selectedCategoryId.value = closestId
  }
}

// Snap to the nearest section when scrolling stops.
const snapToNearestSection = useDebounceFn(() => {
  if (!scrollContainer.value) return
  const toolbarHeight = toolbar.value?.$el?.offsetHeight || 0
  const currentScroll = scrollContainer.value.scrollTop
  let closestId: number | null = null
  let minDistance = Infinity

  Object.entries(categoryRefs.value).forEach(([id, el]) => {
    const distance = Math.abs((el.offsetTop - toolbarHeight) - currentScroll)
    if (distance < minDistance) {
      minDistance = distance
      closestId = Number(id)
    }
  })

  if (closestId !== null && scrollContainer.value && categoryRefs.value[closestId]) {
    scrollContainer.value.scrollTo({
      top: categoryRefs.value[closestId].offsetTop - toolbarHeight,
      behavior: 'smooth'
    })
    selectedCategoryId.value = closestId
  }
}, 300)

// Watch the vertical scroll to update active category and eventually snap to section.
watch(scrollY, () => {
  updateActiveCategory()
  snapToNearestSection()
})

/* --- Touch Event Handling for Horizontal Swipe Gestures --- */
// (These allow swiping left/right to jump between sections.)
let touchStartX = 0, touchStartY = 0

function handleTouchStart(e: TouchEvent) {
  if (e.touches.length > 0) {
    const touch = e.touches[0]
    touchStartX = touch.clientX
    touchStartY = touch.clientY
  }
}

function handleTouchEnd(e: TouchEvent) {
  if (e.changedTouches.length > 0) {
    const touch = e.changedTouches[0]
    const deltaX = touch.clientX - touchStartX
    const deltaY = touch.clientY - touchStartY
    const threshold = 50 // pixels
    // If horizontal swipe is dominant:
    if (Math.abs(deltaX) > threshold && Math.abs(deltaX) > Math.abs(deltaY)) {
      if (deltaX < 0) {
        goToNextSection()
      } else {
        goToPreviousSection()
      }
    }
  }
}

function goToNextSection() {
  const currentIndex = categories.value.findIndex(cat => cat.id === selectedCategoryId.value)
  if (currentIndex >= 0 && currentIndex < categories.value.length - 1) {
    const nextCat = categories.value[currentIndex + 1]
    onUpdateCategory(nextCat.id)
  }
}

function goToPreviousSection() {
  const currentIndex = categories.value.findIndex(cat => cat.id === selectedCategoryId.value)
  if (currentIndex > 0) {
    const prevCat = categories.value[currentIndex - 1]
    onUpdateCategory(prevCat.id)
  }
}

/* --- Register Category Header References --- */
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

/* --- Toolbar Interaction: Scroll to a Specific Category Section --- */
function onUpdateCategory(categoryId: number) {
  // Clear search mode and update active category.
  searchTerm.value = ''
  selectedCategoryId.value = categoryId
  const el = categoryRefs.value[categoryId]
  if (el) {
    const headerOffset = toolbar.value?.$el?.offsetHeight || 0
    scrollTo({
      top: el.offsetTop - headerOffset,
      behavior: 'smooth'
    })
  }
}

/* --- Search Term Interaction --- */
const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  searchTerm.value = newTerm.trim() !== '' ? newTerm : ''
}, 500)

function onUpdateSearchTerm(newSearchTerm: string) {
  debouncedEmitSearchTerm(newSearchTerm)
}

/* --- Setup Touch Listeners on the Scroll Container --- */
onMounted(() => {
  if (scrollContainer.value) {
    scrollContainer.value.addEventListener('touchstart', handleTouchStart)
    scrollContainer.value.addEventListener('touchend', handleTouchEnd)
  }
})
onBeforeUnmount(() => {
  if (scrollContainer.value) {
    scrollContainer.value.removeEventListener('touchstart', handleTouchStart)
    scrollContainer.value.removeEventListener('touchend', handleTouchEnd)
  }
})
</script>

<style scoped lang="scss">
.no-scrollbar {
  scrollbar-width: none; /* Firefox */
}
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
</style>
