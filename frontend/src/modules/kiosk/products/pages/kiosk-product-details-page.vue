<template>
	<div class="relative">
		<!-- Back Button -->
		<button
			@click="goBack"
			class="fixed left-4 top-20 rounded-full bg-slate-800/70 text-white backdrop-blur-md p-2 sm:p-4 z-10"
		>
			<Icon
				icon="mingcute:left-line"
				class="text-3xl"
			/>
		</button>

		<!-- Product Image -->
		<img
			:src="product.image"
			alt="Product Image"
			class="w-full h-[500px] sm:h-[600px] object-cover rounded-3xl"
		/>

		<!-- Gradient Overlay -->
		<div
			class="absolute inset-0 h-[500px] sm:h-[600px] bg-gradient-to-t from-[#F5F5F7] to-transparent pointer-events-none"
		></div>

		<!-- Product Details -->
		<div class="w-full p-4 sm:p-8 sm:-mt-24 relative">
			<h1 class="text-2xl sm:text-4xl font-medium">{{ product.title }}</h1>
			<p class="text-base sm:text-xl mt-1 sm:mt-3">{{ product.description }}</p>

			<!-- Size Selection and Add to Cart -->
			<div class="flex items-center gap-4 justify-between mt-5 sm:mt-10">
				<!-- Size Options -->
				<div class="flex items-center gap-2">
					<KioskDetailsSizes
						v-for="size in sizes"
						:key="size.id"
						:size="size"
						:selected-size="selectedSize"
						@click:size="onSizeClick"
					/>
				</div>

				<!-- Price and Add to Cart Button -->
				<div class="flex items-center gap-4 sm:gap-6">
					<p class="text-2xl sm:text-3xl font-medium">
						{{ formatPrice(totalPrice) }}
					</p>
					<button
						@click="handleAddToCart"
						class="rounded-full bg-primary text-primary-foreground p-2 sm:p-4"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-2xl sm:text-3xl"
						/>
					</button>
				</div>
			</div>

			<!-- Additives Selection -->
			<div class="mt-6 sm:mt-10">
				<p class="text-lg sm:text-2xl">Добавки</p>
				<div class="flex overflow-x-auto no-scrollbar gap-1 mt-2 sm:mt-6">
					<KioskDetailsAdditives
						v-for="additive in additives"
						:key="additive.id"
						:selected-additives="selectedAdditives"
						:additive="additive"
						@click:additive="onAdditiveClick"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
 import { ref, computed, type Ref } from 'vue';
 import { useRouter } from 'vue-router';
 import { Icon } from '@iconify/vue';
 import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue';
 import KioskDetailsAdditives from '@/modules/kiosk/products/components/details/kiosk-details-additives.vue';
 import { formatPrice } from '@/core/utils/price.utils';
 import type { Products, ProductSizes } from '@/modules/products/models/product.model';
 import type { Additives } from '@/modules/additives/models/additive.model';
import { useCartStore } from "@/modules/kiosk/cart/stores/cart.store";

 // Initialize Router
 const router = useRouter();

 // Access Cart Store
 const cartStore = useCartStore();

 // Product Data (Ideally fetched from an API or passed as a prop)
 const product: Ref<Products> = ref({
id: 1,
title: 'Какао',
description: 'Классический какао с молоком на выбор',
image:
  'https://media.istockphoto.com/id/1409579518/photo/perfect-iced-cocoa-aka-iced-cacao-stills-with-ice-cubes.jpg?s=612x612&w=0&k=20&c=ZUCCIGQaKSRWd_sPv27wT3msZR0yHwN7ot72jAKjbi8=',
category: 'Beverages',
startPrice: 1000
 });

 // Available Sizes (Include price in each size)
 const sizes: Ref<ProductSizes[]> = ref([
{ id: 1, label: 'S', volume: 250, price: 1000 },
{ id: 2, label: 'M', volume: 350, price: 1100 },
{ id: 3, label: 'L', volume: 450, price: 1200 },
 ]);

 // Selected Size (Default to first size)
 const selectedSize = ref<ProductSizes>(sizes.value[0]);

 // Available Additives
 const additives: Ref<Additives[]> = ref([
{
  id: 1,
  name: 'Сырная пенка',
  price: 400,
  imageUrl: 'https://example.com/images/additive1.png',
},
{
  id: 2,
  name: 'Шоколадная крошка',
  price: 300,
  imageUrl: 'https://example.com/images/additive2.png',
},
{
  id: 3,
  name: 'Взбитые сливки',
  price: 350,
  imageUrl: 'https://example.com/images/additive3.png',
},
// Add more additives as needed
 ]);

 // Selected Additives
 const selectedAdditives = ref<Additives[]>([]);

 // Quantity (Default to 1)
 const quantity = ref<number>(1);

 // Computed Total Price
 const totalPrice = computed<number>(() => {
const basePrice = selectedSize.value.price;
const additivesPrice = selectedAdditives.value.reduce(
  (sum, additive) => sum + additive.price,
  0
);

return (basePrice + additivesPrice) * quantity.value;
 });

 // Handlers
 const onSizeClick = (newSize: ProductSizes) => {
selectedSize.value = newSize;
 };

 const onAdditiveClick = (additive: Additives) => {
const index = selectedAdditives.value.findIndex(
  (item) => item.id === additive.id
);
if (index !== -1) {
  // Remove additive if already selected
  selectedAdditives.value.splice(index, 1);
} else {
  // Add additive if not selected
  selectedAdditives.value.push(additive);
}
 };

 const handleAddToCart = () => {
cartStore.addToCart(
  product.value,
  selectedSize.value,
  selectedAdditives.value,
  quantity.value
);
// Optional: Provide user feedback (e.g., notification)
 };

 const goBack = () => {
router.back();
 };
</script>

<style lang="scss">
/* Add your styles here */
</style>
