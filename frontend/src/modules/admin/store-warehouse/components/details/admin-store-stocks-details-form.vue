<template>
	<div class="flex flex-col gap-6 mx-auto w-full md:w-2/3">
		<Card>
			<CardHeader>
				<CardTitle>Обновить складские запасы</CardTitle>
				<CardDescription>
					Заполните форму ниже, чтобы обновить информацию о складских запасах.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					@submit="submitForm"
					class="gap-6 grid"
				>
					<!-- Quantity -->
					<FormField
						name="quantity"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Количество</FormLabel>
							<FormControl>
								<Input
									type="number"
									v-bind="componentField"
									placeholder="Введите количество"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Low Stock Threshold -->
					<FormField
						name="lowStockThreshold"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Порог малого запаса</FormLabel>
							<FormControl>
								<Input
									type="number"
									v-bind="componentField"
									placeholder="Введите порог малого запаса"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Action Buttons -->
					<div class="flex gap-4 mt-6">
						<Button
							type="submit"
							class="flex-1"
						>
							Обновить
						</Button>
						<Button
							variant="outline"
							class="flex-1"
							@click="handleCancel"
						>
							Отмена
						</Button>
					</div>
				</form>
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

import type { StoreStocks, UpdateStoreStock } from '@/modules/admin/store-warehouse/models/store-stock.model'

// Props
const props = defineProps<{
	initialData: StoreStocks
}>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: UpdateStoreStock): void
	(e: 'onCancel'): void
}>()

// Define Zod schema
const schema = toTypedSchema(
	z.object({
		quantity: z.coerce
			.number()
			.min(1, 'Количество должно быть не менее 1')
			.refine((value) => Number.isInteger(value), 'Количество должно быть целым числом'),
		lowStockThreshold: z.coerce.number()
			.min(0, 'Порог малого запаса не может быть отрицательным')
			.refine((value) => Number.isInteger(value), 'Порог должен быть целым числом'),
	})
)

// Initialize form
const { handleSubmit } = useForm<UpdateStoreStock>({
	validationSchema: schema,
	initialValues: props.initialData,
})

// Submit form
const submitForm = handleSubmit((formValues) => {
	emit('onSubmit', formValues)
})

// Handle cancel
const handleCancel = () => {
	emit('onCancel')
}
</script>
