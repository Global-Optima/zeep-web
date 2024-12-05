<script setup lang="ts">
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import QRCodeVue3 from 'qrcode-vue3'

interface Props {
  isOpen: boolean;
  qrCodeUrl: string;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'proceed'): void;
  (e: 'back'): void;
}>();
</script>

<template>
	<Dialog
		:open="props.isOpen"
		@update:open="emit('proceed')"
	>
		<DialogContent
			:include-close-button="false"
			class="space-y-8 bg-white sm:p-12 !rounded-[40px] sm:max-w-3xl text-black"
		>
			<DialogHeader class="flex flex-col justify-center items-center w-full">
				<DialogTitle class="text-center text-primary sm:text-3xl">Заказ подтвержден</DialogTitle>
				<p class="mt-3 text-center text-gray-600 sm:text-2xl">
					Ваш заказ готовится! Пожалуйста, используйте QR-код ниже для получения информации о
					заказе.
				</p>
			</DialogHeader>

			<!-- QR Code -->
			<div class="flex justify-center my-8">
				<QRCodeVue3
					:value="props.qrCodeUrl"
					:width="350"
					:height="350"
					:dotsOptions="{ type: 'dots', color: '#000000' }"
					:backgroundOptions="{ color: '#ffffff' }"
				/>
			</div>

			<!-- Footer -->
			<div class="flex justify-center !mt-8 sm:!mt-12">
				<button
					@click="emit('proceed')"
					class="bg-primary px-4 sm:px-8 py-2 sm:py-4 rounded-lg text-primary-foreground sm:text-2xl"
				>
					Вернуться в меню
				</button>
			</div>
		</DialogContent>
	</Dialog>
</template>
