<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { useToast } from '@/core/components/ui/toast'
import { getKaspiConfig, KaspiService } from '@/core/integrations/kaspi.service'
import type { OrderDTO, TransactionDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { PaymentMethod, type PaymentOption } from '@/modules/kiosk/cart/models/kiosk-cart.models'
import { useMutation } from '@tanstack/vue-query'
import { QrCode } from 'lucide-vue-next'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, toRefs } from 'vue'

// ----- Props & Emits -----
const props = defineProps<{
  isOpen: boolean
  selectedPayment: PaymentMethod | null
  order: OrderDTO | null
}>()

const { order } = toRefs(props)

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

const PAYMENT_TIMEOUT = 120_000 // 2 minutes
const STORAGE_KEY = 'ZEEP_CART_PAYMENT_COUNTDOWN'

const isTest: boolean = import.meta.env.VITE_TEST_PAYMENT === 'true'

	const { toast } = useToast()

	// Reactive state
	const selectedMethod = ref<PaymentMethod | null>(null)
	const isPaying = ref(false)
	const errorMessage = ref('')
	const countdown = ref(PAYMENT_TIMEOUT / 1000)
	const hasTimerStarted = ref(false)
	const isTimeoutReached = ref(false)

	// Store interval ID to ensure single timer
	const timerId = ref<number | null>(null)

	// Kaspi mock / real integration

	/* -------------------------------------
	 * Mutations
	 * ------------------------------------- */
	const successOrderPaymentMutation = useMutation({
		mutationFn: ({
			order,
			transactionDTO,
		}: {
			order: OrderDTO
			transactionDTO: TransactionDTO
		}) => {
			return ordersService.successOrderPayment(order.id, transactionDTO)
		},
		onSuccess: async () => {
			console.log('Calling onProceed...')

			try {
        emit("proceed")
        console.log('onProceed executed successfully')
			} catch (error) {
				console.error('onProceed error:', error)
			}
			await nextTick()

			resetPaymentFlow(false)
		},
		onError: err => {
			console.error('Payment confirmation failed:', err)
			toast({
				title: 'Ошибка подтверждения',
				description: 'Не удалось подтвердить платеж. Обратитесь к персоналу',
				variant: 'destructive',
			})
			resetPaymentFlow(false)
		},
	})

	const awaitPayment = async (order: OrderDTO) => {
		try {
			if (!order) throw new Error('Order not found')

			const kaspiConfig = getKaspiConfig()
			if (!kaspiConfig) throw new Error('Kaspi config not found')

			const kaspiService = new KaspiService(kaspiConfig)
			const transaction = isTest
				? await kaspiService.awaitPaymentTest()
				: await kaspiService.awaitPayment(order.total)

			console.log('Transaction Data:', transaction)

			if (!transaction || !transaction.transactionId) {
				console.error('Transaction data missing:', transaction)
				throw new Error('Transaction data missing: ' + transaction)
			}

			await successOrderPaymentMutation.mutateAsync({
				order,
				transactionDTO: transaction,
			})
		} catch (err) {
			console.error('Payment wait failed:', err)
			toast({
				title: 'Ошибка оплаты',
				description: 'Не удалось дождаться оплаты. Попробуйте снова.',
				variant: 'destructive',
			})
			resetPaymentFlow(true)
		}
	}

	const failPaymentMutation = useMutation({
		mutationFn: async (orderId: number) => {
			return ordersService.failOrderPayment(orderId)
		},
		onSuccess: () => {
			errorMessage.value = 'Платеж не завершен вовремя. Пожалуйста, попробуйте снова.'
			resetPaymentFlow(false)
		},
		onError: () => {
			resetPaymentFlow(false)
		},
	})

	/* -------------------------------------
	 * Methods
	 * ------------------------------------- */
	async function selectMethod(method: PaymentMethod) {
		selectedMethod.value = method
		errorMessage.value = ''
		startPaymentProcess()
	}

	async function startPaymentProcess() {
		if (!order.value || !selectedMethod.value) {
			errorMessage.value = 'Номер заказа или способ оплаты отсутствуют.'
			return
		}
		initializeTimer(true)
		isPaying.value = true
		errorMessage.value = ''

		await awaitPayment(order.value)
	}

	async function failPayment() {
		if (order.value) {
			await failPaymentMutation.mutateAsync(order.value.id)
		}
	}

	function initializeTimer(forceRestart = false) {
		clearTimer()

		if (forceRestart) {
			sessionStorage.setItem(STORAGE_KEY, Date.now().toString())
		}

		const storedTime = sessionStorage.getItem(STORAGE_KEY)
		if (!storedTime) return

		const elapsed = Date.now() - Number(storedTime)
		const remaining = PAYMENT_TIMEOUT - elapsed

		if (remaining <= 0) {
			isTimeoutReached.value = true
			failPayment()
			return
		}

		countdown.value = Math.ceil(remaining / 1000)
		hasTimerStarted.value = true

		timerId.value = window.setInterval(() => {
			countdown.value -= 1
			if (countdown.value <= 0) {
				clearTimer()
				isTimeoutReached.value = true
				failPayment()
			}
		}, 1000)
	}

	function clearTimer() {
		if (timerId.value) {
			clearInterval(timerId.value)
			timerId.value = null
		}
	}

	function resetPaymentFlow(goBack: boolean) {
		// Stop countdown
		clearTimer()
		sessionStorage.removeItem(STORAGE_KEY)

		// Reset reactive state
		isPaying.value = false
		hasTimerStarted.value = false
		countdown.value = 120
		selectedMethod.value = null
		errorMessage.value = ''
		isTimeoutReached.value = false

		if (goBack) {
      emit('back')
		}
	}

	// Computed
	const isCountdownVisible = computed(() => hasTimerStarted.value && isPaying.value)
	const displayCountdown = computed(() => {
		const mm = Math.floor(countdown.value / 60)
		const ss = countdown.value % 60
		return `${mm.toString().padStart(2, '0')}:${ss.toString().padStart(2, '0')}`
	})

  onMounted(() => {
    if (order.value) {
      initializeTimer(false)
    }
  })

	onBeforeUnmount(() => {
		clearTimer()
	})

// ----- Handlers -----
async function onSelectMethod(method: PaymentMethod) {
  emit('update:selectedPayment', method)
  await selectMethod(method)
}

function handleClose() {
  resetPaymentFlow(false)
  emit('close')
}

function handleBack() {
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
