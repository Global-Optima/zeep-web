<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
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
import AdminSelectAdditiveDialog from '@/modules/admin/additives/components/admin-select-additive-dialog.vue'
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { ProductSizeMeasures, ProductSizeNames, type ProductSizeDetailsDTO } from '@/modules/kiosk/products/models/product.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'

interface SelectedAdditiveTypesDTO {
  additiveId: number
  isDefault: boolean
  name: string
  categoryName: string
  imageUrl: string
}

export interface UpdateProductSizeFormSchema {
  name: ProductSizeNames
  measure: ProductSizeMeasures
  basePrice: number
  size: number
  additives: SelectedAdditiveTypesDTO[]
}

const {productSize} = defineProps<{productSize: ProductSizeDetailsDTO}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateProductSizeFormSchema]
  onCancel: []
}>()

/**
 * 4. Create Zod schema using the enums
 */
const createProductSizeSchema = toTypedSchema(
  z.object({
    name: z.nativeEnum(ProductSizeNames).describe('Выберите корректный вариант'),
    measure: z.nativeEnum(ProductSizeMeasures).describe('Выберите корректную единицу'),
    basePrice: z.number().min(0, 'Введите корректную цену'),
    size: z.number().min(1, 'Введите корректный размер'),
  })
)

/**
 * 5. Additional validation for additives
 */
const validateAdditives = (additives: SelectedAdditiveTypesDTO[]) => {
  if (!additives.length) {
    return 'Необходимо добавить хотя бы одну добавку.'
  }
  return null
}

/**
 * 6. Form Setup
 */
const { handleSubmit, isSubmitting } = useForm<UpdateProductSizeFormSchema>({
  validationSchema: createProductSizeSchema,
  initialValues: {
    name: productSize.name as ProductSizeNames ,
    measure: productSize.measure as ProductSizeMeasures,
    basePrice: productSize.basePrice,
    size: productSize.size,
    additives: productSize.additives.map(a => ({
      additiveId: a.id,
      isDefault: a.isDefault,
      name: a.name,
      categoryName: a.category.name,
      imageUrl: a.imageUrl
    })),
  }
})

/**
 * 7. Manage Additives
 */
const additives = ref<SelectedAdditiveTypesDTO[]>(productSize.additives.map(a => ({
      additiveId: a.id,
      isDefault: a.isDefault,
      name: a.name,
      categoryName: a.category.name,
      imageUrl: a.imageUrl
    }))
  )
const additivesError = ref<string | null>(null)
const openAdditiveDialog = ref(false)

function addAdditive(additive: AdditiveDTO) {
  if (!additives.value.some((item) => item.additiveId === additive.id)) {
    additives.value.push({
      additiveId: additive.id,
      isDefault: false,
      name: additive.name,
      categoryName: additive.category.name,
      imageUrl: additive.imageUrl,
    })
  }
}

function removeAdditive(index: number) {
  additives.value.splice(index, 1)
}

function toggleDefault(index: number) {
  additives.value[index].isDefault = !additives.value[index].isDefault
}

function sortAdditives() {
  return [...additives.value].sort((a, b) => Number(b.isDefault) - Number(a.isDefault))
}

/**
 * 8. Submit & Cancel
 */
const onSubmit = handleSubmit((formValues) => {
  additivesError.value = validateAdditives(additives.value)
  if (additivesError.value) {
    return
  }
  const finalDTO: UpdateProductSizeFormSchema = {
    ...formValues,
    additives: additives.value,
  }
  emits('onSubmit', finalDTO)
})

const onCancel = () => {
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
		<div class="mt-6">
			<!-- Variant Details -->
			<Card>
				<CardHeader>
					<CardTitle>Детали варианта</CardTitle>
					<CardDescription>Укажите название, размер и цену варианта.</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="gap-6 grid">
						<!-- Name -->
						<FormField
							name="name"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Название</FormLabel>
								<FormControl>
									<!-- This Select uses our VariantName enum -->
									<Select v-bind="componentField">
										<SelectTrigger>
											<SelectValue placeholder="Выберите название" />
										</SelectTrigger>
										<SelectContent>
											<!-- Iterate over VariantName enum -->
											<SelectItem
												v-for="(value, key) in ProductSizeNames"
												:key="key"
												:value="value"
											>
												{{ value }}
											</SelectItem>
										</SelectContent>
									</Select>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<!-- Measure and Size -->
						<div class="flex gap-4">
							<FormField
								name="size"
								v-slot="{ componentField }"
							>
								<FormItem class="flex-1">
									<FormLabel>Размер</FormLabel>
									<FormControl>
										<Input
											type="number"
											v-bind="componentField"
											placeholder="Введите размер"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="measure"
								v-slot="{ componentField }"
							>
								<FormItem class="flex-1">
									<FormLabel>Единица измерения</FormLabel>
									<FormControl>
										<!-- This Select uses our VariantMeasure enum (in Russian) -->
										<Select v-bind="componentField">
											<SelectTrigger>
												<SelectValue placeholder="Выберите единицу" />
											</SelectTrigger>
											<SelectContent>
												<!-- Iterate over VariantMeasure enum -->
												<SelectItem
													v-for="(value, key) in ProductSizeMeasures"
													:key="key"
													:value="value"
												>
													{{ value }}
												</SelectItem>
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
										v-bind="componentField"
										placeholder="Введите цену"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>
				</CardContent>
			</Card>

			<!-- Additives -->
			<Card class="mt-4">
				<CardHeader>
					<div class="flex justify-between items-start">
						<div>
							<CardTitle>Добавки</CardTitle>
							<CardDescription class="mt-2"> Выберите добавки для варианта. </CardDescription>
						</div>
						<Button
							variant="outline"
							@click="openAdditiveDialog = true"
						>
							Добавить
						</Button>
					</div>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead></TableHead>
								<TableHead>Название</TableHead>
								<TableHead>Категория</TableHead>
								<TableHead>По умолчанию</TableHead>
								<TableHead></TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							<TableRow
								v-for="(additive, index) in sortAdditives()"
								:key="additive.additiveId"
							>
								<TableCell>
									<img
										:src="additive.imageUrl"
										class="bg-gray-100 p-1 rounded-md w-16 h-16 object-contain"
									/>
								</TableCell>
								<TableCell>{{ additive.name }}</TableCell>
								<TableCell>{{ additive.categoryName }}</TableCell>
								<TableCell class="text-center">
									<Input
										type="checkbox"
										class="shadow-none h-5"
										:checked="additive.isDefault"
										@change="toggleDefault(index)"
									/>
								</TableCell>
								<TableCell class="text-center">
									<Button
										variant="ghost"
										size="icon"
										@click="removeAdditive(index)"
									>
										<Trash class="w-6 h-6 text-red-500" />
									</Button>
								</TableCell>
							</TableRow>
						</TableBody>
					</Table>
					<div
						v-if="additivesError"
						class="mt-2 text-red-500 text-sm"
					>
						{{ additivesError }}
					</div>
				</CardContent>
			</Card>
		</div>

		<!-- Additive Dialog -->
		<AdminSelectAdditiveDialog
			:open="openAdditiveDialog"
			@close="openAdditiveDialog = false"
			@select="addAdditive"
		/>
	</div>
</template>
