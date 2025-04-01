<script setup lang="ts">
import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import { getKaspiConfig, KaspiService } from '@/core/integrations/kaspi.service'
import {
  OrderStatus,
  type OrderDetailsDTO,
  type TransactionDTO
} from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import {
  PaymentMethod,
  type PaymentOption
} from '@/modules/kiosk/cart/models/kiosk-cart.models'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { TabletSmartphone } from 'lucide-vue-next'
import {
  computed,
  onBeforeUnmount,
  onMounted,
  reactive,
  watch
} from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const { toast } = useToast()

const orderId = route.params.orderId as string

// Payment configuration
const paymentOptions: PaymentOption[] = [
  { id: PaymentMethod.KASPI, label: 'Kaspi QR', icon: TabletSmartphone },
]

// Constants
const PAYMENT_TIMEOUT = 120_000 // 2 minutes
const STORAGE_KEY = 'ZEEP_CART_PAYMENT_COUNTDOWN'
const isTest = import.meta.env.VITE_TEST_PAYMENT === 'true'

// Group payment state
const paymentState = reactive({
  selectedMethod: null as PaymentMethod | null,
  isPaying: false,
  errorMessage: '',
  countdown: PAYMENT_TIMEOUT / 1000,
  hasTimerStarted: false,
  isTimeoutReached: false,
})

// Non-reactive timer ref
let timerId: number | null = null

// Query: Fetch order details
const {
  data: order,
  isPending: isOrderPending,
  isError: isOrderError,
  isSuccess: isOrderSuccess
} = useQuery({
  queryKey: computed(() => ['payment-order', orderId]),
  queryFn: () => ordersService.getOrderById(Number(orderId)),
  enabled: computed(() => !isNaN(Number(orderId))),
  retry: 1,
})

// Allowed statuses for continuing payment
const allowOrderStatuses: OrderStatus[] = [OrderStatus.WAITING_FOR_PAYMENT]

// Watch for order query results (but no longer start the timer here)
watch(
  [isOrderSuccess, isOrderError, order],
  ([success, error, orderData]) => {
    if (success) {
      if (!orderData) {
        toast({
          title: 'Ошибка',
          description: 'Заказ не найден',
          variant: 'destructive'
        })
        router.push({ name: getRouteName('KIOSK_HOME') })
      } else if (!allowOrderStatuses.includes(orderData.status)) {
        toast({
          title: 'Заказ уже оплачен',
          description: 'Платеж уже подтвержден, возвращаем в меню',
          variant: 'destructive'
        })
        router.push({ name: getRouteName('KIOSK_HOME') })
      }
    }

    if (error) {
      toast({
        title: 'Ошибка загрузки',
        description: 'Не удалось загрузить данные заказа',
        variant: 'destructive'
      })
      router.push({ name: getRouteName('KIOSK_HOME') })
    }
  }
)

// Mutation: Confirm successful order payment
const successOrderPaymentMutation = useMutation({
  mutationFn: ({ order, transactionDTO }: { order: OrderDetailsDTO, transactionDTO: TransactionDTO }) =>
    ordersService.successOrderPayment(order.id, transactionDTO),
  onSuccess: (_, variables) => {
    router.push({
      name: getRouteName('KIOSK_CART_PAYMENT_SUCCESS'),
      params: { orderId: variables.order.id }
    })
  },
  onError: () => {
    toast({
      title: 'Ошибка подтверждения',
      description: 'Не удалось подтвердить платеж. Обратитесь к персоналу',
      variant: 'destructive'
    })
    resetPaymentFlow(false)
  }
})

// Mutation: Mark payment as failed
const failPaymentMutation = useMutation({
  mutationFn: (orderId: number) => ordersService.failOrderPayment(orderId),
  onSuccess: () => {
    paymentState.errorMessage = 'Платеж не завершен вовремя. Попробуйте снова.'
    resetPaymentFlow(false)
  }
})

// Timer always starts on mount — ignoring any existing order checks
onMounted(() => {
  // Force a fresh timer start each time the page opens
  sessionStorage.setItem(STORAGE_KEY, Date.now().toString())
  initializeTimer(false)
})

// Payment routine for Kaspi
const awaitPayment = async () => {
  if (!order.value) return

  try {
    const kaspiConfig = getKaspiConfig()
    if (!kaspiConfig) throw new Error('Kaspi config not found')

    const kaspiService = new KaspiService(kaspiConfig)
    const transaction = isTest
      ? await kaspiService.awaitPaymentTest()
      : await kaspiService.awaitPayment(order.value.total)

    if (!transaction?.transactionId) {
      throw new Error('Invalid transaction data')
    }

    successOrderPaymentMutation.mutate({
      order: order.value,
      transactionDTO: transaction
    })
  } catch (error) {
    console.error('PAYMENT ERROR: ', error)
    toast({
      title: 'Ошибка оплаты',
      description: 'Не удалось дождаться оплаты. Попробуйте снова.',
      variant: 'destructive'
    })
    resetPaymentFlow(true)
  }
}

