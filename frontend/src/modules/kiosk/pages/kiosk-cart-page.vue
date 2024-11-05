<template>
	<div class="px-4 relative">
		<section class="flex items-center justify-between gap-4 py-5 sm:py-7">
			<button @click="$router.back">
				<Icon
					icon="mingcute:left-line"
					class="text-3xl sm:text-4xl"
				/>
			</button>

			<p class="text-xl sm:text-2xl font-medium">Корзина</p>

			<button>
				<Icon
					icon="mingcute:delete-line"
					class="text-2xl sm:text-4xl text-destructive"
				/>
			</button>
		</section>

		<section class="flex flex-col gap-1 sm:gap-2">
			<KioskCartItem
				v-for="cartItem in cartItems"
				:item="cartItem"
				:key="cartItem.product.id"
			/>
		</section>

		<div class="mt-4 sm:mt-6">
			<KioskDetailsEnergy />
		</div>

		<section class="mt-6 sm:mt-8">
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
import KioskCartItem from '@/modules/kiosk/components/cart/kiosk-cart-item.vue'
import KioskCartSuggestProduct from '@/modules/kiosk/components/cart/kiosk-cart-suggest-product.vue'
import KioskDetailsEnergy from '@/modules/kiosk/components/details/kiosk-details-energy.vue'
import type { CartItems } from '@/modules/kiosk/models/cart.models'
import type { Products } from '@/modules/products/models/product.model'
import { Icon } from '@iconify/vue'
import { ref } from 'vue'

const cartItems = ref<CartItems[]>([
	{
		product: {
      id: 1,
		title: 'Латте с молоком',
		price: 2323,
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/394/931/small_2x/cup-of-hot-cappuccino-coffee-illustration-generative-ai-png.png',
      category: "Кофе"
  },
		quantity: 1,
    additives: [
      {
        id: 1,
        name: 'Карамель',
        price: 300,
        imageUrl: ''
      },
      {
        id: 2,
        name: 'Ваниль',
        price: 400,
        imageUrl: ''
      }
    ],
	},
	{
		product: {
      id: 2,
		title: 'Эспрессо',
		price: 1450,
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/394/931/small_2x/cup-of-hot-cappuccino-coffee-illustration-generative-ai-png.png',
      category: "Кофе"
  },
		quantity: 1,
    additives: [],
	},
	{
		product: {
      id: 1,
		title: 'Мокко с шоколадной пенкой и посыпкой из карамели',
		price: 3000,
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/394/931/small_2x/cup-of-hot-cappuccino-coffee-illustration-generative-ai-png.png',
      category: "Кофе"
  },
		quantity: 1,
    additives: [
      {
        id: 1,
        name: 'Шоколад',
        price: 300,
        imageUrl: ''
      },
      {
        id: 2,
        name: 'Карамель',
        price: 400,
        imageUrl: ''
      },
      {
        id: 3,
        name: 'Сырная пенка',
        price: 400,
        imageUrl: ''
      },
      {
        id: 4,
        name: 'Фисташковый сироп',
        price: 500,
        imageUrl: ''
      }
    ],
	},
]);

const suggestedProducts = ref<Products[]>([
  {
    id: 1,
    title: 'Круассан с шоколадом',
    price: 150,
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/340/374/small_2x/fresh-cooked-yellow-croissant-png.png',
    category: 'Еда'
  },
  {
    id: 2,
    title: 'Круассан с курицей',
    price: 150,
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/340/374/small_2x/fresh-cooked-yellow-croissant-png.png',
    category: 'Еда'
  },
  {
    id: 3,
    title: 'Круассан с семгой',
    price: 150,
    image: 'https://static.vecteezy.com/system/resources/thumbnails/017/340/374/small_2x/fresh-cooked-yellow-croissant-png.png',
    category: 'Еда'
  },
])
</script>

<style lang="scss"></style>
