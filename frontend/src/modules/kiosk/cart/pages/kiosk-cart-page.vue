<template>
	<div class="relative px-12 pt-safe pb-44 no-scrollbar">
		<!-- Header Section -->
		<section
			class="top-0 left-0 fixed flex justify-between items-center gap-4 bg-white/30 backdrop-blur-md px-12 py-5 sm:py-7 w-full"
		>
			<button @click="goBack">
				<Icon
					icon="mingcute:left-line"
					class="text-3xl sm:text-5xl"
				/>
			</button>

			<p class="font-medium text-xl sm:text-2xl">Корзина</p>

			<button @click="clearCart">
				<Icon
					icon="mingcute:delete-line"
					class="text-2xl text-destructive sm:text-4xl"
				/>
			</button>
		</section>

		<!-- Cart Items -->
		<section class="flex flex-col gap-1 sm:gap-2 pt-28">
			<div
				v-for="cartItem in cartItemsArray"
				:key="cartItem.key"
				@click="openUpdateDialog(cartItem)"
			>
				<KioskCartItem
					:item="cartItem"
					:key="cartItem.key"
				/>
			</div>

			<!-- Update Cart Dialog -->
			<KioskCartUpdateItem
				v-if="selectedCartItem"
				:is-open="!!selectedCartItem"
				:cart-item="selectedCartItem"
				@close="closeUpdateDialog"
				@update="handleUpdate"
			/>
		</section>

		<!-- Total Price & Checkout -->
		<div
			v-if="cartItemsArray.length"
			class="bottom-10 left-0 fixed flex justify-center w-full"
		>
			<div
				class="flex items-center gap-12 bg-primary backdrop-blur-md sm:px-4 sm:py-4 p-3 pl-8 sm:pl-10 rounded-full text-white p"
			>
				<div>
					<p class="text-sm sm:text-lg">Итого ({{ cartStore.totalItems }})</p>
					<p class="font-medium text-xl sm:text-3xl">
						{{ formatPrice(totalPrice) }}
					</p>
				</div>

				<KioskCartCheckout />
			</div>
		</div>

		<!-- Empty Cart -->
		<div
			v-else
			class="py-10 text-center"
		>
			<p class="font-medium text-xl sm:text-2xl">Ваша корзина пуста</p>
			<button
				@click="goBack"
				class="bg-primary mt-4 px-6 py-4 rounded-2xl text-base text-primary-foreground sm:text-xl"
			>
				Вернуться к покупкам
			</button>
		</div>

		<!-- Suggested Products -->
		<!-- TODO: add food here for suggestions -->
		<!-- <section
			class="mt-6 sm:mt-8"
			v-if="suggestedProducts.length"
		>
			<p class="font-medium text-xl sm:text-2xl">Добавьте к заказу</p>
			<div class="gap-2 grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 mt-4">
				<KioskCartSuggestProduct
					v-for="product in suggestedProducts"
					:key="product.id"
					:product="product"
				/>
			</div>
		</section> -->
	</div>
</template>

<script setup lang="ts">
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import KioskCartCheckout from '@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout.vue'
import KioskCartItem from '@/modules/kiosk/cart/components/kiosk-cart-item.vue'
import { useCartStore, type CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import { Icon } from '@iconify/vue'
import { computed, defineAsyncComponent, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter();
const cartStore = useCartStore();

const KioskCartUpdateItem = defineAsyncComponent(() =>
  import('@/modules/kiosk/cart/components/kiosk-cart-update-item.vue')
);

const cartItemsArray = computed(() => Object.values(cartStore.cartItems));
const totalPrice = computed(() => cartStore.totalPrice);


const selectedCartItem = ref<CartItem | null>(null);

const goBack = () => {
  router.back();
};

const clearCart = () => {
  cartStore.clearCart();
};

const openUpdateDialog = (cartItem: CartItem) => {
  selectedCartItem.value = cartItem;
};

const closeUpdateDialog = () => {
  selectedCartItem.value = null;
};

const handleUpdate = (updatedSize: StoreProductSizeDetailsDTO, updatedAdditives: StoreAdditiveCategoryItemDTO[]) => {
  if (!selectedCartItem.value) return;
  cartStore.updateCartItem(selectedCartItem.value.key, {
    size: updatedSize,
    additives: updatedAdditives,
  });
  closeUpdateDialog();
};
</script>

<style scoped>
/* Add your styles here */
</style>
