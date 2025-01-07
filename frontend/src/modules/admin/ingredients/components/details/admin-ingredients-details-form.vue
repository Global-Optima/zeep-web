<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import type { IngredientsDTO, UpdateIngredientDTO } from '@/modules/admin/ingredients/models/ingredients.model'
import { ChevronLeft } from 'lucide-vue-next'

// Props
const props = defineProps<{
  ingredient: IngredientsDTO
}>()

const emits = defineEmits<{
  onSubmit: [dto: UpdateIngredientDTO]
  onCancel: []
}>()



// Validation Schema
const updateIngredientSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название ингредиента'),
    calories: z.number().min(0, 'Введите корректное значение калорий'),
    fat: z.number().min(0, 'Введите корректное значение жиров'),
    carbs: z.number().min(0, 'Введите корректное значение углеводов'),
    proteins: z.number().min(0, 'Введите корректное значение белков'),
    expiresAt: z.string().optional(),
  })
)

// Form Setup
const { handleSubmit, resetForm } = useForm<UpdateIngredientDTO>({
  validationSchema: updateIngredientSchema,
  initialValues: {
    name: props.ingredient.name || '',
    calories: props.ingredient.calories || 0,
    fat: props.ingredient.fat || 0,
    carbs: props.ingredient.carbs || 0,
    proteins: props.ingredient.proteins || 0,
    expiresAt: props.ingredient.expiresAt || '',
  },
})

// Handlers
const onSubmit = handleSubmit((formValues) => {
  emits('onSubmit', formValues)
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
				{{ ingredient.name }}
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

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Детали ингредиента</CardTitle>
				<CardDescription>Заполните информацию об ингредиенте.</CardDescription>
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
									placeholder="Введите название ингредиента"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Calories, Fat, Carbs, Proteins -->
					<div class="flex gap-4">
						<FormField
							name="calories"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Калории (ккал)</FormLabel>
								<FormControl>
									<Input
										id="calories"
										type="number"
										v-bind="componentField"
										placeholder="Введите калории"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
						<FormField
							name="fat"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Жиры (грамм)</FormLabel>
								<FormControl>
									<Input
										id="fat"
										type="number"
										v-bind="componentField"
										placeholder="Введите жиры"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<div class="flex gap-4">
						<FormField
							name="carbs"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Углеводы (грамм)</FormLabel>
								<FormControl>
									<Input
										id="carbs"
										type="number"
										v-bind="componentField"
										placeholder="Введите углеводы"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
						<FormField
							name="proteins"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Белки (грамм)</FormLabel>
								<FormControl>
									<Input
										id="proteins"
										type="number"
										v-bind="componentField"
										placeholder="Введите белки"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<!-- Expiration Date -->
					<FormField
						name="expiresAt"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Срок годности</FormLabel>
							<FormControl>
								<Input
									id="expiresAt"
									type="date"
									v-bind="componentField"
									placeholder="Введите срок годности"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
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
