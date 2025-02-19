<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import { useCartStore } from "@/modules/kiosk/cart/stores/cart.store"
import { useSelectedProductStore } from "@/modules/kiosk/products/stores/current-product.store"
import { useNetwork } from '@vueuse/core'
import { WifiOff } from 'lucide-vue-next'; // ✅ Import Lucide icon
import { defineAsyncComponent, onBeforeUnmount, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const KioskDetailsModal = defineAsyncComponent(() =>
  import('@/modules/kiosk/products/components/details/kiosk-details-modal.vue')
);

const router = useRouter()
const productStore = useSelectedProductStore()
const cartStore = useCartStore()

let inactivityTimeout: ReturnType<typeof setTimeout> | null = null
const inactivityDuration = 300 * 1000 // 5 minutes

const { isOnline } = useNetwork()

const resetAppStates = () => {
  productStore.closeBottomSheet()
  cartStore.clearCart()
}

const handleInactivity = () => {
  resetAppStates()

  if (isOnline.value) {
    router.push({ name: getRouteName('KIOSK_LANDING') })
  }
}

const resetInactivityTimer = () => {
  if (inactivityTimeout) {
    clearTimeout(inactivityTimeout)
  }
  inactivityTimeout = setTimeout(handleInactivity, inactivityDuration)
}

const activityEvents = ['mousemove', 'keydown', 'click', 'touchstart']

onMounted(() => {
  resetInactivityTimer()
  activityEvents.forEach(event => {
    window.addEventListener(event, resetInactivityTimer)
  })
})

onBeforeUnmount(() => {
  if (inactivityTimeout) {
    clearTimeout(inactivityTimeout)
  }
  activityEvents.forEach(event => {
    window.removeEventListener(event, resetInactivityTimer)
  })
})
</script>

<template>
	<div class="relative flex justify-center items-center bg-gray-100 w-full min-h-screen">
		<!-- ✅ Show only this screen when offline -->
		<div
			v-if="!isOnline"
			class="absolute inset-0 flex flex-col justify-center items-center bg-white p-6 text-center"
		>
			<WifiOff class="size-20 text-red-600" />
			<!-- Lucide icon -->
			<h1 class="mt-6 font-bold text-gray-900 text-3xl">Сеть недоступна</h1>
			<p class="mt-4 max-w-md text-gray-600 text-lg">
				К сожалению, терминал временно не может принимать заказы. Пожалуйста, обратитесь к
				персоналу.
			</p>
		</div>

		<!-- ✅ Show the normal app content only if online -->
		<template v-else>
			<main class="relative w-full h-full">
				<router-view v-slot="{ Component }">
					<transition
						name="fade-slide"
						mode="out-in"
					>
						<component :is="Component" />
					</transition>
				</router-view>
			</main>

			<KioskDetailsModal />
		</template>
	</div>
</template>

<style lang="scss">
/* ✅ Smooth transition effect */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.15s ease-in-out;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
}

.fade-slide-enter-to,
.fade-slide-leave-from {
  opacity: 1;
}
</style>
