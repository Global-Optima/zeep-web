<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/core/components/ui/dropdown-menu'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { useToast } from '@/core/components/ui/toast'
import { TechnicalMapEntity, useCopyTechnicalMap } from '@/core/hooks/use-copy-technical-map.hooks'
import AdminIngredientsSelectDialog from "@/modules/admin/ingredients/components/admin-ingredients-select-dialog.vue"
import type { IngredientsDTO } from "@/modules/admin/ingredients/models/ingredients.model"
import type { ProvisionDetailsDTO, UpdateProvisionDTO } from "@/modules/admin/provisions/models/provision.models"
import AdminSelectUnit from '@/modules/admin/units/components/admin-select-unit.vue'
import type { UnitDTO } from '@/modules/admin/units/models/units.model'
import { ChevronDown, ChevronLeft, EllipsisVertical, Trash } from 'lucide-vue-next'

const {provision, readonly = false } = defineProps<{provision: ProvisionDetailsDTO, readonly?: boolean}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateProvisionDTO]
  onCancel: []
}>()

const { toast } = useToast()
const { setTechnicalMapReference, fetchTechnicalMap } = useCopyTechnicalMap()

// Validation Schema for Create Provision
const createProvisionSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название заготовки'),
    absoluteVolume: z.number().min(0, 'Введите абсолютное значение объема'),
    preparationInMinutes: z.number().min(0, 'Введите корректное значение времени приготовления'),
    defaultExpirationInMinutes: z.number().min(0, 'Введите корректное значение времени истечения срока годности'),
    netCost: z.number().min(0, 'Введите корректное значение себестоимости'),
    limitPerDay: z.number().min(0, 'Введите корректное значение лимита по созданию в день'),
    unitId: z.coerce.number().min(1, 'Выберите корректную единицу измерения'),
  })
)

// Form Setup
const { handleSubmit, resetForm, setFieldValue } = useForm({
  validationSchema: createProvisionSchema,
  initialValues: {
    name: provision.name,
    absoluteVolume: provision.absoluteVolume,
    preparationInMinutes: provision.preparationInMinutes,
    defaultExpirationInMinutes: provision.defaultExpirationInMinutes,
    netCost: provision.netCost,
    limitPerDay: provision.limitPerDay,
    unitId: provision.unit.id,
  }
})

// Manage selected ingredients for the provision
interface SelectedIngredientDisplay {
  ingredientId: number
  name: string
  category: string
  unit: string
  quantity: number
}

const selectedIngredients = ref<SelectedIngredientDisplay[]>(provision.ingredients.map(i => ({
  ingredientId: i.ingredient.id,
  name: i.ingredient.name,
  category: i.ingredient.category.name,
  unit: i.ingredient.unit.name,
  quantity: i.quantity,
})))
const openIngredientsDialog = ref(false)

function removeIngredient(index: number) {
  selectedIngredients.value.splice(index, 1)
}

function addIngredient(ingredient: IngredientsDTO) {
  if (readonly) return

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
  if (readonly) return

  const dto: UpdateProvisionDTO = {
    name: formValues.name,
    absoluteVolume: formValues.absoluteVolume,
    preparationInMinutes: formValues.preparationInMinutes,
    defaultExpirationInMinutes: formValues.defaultExpirationInMinutes,
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
const selectedUnit = ref<UnitDTO | null>(provision.unit)

function selectUnit(unit: UnitDTO) {
  if (readonly) return

  selectedUnit.value = unit
  openUnitDialog.value = false
  setFieldValue('unitId', unit.id)
}

const onCopyTechMapClick = () => {
  try {
    setTechnicalMapReference(TechnicalMapEntity.PROVISION, provision.id)
    toast({ description: "Технологическая карта успешно скопирована", variant: "success" })

  } catch {
    toast({ description: "Ошибка при копировании технологической карты", variant: "destructive" })
  }
}

const onPasteTechMapClick = async () => {
  try {
    const techMap = await fetchTechnicalMap()
    if (!techMap) {
      toast({ description: "Технологическая карта не найдена" })
      return
    }

    selectedIngredients.value = techMap.map(t => ({
      ingredientId: t.ingredient.id,
      name: t.ingredient.name,
      unit: t.ingredient.unit.name,
      category: t.ingredient.category.name,
      quantity: t.quantity,
    }))

    toast({ description: "Технологическая карта успешно вставлена", variant: "success" })

  } catch {
    toast({ description: "Ошибка при вставке технологической карты", variant: "destructive" })
  }
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
				Заготовка {{ provision.name }}
			</h1>

			<div
				v-if="!readonly"
				class="hidden md:flex items-center gap-2 md:ml-auto"
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

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Детали заготовки</CardTitle>
				<CardDescription v-if="!readonly">Заполните информацию о заготовке.</CardDescription>
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
									:readonly="readonly"
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
								@click="!readonly && (openUnitDialog = true)"
								class="flex justify-between items-center gap-4 px-4 py-2 border rounded-md text-sm"
								:class="{ 'cursor-pointer': !readonly }"
							>
								{{ selectedUnit?.name || 'Не выбрана' }}
								<ChevronDown
									v-if="!readonly"
									class="w-5 h-5 text-gray-500"
								/>
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
										:readonly="readonly"
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
										:readonly="readonly"
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
							<FormLabel>Время приготовления (минут)</FormLabel>
							<FormControl>
								<Input
									id="preparationInMinutes"
									type="number"
									v-bind="componentField"
									:readonly="readonly"
									placeholder="Введите время приготовления"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

          <!-- Default Expiration Time (Minutes) -->
          <FormField
            name="defaultExpirationInMinutes"
            v-slot="{ componentField }"
          >
            <FormItem>
              <FormLabel>Срок годности по умолчанию (минут)</FormLabel>
              <FormControl>
                <Input
                  id="defaultExpirationInMinutes"
                  type="number"
                  v-bind="componentField"
                  :readonly="readonly"
                  placeholder="Введите срок годности"
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
									:readonly="readonly"
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
						<CardDescription
							class="mt-2"
							v-if="!readonly"
						>
							Выберите ингредиенты и их количество
						</CardDescription>
					</div>
					<div
						v-if="!readonly"
						class="flex items-center gap-2"
					>
						<DropdownMenu>
							<DropdownMenuTrigger class="p-2 border rounded-md">
								<EllipsisVertical class="size-4" />
							</DropdownMenuTrigger>
							<DropdownMenuContent>
								<DropdownMenuItem @click="onCopyTechMapClick">Копировать</DropdownMenuItem>
								<DropdownMenuItem @click="onPasteTechMapClick">Вставить</DropdownMenuItem>
							</DropdownMenuContent>
						</DropdownMenu>

						<Button
							variant="outline"
							@click="openIngredientsDialog = true"
						>
							Добавить
						</Button>
					</div>
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
							<TableHead v-if="!readonly"></TableHead>
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
							<TableCell
								v-if="!readonly"
								class="text-center"
							>
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
		<div
			v-if="!readonly"
			class="md:hidden flex justify-center items-center gap-2"
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
