<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { CreditCard, Loader, QrCode } from 'lucide-vue-next'
import { computed, defineEmits, defineProps, ref } from 'vue'

// Props & Emits
defineProps<{
  isOpen: boolean;
  selectedPayment: string;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'back'): void;
  (e: 'update:selectedPayment', value: string): void;
  (e: 'proceed', data: { paymentMethod: string }): void;
}>();

// Payment Methods
const paymentOptions = [
  { id: 'kaspi', label: 'Kaspi QR', icon: QrCode },
  { id: 'card', label: 'Оплата картой', icon: CreditCard },
];

// State
const isLoading = ref(false);
const errorMessage = ref('');

// Computed Properties
const hasError = computed(() => !!errorMessage.value);

// Methods
const selectPaymentMethod = async (method: string) => {
  emit('update:selectedPayment', method);
  await processPayment(method);
};

const processPayment = async (method: string) => {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    await new Promise(resolve => setTimeout(resolve, 2000));
    const success = Math.random() > 0.3;

    if (!success) {
      emit('update:selectedPayment', '');
      throw new Error('Ошибка при обработке, попробуйте снова!');
    }

    emit('proceed', { paymentMethod: method });
  } catch (err) {
    errorMessage.value = (err as Error).message;
  } finally {
    isLoading.value = false;
  }
};
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
				<DialogTitle class="font-medium text-gray-900 sm:text-4xl text-center">
					Выберите способ оплаты
				</DialogTitle>
			</DialogHeader>

			<!-- Payment Options -->
			<div
				v-if="!isLoading"
				class="gap-4 sm:gap-6 grid grid-cols-2"
			>
				<div
					v-for="option in paymentOptions"
					:key="option.id"
					class="flex flex-col justify-center items-center px-6 py-10 rounded-lg cursor-pointer"
					:class="{
            'bg-primary text-white': selectedPayment === option.id,
            'bg-gray-100': selectedPayment !== option.id,
          }"
					@click="selectPaymentMethod(option.id)"
				>
					<component
						:is="option.icon"
						class="mb-4 sm:mb-6 w-14 sm:w-24 h-14 sm:h-24"
					/>
					<p class="font-medium sm:text-2xl text-center">{{ option.label }}</p>
				</div>
			</div>

			<!-- Loading State -->
			<div
				v-if="isLoading"
				class="flex flex-col justify-center items-center space-y-4"
			>
				<Loader class="w-16 h-16 text-primary animate-spin" />
				<p class="text-gray-600 text-2xl">Обработка платежа...</p>
			</div>

			<!-- Error Message -->
			<p
				v-if="hasError && !isLoading"
				class="mt-2 text-red-500 text-base sm:text-2xl text-center"
			>
				{{ errorMessage }}
			</p>

			<!-- Footer Buttons -->
			<div class="flex justify-center items-center mt-4 w-full">
				<Button
					variant="ghost"
					:disabled="isLoading"
					@click="emit('back')"
					class="sm:px-6 sm:py-8 sm:text-2xl"
				>
					Назад
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>
