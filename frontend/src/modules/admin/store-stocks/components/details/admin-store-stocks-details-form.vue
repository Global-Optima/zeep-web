<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				{{ initialData.name }}
			</h1>

			<div
				class="hidden md:flex items-center gap-2 md:ml-auto"
				v-if="!readonly"
			>
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					:disabled="readonly"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					:disabled="readonly"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Обновить складские запасы</CardTitle>
				<CardDescription>
					Заполните форму ниже, чтобы обновить информацию о складских запасах.
				</CardDescription>
			</CardHeader>

			<CardContent>
				<form
					@submit="onSubmit"
					class="gap-6 grid"
				>
					<!-- Quantity -->
					<FormField
						name="quantity"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>
								Количество ({{ initialData.ingredient.unit.name.toLowerCase() }})</FormLabel
							>
							<FormControl>
								<Input
									type="number"
									v-bind="componentField"
									placeholder="Введите количество"
									step="0.5"
									:readonly="readonly"
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
							<FormLabel>
								Порог малого запаса ({{ initialData.ingredient.unit.name.toLowerCase()}})
							</FormLabel>
							<FormControl>
								<Input
									type="number"
									v-bind="componentField"
									placeholder="Введите порог малого запаса"
									:readonly="readonly"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</form>
			</CardContent>
		</Card>

		<!-- Footer -->
		<div
			class="md:hidden flex justify-center items-center gap-2"
			v-if="!readonly"
		>
			<Button
				variant="outline"
				@click="onCancel"
				:disabled="readonly"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
				:disabled="readonly"
				>Сохранить</Button
			>
		</div>
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
import type { StoreWarehouseStockDTO, UpdateStoreWarehouseStockDTO } from '@/modules/admin/store-stocks/models/store-stock.model'
import { toTypedSchema } from '@vee-validate/zod'
import { ChevronLeft } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import * as z from 'zod'
// Props
const props = defineProps<{
	initialData: StoreWarehouseStockDTO,
  readonly?: boolean
}>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: UpdateStoreWarehouseStockDTO): void
	(e: 'onCancel'): void
}>()

// Define Zod schema
const schema = toTypedSchema(
	z.object({
		quantity: z.coerce
			.number()
			.min(0, 'Количество должно быть не менее 0'),
		lowStockThreshold: z.coerce.number()
			.min(0, 'Порог малого запаса не может быть отрицательным')
			.refine((value) => Number.isInteger(value), 'Порог должен быть целым числом'),
	})
)

// Initialize form
const { handleSubmit } = useForm<UpdateStoreWarehouseStockDTO>({
	validationSchema: schema,
	initialValues: props.initialData,
})

// Submit form
const onSubmit = handleSubmit((formValues) => {
	emit('onSubmit', formValues)
})

// Handle cancel
const onCancel = () => {
	emit('onCancel')
}
</script>
