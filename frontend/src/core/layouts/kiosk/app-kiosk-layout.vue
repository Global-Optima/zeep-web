<script setup lang="ts">
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import { useCartStore } from "@/modules/kiosk/cart/stores/cart.store"
import { useSelectedProductStore } from "@/modules/kiosk/products/stores/current-product.store"
import { defineAsyncComponent, onBeforeUnmount, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const KioskDetailsModal = defineAsyncComponent(() =>
  import('@/modules/kiosk/products/components/details/kiosk-details-modal.vue')
);

const pagesOmitRedirect: RouteKey[] = [ 'KIOSK_ORDERS', 'KIOSK_ORDERS_DISPLAY']

const router = useRouter()
const productStore = useSelectedProductStore()
const cartStore = useCartStore()
let inactivityTimeout: ReturnType<typeof setTimeout> | null = null
const inactivityDuration = 300 * 1000

const resetAppStates = () => {
  productStore.closeBottomSheet()
  cartStore.clearCart()
}

const resetInactivityTimer = () => {
  if (inactivityTimeout) {
    clearTimeout(inactivityTimeout)
  }
  inactivityTimeout = setTimeout(() => {
    const currentRouteName = router.currentRoute.value.name
    if (pagesOmitRedirect.some(page => page === currentRouteName)) return

    resetAppStates()
    router.push({ name: getRouteName('KIOSK_LANDING') })
  }, inactivityDuration)
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
	<div class="w-full min-h-screen overflow-x-hidden no-scrollbar">
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
	</div>
</template>

<style lang="scss">
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