// Timer management
const initializeTimer = (forceRestart = false) => {
  clearTimer()

  if (forceRestart) {
    sessionStorage.setItem(STORAGE_KEY, Date.now().toString())
  }

  const storedTime = sessionStorage.getItem(STORAGE_KEY)
  if (!storedTime) return

  const elapsed = Date.now() - Number(storedTime)
  const remaining = PAYMENT_TIMEOUT - elapsed

  if (remaining <= 0) {
    paymentState.isTimeoutReached = true
    failPayment()
    return
  }

  paymentState.countdown = Math.ceil(remaining / 1000)
  paymentState.hasTimerStarted = true

  timerId = window.setInterval(() => {
    paymentState.countdown -= 1
    if (paymentState.countdown <= 0) {
      clearTimer()
      paymentState.isTimeoutReached = true
      failPayment()
    }
  }, 1000)
}

const clearTimer = () => {
  if (timerId) {
    window.clearInterval(timerId)
    timerId = null
  }
}

// Navigation & Payment Flow Handlers
const handleBack = async () => {
  await failPayment()
  resetPaymentFlow(false)
  router.push({ name: getRouteName('KIOSK_HOME') })
}

const selectMethod = (method: PaymentMethod) => {
  if (!order.value) {
    paymentState.errorMessage = 'Номер заказа отсутствует'
    return
  }
  paymentState.selectedMethod = method
  paymentState.errorMessage = ''
  startPaymentProcess()
}

const startPaymentProcess = async () => {
  if (!order.value || !paymentState.selectedMethod) {
    paymentState.errorMessage = 'Номер заказа или способ оплаты отсутствуют'
    return
  }
  // Re-initialize timer from scratch
  initializeTimer(true)
  paymentState.isPaying = true
  paymentState.errorMessage = ''
  await awaitPayment()
}

const failPayment = async () => {
  if (order.value) {
    await failPaymentMutation.mutateAsync(order.value.id)
  }
}

const resetPaymentFlow = (goBack: boolean) => {
  clearTimer()
  sessionStorage.removeItem(STORAGE_KEY)
  paymentState.isPaying = false
  paymentState.hasTimerStarted = false
  paymentState.countdown = PAYMENT_TIMEOUT / 1000
  paymentState.selectedMethod = null
  paymentState.errorMessage = ''
  paymentState.isTimeoutReached = false
  if (goBack) {
    router.push({ name: getRouteName('KIOSK_HOME') })
  }
}

// Computed for UI
const displayCountdown = computed(() => {
  const mm = Math.floor(paymentState.countdown / 60)
  const ss = paymentState.countdown % 60
  return `${mm.toString().padStart(2, '0')}:${ss.toString().padStart(2, '0')}`
})

const isCountdownVisible = computed(() =>
  paymentState.hasTimerStarted && paymentState.isPaying
)

onBeforeUnmount(() => {
  clearTimer()
})
</script>

<template>
	<div class="flex justify-center items-center bg-slate-100 h-screen">
		<PageLoader v-if="isOrderPending" />
		<div
			v-else
			class="p-6 rounded-3xl w-full max-w-3xl"
		>
			<!-- Header -->
			<div class="text-center">
				<h1
					v-if="!paymentState.isPaying"
					class="font-medium text-5xl"
				>
					Выберите способ оплаты
				</h1>
				<h1
					v-else
					class="font-medium text-5xl"
				>
					Идет обработка платежа
				</h1>
			</div>

			<!-- Payment Options -->
			<div
				v-if="!paymentState.isPaying"
				class="flex justify-center gap-4 mt-12"
			>
				<div
					v-for="option in paymentOptions"
					:key="option.id"
					class="flex flex-col flex-1 justify-center items-center px-6 py-10 border rounded-3xl max-w-xs cursor-pointer"
					:class="paymentState.selectedMethod === option.id
            ? 'bg-primary text-white border-primary'
            : 'bg-white border-slate-200'"
					@click="selectMethod(option.id)"
				>
					<component
						:is="option.icon"
						class="size-16"
						stroke-width="1.6"
					/>
					<p class="mt-6 font-medium text-3xl">{{ option.label }}</p>
				</div>
			</div>

			<!-- Payment Processing -->
			<div
				v-else
				class="flex flex-col items-center gap-6 mt-16 text-center"
			>
				<p
					v-if="isCountdownVisible"
					class="font-semibold text-green-600 text-5xl"
				>
					{{ displayCountdown }}
				</p>
				<p class="text-3xl">
					{{
            paymentState.selectedMethod === PaymentMethod.KASPI
              ? 'Откройте приложение Kaspi и просканируйте QR'
              : 'Идет обработка...'
					}}
				</p>
			</div>

			<!-- Error Message -->
			<p
				v-if="paymentState.errorMessage"
				class="mt-12 text-red-500 text-3xl text-center"
			>
				{{ paymentState.errorMessage }}
			</p>

			<!-- Action Buttons -->
			<div class="flex justify-center gap-4 mt-16">
				<Button
					variant="ghost"
					@click="handleBack"
					:disabled="paymentState.isPaying"
					class="hover:bg-slate-200 text-3xl"
				>
					Назад
				</Button>
			</div>
		</div>
	</div>
</template>
