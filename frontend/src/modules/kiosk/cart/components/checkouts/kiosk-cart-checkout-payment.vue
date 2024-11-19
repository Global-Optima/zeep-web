<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import { CreditCard, Loader, QrCode } from 'lucide-vue-next'
import { ref } from 'vue'

interface Props {
  isOpen: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'back'): void;
  (e: 'proceed', paymentMethod: string): void;
}>();

const selectedPayment = ref('');
const error = ref('');
const isLoading = ref(false);

const paymentOptions = [
  { id: 'kaspi', label: 'Kaspi QR', icon: QrCode },
  { id: 'card', label: 'Оплата картой', icon: CreditCard },
];

// Simulate Payment Process
const processPayment = async () => {
  isLoading.value = true;
  error.value = '';

  try {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    const success = Math.random() > 0.3;

    if (!success) {
      throw new Error('Ошибка при обработке, попробуйте снова!');
    }

    emit('proceed', selectedPayment.value);
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    isLoading.value = false;
  }
};

const proceedToPayment = () => {
  if (!selectedPayment.value) {
    error.value = 'Выберите способ оплаты.';
    return;
  }

  processPayment();
};
</script>

<template>
	<Dialog
		:open="props.isOpen"
		@update:open="emit('close')"
	>
		<DialogContent
			:include-close-button="false"
			class="space-y-8 bg-white shadow-lg mx-auto p-10 rounded-lg max-w-3xl"
		>
			<DialogHeader>
				<DialogTitle class="text-gray-800 sm:text-3xl">Выберите способ оплаты</DialogTitle>
				<p class="mt-3 text-gray-600 sm:text-2xl">
					Пожалуйста, выберите предпочитаемый способ оплаты.
				</p>
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
					@click="selectedPayment = option.id"
				>
					<component
						:is="option.icon"
						class="mb-4 sm:mb-6 w-14 sm:w-20 h-14 sm:h-20"
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
			<DialogFooter class="flex justify-between space-x-4 sm:!mt-14">
				<Button
					variant="ghost"
					@click="emit('back')"
					:disabled="isLoading"
					class="sm:px-6 sm:py-8 sm:text-2xl"
				>
					Назад
				</Button>
				<Button
					@click="proceedToPayment"
					:disabled="isLoading"
					class="bg-primary sm:px-6 sm:py-8 text-white sm:text-2xl"
				>
					{{ isLoading ? 'Ожидание...' : 'Продолжить' }}
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
