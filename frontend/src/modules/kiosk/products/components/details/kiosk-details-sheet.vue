<template>
	<div class="w-full p-4 sm:p-8">
		<div class="relative">
			<!-- Close Button -->
			<button
				@click="$emit('close')"
				class="absolute top-4 right-4 rounded-full bg-gray-200 p-2 sm:p-3"
			>
				<Icon
					icon="mingcute:close-line"
					class="text-xl sm:text-2xl"
				/>
			</button>

			<!-- Product Image -->
			<img
				:src="product.image"
				alt="Product Image"
				class="w-full h-64 sm:h-80 object-cover rounded-2xl"
			/>

			<!-- Product Title and Description -->
			<div class="mt-4">
				<h1 class="text-2xl sm:text-4xl font-medium">{{ product.title }}</h1>
				<p class="text-base sm:text-xl mt-2 sm:mt-4">{{ product.description }}</p>
			</div>

			<!-- Size Selection -->
			<div class="mt-6 sm:mt-10">
				<p class="text-lg sm:text-2xl font-medium">Выберите размер</p>
				<div class="flex items-center gap-4 mt-2 sm:mt-4">
					<KioskDetailsSizes
						v-for="size in sizes"
						:key="size.id"
						:size="size"
						:selected-size="selectedSize"
						@click:size="onSizeClick"
					/>
				</div>
			</div>

			<!-- Additives Selection -->
			<div class="mt-6 sm:mt-10">
				<p class="text-lg sm:text-2xl font-medium">Добавки</p>
				<div class="flex flex-wrap gap-2 mt-2 sm:mt-4">
					<KioskDetailsAdditives
						v-for="additive in additives"
						:key="additive.id"
						:additive="additive"
						:selected-additives="selectedAdditives"
						@click:additive="onAdditiveClick"
					/>
				</div>
			</div>

			<!-- Total Price and Add to Cart -->
			<div class="flex items-center justify-between mt-6 sm:mt-10">
				<p class="text-2xl sm:text-3xl font-medium">{{ formatPrice(totalPrice) }}</p>
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
	</div>
</template>

<script setup lang="ts">
 import { computed, ref, watch } from 'vue';
 import { Icon } from '@iconify/vue';
 import KioskDetailsSizes from '@/modules/kiosk/products/components/details/kiosk-details-sizes.vue';
 import KioskDetailsAdditives from '@/modules/kiosk/products/components/details/kiosk-details-additives.vue';
 import { formatPrice } from '@/core/utils/price.utils';
 import type { Products, ProductSizes } from '@/modules/products/models/product.model';
 import type { Additives } from '@/modules/additives/models/additive.model';
import { useCartStore } from "@/modules/kiosk/cart/stores/cart.store";

 // Define props
 const props = defineProps<{
product: Products;
 }>();

 const cartStore = useCartStore();

 // Available Sizes (if dynamic, otherwise pass via props or fetch)
 const sizes = ref<ProductSizes[]>([
{ id: 1, label: 'S', volume: 250, price: 1000 },
{ id: 2, label: 'M', volume: 350, price: 1100 },
{ id: 3, label: 'L', volume: 450, price: 1200 },
 ]);

 // Available Additives (if dynamic, otherwise pass via props or fetch)
 const additives = ref<Additives[]>([
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

 // Selected Size (default to first size)
 const selectedSize = ref<ProductSizes>(sizes.value[0]);

 // Selected Additives
 const selectedAdditives = ref<Additives[]>([]);

 // Quantity (default to 1)
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
  selectedAdditives.value.splice(index, 1);
} else {
  selectedAdditives.value.push(additive);
}
 };

 const handleAddToCart = () => {
cartStore.addToCart(
  props.product,
  selectedSize.value,
  selectedAdditives.value,
  quantity.value
);
 };

 // Watch for product prop changes to reset selections
 watch(
() => props.product,
() => {
  selectedSize.value = sizes.value[0];
  selectedAdditives.value = [];
  quantity.value = 1;
},
{ immediate: true }
 );
</script>

<style scoped>
/* Add your styles here */
</style>
