<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import type { RegionDTO, UpdateRegionDTO } from '@/modules/admin/regions/models/regions.model'
import { ChevronLeft } from 'lucide-vue-next'

// Props & Events
const props = defineProps<{
  region: RegionDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  (e: 'onSubmit', dto: UpdateRegionDTO): void
  (e: 'onCancel'): void
}>()


// Validation Schema
const updateRegionSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название региона'),
  })
)

// Form Setup
const { handleSubmit, resetForm } = useForm({
  validationSchema: updateRegionSchema,
  initialValues: { name: props.region.name }
})


// Handlers
const onSubmit = handleSubmit(async (formValues) => {
  if (props.readonly) return
  emits('onSubmit', formValues)
})

const onCancel = () => {
  resetForm()
  emits('onCancel')
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
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
				{{ region.name }}
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

		<!-- Main Content -->
		<Card>
			<CardHeader>
				<CardTitle>Редактирование региона</CardTitle>
				<CardDescription>Измените детали региона.</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					@submit="onSubmit"
					class="gap-6 grid"
				>
					<FormField
						name="name"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Название региона</FormLabel>
							<FormControl>
								<Input
									id="name"
									type="text"
									v-bind="componentField"
									:readonly="readonly"
									placeholder="Введите название региона"
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
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>
	</div>
</template>
