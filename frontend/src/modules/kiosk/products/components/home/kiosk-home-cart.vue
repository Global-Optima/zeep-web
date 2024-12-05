<template>
	<button
		@click="onCartClick"
		class="flex items-center gap-3 bg-primary shadow-xl backdrop-blur-md px-6 sm:px-8 py-6 rounded-full text-primary-foreground"
	>
		<p class="text-xl sm:text-4xl">
			{{ formatPrice(cartStore.totalPrice) }}
		</p>

		<ChevronRight class="w-8 h-8" />
	</button>
</template>

<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { ChevronRight } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const cartStore = useCartStore();
const router  = useRouter()

const onCartClick = () => {
  router.push({ name: getRouteName('KIOSK_CART') });
};
</script>

<style scoped></style>
