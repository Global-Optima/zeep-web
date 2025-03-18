<template>
	<div class="flex flex-col justify-center items-center p-4 min-h-screen">
		<PageLoader v-if="isOrderPending" />

		<div
			v-else-if="order"
			class="bg-white shadow-2xl hover:shadow-2xl p-8 rounded-3xl w-full max-w-md text-center transition-all duration-500 transform"
		>
			<div class="flex justify-center items-center bg-green-100 mx-auto mb-6 rounded-full size-20">
				<Check class="size-10 text-green-500" />
			</div>

			<!-- Header -->
			<h1 class="mb-4 font-bold text-gray-800 text-3xl">Заказ подтвержден</h1>
			<p class="mb-8 text-gray-600 text-lg">
				Спасибо за покупку! Ваш заказ принят и обрабатывается.
			</p>

			<!-- Order Display Number -->
			<div class="mb-10">
				<p class="font-extrabold text-emerald-600 text-6xl">
					{{ order.displayNumber }}
				</p>

				<p class="mt-2 text-gray-500 text-base">Номер заказа</p>
			</div>

			<!-- Footer Button -->
			<Button
				@click="handleProceed"
				class="px-6 rounded-2xl h-12 font-semibold text-white text-xl"
			>
				Вернуться в меню
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useQuery } from '@tanstack/vue-query'
import { Check } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// Router setup
const router = useRouter()
const route = useRoute()
const cartStore = useCartStore()

const orderId = route.params.orderId as string

const { data: order, isPending: isOrderPending } = useQuery({
  queryKey: computed(() => ["success-orders", orderId]),
  queryFn: () => ordersService.getOrderById(Number(orderId)),
  enabled: computed(() => !isNaN(Number(orderId)))
})

const handleProceed = () => {
  cartStore.clearCart()
  router.push({ name: getRouteName("KIOSK_HOME") })
}
</script>

<style scoped lang="scss">
/* Additional styles can be added here if needed */
</style>
