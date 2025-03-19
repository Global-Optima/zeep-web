<template>
	<button
		@click="onCartClick"
		class="flex items-center gap-4 bg-primary shadow-xl backdrop-blur-md px-6 sm:px-8 py-6 rounded-full text-primary-foreground transition-all duration-300 ease-out pointer-events-auto"
		:class="{ 'scale-110 shadow-2xl': isAnimating }"
	>
		<div>
			<ShoppingBasket
				class="size-10 transition-all duration-300 ease-out"
				stroke-width="1.5"
			/>
		</div>

		<p class="font-medium text-xl sm:text-4xl transition-all duration-300 ease-out">
			{{ formatPrice(cartStore.totalPrice) }}
		</p>
	</button>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { ShoppingBasket } from 'lucide-vue-next'
import { ref, watch } from 'vue'

const cartStore = useCartStore();
const isAnimating = ref(false);

const onCartClick = () => {
  cartStore.toggleModal()
};

// Watch for changes in cart items to trigger animation
watch(() => cartStore.totalItems, (newVal, oldVal) => {
  if (newVal > oldVal) {
    isAnimating.value = true;
    setTimeout(() => {
      isAnimating.value = false;
    }, 500);
  }
});
</script>
