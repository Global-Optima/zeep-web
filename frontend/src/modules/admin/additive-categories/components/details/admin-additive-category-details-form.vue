<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import Switch from '@/core/components/ui/switch/Switch.vue'
import type { AdditiveCategoryDetailsDTO, UpdateAdditiveCategoryDTO } from '@/modules/admin/additives/models/additives.model'
import { ChevronLeft } from 'lucide-vue-next'

// Props
const props = defineProps<{
  category: AdditiveCategoryDetailsDTO
  readonly?: boolean
}>()

// Emits
const emits = defineEmits<{
  onSubmit: [dto: UpdateAdditiveCategoryDTO]
  onCancel: []
}>()

// Validation Schema
const createAdditiveCategorySchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название категории'),
    description: z.string().optional(),
    isMultipleSelect: z.boolean().optional().describe('Можно ли выбирать несколько модификаторов в этой категории'),
  })
)

// Form Setup
const { handleSubmit, resetForm } = useForm<UpdateAdditiveCategoryDTO>({
  validationSchema: createAdditiveCategorySchema,
  initialValues: props.category
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  if (props.readonly) return // Prevent submission in readonly mode
  emits('onSubmit', { ...formValues, isMultipleSelect: formValues.isMultipleSelect ?? false })
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
				{{ category.name }}
			</h1>

			<div
				class="md:flex items-center gap-2 hidden md:ml-auto"
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
				<CardTitle>Детали категории модификатора</CardTitle>
				<CardDescription>Заполните информацию о категории модификатора.</CardDescription>
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
									placeholder="Введите название категории модификатора"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Description -->
					<FormField
						name="description"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Описание</FormLabel>
							<FormControl>
								<Input
									id="description"
									type="text"
									v-bind="componentField"
									:readonly="readonly"
									placeholder="Введите описание категории модификатора"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Is Multiple Select -->
					<FormField
						v-slot="{ value, handleChange }"
						name="isMultipleSelect"
					>
						<FormItem
							class="flex flex-row justify-between items-center gap-12 p-4 border rounded-lg"
						>
							<div class="flex flex-col space-y-0.5">
								<FormLabel class="font-medium text-base">Множественный выбор</FormLabel>
								<FormDescription class="text-sm">
									Укажите можно ли выбрать несколько модификаторов в этой категории при заказе
								</FormDescription>
							</div>

							<FormControl>
								<Switch
									:checked="value"
									@update:checked="handleChange"
									:disabled="readonly"
								/>
							</FormControl>
						</FormItem>
					</FormField>
				</div>
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
