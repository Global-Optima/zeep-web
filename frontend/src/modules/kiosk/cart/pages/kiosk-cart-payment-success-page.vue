<template>
	<div class="flex justify-center items-center bg-slate-100 h-screen">
		<PageLoader v-if="isOrderPending" />

		<div
			v-else-if="order"
			class="bg-white shadow-2xl hover:shadow-2xl p-8 rounded-3xl w-full max-w-md text-center transition-all duration-500 transform"
		>
			<div
				class="flex justify-center items-center bg-green-100 mx-auto mb-6 rounded-full w-20 h-20"
			>
				<Check class="w-10 h-10 text-green-500" />
			</div>

			<!-- Header -->
			<h1 class="mb-4 font-bold text-gray-800 text-3xl">Заказ подтвержден</h1>
			<p class="mb-8 text-gray-600 text-lg">
				Спасибо за покупку! Ваш заказ принят и обрабатывается.
			</p>

			<!-- Order Display Number -->
			<div class="mb-8">
				<p class="font-extrabold text-emerald-600 text-7xl">
					{{ order.displayNumber }}
				</p>
				<p class="mt-4 text-gray-500 text-lg">Номер заказа</p>
			</div>

			<!-- Footer Button -->
			<Button
				@click="handleProceed"
				class="px-6 rounded-2xl h-12 font-normal text-white text-xl"
			>
				Вернуться в меню
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useQuery } from '@tanstack/vue-query'
import { Check } from 'lucide-vue-next'
import { computed, onBeforeUnmount, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const cartStore = useCartStore()
const { toast } = useToast()

const orderId = route.params.orderId as string

// Fetch order details using Vue Query.
const { data: order, isPending: isOrderPending, isError: isOrderError } = useQuery({
  queryKey: computed(() => ['success-orders', orderId]),
  queryFn: () => ordersService.getOrderById(Number(orderId)),
  enabled: computed(() => !isNaN(Number(orderId))),
  retry: 1,
})

// Watch the query state to handle errors or missing order data.
watch([isOrderError, order], ([isError, orderData]) => {
  if (isError || (!isOrderPending.value && !orderData)) {
    toast({
      title: 'Ошибка',
      description: 'Не удалось загрузить данные заказа',
      variant: 'destructive'
    })
    router.push({ name: getRouteName("KIOSK_HOME") })
  }
})

// Auto-redirect after 60 seconds.
let autoRedirectTimeout: number | null = null
const redirectDelay = 60

onMounted(() => {
  autoRedirectTimeout = window.setTimeout(() => {
    handleProceed()
  }, redirectDelay * 1000)
})

onBeforeUnmount(() => {
  if (autoRedirectTimeout) clearTimeout(autoRedirectTimeout)
})

// When the user confirms, clear the cart and navigate home.
const handleProceed = () => {
  if (autoRedirectTimeout) clearTimeout(autoRedirectTimeout)
  cartStore.clearCart()
  router.push({ name: getRouteName("KIOSK_HOME") })
}
</script>

<style scoped lang="scss">
/* Additional styles can be added here if needed */
</style>
