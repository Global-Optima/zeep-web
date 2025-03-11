<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { useToast } from '@/core/components/ui/toast'
import { useMutation } from '@tanstack/vue-query'
import { computed, onMounted, ref } from 'vue'

// Services
import { KaspiService } from '@/core/integrations/kaspi.service'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'

// Types
import { PaymentMethod, type PaymentOption } from '@/modules/kiosk/cart/models/kiosk-cart.models'
import { CreditCard, QrCode } from 'lucide-vue-next'; // Loader2 is a spinning icon in lucide-vue-next (optional if you prefer another)

const props = defineProps<{
  isOpen: boolean;
  selectedPayment: PaymentMethod | null;
  orderId: number | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;        // Close the entire checkout flow
  (e: 'back'): void;         // Back to previous step (customer info)
  (e: 'update:selectedPayment', val: PaymentMethod | null): void;
  (e: 'proceed'): void;      // Payment success => proceed to next step
}>();

const { toast } = useToast();

// Payment Options
const paymentOptions: PaymentOption[] = [
  { id: PaymentMethod.KASPI, label: 'Kaspi QR', icon: QrCode },
  { id: PaymentMethod.CARD, label: 'Оплата картой', icon: CreditCard },
];

// Reactive State
const selectedMethod = ref<PaymentMethod | null>(props.selectedPayment); // which method user selected
const isPaying = ref(false);         // toggles "awaiting payment" UI
const errorMessage = ref('');
const countdown = ref(120);         // timer in seconds
const hasTimerStarted = ref(false); // for controlling UI of countdown
const isTimeoutReached = ref(false);

// Payment Timer Constants
const PAYMENT_TIMEOUT = 120_000;  // 2 minutes in ms
const STORAGE_KEY = 'paymentStartTime';

// Kaspi mock / real integration
const kaspiService = new KaspiService({
  posIpAddress: '',
  integrationName: '',
});

/* -------------------------------------
   MUTATIONS
   ------------------------------------- */
// 1) AWAIT / Poll Payment
const awaitPaymentMutation = useMutation({
  mutationFn: async ({ method }: {orderId: number, method: PaymentMethod}) => {
    // Simulate or call real polling
    return kaspiService.awaitPayment(method);
  },
  onSuccess: async (transaction, variables) => {
    // Mark order as success on the server
    await ordersService.successOrderPayment(variables.orderId, transaction);
    // Payment successful → proceed to next step
    emit('proceed');
  },
  onError: async () => {
    errorMessage.value = 'Ошибка при обработке платежа. Попробуйте снова.';
    toast({ title: 'Ошибка', description: errorMessage.value, variant: 'destructive' });
    // If the payment fails, handle local reset
    resetPaymentFlow(true);
  },
});

// 2) FAIL Payment (Timeout)
const failPaymentMutation = useMutation({
  mutationFn: async (orderId: number) => {
    return ordersService.failOrderPayment(orderId);
  },
  onSuccess: () => {
    errorMessage.value = 'Платеж не завершен вовремя. Пожалуйста, попробуйте снова.';
    isPaying.value = false;
    resetPaymentFlow(false);
  },
  onError: () => {
    // Even if failing the payment on the server has an error,
    // we still reset the local state to let user try again.
    resetPaymentFlow(false);
  },
});

/* -------------------------------------
   METHODS
   ------------------------------------- */

/**
 * The user picks a payment method by clicking on it.
 * We store the selection, reset error, and begin the process.
 */
function selectMethod(method: PaymentMethod) {
  selectedMethod.value = method;
  errorMessage.value = '';

  emit('update:selectedPayment', method);
  // Begin Payment
  startPaymentProcess();
}

/**
 * Begin Payment Process:
 * 1) Set/Restore the timer
 * 2) Trigger the polling mutation
 */
function startPaymentProcess() {
  if (!props.orderId || !selectedMethod.value) {
    errorMessage.value = 'Номер заказа или способ оплаты отсутствуют.';
    return;
  }
  initializeTimer(true);

  isPaying.value = true;
  errorMessage.value = '';

  console.log({ orderId: props.orderId, method: selectedMethod.value })
  awaitPaymentMutation.mutate({ orderId: props.orderId, method: selectedMethod.value });
}

/**
 * When the countdown hits 0 or a critical error occurs,
 * we mark the payment as failed on the server.
 */
async function failPayment() {
  if (!props.orderId) return;
  await failPaymentMutation.mutateAsync(props.orderId);
}

/**
 * Timer Initialization
 * forceRestart = true => starts a new 120s window
 * forceRestart = false => tries to restore from sessionStorage
 */
