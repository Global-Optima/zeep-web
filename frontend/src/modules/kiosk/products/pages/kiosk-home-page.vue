<template>
	<!-- Toolbar with Categories and Search -->
	<KioskHomeToolbar
		v-if="!categoriesLoading"
		:categories="categories"
		:selected-category-id="selectedCategoryId"
		:search-term="searchTerm"
		@update:category="onUpdateCategory"
		@update:search-term="onUpdateSearchTerm"
	/>

	<div
		v-else
		class="w-full py-4 sm:py-6 px-4 flex items-center gap-2 overflow-x-auto no-scrollbar sticky top-0 z-10"
	>
		<Skeleton
			v-for="n in 4"
			:key="n"
			class="h-16 w-32 rounded-full bg-gray-200"
		/>
	</div>

	<!-- Products Grid -->
	<section class="w-full px-4 sm:px-6 overflow-y-auto">
		<div
			v-if="productsLoading"
			class="grid grid-cols-2 sm:grid-cols-3 gap-2 sm:gap-4"
		>
			<Skeleton
				v-for="n in 6"
				:key="n"
				class="rounded-3xl w-full h-60 bg-gray-200"
			/>
		</div>

		<div
			v-else-if="products.length === 0"
			class="flex items-center justify-center h-20"
		>
			<p class="text-lg text-gray-400">Ничего не найдено</p>
		</div>

		<div
			v-else
			class="grid grid-cols-2 sm:grid-cols-3 gap-2 sm:gap-4"
		>
			<KioskHomeProductCard
				v-for="product in products"
				:key="product.id"
				:product="product"
				@select-product="openProductSheet"
			/>
		</div>
	</section>

	<div
		v-if="!isCartEmpty"
		class="fixed bottom-10 left-0 w-full flex justify-center"
	>
		<button
			@click="onCartClick"
			class="rounded-3xl px-6 py-4 sm:px-8 sm:py-6 bg-slate-800/70 text-white backdrop-blur-md"
		>
			<div class="flex items-center gap-6">
				<p class="text-lg sm:text-2xl">Корзина ({{ cartTotalItems }})</p>
				<p class="text-lg sm:text-3xl font-medium">
					{{ formatPrice(cartTotalPrice) }}
				</p>
			</div>
		</button>
	</div>

	<Dialog
		:open="isSheetOpen"
		@update:open="closeProductSheet"
	>
		<DialogContent
			class="bg-[#F5F5F7] p-0 max-w-[95vw] rounded-3xl sm:rounded-[36px] overflow-clip"
		>
			<div class="overflow-auto max-h-[82vh] sm:max-h-[82vh] no-scrollbar">
				<KioskDetailsSheetContent
					v-if="selectedProductId"
					:selected-product-id="selectedProductId"
				/>
			</div>
		</DialogContent>
	</Dialog>
</template>

<script setup lang="ts">

import { useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import {
  Dialog,
  DialogContent
} from '@/core/components/ui/dialog'

import { Skeleton } from '@/core/components/ui/skeleton'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsSheetContent from '@/modules/kiosk/products/components/details/kiosk-details-sheet-content.vue'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeToolbar from '@/modules/kiosk/products/components/home/kiosk-home-toolbar.vue'
import type { ProductCategory, StoreProducts } from '@/modules/kiosk/products/models/product.model'
import { productService } from '@/modules/kiosk/products/services/products.service'

// Initialize router and cart store
const router = useRouter();
const cartStore = useCartStore();

const selectedCategoryId = ref<number | null>(null);
const searchTerm = ref('');

// State for products and selected product
const selectedProductId = ref<number | null>(null);

// Reactive queryKey for products
const productsQueryKey = computed(() => [
  'products',
  { categoryId: selectedCategoryId.value, searchTerm: searchTerm.value },
]);

// Fetch products based on selected category and search term
const { data: products, isLoading: productsLoading } = useQuery<StoreProducts[]>({
  queryKey: productsQueryKey,
  queryFn: () =>
    productService.getStoreProducts(
      selectedCategoryId.value!,
      searchTerm.value
    ),
  enabled: computed(() => Boolean(selectedCategoryId.value)),
  initialData: []
});

const { data: categories, isLoading: categoriesLoading } = useQuery<ProductCategory[]>({
  queryKey: ['categories'],
  queryFn: () => productService.getStoreCategories(),
  initialData: []
});

watch(
  categories,
  (newCategories) => {
    if (newCategories && newCategories.length > 0 && !selectedCategoryId.value) {
      selectedCategoryId.value = newCategories[0].id;
    }
  },
  { immediate: true }
);


const cartTotalItems = computed(() => cartStore.totalItems);
const cartTotalPrice = computed(() => cartStore.totalPrice);
const isCartEmpty = computed(() => cartStore.isEmpty);

const isSheetOpen = computed(() => selectedProductId.value !== null);

function onUpdateCategory(categoryId: number) {
  selectedCategoryId.value = categoryId;
}

function onUpdateSearchTerm(newSearchTerm: string) {
  searchTerm.value = newSearchTerm;
}

const onCartClick = () => {
  router.push({ name: getRouteName('KIOSK_CART') });
};

const openProductSheet = (productId: number) => {
  selectedProductId.value = productId;
};

const closeProductSheet = () => {
  selectedProductId.value = null;
};
</script>

<style scoped lang="scss"></style>
