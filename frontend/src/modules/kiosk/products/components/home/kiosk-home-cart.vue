<template>
	<button
		@click="onCartClick"
		class="flex items-center gap-4 bg-primary shadow-xl backdrop-blur-md px-6 sm:px-8 py-6 rounded-full text-primary-foreground"
	>
		<div>
			<ShoppingBasket
				class="size-10"
				stroke-width="1.5"
			/>
		</div>

		<p class="font-medium text-xl sm:text-4xl">
			{{ formatPrice(cartStore.totalPrice) }}
		</p>
	</button>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { ShoppingBasket } from 'lucide-vue-next'

const cartStore = useCartStore();

const onCartClick = () => {
  cartStore.toggleModal()
};
</script>

<style scoped></style>
