<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
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
import type { CreateSupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { ChevronLeft } from 'lucide-vue-next'
import { phoneValidationSchema } from '@/core/validators/phone.validator'

// Props & Events
const emits = defineEmits<{
  (e: 'onSubmit', dto: CreateSupplierDTO): void
  (e: 'onCancel'): void
}>()

// Validation Schema
const createSupplierSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название поставщика'),
    contactEmail: z.string().email('Введите корректный email').min(1, 'Введите почтовый адрес'),
    contactPhone: phoneValidationSchema,
    city: z.string().min(1, 'Введите город'),
    address: z.string().min(1, 'Введите адрес'),
  })
)

// Form Setup
const { handleSubmit, resetForm } = useForm({
  validationSchema: createSupplierSchema,
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
				Создать поставщика
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
				<CardTitle>Информация о поставщике</CardTitle>
				<CardDescription>Заполните основные данные о поставщике.</CardDescription>
			</CardHeader>
			<CardContent>
				<div class="gap-6 grid">
					<!-- Name -->
					<FormField
						name="name"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Название поставщика</FormLabel>
							<FormControl>
								<Input
									id="name"
									type="text"
									v-bind="componentField"
									placeholder="Введите название"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Contact Details -->
					<div class="flex gap-4">
						<FormField
							name="contactEmail"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Email</FormLabel>
								<FormControl>
									<Input
										id="contactEmail"
										type="email"
										v-bind="componentField"
										placeholder="Введите email"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
						<FormField
							name="contactPhone"
							v-slot="{ componentField }"
						>
							<FormItem class="flex-1">
								<FormLabel>Телефон</FormLabel>
								<FormControl>
									<Input
										id="contactPhone"
										type="text"
										v-bind="componentField"
										placeholder="Введите телефон"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<!-- City -->
					<FormField
						name="city"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Город</FormLabel>
							<FormControl>
								<Input
									id="city"
									type="text"
									v-bind="componentField"
									placeholder="Введите город"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Address -->
					<FormField
						name="address"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Адрес</FormLabel>
							<FormControl>
								<Input
									id="address"
									type="text"
									v-bind="componentField"
									placeholder="Введите адрес"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
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
				>Сохранить</Button
			>
		</div>
	</div>
</template>
