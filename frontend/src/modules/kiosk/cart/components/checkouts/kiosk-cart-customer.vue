<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { KIOSK_CUSTOMER_ADJECTIVES, KIOSK_CUSTOMER_NOUNS } from '@/modules/kiosk/cart/components/checkouts/customer-names-dictionary'
import KioskCartItem from '@/modules/kiosk/cart/components/kiosk-cart-item.vue'
import KioskCartUpdateItem from '@/modules/kiosk/cart/components/kiosk-cart-update-item.vue'
import { useCartStore, type CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import { ShoppingBasket, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

// Props & Emits
const props = defineProps<{
  open: boolean;
  customerName: string;
}>();

const emits = defineEmits<{
  (e: 'toggle', value: boolean): void;
  (e: 'next'): void;
  (e: 'update:customerName', value: string): void;
}>();

// Cart Store
const cartStore = useCartStore();
const router = useRouter()

const cartItemsArray = computed(() => Object.values(cartStore.cartItems));
const totalPrice = computed(() => cartStore.totalPrice);
const showCartItems = computed(() => cartStore.totalItems > 0)

// Cart Item Management
const selectedCartItem = ref<CartItem | null>(null);

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

// Unique Name Generator
const generateUniqueName = () => {
  const randomAdjective = KIOSK_CUSTOMER_ADJECTIVES[Math.floor(Math.random() * KIOSK_CUSTOMER_ADJECTIVES.length)];
  const randomNoun = KIOSK_CUSTOMER_NOUNS[Math.floor(Math.random() * KIOSK_CUSTOMER_NOUNS.length)];

  emits('update:customerName', `${randomAdjective} ${randomNoun}`);
};

// Validation State
const showValidationError = ref(false);

// Input Validation
const errorMessage = computed(() => {
  if (!showValidationError.value) return ''; // Do not show errors initially
  if (!props.customerName.trim()) return 'Пожалуйста, введите ваше имя.';
  if (props.customerName.length < 3) return 'Имя должно содержать минимум 3 символа.';
  return '';
});

// Proceed to next step with validation
const handleProceed = () => {
  showValidationError.value = true; // Only show error after attempt to proceed

  if (!errorMessage.value) {
    emits('next'); // Proceed only if validation passes
  }
};

const returnToMenu = () => {
  emits('toggle', false)
  router.push({name: getRouteName("KIOSK_HOME")})
}
</script>

<template>
	<Dialog
		:open="open"
		@update:open="v => emits('toggle', v)"
	>
		<DialogContent
			v-if="!showCartItems"
			:include-close-button="false"
			class="p-6 sm:rounded-[36px] min-h-[300px]"
		>
			<div class="flex flex-col justify-center items-center space-y-6 h-full">
				<!-- Icon and message -->
				<div class="flex flex-col items-center space-y-6 text-center">
					<ShoppingBasket class="size-16 text-gray-400" />
					<h2 class="mt-6 font-semibold text-slate-800 text-3xl">Ваша корзина пуста</h2>
					<p class="mt-6 text-slate-500 text-xl">
						Добавьте товары из меню, чтобы начать оформление заказа
					</p>
				</div>

				<!-- Return button -->
				<Button
					class="!mt-8 px-6 h-14 text-xl"
					@click="returnToMenu"
				>
					Вернуться к меню
				</Button>
			</div>
		</DialogContent>

		<DialogContent
			v-else
			:include-close-button="false"
			class="grid-rows-[auto_minmax(0,1fr)_auto] p-4 sm:rounded-[36px] max-w-2xl max-h-[80dvh]"
		>
			<!-- Header -->
			<DialogHeader class="p-6 pb-0">
				<div class="flex justify-between items-start gap-4">
					<DialogTitle class="text-4xl">Детали заказа</DialogTitle>
					<Button
						variant="ghost"
						size="icon"
						@click="clearCart"
					>
						<Trash
							stroke-width="1.6"
							class="size-9 text-red-500"
						/>
					</Button>
				</div>
			</DialogHeader>

			<!-- Cart Items -->
			<div class="gap-4 grid px-6 py-4 overflow-y-auto no-scrollbar">
				<div class="flex flex-col gap-6 h-[50dvh]">
					<div
						v-for="cartItem in cartItemsArray"
						:key="cartItem.key"
						@click="openUpdateDialog(cartItem)"
					>
						<KioskCartItem :item="cartItem" />
					</div>

					<!-- Update Cart Dialog -->
					<KioskCartUpdateItem
						v-if="selectedCartItem"
						:is-open="!!selectedCartItem"
						:cart-item="selectedCartItem"
						@close="closeUpdateDialog"
						@update="handleUpdate"
					/>
				</div>
			</div>

			<!-- Footer -->
			<DialogFooter class="block px-6 pb-6 w-full">
				<div>
					<Input
						:model-value="customerName"
						@update:model-value="v => emits('update:customerName', v.toString())"
						class="bg-slate-100 px-6 py-8 border-none rounded-xl w-full text-xl"
						placeholder="Введите ваше имя"
					/>

					<p
						v-if="errorMessage"
						class="mt-2 px-6 text-red-500 text-lg"
					>
						{{ errorMessage }}
					</p>

					<div class="flex justify-center mt-4">
						<Button
							type="button"
							variant="ghost"
							class="font-medium text-blue-500 hover:text-blue-500 text-xl"
							@click="generateUniqueName"
						>
							Сгенерировать имя
						</Button>
					</div>
				</div>

				<div class="flex justify-between items-center gap-4 mt-6">
					<p class="font-semibold text-primary text-4xl">{{ formatPrice(totalPrice) }}</p>

					<Button
						size="lg"
						class="py-6 text-xl"
						@click="handleProceed"
					>
						Оплатить
					</Button>
				</div>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
