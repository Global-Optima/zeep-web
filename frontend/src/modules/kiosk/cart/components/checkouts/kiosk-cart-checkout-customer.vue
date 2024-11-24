<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/core/components/ui/dialog'
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

interface Props {
  isOpen: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'back'): void;
  (e: 'proceed', customerName: string): void;
}>();

// Form Schema with Zod
const formSchema = toTypedSchema(
  z.object({
    customerName: z.string()
      .min(2, { message: 'Имя должно содержать не менее 2 символов' })
      .max(50, { message: 'Имя не должно превышать 50 символов' }),
  })
);

// Unique Name Generator
const generateUniqueName = () => {
  const uniqueName = `Покупатель-${Math.floor(Math.random() * 10000)}`;
  form.setFieldValue('customerName', uniqueName);
};

// Form Hook
const form = useForm({
  validationSchema: formSchema,
  initialValues: { customerName: '' },
});

const handleSubmit = form.handleSubmit((values) => {
  emit('proceed', values.customerName);
});
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
					Укажите ваше имя для заказа
				</DialogTitle>
			</DialogHeader>

			<form
				@submit="handleSubmit"
				class="space-y-6 !mt-12"
			>
				<!-- Customer Name Field -->
				<FormField
					name="customerName"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormControl>
							<Input
								type="text"
								placeholder="Введите ваше имя"
								class="border-gray-300 sm:px-6 sm:py-8 border rounded-lg w-full sm:text-2xl"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage class="text-red-500 sm:text-xl" />
					</FormItem>
				</FormField>

				<!-- Generate Unique Name Button -->
				<Button
					type="button"
					variant="link"
					@click="generateUniqueName"
					class="hover:bg-transparent sm:!mt-6 sm:px-6 sm:py-8 w-full text-gray-600 sm:text-2xl"
				>
					Сгенерировать уникальное имя
				</Button>

				<!-- Footer Buttons -->
				<div class="flex justify-between items-center !mt-12 w-full">
					<Button
						variant="ghost"
						@click="$emit('back')"
						class="sm:px-6 sm:py-8 sm:text-2xl"
					>
						Назад
					</Button>
					<Button
						type="submit"
						class="bg-primary sm:px-6 sm:py-8 text-white sm:text-2xl"
					>
						Продолжить
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>

<style scoped>
/* Add any custom styles if necessary */
</style>