function initializeTimer(forceRestart = false) {
  if (forceRestart) {
    sessionStorage.setItem(STORAGE_KEY, Date.now().toString());
  }

  const storedTime = sessionStorage.getItem(STORAGE_KEY);
  if (!storedTime) return;

  hasTimerStarted.value = true;
  isTimeoutReached.value = false; // Reset timeout state

  const elapsed = Date.now() - Number(storedTime);
  const remaining = PAYMENT_TIMEOUT - elapsed;
  countdown.value = Math.ceil(remaining / 1000);

  if (countdown.value <= 0) {
    isTimeoutReached.value = true;
    failPayment();
    return;
  }

  // Start countdown
  const interval = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      clearInterval(interval);
      isTimeoutReached.value = true;
      failPayment();
    }
  }, 1000);
}



/**
 * Reset Payment Flow (local).
 * If we want to let the user pick again, we:
 *  - Clear isPaying
 *  - Clear selectedMethod (optional, if you want them to reselect)
 *  - Clear the timer
 *  - Optionally call failPayment() if the order is truly done
 *  - Step them "back" to the Payment Option UI
 *
 * If `goBack` is true, we also call `emit('back')` to go to the previous step (e.g. Customer).
 * Otherwise, we just remain on the Payment step with the Payment Options visible again.
 */
function resetPaymentFlow(goBack: boolean) {
  // Clear local states
  isPaying.value = false;
  hasTimerStarted.value = false;
  countdown.value = 120;
  sessionStorage.removeItem(STORAGE_KEY);

  // If you'd like the user to explicitly re-select the method:
  selectedMethod.value = null;
  emit('update:selectedPayment', null);

  // If you want to move the user all the way back to customer step,
  // set `goBack = true`. Or you can simply remain on Payment step
  // so the user can choose a payment method again.
  if (goBack) {
    // This triggers your parent to set currentStep = CheckoutStep.CUSTOMER
    emit('back');
  }
}

/* -------------------------------------
   COMPUTED
   ------------------------------------- */
const isCountdownVisible = computed(() => hasTimerStarted.value && isPaying.value);
const displayCountdown = computed(() => {
  const mm = Math.floor(countdown.value / 60);
  const ss = countdown.value % 60;
  return `${mm.toString().padStart(2, '0')}:${ss.toString().padStart(2, '0')}`;
});

/* -------------------------------------
   LIFECYCLE
   ------------------------------------- */
onMounted(() => {
  // If user reloaded mid-payment, restore the timer
  if (!hasTimerStarted.value && props.orderId) {
    initializeTimer(false);
  }
});
</script>

<template>
	<Dialog
		:open="isOpen"
		@update:open="emit('close')"
	>
		<DialogContent
			:include-close-button="false"
			class="space-y-8 bg-white sm:p-12 !rounded-[40px] sm:max-w-3xl text-black"
		>
			<!-- Header -->
			<DialogHeader>
				<!-- Show Title if we haven't started paying yet -->
				<DialogTitle
					v-if="!isPaying"
					class="font-medium text-gray-900 sm:text-4xl text-center"
				>
					Выберите способ оплаты
				</DialogTitle>

				<!-- Show 'Payment in progress' message if paying -->
				<DialogTitle
					v-else
					class="font-medium text-gray-900 sm:text-3xl text-center"
				>
					Идет обработка платежа
				</DialogTitle>
			</DialogHeader>

			<!-- Payment Options (Disable if Timeout Reached) -->
			<div
				v-if="!isPaying"
				class="gap-4 grid grid-cols-2 mb-4"
			>
				<div
					v-for="option in paymentOptions"
					:key="option.id"
					class="flex flex-col justify-center items-center px-6 py-8 border-2 border-transparent rounded-2xl cursor-pointer"
					:class="{
            'bg-primary text-white': selectedMethod === option.id,
            'bg-gray-100': selectedMethod !== option.id,
            'cursor-not-allowed opacity-40': isTimeoutReached, // Disable style
          }"
					@click="!isTimeoutReached && selectMethod(option.id)"
				>
					<component
						:is="option.icon"
						class="mb-4 w-14 h-14"
					/>
					<p class="font-medium sm:text-2xl text-center">{{ option.label }}</p>
				</div>
			</div>

			<!-- Payment Processing Screen -->
			<div
				v-else
				class="flex flex-col items-center gap-6 mb-5 text-center"
			>
				<!-- Countdown + Possibly a spinner -->
				<div v-if="isCountdownVisible">
					<p class="font-bold text-green-600 text-3xl">{{ displayCountdown }}</p>
				</div>

				<!-- Instruction based on selected method -->
				<p
					v-if="selectedMethod === 'card'"
					class="text-2xl"
				>
					Приложите карту к терминалу для оплаты
				</p>
				<p
					v-else-if="selectedMethod === 'kaspi'"
					class="text-2xl"
				>
					Откройте приложение Kaspi и просканируйте QR
				</p>
				<p
					v-else
					class="text-2xl"
				>
					Идет обработка платежа...
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
				<!-- Back to Customer Info Step -->
				<Button
					variant="ghost"
					@click="emit('back')"
					:disabled="isPaying"
					class="sm:px-6 sm:py-8 sm:text-2xl"
				>
					Назад
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>
