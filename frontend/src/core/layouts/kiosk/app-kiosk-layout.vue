<script setup lang="ts">
import { useNetwork } from '@vueuse/core'
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// Kiosk & Services
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'

// UI Components
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { useInactivityTimer } from '@/core/hooks/use-inactivity-timer.hooks'
import { getKaspiConfig } from '@/core/integrations/kaspi.service'
import { cn } from '@/core/utils/tailwind.utils'
import KioskHomeCart from '@/modules/kiosk/products/components/home/kiosk-home-cart.vue'
import { WifiOff } from 'lucide-vue-next'
import { defineAsyncComponent } from 'vue'

// Lazy load the cart dialog
const KioskCartDialog = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-wrapper.vue')
)

// Store & Router
const cartStore = useCartStore()
const router = useRouter()

// Network
const { isOnline } = useNetwork()

// Inactivity constants
const INACTIVITY_MS = 5 * 60 * 1000 // 5 minutes

// On inactivity, reset store states and go to landing
function handleInactivity() {
  cartStore.clearCart()
  cartStore.closeModal()

  // Only navigate if online (use your logic if offline)
  if (isOnline.value) {
    router.push({ name: getRouteName('KIOSK_LANDING') })
  }
}

// We set up the composable
useInactivityTimer(INACTIVITY_MS, handleInactivity, true)

// Some routes should NOT show the cart
const omitShowCartPages: RouteKey[] = ['KIOSK_LANDING', 'KIOSK_CHECKLIST', "KIOSK_CART_PAYMENT", "KIOSK_CART_PAYMENT_SUCCESS"]

const showCart = computed(() => {
  if (cartStore.isEmpty) return false
  const currentRouteName = router.currentRoute.value.name as RouteKey
  return !omitShowCartPages.includes(currentRouteName)
})

// Check kiosk config on mount
onMounted(() => {
  if (!Boolean(getKaspiConfig())) {
    router.replace({ name: getRouteName('KIOSK_CHECKLIST') })
  }
})
</script>

<template>
	<div class="relative bg-slate-100 w-full min-h-screen no-scrollbar">
		<!-- If offline, show overlay -->
		<div
			v-if="!isOnline"
			class="absolute inset-0 flex flex-col justify-center items-center bg-white p-6 text-center"
		>
			<WifiOff class="size-20 text-red-600" />
			<h1 class="mt-6 font-bold text-gray-900 text-3xl">Сеть недоступна</h1>
			<p class="mt-4 max-w-md text-gray-600 text-lg">
				К сожалению, терминал временно не может принимать заказы. Пожалуйста, обратитесь к
				персоналу.
			</p>
		</div>

		<!-- If online, show kiosk UI -->
		<template v-else>
			<main :class="cn('relative w-full h-full no-scrollbar', showCart && 'pb-40')">
				<RouterView v-slot="{ Component, route }">
					<template v-if="Component">
						<Transition
							name="fade-slide"
							mode="out-in"
						>
							<Suspense>
								<div :key="route.matched[0]?.path">
									<component :is="Component" />
								</div>
								<template #fallback>
									<PageLoader />
								</template>
							</Suspense>
						</Transition>
					</template>
				</RouterView>
			</main>

			<!-- Show the cart button if conditions pass -->
			<div
				v-if="showCart"
				class="bottom-8 z-50 fixed flex justify-center w-full pointer-events-none"
				aria-label="Открыть корзину"
			>
				<KioskHomeCart />
			</div>

			<!-- Cart dialog -->
			<KioskCartDialog v-if="cartStore.isModalOpen" />
		</template>
	</div>
</template>

<style lang="scss">
/* Smooth transition effect */
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
