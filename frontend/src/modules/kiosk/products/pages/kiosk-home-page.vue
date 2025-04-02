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
		<section class="mt-8 px-8 pb-4">
			<!-- Search Mode -->
			<div v-if="searchTerm">
				<KioskHomeSearchProducts
					:products="sortedSearchProducts"
					:isLoading="isSearchProductsPending"
				/>
			</div>

			<!-- Category Mode -->
			<div v-else>
				<KioskHomeCategoryProducts
					v-for="(category, index) in categories"
					:key="category.id"
					:category="category"
					:products="categoryProducts[category.id]"
					:isLoading="categoryProductsQueries[index].isPending"
					:setCategoryRef="setCategoryRef"
				/>
			</div>
		</section>
	</div>
</template>

<script setup lang="ts">
import type { StoreProductDTO } from '@/modules/admin/store-products/models/store-products.model'
import { storeProductsService } from '@/modules/admin/store-products/services/store-products.service'
import KioskHomeCategoryProducts from '@/modules/kiosk/products/components/home/list/kiosk-home-category-products.vue'
import KioskHomeSearchProducts from '@/modules/kiosk/products/components/home/list/kiosk-home-search-products.vue'
import KioskHomeToolbarMobile from '@/modules/kiosk/products/components/home/toolbar/kiosk-home-toolbar-mobile.vue'
import { useQueries, useQuery } from '@tanstack/vue-query'
import { useDebounceFn, useScroll } from '@vueuse/core'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, shallowRef, watch, watchEffect } from 'vue'

/* --- Reactive State --- */
const scrollContainer = shallowRef<HTMLElement | null>(null)
const toolbar = shallowRef<InstanceType<typeof KioskHomeToolbarMobile> | null>(null)

const searchTerm = ref('')
const selectedCategoryId = ref<number | null>(null)
const categoryRefs = ref<Record<number, HTMLElement>>({})
const isManualScroll = ref(false)
const scrollEndTimeout = ref<number | null>(null)

/* --- Fetch Categories --- */
const { data: categories, isPending: categoriesLoading } = useQuery({
  queryKey: ['store-product-categories'],
  queryFn: () => storeProductsService.getStoreProductCategories(),
  refetchInterval: 20000,
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
  if (!searchProducts.value?.data) return []
  return searchProducts.value.data.slice().sort((a, b) => {
    if (a.isOutOfStock === b.isOutOfStock) return 0
    return a.isOutOfStock ? 1 : -1
  })
})

/* --- Scroll Handling --- */
const { y: scrollY } = useScroll(() => scrollContainer.value)

const updateActiveCategory = () => {
  if (searchTerm.value || isManualScroll.value) return

  const toolbarHeight = toolbar.value?.$el?.offsetHeight || 0
  let closestId: number | null = null
  let minDistance = Infinity

  Object.entries(categoryRefs.value).forEach(([id, el]) => {
    const rect = el.getBoundingClientRect()
    const distance = Math.abs(rect.top - toolbarHeight)
    if (distance < minDistance) {
      minDistance = distance
      closestId = Number(id)
    }
  })

  if (closestId !== null && closestId !== selectedCategoryId.value) {
    selectedCategoryId.value = closestId
  }
}

watch(scrollY, () => {
  if (!isManualScroll.value) updateActiveCategory()
})

/* --- Category Reference Management --- */
const setCategoryRef = (id: number, el: Element | { $el: HTMLElement } | null) => {
  let element: HTMLElement | null = null
  if (el) {
    element = el instanceof HTMLElement ? el : ('$el' in el && el.$el instanceof HTMLElement ? el.$el : null)
  }
  if (element) {
    categoryRefs.value[id] = element
  } else {
    delete categoryRefs.value[id]
  }
}

/* --- Toolbar Interaction --- */
const onUpdateCategory = (categoryId: number) => {
  searchTerm.value = ''
  isManualScroll.value = true
  selectedCategoryId.value = categoryId

  if (scrollEndTimeout.value !== null) {
    window.clearTimeout(scrollEndTimeout.value)
  }

  const el = categoryRefs.value[categoryId]
  if (el && scrollContainer.value) {
    const toolbarHeight = toolbar.value?.$el?.offsetHeight || 0
    scrollContainer.value.scrollTo({
      top: el.offsetTop - toolbarHeight - 20,
      behavior: 'smooth'
    })

    scrollEndTimeout.value = window.setTimeout(() => {
      isManualScroll.value = false
      scrollEndTimeout.value = null
    }, 1000)
  } else {
    isManualScroll.value = false
  }
}

/* --- Search Term Interaction --- */
const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  searchTerm.value = newTerm.trim() !== '' ? newTerm : ''
}, 800)

const onUpdateSearchTerm = (newSearchTerm: string) => {
  debouncedEmitSearchTerm(newSearchTerm)
}

/* --- Touch Event Handling for Horizontal Swipe Gestures --- */
const touchStartX = ref(0)
const touchStartY = ref(0)

const handleTouchStart = (e: TouchEvent) => {
  if ((e.target as HTMLElement).closest('[data-testid="category-buttons"]')) {
    return;
  }
  if (e.touches.length > 0) {
    const touch = e.touches[0];
    touchStartX.value = touch.clientX;
    touchStartY.value = touch.clientY;
  }
}

const handleTouchEnd = (e: TouchEvent) => {
  // Ignore touches that end inside the category buttons container.
  if ((e.target as HTMLElement).closest('[data-testid="category-buttons"]')) {
    return;
  }
  if (e.changedTouches.length > 0) {
    const touch = e.changedTouches[0];
    const deltaX = touch.clientX - touchStartX.value;
    const deltaY = touch.clientY - touchStartY.value;
    const threshold = 50;
    if (Math.abs(deltaX) > threshold && Math.abs(deltaX) > Math.abs(deltaY)) {
      if (deltaX < 0) {
        goToNextSection()
      } else {
        goToPreviousSection()
      }    }
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
  if (scrollContainer.value) {
    scrollContainer.value.addEventListener('touchstart', handleTouchStart)
    scrollContainer.value.addEventListener('touchend', handleTouchEnd)
  }
  nextTick(updateActiveCategory)
})

onBeforeUnmount(() => {
  if (scrollContainer.value) {
    scrollContainer.value.removeEventListener('touchstart', handleTouchStart)
    scrollContainer.value.removeEventListener('touchend', handleTouchEnd)
  }
  if (scrollEndTimeout.value !== null) {
    window.clearTimeout(scrollEndTimeout.value)
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
