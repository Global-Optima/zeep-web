<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed } from 'vue'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import { ChevronLeft } from 'lucide-vue-next'

// Types
import type { ProvisionDetailsDTO } from '@/modules/admin/provisions/models/provision.models'
import type {
  StoreProvisionDetailsDTO,
  StoreProvisionDTO,
  UpdateStoreProvisionDTO
} from '@/modules/admin/store-provisions/models/store-provision.models'

// Props:
// - storeProvision: current store provision (with volume and expirationInMinutes values)
// - provision: the template provision with full details including ingredients.
const props = defineProps<{
  storeProvision: StoreProvisionDetailsDTO
}>()

// Emits: onSubmit and onCancel events.
const emits = defineEmits<{
  onSubmit: [dto: UpdateStoreProvisionDTO]
  onCancel: []
}>()

// Validation Schema for updating the store provision.
// Only volume and expirationInMinutes are updated.
const updateStoreProvisionSchema = toTypedSchema(
  z.object({
    volume: z.number().min(0.0001, 'Введите объем заготовки'),
    expirationInMinutes: z.number().min(0, 'Введите срок годности в минутах').default(60),
  })
)

// Initialize the form with the current values from the storeProvision prop.
const { handleSubmit, resetForm, values } = useForm({
  validationSchema: updateStoreProvisionSchema,
  initialValues: {
    volume: props.storeProvision.volume,
    expirationInMinutes: props.storeProvision.expirationInMinutes,
  }
})

// Compute the technical map scaled to the updated store provision’s volume.
// The scaling ratio is: newVolume / templateAbsoluteVolume.
const scaledTechnicalMap = computed(() => {
  if (!props.storeProvision || !values.volume) return []
  const ratio = values.volume / props.storeProvision.provision.absoluteVolume
  return props.storeProvision.ingredients.map(item => ({
    ingredientId: item.ingredient.id,
    name: item.ingredient.name,
    category: item.ingredient.category,
    quantity: item.quantity,
    scaledQuantity: item.quantity * ratio,
    unit: item.ingredient.unit.name,
  }))
})

// Handle form submission: build and emit the UpdateStoreProvisionDTO.
const onSubmit = handleSubmit((formValues) => {
  const dto: UpdateStoreProvisionDTO = {
    volume: formValues.volume,
    expirationInMinutes: formValues.expirationInMinutes,
  }
  emits('onSubmit', dto)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
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
				Обновить заготовку
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

		<!-- Update Details Form -->
		<Card>
			<CardHeader>
				<CardTitle>Детали заготовки</CardTitle>
				<CardDescription> Укажите требуемый объем и срок годности. </CardDescription>
			</CardHeader>
			<CardContent>
				<div class="gap-6 grid">
					<FormField
						name="volume"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Объем</FormLabel>
							<FormControl>
								<Input
									id="volume"
									type="number"
									step="0.5"
									v-bind="componentField"
									placeholder="Введите объем заготовки"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
					<FormField
						name="expirationInMinutes"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Срок годности (минут)</FormLabel>
							<FormControl>
								<Input
									id="expirationInMinutes"
									type="number"
									v-bind="componentField"
									placeholder="Введите срок годности в минутах"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<!-- Technical Map / Ingredients Table -->
		<Card>
			<CardHeader>
				<CardTitle>Сырье</CardTitle>
				<CardDescription> Технологическая карта на заданный обьем </CardDescription>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Категория</TableHead>
							<TableHead>Абсолютный обьем</TableHead>
							<TableHead>Итоговый обьем</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="ingredient in scaledTechnicalMap"
							:key="ingredient.ingredientId"
						>
							<TableCell>{{ ingredient.name }}</TableCell>
							<TableCell>{{ ingredient.category.name }}</TableCell>
							<TableCell>{{ ingredient.quantity }} {{ ingredient.unit.toLowerCase() }}</TableCell>
							<TableCell class="font-semibold">
								{{ ingredient.scaledQuantity.toFixed(2) }}
								{{ ingredient.unit.toLowerCase() }}
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

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
