<template>
	<!-- Toolbar with Categories and Search -->
	<KioskHomeToolbar
		:categories="categories"
		:selected-category="selectedCategory"
		:search-term="searchTerm"
		@update:category="onUpdateCategory"
		@update:search-term="onUpdateSearchTerm"
	/>

	<!-- Products Grid -->
	<section class="w-full px-4 sm:px-6 overflow-y-auto">
		<div class="grid grid-cols-2 sm:grid-cols-3 gap-2 sm:gap-4">
			<KioskHomeProductCard
				v-for="product in filteredProducts"
				:key="product.id"
				:product="product"
				@select-product="openProductSheet"
			/>
		</div>
	</section>

	<!-- Cart Button -->
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

	<!-- Product Details Sheet -->
	<Sheet
		:open="isSheetOpen"
		@update:open="closeProductSheet"
	>
		<SheetContent
			side="bottom"
			class="max-h-[92vh] overflow-auto p-0 pb-14 bg-[#F5F5F7] border-none no-scrollbar rounded-t-3xl"
		>
			<KioskDetailsSheet
				v-if="selectedProduct"
				:product="selectedProduct"
			/>
		</SheetContent>
	</Sheet>
</template>

<script setup lang="ts">
 import {
  Sheet,
  SheetContent,
} from '@/core/components/ui/sheet'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from "@/modules/kiosk/cart/stores/cart.store"
import KioskDetailsSheet from '@/modules/kiosk/products/components/details/kiosk-details-sheet-content.vue'
import KioskHomeProductCard from '@/modules/kiosk/products/components/home/kiosk-home-product-card.vue'
import { products } from '@/modules/kiosk/products/components/home/kiosk-home-products-list'
import KioskHomeToolbar from '@/modules/kiosk/products/components/home/kiosk-home-toolbar.vue'
import type { Products } from "@/modules/products/models/product.model"
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

 const router = useRouter();
 const cartStore = useCartStore();


 const searchTerm = ref('');
 const categories = [
'Популярное',
'Новинки',
'Сезонное',
'Кофе',
'Чай',
'Холодные напитки',
'Мороженое',
 ];
 const selectedCategory = ref('Популярное');

 const cartTotalItems = computed(() => cartStore.totalItems);
 const cartTotalPrice = computed(() => cartStore.totalPrice);
 const isCartEmpty = computed(() => cartStore.isEmpty);

 const selectedProduct = ref<Products | null>(null);

 const isSheetOpen = computed(() => selectedProduct.value !== null);

const filteredProducts = computed(() => {
  if (searchTerm.value) {
    return products.filter((product) =>
    product.title.toLowerCase().includes(searchTerm.value.toLowerCase())
    );
  }
  return products.filter(
    (product) => product.category === selectedCategory.value
  );
});

const onCartClick = () => {
  router.push({ name: getRouteName('KIOSK_CART') });
};

const onUpdateCategory = (category: string) => {
  selectedCategory.value = category;
};

const onUpdateSearchTerm = (newSearchTerm: string) => {
  searchTerm.value = newSearchTerm;
};

const openProductSheet = (product: Products) => {
  selectedProduct.value = product;
};

const closeProductSheet = () => {
  selectedProduct.value = null;
};
</script>

<style lang="scss"></style>
