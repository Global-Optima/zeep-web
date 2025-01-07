<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'

// Icons
import { ProductSizeMeasures, ProductSizeNames } from '@/modules/kiosk/products/models/product.model'
import { ChevronLeft } from 'lucide-vue-next'

// DTO
interface CreateProductSizeDTO {
  productId: number
  name: string
  measure: string
  basePrice: number
  size: number
  isDefault?: boolean
  additives?: SelectedAdditiveTypesDTO[]
  ingredientIds?: number[]
}

interface SelectedAdditiveTypesDTO {
  additiveId: number
  isDefault: boolean
}

// Emits
const emits = defineEmits<{
  onSubmit: [dto: CreateProductSizeDTO]
  onCancel: []
}>()

// Router
const router = useRouter()

// Default Additives (Mock Data)
const defaultAdditives = ref([
  { id: 1, name: 'Карамельный сироп', size: '500 мл', imageUrl: 'caramel.jpg' },
  { id: 2, name: 'Ванильный сироп', size: '500 мл', imageUrl: 'vanilla.jpg' },
])

// Technical Card Ingredients (Mock Data)
const technicalCardIngredients = ref([
  { id: 1, imageUrl: 'ingredient1.jpg', name: 'Кофе', unit: 'грамм', weight: 30 },
  { id: 2, imageUrl: 'ingredient2.jpg', name: 'Молоко', unit: 'мл', weight: 100 },
])

// Zod Schema for Validation
const createProductSizeSchema = toTypedSchema(
  z.object({
    productId: z.number().min(1, 'ID продукта обязателен'),
    name: z.enum([ProductSizeNames.S, ProductSizeNames.M, ProductSizeNames.L, ProductSizeNames.XL], {message: 'Выберите корректное название варианта'}),
    size: z.number().min(1, 'Размер должен быть положительным'),
    measure: z.enum([ProductSizeMeasures.G, ProductSizeMeasures.ML, ProductSizeMeasures.PIECE], {message:'Выберите корректную единицу измерения'}),
    basePrice: z.number().min(0, 'Цена должна быть положительным числом'),
    isDefault: z.boolean().optional(),
    additives: z.array(
      z.object({
        additiveId: z.number().min(1),
        isDefault: z.boolean(),
      })
    ).optional(),
    ingredientIds: z.array(z.number()).optional(),
  })
)

// Form Setup
const { handleSubmit, resetForm, isSubmitting } = useForm<CreateProductSizeDTO>({
  validationSchema: createProductSizeSchema,
  initialValues: {
    name: '',
    measure: '',
    basePrice: 0,
    size: 0,
    isDefault: false,
    additives: [],
    ingredientIds: [],
  },
})

// Handlers
const onSubmit = handleSubmit((values) => {
  emits('onSubmit', values)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}
</script>

<template>
	<div class="mx-auto w-full max-w-2xl">
		<!-- Header -->
		<div class="flex justify-between items-center gap-4 w-full">
			<div class="flex items-center gap-4">
				<Button
					variant="outline"
					size="icon"
					@click="onCancel"
				>
					<ChevronLeft class="w-5 h-5" />
					<span class="sr-only">Назад</span>
				</Button>
				<h1
					class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0"
				>
					Детали варианта
				</h1>
			</div>

			<div class="flex items-center gap-2">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					:disabled="isSubmitting"
					@click="onSubmit"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- Main Content -->
		<div class="items-start gap-4 grid lg:col-span-2 auto-rows-max mt-6">
			<!-- Variant Details -->
			<Card>
				<CardHeader>
					<CardTitle>Детали варианта</CardTitle>
					<CardDescription>Укажите название, размер и цену варианта.</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="gap-6 grid">
						<!-- Name Selector -->
						<FormField
							name="name"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Название</FormLabel>
								<FormControl>
									<Select v-bind="componentField">
										<SelectTrigger>
											<SelectValue placeholder="Выберите название" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="S">S</SelectItem>
											<SelectItem value="M">M</SelectItem>
											<SelectItem value="L">L</SelectItem>
											<SelectItem value="XL">XL</SelectItem>
										</SelectContent>
									</Select>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<!-- Measure and Unit -->
						<div class="flex items-start gap-4">
							<FormField
								name="size"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Размер (число)</FormLabel>
									<FormControl>
										<Input
											type="number"
											placeholder="Например, 250"
											v-bind="componentField"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="measure"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Единица измерения</FormLabel>
									<FormControl>
										<Select v-bind="componentField">
											<SelectTrigger>
												<SelectValue placeholder="Выберите единицу" />
											</SelectTrigger>
											<SelectContent>
												<SelectItem value="грамм">Грамм</SelectItem>
												<SelectItem value="мл">Миллилитры</SelectItem>
											</SelectContent>
										</Select>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<!-- Price -->
						<FormField
							name="basePrice"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Начальная цена</FormLabel>
								<FormControl>
									<Input
										type="number"
										placeholder="Введите начальную цену"
										v-bind="componentField"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>
				</CardContent>
			</Card>

			<!-- Default Toppings -->
			<Card>
				<CardHeader>
					<CardTitle>Топпинги по умолчанию</CardTitle>
					<CardDescription>
						Добавьте топпинги, которые идут по умолчанию с этим вариантом.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead class="hidden w-[100px] sm:table-cell"> </TableHead>
								<TableHead>Название</TableHead>
								<TableHead>Размер</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="additive in defaultAdditives"
								:key="additive.id"
							>
								<TableCell class="hidden sm:table-cell">
									<img
										:src="additive.imageUrl"
										alt="Изображение"
										class="bg-gray-100 rounded-md aspect-square object-contain"
										height="64"
										width="64"
									/>
								</TableCell>
								<TableCell class="font-medium">{{ additive.name }}</TableCell>
								<TableCell>{{ additive.size }}</TableCell>
							</TableRow>
						</TableBody>
					</Table>
				</CardContent>
			</Card>
		</div>
	</div>
</template>
