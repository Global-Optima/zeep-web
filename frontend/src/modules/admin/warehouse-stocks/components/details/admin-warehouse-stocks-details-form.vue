<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed } from 'vue'
import * as z from 'zod'

// UI Components
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
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { usePrinter } from '@/core/hooks/use-print.hook'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import type {
  UpdateWarehouseStockDTO,
  WarehouseStockMaterialDetailsDTO,
} from '@/modules/admin/warehouse-stocks/models/warehouse-stock.model'
import { ChevronLeft, Printer } from 'lucide-vue-next'

// Props
const props = defineProps<{
  initialData: WarehouseStockMaterialDetailsDTO
  readonly?: boolean
}>()

const emit = defineEmits<{
  (e: 'onSubmit', formValues: UpdateWarehouseStockDTO): void
  (e: 'onCancel'): void
}>()

// Predefined Material Info Array
const materialInfo = computed(() => [
  { label: 'Категория', value: props.initialData.stockMaterial.category.name },
  { label: 'Упаковка', value: `${props.initialData.stockMaterial.size} ${props.initialData.stockMaterial.unit.name}` },
  {
    label: 'Срок годности',
    value: `${props.initialData.stockMaterial.expirationPeriodInDays} дней`
  },
  {
    label: 'Ранняя дата истечения срока годности',
    value: props.initialData.earliestExpirationDate ? formatDate(new Date(props.initialData.earliestExpirationDate)) : "Доставки товара отсутствуют",
  },
  { label: 'Штрихкод', value: props.initialData.stockMaterial.barcode },
])

const {print} = usePrinter()

const onPrintBarcode = async () => {
  if (props.readonly) return
  const blob = await stockMaterialsService.getStockMaterialsBarcodeFile(props.initialData.stockMaterial.id);
  print(blob)
}

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
    expirationDate: z.string().optional()
  })
)

// Form Handling
const { handleSubmit } = useForm({
  validationSchema: schema,
  initialValues: {
    quantity: props.initialData.quantity,
    expirationDate: props.initialData.earliestExpirationDate?.split('T')[0],
  },
})

// Event Handlers
const onSubmit = handleSubmit((formValues) => {
  if (props.readonly) return

  const dto: UpdateWarehouseStockDTO = {
    quantity: formValues.quantity,
    expirationDate: formValues.expirationDate ? new Date(formValues.expirationDate) : undefined,
  }

  emit('onSubmit', dto)
})

const onCancel = () => {
  emit('onCancel')
}
</script>

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
				{{ initialData.stockMaterial.name }}
			</h1>

			<div
				class="md:flex items-center gap-2 hidden md:ml-auto"
				v-if="!readonly"
			>
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
				<div class="flex justify-between items-start gap-4">
					<div>
						<CardTitle>Информация о материале</CardTitle>
						<CardDescription class="mt-1.5">Полная информация о материале</CardDescription>
					</div>
					<Button
						variant="outline"
						@click="onPrintBarcode"
						class="gap-2"
					>
						<Printer class="text-gray-800 size-4" />
						<span>Штрихкод</span>
					</Button>
				</div>
			</CardHeader>
			<CardContent>
				<ul class="space-y-2">
					<li
						v-for="info in materialInfo"
						:key="info.label"
					>
						<span>{{ info.label }}:</span> <span class="font-medium">{{ info.value }}</span>
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
					<FormField name="safetyStock">
						<FormItem>
							<FormLabel>Безопасный запас упаковок</FormLabel>
							<FormControl>
								<Input
									type="number"
									placeholder="Введите безопасный запас упаковок"
									:model-value="initialData.stockMaterial.safetyStock"
									readonly
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

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
			class="flex justify-center items-center gap-2 md:hidden"
			v-if="!readonly"
		>
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
