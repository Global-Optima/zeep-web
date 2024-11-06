<template>
	<KioskHomeToolbar
		:categories="categories"
		:selected-category="selectedCategory"
		:search-term="searchTerm"
		@update:category="onUpdateCategory"
		@update:search-term="onUpdateSearchTerm"
	/>

	<section class="w-full px-4 sm:px-6 overflow-y-auto">
		<div class="grid grid-cols-2 sm:grid-cols-3 gap-2 sm:gap-4">
			<KioskHomeProductCard
				v-for="product in filteredProducts"
				:key="product.id"
				:product="product"
			/>
		</div>
	</section>

	<div class="fixed bottom-10 left-0 w-full flex justify-center">
		<button
			@click="onCartClick"
			class="rounded-full px-6 py-4 bg-slate-800/70 text-white backdrop-blur-md"
		>
			<div class="flex items-center gap-6">
				<p class="text-lg sm:text-xl">Корзина</p>
				<p class="text-lg sm:text-2xl font-medium">{{ formatPrice(5950) }}</p>
			</div>
		</button>
	</div>

	<Sheet>
		<SheetTrigger>
			<button class="px-4 mt-5">Test Popup</button>
		</SheetTrigger>
		<SheetContent
			side="bottom"
			class="max-h-[92vh] overflow-auto p-0 pb-14 bg-[#F5F5F7] border-none no-scrollbar rounded-t-3xl"
		>
			<KioskDetailsSheet />
		</SheetContent>
	</Sheet>
</template>

<script setup lang="ts">
import {
  Sheet,
  SheetContent,
  SheetTrigger
} from '@/core/components/ui/sheet'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import KioskDetailsSheet from '@/modules/kiosk/components/details/kiosk-details-sheet.vue'
import KioskHomeProductCard from '@/modules/kiosk/components/home/kiosk-home-product-card.vue'
import { products } from '@/modules/kiosk/components/home/kiosk-home-products-list'
import KioskHomeToolbar from '@/modules/kiosk/components/home/kiosk-home-toolbar.vue'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const searchTerm = ref('');
const categories = [
	'Популярное',
	'Новинки',
	'Сезонное',
	'Кофе',
	'Чай',
	'Холодные напитки',
	'Мороженое'
];
const selectedCategory = ref("Популярное");

const onCartClick = () => {
  router.push({name: getRouteName("KIOSK_CART")})
}

const onUpdateCategory = (category: string) => {
  selectedCategory.value = category
}

const onUpdateSearchTerm = (newSearchTerm: string) => {
  searchTerm.value = newSearchTerm
}

const filteredProducts = computed(() => {
	if (searchTerm.value) {
		return products.filter(product =>
			product.title.toLowerCase().includes(searchTerm.value.toLowerCase())
		);
	}

		return products.filter(product => product.category === selectedCategory.value);
});
</script>

<style lang="scss"></style>
