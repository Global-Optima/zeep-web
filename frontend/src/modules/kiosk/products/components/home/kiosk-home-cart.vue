<template>
	<button
		@click="onCartClick"
		class="flex items-center gap-3 rounded-full px-6 sm:px-8 py-4 bg-primary
				text-primary-foreground shadow-xl backdrop-blur-md"
	>
		<p class="text-xl sm:text-3xl">
			{{ formatPrice(cartStore.totalPrice) }}
		</p>

		<ChevronRight class="w-7 h-7" />
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
