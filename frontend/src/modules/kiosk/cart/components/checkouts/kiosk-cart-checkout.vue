<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import { useResetKioskState } from '@/modules/kiosk/hooks/use-reset-kiosk.hook'
import { ChevronRight } from 'lucide-vue-next'
import { defineAsyncComponent, ref } from 'vue'
import { useRouter } from 'vue-router'

interface StepState {
  customerName: string;
  selectedPayment: string;
  qrCodeUrl: string;
}

const stepState = ref<StepState>({
  customerName: '',
  selectedPayment: '',
  qrCodeUrl: 'https://urbocoffee.kz/',
});


interface StepConfig {
  name: string;
  component: ReturnType<typeof defineAsyncComponent>;
  onProceed?: (data: StepState) => void;
  onBack?: string | null;
}

const router = useRouter();
const resetKioskState = useResetKioskState();

const stepsConfig: StepConfig[] = [
  {
    name: 'customer',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-customer.vue')
    ),
    onProceed: (data: { customerName: string }) => {
      return stepState.value.customerName = data.customerName
    },
    onBack: null,
  },
  {
    name: 'payment',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-payment.vue')
    ),
    onProceed: (data: { selectedPayment: string }) => {
      stepState.value.selectedPayment = data.selectedPayment;
    },
    onBack: 'customer',
  },
  {
    name: 'confirmation',
    component: defineAsyncComponent(() =>
      import('@/modules/kiosk/cart/components/checkouts/kiosk-cart-checkout-confirm.vue')
    ),
    onProceed: () => {
      router.push({ name: getRouteName('KIOSK_LANDING') });
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
    qrCodeUrl: 'https://urbocoffee.kz/',
  });
  resetKioskState.resetAll();
};

const handleProceed = (stepName: string, data: StepState) => {
  const step = stepsConfig.find((s) => s.name === stepName);
  if (step?.onProceed) {
    step.onProceed(data);
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
		@proceed="(data: StepState) => handleProceed(step.name, data)"
		@back="() => handleBack(step.name)"
		v-bind="stepState"
	/>
</template>
