import { useToast } from '@/core/components/ui/toast'
import { useMutation } from '@tanstack/vue-query'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

// Services
import { getKaspiConfig, KaspiService } from '@/core/integrations/kaspi.service'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'

// Types
import type { OrderDTO, TransactionDTO } from '@/modules/admin/store-orders/models/orders.models'
import { PaymentMethod } from '@/modules/kiosk/cart/models/kiosk-cart.models'

const PAYMENT_TIMEOUT = 120_000 // 2 minutes
const STORAGE_KEY = 'ZEEP_CART_PAYMENT_COUNTDOWN'

export function useCartPayment(order: OrderDTO | null, onProceed: () => void, onBack: () => void) {
	const isTest: boolean =
		import.meta.env.VITE_TEST_PAYMENT === undefined
			? false
			: import.meta.env.VITE_TEST_PAYMENT === 'true'
				? true
				: false

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
		mutationFn: async ({
			order,
			transactionDTO,
		}: {
			order: OrderDTO
			transactionDTO: TransactionDTO
		}) => ordersService.successOrderPayment(order.id, transactionDTO),
		onSuccess: () => {
			onProceed()
			resetPaymentFlow(false)
		},
		onError: err => {
			console.error(err)
			errorMessage.value = 'Ошибка при отправки платежа. Просим обратиться к персоналу'
			toast({ title: 'Ошибка', description: errorMessage.value, variant: 'destructive' })
		},
	})

	const awaitPaymentMutation = useMutation({
		mutationFn: async ({ order }: { order: OrderDTO }) => {
			const kaspiConfig = getKaspiConfig()

			if (!kaspiConfig) throw new Error('Failed to get kaspi config')

			const kaspiService = new KaspiService(kaspiConfig)

			if (isTest) return await kaspiService.awaitPaymentTest()

			return await kaspiService.awaitPayment(order.total)
		},
		onSettled: (transaction, error, variables) => {
			if (transaction && !error) {
				successOrderPaymentMutation.mutate({ order: variables.order, transactionDTO: transaction })
			}
		},
		onError: err => {
			console.error(err)
			errorMessage.value = 'Ошибка при обработке платежа. Попробуйте снова.'
			toast({ title: 'Ошибка', description: errorMessage.value, variant: 'destructive' })
			resetPaymentFlow(true)
		},
	})

	const failPaymentMutation = useMutation({
		mutationFn: async (orderId: number) => {
			return ordersService.failOrderPayment(orderId)
		},
		onSuccess: () => {
			errorMessage.value = 'Платеж не завершен вовремя. Пожалуйста, попробуйте снова.'
			isPaying.value = false
			resetPaymentFlow(false)
		},
		onError: () => {
			resetPaymentFlow(false)
		},
	})

	/* -------------------------------------
	 * Methods
	 * ------------------------------------- */
	function selectMethod(method: PaymentMethod) {
		selectedMethod.value = method
		errorMessage.value = ''
		startPaymentProcess()
	}

	function startPaymentProcess() {
		if (!order || !selectedMethod.value) {
			errorMessage.value = 'Номер заказа или способ оплаты отсутствуют.'
			return
		}
		initializeTimer(true)
		isPaying.value = true
		errorMessage.value = ''

		awaitPaymentMutation.mutate({ order })
	}

	async function failPayment() {
		if (order) {
			await failPaymentMutation.mutateAsync(order.id)
		}
	}

	function initializeTimer(forceRestart = false) {
		// Clear any existing timer
		clearTimer()

		// If forcing a restart, store new start time
		if (forceRestart) {
			sessionStorage.setItem(STORAGE_KEY, Date.now().toString())
		}

		const storedTime = sessionStorage.getItem(STORAGE_KEY)
		if (!storedTime) return

		hasTimerStarted.value = true
		isTimeoutReached.value = false

		const elapsed = Date.now() - Number(storedTime)
		const remaining = PAYMENT_TIMEOUT - elapsed
		countdown.value = Math.ceil(remaining / 1000)

		if (countdown.value <= 0) {
			isTimeoutReached.value = true
			failPayment()
			return
		}

		// Launch a single setInterval
		timerId = window.setInterval(() => {
			countdown.value--
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
