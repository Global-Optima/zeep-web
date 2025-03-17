import { useToast } from '@/core/components/ui/toast'
import { useMutation } from '@tanstack/vue-query'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

// Services
import { getKaspiConfig, KaspiService } from '@/core/integrations/kaspi.service'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'

// Types
import type { OrderDTO, TransactionDTO } from '@/modules/admin/store-orders/models/orders.models'
import { PaymentMethod } from '@/modules/kiosk/cart/models/kiosk-cart.models'

const PAYMENT_TIMEOUT = 120_000 // 2 minutes
const STORAGE_KEY = 'ZEEP_CART_PAYMENT_COUNTDOWN'

export function useCartPayment(order: OrderDTO | null, onProceed: () => void, onBack: () => void) {
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
	let timerId: number | null = null

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
				onProceed()
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
			console.log('Here 4')
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
		if (!order || !selectedMethod.value) {
			errorMessage.value = 'Номер заказа или способ оплаты отсутствуют.'
			return
		}
		initializeTimer(true)
		isPaying.value = true
		errorMessage.value = ''

		await awaitPayment(order)
	}

	async function failPayment() {
		if (order) {
			await failPaymentMutation.mutateAsync(order.id)
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

		timerId = window.setInterval(() => {
			countdown.value -= 1
			if (countdown.value <= 0) {
				clearTimer()
				isTimeoutReached.value = true
				failPayment()
			}
		}, 1000)
	}

	function clearTimer() {
		if (timerId) {
			clearInterval(timerId)
			timerId = null
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
			onBack()
		}
	}

	// Computed
	const isCountdownVisible = computed(() => hasTimerStarted.value && isPaying.value)
	const displayCountdown = computed(() => {
		const mm = Math.floor(countdown.value / 60)
		const ss = countdown.value % 60
		return `${mm.toString().padStart(2, '0')}:${ss.toString().padStart(2, '0')}`
	})

	// Lifecycle
	onMounted(() => {
		if (order) {
			initializeTimer(false)
		}
	})

	onBeforeUnmount(() => {
		clearTimer()
	})

	return {
		// State
		selectedMethod,
		isPaying,
		errorMessage,
		countdown,
		isCountdownVisible,
		displayCountdown,
		isTimeoutReached,

		// Methods
		selectMethod,
		resetPaymentFlow,
		startPaymentProcess,
	}
}
