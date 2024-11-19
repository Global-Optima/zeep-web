<template>
	<div class="flex gap-4 sm:gap-8 bg-white p-4 sm:p-6 rounded-3xl">
		<!-- Product Image -->
		<img
			:src="item.product.imageUrl"
			alt="Product Image"
			class="rounded w-16 sm:w-20 h-16 sm:h-20 object-contain"
		/>

		<!-- Product Details -->
		<div class="flex-1">
			<p class="text-base sm:text-2xl">{{ item.product.name }}, {{ item.size.name }}</p>

			<div class="mt-2">
				<div
					v-if="item.additives.length > 0"
					class="flex flex-col gap-1 sm:gap-2 text-gray-600"
				>
					<p class="text-xs sm:text-lg">{{ additivesList }}</p>
				</div>
			</div>

			<div class="flex justify-between items-center mt-2">
				<p class="font-medium text-lg sm:text-2xl">
					{{ formatPrice(itemTotalPrice) }}
				</p>

				<!-- Quantity Controls -->
				<div class="flex items-center gap-2">
					<button
						@click="decrement"
						class="bg-gray-200 p-1 sm:p-3 rounded-xl"
					>
						<Icon
							icon="mingcute:minimize-line"
							class="text-lg sm:text-2xl"
						/>
					</button>

					<span class="mx-1 sm:mx-2 text-base sm:text-2xl">
						{{ item.quantity }}
					</span>

					<button
						@click="increment"
						class="bg-gray-200 p-1 sm:p-3 rounded-xl"
					>
						<Icon
							icon="mingcute:add-line"
							class="text-lg sm:text-xl"
						/>
					</button>
				</div>
			</div>

			<!-- Additives List -->
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore, type CartItem } from "@/modules/kiosk/cart/stores/cart.store"
import { Icon } from '@iconify/vue'
import { computed } from 'vue'

const props = defineProps<{
  item: CartItem;
}>();

const cartStore = useCartStore();

// Reactive methods
const increment = (e: Event) => {
  e.stopPropagation()
  cartStore.incrementQuantity(props.item.key);
};

const decrement = (e: Event) => {
  e.stopPropagation()
  cartStore.decrementQuantity(props.item.key);
};

// Computed property for item total price
const itemTotalPrice = computed(() => {
  const basePrice = props.item.size.basePrice;
  const additivesPrice = props.item.additives.reduce(
    (sum, additive) => sum + additive.price,
    0
  );
  return (basePrice + additivesPrice) * props.item.quantity;
});

// Computed property for additives list as a comma-separated string
const additivesList = computed(() => {
  return props.item.additives.map(additive => additive.name).join(', ');
});
</script>

<style scoped></style>
