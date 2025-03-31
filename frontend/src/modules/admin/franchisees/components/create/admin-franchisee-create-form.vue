<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Textarea } from '@/core/components/ui/textarea'
import { ChevronLeft } from 'lucide-vue-next'

// API Service & DTOs
import type { CreateFranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'

// Emits
const emits = defineEmits<{
  (e: 'onSubmit', dto: CreateFranchiseeDTO): void
  (e: 'onCancel'): void
}>()

// Validation Schema
const createFranchiseeSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название франчайзи'),
    description: z.string().optional(),
  })
)

// Form Setup
const { handleSubmit, resetForm } = useForm({
  validationSchema: createFranchiseeSchema,
})

// Handlers
const onSubmit = handleSubmit(async (formValues) => {
  emits('onSubmit', formValues)
  resetForm()
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
				Создать франчайзи
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
				<CardTitle>Детали франчайзи</CardTitle>
				<CardDescription>Введите название и описание нового франчайзи.</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					@submit="onSubmit"
					class="gap-6 grid"
				>
					<!-- Name -->
					<FormField
						name="name"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Название франчайзи</FormLabel>
							<FormControl>
								<Input
									id="name"
									type="text"
									v-bind="componentField"
									placeholder="Введите название франчайзи"
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
								<Textarea
									id="description"
									v-bind="componentField"
									placeholder="Введите описание"
									class="min-h-32"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</form>
			</CardContent>
		</Card>

		<!-- Footer -->
		<div class="md:hidden flex justify-center items-center gap-2">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>
	</div>
</template>
