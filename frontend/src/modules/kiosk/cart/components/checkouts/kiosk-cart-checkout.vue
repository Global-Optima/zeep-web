<script setup lang="ts">
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import type { CreateOrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useCartStore } from '@/modules/kiosk/cart/stores/cart.store'
import { useResetKioskState } from '@/modules/kiosk/hooks/use-reset-kiosk.hook'
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
const {toast} = useToast()
const resetKioskState = useResetKioskState();
const { cartItems, clearCart } = useCartStore();

// Mutation for creating the order
const createOrderMutation = useMutation({
  mutationFn: (orderDTO: CreateOrderDTO) => ordersService.createOrder(orderDTO),
  onError: () => {
    toast({description: "Ошибка при создании заказа"});
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
      try {

        const orderDTO: CreateOrderDTO = {
          customerName: stepState.value.customerName,
          subOrders: Object.entries(cartItems).map(([_, item]) => ({
            storeProductSizeId: item.size.id,
            quantity: item.quantity,
            storeAdditivesIds: item.additives.map((add) => add.id),
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
		<ChevronRight class="size-10" />
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
