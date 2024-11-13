<template>
	<div class="flex h-screen sm:flex-row flex-col pt-safe">
		<!-- Sidebar for tablet and larger screens -->
		<aside class="hidden sm:flex py-8 px-6 bg-white h-full overflow-hidden flex-col max-w-[250px]">
			<p class="text-3xl font-semibold">ZEEP</p>
			<ul class="flex flex-col justify-center gap-4 mt-8">
				<li
					v-for="category in categories"
					:key="category.id"
					@click="onUpdateCategory(category.id)"
					:class="[
						'cursor-pointer text-xl',
						category.id === selectedCategoryId ? 'text-primary' : 'text-gray-600',
					]"
				>
					{{ category.name }}
				</li>
			</ul>

			<button
				@click="onCartClick"
				class="rounded-2xl px-6 py-4 bg-gray-200 text-black mt-auto"
			>
				<span class="text-xl">Корзина ({{ cartTotalItems }})</span>
			</button>
		</aside>

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
						@select-product="openProductSheet"
					/>
				</div>
			</section>
		</div>

		<!-- Cart Button for mobile -->
		<div
			v-if="!isCartEmpty"
			class="fixed bottom-10 left-0 w-full flex justify-center sm:hidden"
		>
			<button
				@click="onCartClick"
				class="rounded-3xl px-6 py-4 sm:px-8 sm:py-4 bg-slate-800/70 text-white backdrop-blur-md"
			>
				<div class="flex items-center gap-6">
					<p class="text-lg sm:text-2xl">Корзина ({{ cartTotalItems }})</p>
					<p class="text-lg sm:text-3xl font-medium">
						{{ formatPrice(cartTotalPrice) }}
					</p>
				</div>
			</button>
		</div>

		<!-- Product Details Dialog -->
		<Sheet
			:open="isSheetOpen"
			@update:open="closeProductSheet"
		>
			<SheetContent
				side="bottom"
				class="p-0 overflow-clip rounded-t-3xl overflow-y-auto h-[92vh] no-scrollbar bg-[#F5F5F7] border-t-0"
			>
				<KioskDetailsSheetContent :selected-product-id="selectedProductId" />
			</SheetContent>
		</Sheet>
	</div>
</template>

<script setup lang="ts">
import { Sheet, SheetContent } from '@/core/components/ui/sheet'
import { Skeleton } from '@/core/components/ui/skeleton'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import KioskDetailsSheetContent from '@/modules/kiosk/products/components/details/kiosk-details-sheet-content.vue'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import KioskHomeToolbarMobile from '@/modules/kiosk/products/components/home/kiosk-home-toolbar-mobile.vue'
import KioskHomeToolbarTablet from '@/modules/kiosk/products/components/home/kiosk-home-toolbar-tablet.vue'
import type { ProductCategory, StoreProducts } from '@/modules/kiosk/products/models/product.model'
import { productService } from '@/modules/kiosk/products/services/products.service'
import { useProductStore } from "@/modules/kiosk/products/stores/current-product.store"
import { useQuery } from '@tanstack/vue-query'
import { useDebounceFn } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

// Initialize router and cart store
const router = useRouter();
const cartStore = useCartStore();

const selectedCategoryId = ref<number | null>(null);
const searchTerm = ref('');

// State for products and selected product
const selectedProductId = ref<number | null>(null);

// New state variable to control the Sheet's open/close state
const isSheetOpen = ref(false);

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

// Remove the old computed isSheetOpen based on selectedProductId
// const isSheetOpen = computed(() => selectedProductId.value !== null);

function onUpdateCategory(categoryId: number) {
  selectedCategoryId.value = categoryId;
}

const debouncedEmitSearchTerm = useDebounceFn((newTerm: string) => {
  searchTerm.value = newTerm;
}, 500);

function onUpdateSearchTerm(newSearchTerm: string) {
  debouncedEmitSearchTerm(newSearchTerm);
}

const onCartClick = () => {
  router.push({ name: getRouteName('KIOSK_CART') });
};

const productStore = useProductStore()

const openProductSheet = (productId: number) => {
//   selectedProductId.value = productId;
//   isSheetOpen.value = true; // Open the Sheet
	productStore.selectProduct(productId);
};

const closeProductSheet = () => {
  isSheetOpen.value = false; // Close the Sheet
  // Optionally, you can reset selectedProductId here if needed
  // selectedProductId.value = null;
};
</script>

<style scoped lang="scss"></style>
