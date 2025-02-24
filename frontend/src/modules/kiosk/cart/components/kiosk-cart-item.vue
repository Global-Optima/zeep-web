<template>
	<div class="flex gap-4 sm:gap-8 bg-white p-4 sm:p-6 rounded-3xl">
		<!-- Product Image -->
		<img
			:src="item.product.imageUrl"
			alt="Product Image"
			class="rounded w-16 sm:w-28 h-16 sm:h-28 object-contain"
		/>

		<!-- Product Details -->
		<div class="flex-1">
			<div class="flex justify-between items-start gap-4">
				<p class="flex-1 text-xl sm:text-3xl">{{ item.product.name }}, {{ item.size.name }}</p>
				<button class="bg-gray-200 p-2 sm:p-3 rounded-xl">
					<Pencil
						class="size-5 text-gray-600"
						stroke-width="1.6"
					/>
				</button>
			</div>

			<div class="mt-1">
				<div class="flex flex-col gap-1 sm:gap-2 text-gray-600">
					<p class="text-xs sm:text-lg">{{ itemDescription }}</p>
				</div>
			</div>

			<div class="flex justify-between items-start mt-2">
				<p class="font-medium text-xl sm:text-3xl">
					{{ formatPrice(itemTotalPrice) }}
				</p>

				<!-- Quantity Controls -->
				<div class="flex items-center gap-2">
					<button
						@click="decrement"
						class="bg-gray-200 p-2 sm:p-3 rounded-xl"
					>
						<Trash
							v-if="item.quantity === 1"
							stroke-width="1.6"
							class="size-5 text-gray-600"
						/>

						<Minus
							v-if="item.quantity > 1"
							class="size-5 text-gray-600"
						/>
					</button>

					<span class="mx-1 sm:mx-2 text-base sm:text-2xl">
						{{ item.quantity }}
					</span>

					<button
						@click="increment"
						class="bg-gray-200 p-2 sm:p-3 rounded-xl"
					>
						<Plus class="size-5 text-gray-600" />
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore, type CartItem } from "@/modules/kiosk/cart/stores/cart.store"
import { Minus, Pencil, Plus, Trash } from 'lucide-vue-next'
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
  const storePrice = props.item.size.storePrice;

  const additivesPrice = props.item.additives.reduce(
    (sum, additive) => sum + additive.storePrice,
    0
  );
  return (storePrice + additivesPrice) * props.item.quantity;
});

// Computed property for additives list as a comma-separated string
const itemDescription = computed(() => {
  const itemSize = `${props.item.size.size} ${props.item.size.unit.name}`
  if (props.item.additives.length > 0) {
    return `${itemSize}, ${props.item.additives.map(additive => additive.name).join(", ")}`.toLowerCase()
  }

  return itemSize.toLowerCase()
});
</script>

<style scoped></style>
