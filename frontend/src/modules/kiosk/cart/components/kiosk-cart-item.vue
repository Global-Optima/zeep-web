<template>
	<div class="flex gap-4 sm:gap-8 bg-white shadow-lg shadow-slate-100 p-4 sm:p-6 rounded-[36px]">
		<!-- Product Image -->
		<LazyImage
			:src="item.product.imageUrl"
			alt="Изображение продукта"
			class="rounded-md size-24 object-contain"
		/>

		<!-- Product Details -->
		<div class="flex-1">
			<p class="flex-1 text-3xl">{{ item.product.name }}, {{ item.size.name }}</p>
			<p class="mt-2 text-xl">{{ itemDescription }}</p>

			<p class="mt-4 font-semibold text-primary text-4xl">
				{{ formatPrice(itemTotalPrice) }}
			</p>

			<div class="flex justify-between items-center gap-4 mt-6">
				<div class="flex items-center gap-2">
					<button
						@click="decrement"
						class="bg-gray-100 p-4 rounded-2xl"
					>
						<Trash
							v-if="item.quantity === 1"
							stroke-width="1.6"
							class="size-6 text-gray-600"
						/>

						<Minus
							v-if="item.quantity > 1"
							class="size-6 text-gray-600"
						/>
					</button>

					<span class="mx-4 font-medium text-3xl">
						{{ item.quantity }}
					</span>

					<button
						@click="increment"
						class="bg-gray-100 p-4 rounded-2xl"
					>
						<Plus class="size-6 text-gray-600" />
					</button>
				</div>

				<button class="bg-gray-100 p-4 rounded-2xl">
					<Pencil
						class="size-6 text-gray-600"
						stroke-width="1.6"
					/>
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import LazyImage from '@/core/components/lazy-image/LazyImage.vue'
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
    (sum, additive) => {
      if (additive.isDefault) {
				return sum
			}

      return sum + additive.storePrice
    },
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
