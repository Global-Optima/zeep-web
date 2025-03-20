<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import { formatPrice } from '@/core/utils/price.utils'
import type { StoreAdditiveCategoryItemDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import type { CreateOrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import type { StoreProductSizeDetailsDTO } from '@/modules/admin/store-products/models/store-products.model'
import { KIOSK_CUSTOMER_ADJECTIVES, KIOSK_CUSTOMER_NOUNS } from '@/modules/kiosk/cart/components/checkouts/customer-names-dictionary'
import KioskCartItem from '@/modules/kiosk/cart/components/kiosk-cart-item.vue'
import KioskCartUpdateItem from '@/modules/kiosk/cart/components/kiosk-cart-update-item.vue'
import { useCartStore, type CartItem } from '@/modules/kiosk/cart/stores/cart.store'
import { useMutation } from '@tanstack/vue-query'
import { ShoppingBasket, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

// Cart Store
const cartStore = useCartStore()
const router = useRouter()

const {toastLocalizedError} = useAxiosLocaleToast()

// State
const customerName = ref('')
const selectedCartItem = ref<CartItem | null>(null)
const showValidationError = ref(false)
const createOrderMutation = useMutation({
  mutationFn: (orderDTO: CreateOrderDTO) => ordersService.createOrder(orderDTO),
  onError: (err: AxiosLocalizedError) => {
    toastLocalizedError(err, "Ошибка при создании заказа")
  }
})

// Computed
const cartItemsArray = computed(() => Object.values(cartStore.cartItems))
const totalPrice = computed(() => cartStore.totalPrice)
const showCartItems = computed(() => cartStore.totalItems > 0)
const errorMessage = computed(() => {
  if (!showValidationError.value) return ''
  if (!customerName.value.trim()) return 'Пожалуйста, введите ваше имя.'
  if (customerName.value.length < 3) return 'Имя должно содержать минимум 3 символа.'
  return ''
})

// Handlers
const generateUniqueName = () => {
  const adj = KIOSK_CUSTOMER_ADJECTIVES[Math.floor(Math.random() * KIOSK_CUSTOMER_ADJECTIVES.length)]
  const noun = KIOSK_CUSTOMER_NOUNS[Math.floor(Math.random() * KIOSK_CUSTOMER_NOUNS.length)]
  customerName.value = `${adj} ${noun}`
}

const handleProceed = async () => {
  showValidationError.value = true
  if (errorMessage.value) return

  const orderDTO: CreateOrderDTO = {
    customerName: customerName.value,
    subOrders: cartItemsArray.value.map(item => ({
      storeProductSizeId: item.size.id,
      quantity: item.quantity,
      storeAdditivesIds: item.additives.map(a => a.id)
    }))
  }

  const order = await createOrderMutation.mutateAsync(orderDTO)
  cartStore.toggleModal()
  router.push({ name: getRouteName("KIOSK_CART_PAYMENT"), params: { orderId: order.id } })
}

const onBackClick = () => {
  router.push({name: getRouteName("KIOSK_HOME")})
  cartStore.toggleModal()
}


const openUpdateDialog = (item: CartItem) => (selectedCartItem.value = item)
const closeUpdateDialog = () => (selectedCartItem.value = null)
const handleUpdate = (size: StoreProductSizeDetailsDTO, additives: StoreAdditiveCategoryItemDTO[]) => {
  if (selectedCartItem.value) {
    cartStore.updateCartItem(selectedCartItem.value.key, { size, additives })
    closeUpdateDialog()
  }
}
</script>

<template>
	<Dialog
		:open="cartStore.isModalOpen"
		@update:open="cartStore.toggleModal"
	>
		<DialogContent
			v-if="!showCartItems"
			:include-close-button="false"
			class="p-6 sm:rounded-[36px]"
		>
			<div class="flex flex-col items-center space-y-6 h-full">
				<ShoppingBasket class="size-16 text-gray-400" />
				<h2 class="font-semibold text-slate-800 text-3xl">Ваша корзина пуста</h2>
				<p class="text-slate-500 text-xl text-center">
					Добавьте товары из меню, чтобы начать оформление заказа
				</p>
				<Button
					class="!mt-8 px-6 h-14 text-xl"
					@click="onBackClick"
				>
					Вернуться к меню
				</Button>
			</div>
		</DialogContent>

		<DialogContent
			v-else
			:include-close-button="false"
			class="p-4 sm:rounded-[36px] max-w-2xl"
		>
			<DialogHeader class="p-6 pb-0">
				<div class="flex justify-between items-start gap-4">
					<DialogTitle class="text-4xl">Детали заказа</DialogTitle>
					<Button
						variant="ghost"
						size="icon"
						@click="cartStore.clearCart"
					>
						<Trash class="size-9 text-red-500" />
					</Button>
				</div>
			</DialogHeader>

			<div class="px-6 py-4 overflow-y-auto no-scrollbar">
				<div class="flex flex-col gap-6 h-[50dvh]">
					<div
						v-for="item in cartItemsArray"
						:key="item.key"
						@click="openUpdateDialog(item)"
					>
						<KioskCartItem :item="item" />
					</div>
					<KioskCartUpdateItem
						v-if="selectedCartItem"
						:is-open="!!selectedCartItem"
						:cart-item="selectedCartItem"
						@close="closeUpdateDialog"
						@update="handleUpdate"
					/>
				</div>
			</div>

			<DialogFooter class="block px-6 pb-6 w-full">
				<div>
					<Input
						v-model="customerName"
						class="bg-slate-100 px-6 py-8 rounded-xl w-full text-xl"
						placeholder="Введите ваше имя"
					/>
					<p
						v-if="errorMessage"
						class="mt-2 px-5 text-red-500 text-lg"
					>
						{{ errorMessage }}
					</p>
					<div class="flex justify-center mt-4">
						<Button
							variant="ghost"
							class="text-blue-500 text-xl"
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
