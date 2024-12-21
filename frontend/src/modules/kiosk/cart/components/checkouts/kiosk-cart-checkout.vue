<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import { toastError } from '@/core/config/toast.config'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useResetKioskState } from '@/modules/kiosk/hooks/use-reset-kiosk.hook'
import type { CreateOrderDTO } from '@/modules/orders/models/orders.models'
import { orderService } from '@/modules/orders/services/orders.service'
import { useCurrentStoreStore } from '@/modules/stores/store/current-store.store'
import { useMutation } from '@tanstack/vue-query'
import { ChevronRight } from 'lucide-vue-next'
import { defineAsyncComponent, ref } from 'vue'
import { useRouter } from 'vue-router'

interface StepState {
  customerName: string;
  selectedPayment: string;
  qrCodeUrl: string;
  orderId?: number;
}

const stepState = ref<StepState>({
  customerName: '',
  selectedPayment: '',
  qrCodeUrl: '',
  orderId: undefined,
});

const router = useRouter();
const resetKioskState = useResetKioskState();
const { cartItems, clearCart } = useCartStore();
const {currentStoreId} = useCurrentStoreStore()

// Mutation for creating the order
const createOrderMutation = useMutation({
  mutationFn: (orderDTO: CreateOrderDTO) => orderService.createOrder(orderDTO),
  onSuccess: (response) => {
    stepState.value.orderId = response.orderId;
    stepState.value.qrCodeUrl = `/api/v1/orders/${response.orderId}/receipt`;
  },
  onError: (error) => {
    console.error('Failed to create order:', error);
    toastError('An error occurred while creating the order. Please try again.');
  },
});

interface StepConfig {
  name: string;
  component: ReturnType<typeof defineAsyncComponent>;
  onProceed?: (data: StepState) => Promise<void> | void;
  onBack?: string | null;
}

const stepsConfig: StepConfig[] = [
  {
    name: 'customer',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-customer.vue')
    ),
    onProceed: (data: { customerName: string }) => {
      stepState.value.customerName = data.customerName;
      console.log(stepState.value.customerName)
    },
    onBack: null,
  },
  {
    name: 'payment',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-payment.vue')
    ),
    onProceed: async (data: { selectedPayment: string }) => {
      stepState.value.selectedPayment = data.selectedPayment;
      console.log("HERE", stepState.value.customerName)
      try {
        if (!currentStoreId) return

        const orderDTO: CreateOrderDTO = {
          customerName: stepState.value.customerName,
          storeId: currentStoreId,
          subOrders: Object.entries(cartItems).map(([_, item]) => ({
            productSizeId: item.size.id,
            quantity: item.quantity,
            additivesIds: item.additives.map((add) => add.id),
          })),
        };

        await createOrderMutation.mutateAsync(orderDTO);
      } catch (error) {
        throw error;
      }
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

const currentStep = ref<string | null>(null);

const openStep = (stepName: string) => {
  currentStep.value = stepName;
};

const closeAllAndReset = () => {
  currentStep.value = null;
  Object.assign(stepState.value, {
    customerName: '',
    selectedPayment: '',
    qrCodeUrl: '',
    orderId: undefined,
  });
};

const handleProceed = async (stepName: string, data: StepState) => {
  const step = stepsConfig.find((s) => s.name === stepName);
  if (step?.onProceed) {
    try {
      await step.onProceed(data);
    } catch {
      return;
    }
  }
  const nextStepIndex = stepsConfig.findIndex((s) => s.name === stepName) + 1;
  if (stepsConfig[nextStepIndex]) {
    openStep(stepsConfig[nextStepIndex].name);
  } else {
    closeAllAndReset();
  }
};

const handleBack = (stepName: string) => {
  const step = stepsConfig.find((s) => s.name === stepName);
  if (step?.onBack) openStep(step.onBack);
};
</script>

<template>
	<button
		class="flex items-center bg-white p-2 sm:p-4 rounded-full text-primary"
		@click="openStep('customer')"
	>
		<ChevronRight class="w-9 h-9" />
	</button>

	<component
		v-for="step in stepsConfig"
		:key="step.name"
		:is="step.component"
		:isOpen="currentStep === step.name"
		@close="closeAllAndReset"
		@proceed="(data: any) => handleProceed(step.name, data)"
		@back="() => handleBack(step.name)"
		v-bind="stepState"
	/>
</template>
