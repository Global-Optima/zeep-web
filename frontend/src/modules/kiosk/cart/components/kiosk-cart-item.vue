<template>
	<div class="flex gap-4 sm:gap-8 p-4 sm:p-6 rounded-3xl bg-white">
		<!-- Product Image -->
		<img
			:src="item.product.image"
			alt="Product Image"
			class="w-16 h-16 sm:w-fit sm:h-20 object-contain rounded"
		/>

		<!-- Product Details -->
		<div class="flex-1">
			<p class="text-base sm:text-2xl">
				{{ item.product.title }}, {{ item.size.label }} ({{ item.size.volume }}ml)
			</p>

			<div class="flex justify-between items-center mt-2">
				<p class="text-lg sm:text-2xl font-medium">
					{{ formatPrice(itemTotalPrice) }}
				</p>

				<!-- Quantity Controls -->
				<div class="flex items-center gap-2">
					<button
						@click="decrement"
						class="p-1 sm:p-2 bg-gray-200 rounded-xl"
					>
						<Icon
							icon="mingcute:minimize-line"
							class="text-lg sm:text-xl"
						/>
					</button>

					<span class="text-base sm:text-2xl mx-1 sm:mx-2">
						{{ item.quantity }}
					</span>

					<button
						@click="increment"
						class="p-1 sm:p-2 bg-gray-200 rounded-xl"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-lg sm:text-xl"
						/>
					</button>
				</div>
			</div>

			<!-- Additives List -->
			<div class="mt-3">
				<div
					v-if="item.additives.length > 0"
					class="flex flex-col gap-1 sm:gap-2 text-gray-600"
				>
					<div
						v-for="additive in item.additives"
						:key="additive.id"
						class="inline-flex items-center gap-1 sm:gap-2"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-lg sm:text-xl text-primary"
						/>
						<p class="text-xs sm:text-lg">{{ additive.name }}</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
 import { computed } from 'vue';
 import { formatPrice } from '@/core/utils/price.utils';
 import { Icon } from '@iconify/vue';
import { useCartStore, type CartItem } from "@/modules/kiosk/cart/stores/cart.store";

 const props = defineProps<{
item: CartItem;
 }>();

 const cartStore = useCartStore();

 // Reactive methods
 const increment = () => {
cartStore.incrementQuantity(props.item.key);
 };

 const decrement = () => {
cartStore.decrementQuantity(props.item.key);
 };

 // Computed property for item total price
 const itemTotalPrice = computed(() => {
const basePrice = props.item.size.price;
const additivesPrice = props.item.additives.reduce(
  (sum, additive) => sum + additive.price,
  0
);
return (basePrice + additivesPrice) * props.item.quantity;
 });
</script>

<style scoped></style>
