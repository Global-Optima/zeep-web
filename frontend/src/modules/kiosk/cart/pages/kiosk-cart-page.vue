<template>
	<div class="px-4 relative pb-10">
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
				class="flex items-center gap-12 rounded-3xl px-8 py-4 bg-slate-700/70 text-white backdrop-blur-md"
			>
				<div>
					<p class="text-base sm:text-xl">Итого ({{ cartStore.totalItems }})</p>
					<p class="text-2xl sm:text-3xl font-medium sm:mt-1">
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
import type { Products } from "@/modules/products/models/product.model"
import { Icon } from '@iconify/vue'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

// Initialize Router and Cart Store
const router = useRouter();
const cartStore = useCartStore();

// Computed Properties
const cartItemsArray = computed(() => Object.values(cartStore.cartItems));

const totalPrice = computed(() => cartStore.totalPrice);

// Suggested Products (Fetch from API or define as needed)
const suggestedProducts = ref<Products[]>([
  {
    id: 1,
    title: 'Круассан с шоколадом',
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/340/374/small_2x/fresh-cooked-yellow-croissant-png.png',
    category: 'Еда',
    startPrice: 1000
  },
  {
    id: 2,
    title: 'Круассан с курицей',
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/340/374/small_2x/fresh-cooked-yellow-croissant-png.png',
    category: 'Еда',
    startPrice: 1000
  },
  {
    id: 3,
    title: 'Круассан с семгой',
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/340/374/small_2x/fresh-cooked-yellow-croissant-png.png',
    category: 'Еда',
    startPrice: 1000
  },
])
// Methods
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
