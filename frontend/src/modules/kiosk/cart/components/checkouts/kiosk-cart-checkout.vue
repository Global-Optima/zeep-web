<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import type { LocalizedError } from '@/core/models/errors.model'
import type { CreateOrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useResetKioskState } from '@/modules/kiosk/hooks/use-reset-kiosk.hook'
import { useMutation } from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import { ChevronRight } from 'lucide-vue-next'
import { defineAsyncComponent, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

// Define step state interface
interface StepState {
  customerName: string;
  selectedPayment: string;
  qrCodeUrl: string;
  orderId?: number;
}

// Reactive state for steps
const stepState = reactive<StepState>({
  customerName: '',
  selectedPayment: '',
  qrCodeUrl: '',
  orderId: undefined,
});

// Router and toast initialization
const router = useRouter();
const { toast } = useToast();
const resetKioskState = useResetKioskState();
const { cartItems, clearCart } = useCartStore();

// Mutations for API calls
const createOrderMutation = useMutation({
  mutationFn: (orderDTO: CreateOrderDTO) => ordersService.createOrder(orderDTO),
  onSuccess: () => {
    stepState.qrCodeUrl = "https://cabinet.kofd.kz/consumer";
  },
  onError: (error: AxiosError<LocalizedError>) => {
    toast({ title: "Ошибка", description: error.response?.data.message.ru ?? "Ошибка при создании заказа" });
  },
});

const checkCustomerNameMutation = useMutation({
  mutationFn: (name: string) => ordersService.checkCustomerName({ customerName: name }),
  onError: (error: AxiosError<LocalizedError>) => {
    toast({ title: "Ошибка", description: error.response?.data.message.ru ?? "Ошибка при проверке имени клиента" });
  },
});

// Step configuration interface
interface StepConfig {
  name: string;
  component: ReturnType<typeof defineAsyncComponent>;
  onBeforeProceed?: (state: StepState) => Promise<boolean> | boolean;
  onProceed?: (state: StepState) => Promise<void> | void;
  onBack?: string | null;
}

// Steps configuration
const stepsConfig: StepConfig[] = [
  {
    name: 'customer',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-customer.vue')
    ),
    onBeforeProceed: async (state) => {
      try {
        await checkCustomerNameMutation.mutateAsync(state.customerName);
        return !checkCustomerNameMutation.isError.value; // Proceed only if validation succeeds
      } catch (error) {
        throw new Error(checkCustomerNameMutation.error.value?.response?.data.message.ru || "Ошибка при проверке имени клиента");
      }
    },
    onProceed: (state) => {
      stepState.customerName = state.customerName;
    },
  },
  {
    name: 'payment',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-payment.vue')
    ),
    onProceed: async (state) => {
      stepState.selectedPayment = state.selectedPayment;
      const orderDTO: CreateOrderDTO = {
        customerName: stepState.customerName,
        subOrders: Object.entries(cartItems).map(([_, item]) => ({
          storeProductSizeId: item.size.id,
          quantity: item.quantity,
          storeAdditivesIds: item.additives.map((add) => add.id),
        })),
      };
      await createOrderMutation.mutateAsync(orderDTO);
    },
    onBack: 'customer',
  },
  {
    name: 'confirmation',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-confirm.vue')
    ),
    onProceed: () => {
      resetKioskState.resetAll();
      clearCart();
      router.push({ name: getRouteName('KIOSK_HOME') });
    },
    onBack: 'payment',
  },
];

// Current step reference
const currentStep = ref<string | null>(null);

// Open a specific step
const openStep = (stepName: string) => {
  currentStep.value = stepName;
};

// Close all steps and reset state
const closeAllAndReset = () => {
  currentStep.value = null;
  Object.assign(stepState, {
    customerName: '',
    selectedPayment: '',
    qrCodeUrl: '',
    orderId: undefined,
  });
};

// Handle proceed action
const handleProceed = async (stepName: string, data: StepState) => {
  const step = stepsConfig.find((s) => s.name === stepName);
  if (!step) return;

  try {
    // Execute onBeforeProceed hook
    if (step.onBeforeProceed && !(await step.onBeforeProceed(data))) {
      return;
    }

    // Execute onProceed hook
    if (step.onProceed) {
      await step.onProceed(data);
    }
  } catch (error) {
    // Display specific error message from the error object
    if (error instanceof Error) {
      toast({ title: "Ошибка", description: error.message });
    } else {
      toast({ title: "Ошибка", description: "Произошла ошибка при обработке шага" });
    }
    return;
  }

  // Navigate to the next step
  const nextStepIndex = stepsConfig.findIndex((s) => s.name === stepName) + 1;
  if (stepsConfig[nextStepIndex]) {
    openStep(stepsConfig[nextStepIndex].name);
  } else {
    closeAllAndReset();
  }
};

// Handle back action
const handleBack = (stepName: string) => {
  const step = stepsConfig.find((s) => s.name === stepName);
  if (step?.onBack) {
    openStep(step.onBack);
  } else {
    closeAllAndReset();
  }
};
</script>

<template>
	<button
		class="flex items-center bg-white p-2 sm:p-4 rounded-full text-primary"
		@click="openStep('customer')"
	>
		<ChevronRight class="size-10" />
	</button>
	<component
		v-for="step in stepsConfig"
		:key="step.name"
		:is="step.component"
		:isOpen="currentStep === step.name"
		@close="closeAllAndReset"
		@proceed="(data: StepState) => handleProceed(step.name, data)"
		@back="() => handleBack(step.name)"
		v-bind="stepState"
	/>
</template>
