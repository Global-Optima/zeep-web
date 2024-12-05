<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { CreditCard, Loader, QrCode } from 'lucide-vue-next'
import { ref } from 'vue'

interface Props {
  isOpen: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'back'): void;
  (e: 'proceed', data: {paymentMethod: string}): void;
}>();

const selectedPayment = ref<null | string>('');
const error = ref('');
const isLoading = ref(false);

const paymentOptions = [
  { id: 'kaspi', label: 'Kaspi QR', icon: QrCode },
  { id: 'card', label: 'Оплата картой', icon: CreditCard },
];

// Handle Payment Option Selection
const selectPaymentMethod = (method: string) => {
  selectedPayment.value = method;
  processPayment();
};

// Simulate Payment Process
const processPayment = async () => {
  isLoading.value = true;
  error.value = '';

  try {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    const success = Math.random() > 0.3;

    if (!success) {
      selectedPayment.value = null
      throw new Error('Ошибка при обработке, попробуйте снова!');
    }

    if (selectedPayment.value) {
      emit('proceed', {paymentMethod: selectedPayment.value});
    }
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
	<Dialog
		:open="props.isOpen"
		@update:open="emit('close')"
	>
		<DialogContent
			:include-close-button="false"
			class="space-y-8 bg-white sm:p-12 !rounded-[40px] sm:max-w-3xl text-black"
		>
			<DialogHeader>
				<DialogTitle class="font-medium text-center text-gray-900 sm:text-4xl">
					Выберите способ оплаты</DialogTitle
				>
			</DialogHeader>

			<!-- Payment Options -->
			<div
				v-if="!isLoading"
				class="gap-4 sm:gap-6 grid grid-cols-2"
			>
				<div
					v-for="option in paymentOptions"
					:key="option.id"
					class="flex flex-col justify-center items-center p-6 border rounded-lg cursor-pointer"
					:class="{
            'bg-primary text-white': selectedPayment === option.id,
            'border-gray-300': selectedPayment !== option.id,
          }"
					@click="selectPaymentMethod(option.id)"
				>
					<component
						:is="option.icon"
						class="mb-4 sm:mb-6 w-14 sm:w-24 h-14 sm:h-24"
					/>
					<p class="font-medium text-center sm:text-2xl">{{ option.label }}</p>
				</div>
			</div>

			<!-- Loading State -->
			<div
				v-if="isLoading"
				class="flex flex-col justify-center items-center space-y-4"
			>
				<Loader class="w-16 h-16 text-primary animate-spin" />
				<p class="text-2xl text-gray-600">Обработка платежа...</p>
			</div>

			<!-- Error Message -->
			<p
				v-if="error && !isLoading"
				class="mt-2 text-base text-center text-red-500 sm:text-2xl"
			>
				{{ error }}
			</p>

			<!-- Footer Buttons -->
			<div class="flex justify-center items-center !mt-0 w-full">
				<Button
					variant="ghost"
					:disabled="isLoading"
					@click="$emit('back')"
					class="sm:px-6 sm:py-8 sm:text-2xl"
				>
					Назад
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>
