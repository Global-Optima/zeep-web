<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import type { OrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { useCartPayment } from '@/modules/kiosk/cart/hooks/use-cart-payment.hook'
import { PaymentMethod, type PaymentOption } from '@/modules/kiosk/cart/models/kiosk-cart.models'
import { QrCode } from 'lucide-vue-next'

// ----- Props & Emits -----
const {order} = defineProps<{
  isOpen: boolean
  selectedPayment: PaymentMethod | null
  order: OrderDTO | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'back'): void
  (e: 'update:selectedPayment', val: PaymentMethod | null): void
  (e: 'proceed'): void
}>()

// ----- Payment Options -----
const paymentOptions: PaymentOption[] = [
  { id: PaymentMethod.KASPI, label: 'Kaspi QR', icon: QrCode },
]

// ----- Composable Payment Logic -----
const {
  selectedMethod,
  selectMethod,
  isPaying,
  errorMessage,
  isCountdownVisible,
  displayCountdown,
  resetPaymentFlow,
} = useCartPayment(
  order,
  () => emit('proceed'),
  () => emit('back')
)

// ----- Handlers -----
async function onSelectMethod(method: PaymentMethod) {
  // Update parent's v-model
  emit('update:selectedPayment', method)
  // Trigger payment logic
  await selectMethod(method)
}

function handleClose() {
  // If you want to reset on close, do so:
  resetPaymentFlow(false)
  emit('close')
}

function handleBack() {
  // Example: if user cancels payment but hasn't even started
  // a real payment request, or wants to go to previous screen
  resetPaymentFlow(false)
  emit('back')
}
</script>

<template>
	<Dialog
		:open="isOpen"
		@update:open="handleClose"
	>
		<DialogContent
			:include-close-button="false"
			class="space-y-8 bg-white sm:p-12 !rounded-[40px] sm:max-w-3xl text-black"
		>
			<DialogHeader>
				<!-- Toggle Title -->
				<DialogTitle
					v-if="!isPaying"
					class="font-medium text-gray-900 sm:text-4xl text-center"
				>
					Выберите способ оплаты
				</DialogTitle>
				<DialogTitle
					v-else
					class="font-medium text-gray-900 sm:text-3xl text-center"
				>
					Идет обработка платежа
				</DialogTitle>
			</DialogHeader>

			<!-- Payment Options -->
			<div
				v-if="!isPaying"
				class="flex justify-center gap-4 mb-4"
			>
				<div
					v-for="option in paymentOptions"
					:key="option.id"
					class="flex flex-col flex-1 justify-center items-center px-6 py-8 border-2 border-transparent rounded-2xl max-w-sm cursor-pointer"
					:class="{
            'bg-primary text-white': selectedMethod === option.id,
            'bg-gray-100': selectedMethod !== option.id
          }"
					@click="onSelectMethod(option.id)"
				>
					<component
						:is="option.icon"
						class="mb-4 w-14 h-14"
					/>
					<p class="font-medium sm:text-2xl text-center">{{ option.label }}</p>
				</div>
			</div>

			<!-- Payment Processing -->
			<div
				v-else
				class="flex flex-col items-center gap-6 mb-5 text-center"
			>
				<div v-if="isCountdownVisible">
					<p class="font-bold text-green-600 text-3xl">{{ displayCountdown }}</p>
				</div>
				<p class="text-2xl">
					{{
            selectedMethod === PaymentMethod.KASPI
              ? 'Откройте приложение Kaspi и просканируйте QR'
              : 'Идет обработка платежа...'
					}}
				</p>
			</div>

			<!-- Error Message -->
			<p
				v-if="errorMessage"
				class="text-red-500 text-xl text-center"
			>
				{{ errorMessage }}
			</p>

			<!-- Footer Buttons -->
			<div class="flex justify-center items-center gap-4 mt-4 w-full">
				<Button
					variant="ghost"
					@click="handleBack"
					:disabled="isPaying"
					class="sm:px-6 sm:py-8 sm:text-2xl"
				>
					Назад
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>
