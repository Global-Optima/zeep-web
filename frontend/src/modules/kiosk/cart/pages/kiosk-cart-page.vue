<template>
	<div class="px-4 relative pt-safe pb-44 no-scrollbar">
		<section class="flex items-center justify-between gap-4 py-5 sm:py-7">
			<button @click="goBack">
				<Icon
					icon="mingcute:left-line"
					class="text-3xl sm:text-4xl"
				/>
			</button>

			<p class="text-xl sm:text-2xl font-medium">Корзина</p>

			<button @click="clearCart">
				<Icon
					icon="mingcute:delete-line"
					class="text-2xl sm:text-4xl text-destructive"
				/>
			</button>
		</section>

		<section class="flex flex-col gap-1 sm:gap-2">
			<KioskCartItem
				v-for="cartItem in cartItemsArray"
				:item="cartItem"
				:key="cartItem.key"
			/>
		</section>

		<div
			v-if="cartItemsArray.length"
			class="fixed bottom-10 left-0 w-full flex justify-center"
		>
			<div
				class="flex items-center gap-12 rounded-3xl px-6 sm:px-8 py-4 bg-slate-700/70 text-white backdrop-blur-md"
			>
				<div>
					<p class="text-sm sm:text-xl">Итого ({{ cartStore.totalItems }})</p>
					<p class="text-xl sm:text-3xl font-medium sm:mt-1">
						{{ formatPrice(totalPrice) }}
					</p>
				</div>

				<KioskCartCheckout />
			</div>
		</div>

		<div
			v-else
			class="text-center py-10"
		>
			<p class="text-xl sm:text-2xl font-medium">Ваша корзина пуста</p>
			<button
				@click="goBack"
				class="mt-4 px-6 py-4 bg-primary text-primary-foreground rounded-2xl text-base sm:text-xl"
			>
				Вернуться к покупкам
			</button>
		</div>

		<section
			class="mt-6 sm:mt-8"
			v-if="suggestedProducts.length"
		>
			<p class="text-xl sm:text-2xl font-medium">Добавьте к заказу</p>
			<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 mt-4">
				<KioskCartSuggestProduct
					v-for="product in suggestedProducts"
					:key="product.id"
					:product="product"
				/>
			</div>
		</section>
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import KioskCartCheckout from '@/modules/kiosk/cart/components/kiosk-cart-checkout.vue'
import KioskCartItem from '@/modules/kiosk/cart/components/kiosk-cart-item.vue'
import KioskCartSuggestProduct from '@/modules/kiosk/cart/components/kiosk-cart-suggest-product.vue'
import { useCartStore } from "@/modules/kiosk/cart/stores/cart.store"
import type { StoreProductDetails } from '@/modules/kiosk/products/models/product.model'
import { Icon } from '@iconify/vue'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter();
const cartStore = useCartStore();

const cartItemsArray = computed(() => Object.values(cartStore.cartItems));

const totalPrice = computed(() => cartStore.totalPrice);

const suggestedProducts = ref<StoreProductDetails[]>([
  {
    id: 1000,
    name: 'Круассан с шоколадом',
    imageUrl: 'https://lamin8patisserie.com.au/cdn/shop/products/Chocolatecroissant_530x@2x.png?v=1611018458',
    basePrice: 1200,
    description: 'Круассан с шоколадом',
    sizes: [],
    defaultAdditives: [],
    recipeSteps: []
  },
  {
    id: 2,
    name: 'Круассан с курицей',
    imageUrl: 'https://static.vecteezy.com/system/resources/previews/044/308/224/non_2x/croissant-sanwich-isolated-on-transparent-background-png.png',
    basePrice: 1240,
    description: 'Круассан с курицей',
    sizes: [],
    defaultAdditives: [],
    recipeSteps: []
  },
])

const goBack = () => {
  router.back();
};

const clearCart = () => {
  cartStore.clearCart();
};
</script>

<style lang="scss">
/* Add your styles here */
</style>
