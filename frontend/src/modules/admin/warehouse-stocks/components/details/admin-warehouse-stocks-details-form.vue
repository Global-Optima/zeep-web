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

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
					>Отменить</Button
				>
				<Button
					type="submit"
					@click="onSubmit"
					>Сохранить</Button
				>
			</div>
		</div>

		<!-- Stock Material Info -->
		<Card>
			<CardHeader>
				<CardTitle>Информация о материале</CardTitle>
			</CardHeader>
			<CardContent>
				<ul class="space-y-2">
					<li
						v-for="info in materialInfo"
						:key="info.label"
					>
						<span class="font-semibold">{{ info.label }}:</span> {{ info.value }}
					</li>
				</ul>
			</CardContent>
		</Card>

		<!-- Update Form -->
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

					<FormField
						name="expirationDate"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Дата истечения срока годности</FormLabel>
							<FormControl>
								<Input
									type="date"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</form>
			</CardContent>
		</Card>

		<!-- Footer -->
		<div class="flex justify-center items-center gap-2 md:hidden">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
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
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import type {
  UpdateWarehouseStockDTO,
  WarehouseStockDetailsDTO
} from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { toTypedSchema } from '@vee-validate/zod'
import { ChevronLeft } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// Props
const props = defineProps<{
  initialData: WarehouseStockDetailsDTO
}>()

const emit = defineEmits<{
  (e: 'onSubmit', formValues: UpdateWarehouseStockDTO): void
  (e: 'onCancel'): void
}>()

// Predefined Material Info Array
const materialInfo = [
  { label: 'Категория', value: props.initialData.category.name },
  { label: 'Единица измерения', value: props.initialData.unit?.name ?? "TODO" },
  { label: 'Безопасный запас', value: props.initialData.safetyStock },
  {
    label: 'Срок годности',
    value: props.initialData.expirationFlag
      ? `${props.initialData.expirationPeriodInDays} дней`
      : 'Нет',
  },
  { label: 'Штрихкод', value: props.initialData.barcode },
  { label: 'Количество на складе', value: props.initialData.totalQuantity },
  {
    label: 'Ранняя дата истечения срока годности',
    value: formatDate(new Date(props.initialData.earliestExpirationDate)),
  },
]

// Date Formatter Utility
function formatDate(date: Date): string {
  return new Date(date).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

// Zod Schema for Validation
const schema = toTypedSchema(
  z.object({
    quantity: z.coerce
      .number()
      .min(1, 'Количество должно быть не менее 1')
      .refine((value) => Number.isInteger(value), 'Количество должно быть целым числом'),
    expirationDate: z.string().min(1, 'Дата истечения срока годности обязательна'),
  })
)

// Form Handling
const { handleSubmit } = useForm({
  validationSchema: schema,
  initialValues: {
    quantity: props.initialData.totalQuantity,
    expirationDate: props.initialData.earliestExpirationDate.split('T')[0],
  },
})

// Event Handlers
const onSubmit = handleSubmit((formValues) => {
  const dto: UpdateWarehouseStockDTO = {
    quantity: formValues.quantity,
    expirationDate: new Date(formValues.expirationDate),
  }

  emit('onSubmit', dto)
})

const onCancel = () => {
  emit('onCancel')
}
</script>
