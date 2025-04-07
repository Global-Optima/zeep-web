<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import AdminIngredientsSelectDialog from "@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue"
import type { CreateProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import AdminSelectUnit from '@/modules/admin/units/components/admin-select-unit.vue'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ChevronDown, ChevronLeft, Trash } from 'lucide-vue-next'
import type { IngredientsDTO } from "@/modules/admin/ingredients/models/ingredients.model"


const emits = defineEmits<{
  onSubmit: [dto: CreateProvisionDTO]
  onCancel: []
}>()

// Validation Schema for Create Provision
const createProvisionSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название заготовки'),
    absoluteVolume: z.number().min(0, 'Введите абсолютное значение объема'),
    preparationInMinutes: z.number().min(0, 'Введите корректное значение времени подготовки'),
    netCost: z.number().min(0, 'Введите корректное значение себестоимости'),
    limitPerDay: z.number().min(0, 'Введите корректное значение лимита по созданию в день'),
    unitId: z.coerce.number().min(1, 'Выберите корректную единицу измерения'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: createProvisionSchema,
})

// Manage selected ingredients for the provision
interface SelectedIngredientDisplay {
  ingredientId: number
  name: string
  category: string
  unit: string
  quantity: number
}

const selectedIngredients = ref<SelectedIngredientDisplay[]>([])
const openIngredientsDialog = ref(false)

function removeIngredient(index: number) {
  selectedIngredients.value.splice(index, 1)
}

// This dummy function simulates adding an ingredient from a dialog.
// Replace with your actual component/dialog implementation.
function addIngredient(ingredient: IngredientsDTO) {
  if (!selectedIngredients.value.some((item) => item.ingredientId === ingredient.id)) {
    selectedIngredients.value.push({
      ingredientId: ingredient.id,
      name: ingredient.name,
      unit: ingredient.unit.name,
      category: ingredient.category.name,
      quantity: 0
    })
  }
    openIngredientsDialog.value = false
}

// Handlers
const onSubmit = handleSubmit((formValues) => {
  const dto: CreateProvisionDTO = {
    name: formValues.name,
    absoluteVolume: formValues.absoluteVolume,
    preparationInMinutes: formValues.preparationInMinutes,
    netCost: formValues.netCost,
    limitPerDay: formValues.limitPerDay,
    unitId: formValues.unitId,
    ingredients: selectedIngredients.value.map(i => ({
      ingredientId: i.ingredientId,
      quantity: i.quantity,
    })),
  }
  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  selectedIngredients.value = []
  emits('onCancel')
}

// Unit selection
const openUnitDialog = ref(false)
const selectedUnit = ref<UnitDTO | null>(null)

function selectUnit(unit: UnitDTO) {
  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
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
				Создать заготовку
			</h1>
			<div class="hidden md:flex items-center gap-2 md:ml-auto">
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

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Детали заготовки</CardTitle>
				<CardDescription>Заполните информацию о заготовке.</CardDescription>
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
								<Input
									id="name"
									type="text"
									v-bind="componentField"
									placeholder="Введите название заготовки"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="unitId">
						<FormItem>
							<FormLabel>Единица измерения</FormLabel>
							<div
								@click="openUnitDialog = true"
								class="flex justify-between items-center gap-4 px-4 py-2 border rounded-md text-sm"
							>
								{{ selectedUnit?.name || 'Не выбрана' }}

								<ChevronDown class="w-5 h-5 text-gray-500" />
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<div class="flex items-start gap-4 w-full">
						<!-- Absolute Volume -->
						<FormField
							name="absoluteVolume"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>
									Обьем ({{ selectedUnit?.name.toLowerCase() || 'Не выбрана' }})
								</FormLabel>
								<FormControl>
									<Input
										id="absoluteVolume"
										type="number"
										v-bind="componentField"
										placeholder="Введите объем"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<!-- Net Cost -->
						<FormField
							name="netCost"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Себестоимость</FormLabel>
								<FormControl>
									<Input
										id="netCost"
										type="number"
										v-bind="componentField"
										placeholder="Введите себестоимость"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<!-- Preparation Time (Minutes) -->
					<FormField
						name="preparationInMinutes"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Время подготовки (минут)</FormLabel>
							<FormControl>
								<Input
									id="preparationInMinutes"
									type="number"
									v-bind="componentField"
									placeholder="Введите время подготовки"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Limit Per Day -->
					<FormField
						name="limitPerDay"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Лимит в день</FormLabel>
							<FormControl>
								<Input
									id="limitPerDay"
									type="number"
									v-bind="componentField"
									placeholder="Введите лимит"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<!-- Ingredients List Table -->
		<Card class="mt-2">
			<CardHeader>
				<div class="flex justify-between items-start">
					<div>
						<CardTitle>Технологическая карта</CardTitle>
						<CardDescription class="mt-2"> Выберите ингредиенты и их количество </CardDescription>
					</div>
					<Button
						variant="outline"
						@click="openIngredientsDialog = true"
					>
						Добавить
					</Button>
				</div>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Категория</TableHead>
							<TableHead>Количество</TableHead>
							<TableHead>Единица</TableHead>
							<TableHead></TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="(ingredient, index) in selectedIngredients"
							:key="ingredient.ingredientId"
						>
							<TableCell>{{ ingredient.name }}</TableCell>
							<TableCell>{{ ingredient.category }}</TableCell>
							<TableCell class="flex items-center gap-2">
								<Input
									type="number"
									v-model.number="ingredient.quantity"
									min="0"
									placeholder="Введите количество"
									class="w-16"
								/>
							</TableCell>
							<TableCell>{{ ingredient.unit }}</TableCell>
							<TableCell class="text-center">
								<Button
									variant="ghost"
									size="icon"
									@click="removeIngredient(index)"
								>
									<Trash class="w-6 h-6 text-red-500" />
								</Button>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<!-- Dialogs -->
		<AdminSelectUnit
			:open="openUnitDialog"
			@close="openUnitDialog = false"
			@select="selectUnit"
		/>

		<AdminIngredientsSelectDialog
			:open="openIngredientsDialog"
			@close="openIngredientsDialog = false"
			@select="addIngredient"
		/>
		<!-- Mobile Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
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
