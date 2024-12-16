<template>
	<div class="flex flex-col gap-6 mx-auto w-full md:w-2/3">
		<Card>
			<CardHeader>
				<CardTitle>Обновить Поставщика</CardTitle>
				<CardDescription>Заполните форму ниже, чтобы обновить поставщика.</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					@submit="submitForm"
					class="gap-6 grid"
				>
					<!-- Supplier Name -->
					<FormField
						name="name"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Название поставщика</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="Введите название поставщика"
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
							<FormLabel>Адрес поставщика</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="Введите адрес поставщика"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Contact Phone -->
					<FormField
						name="contactPhone"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Контактный телефон</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="+7 (___) ___-__-__"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Contact Email -->
					<FormField
						name="contactEmail"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Контактный Email</FormLabel>
							<FormControl>
								<Input
									type="email"
									v-bind="componentField"
									placeholder="example@example.com"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Action Buttons -->
					<div class="flex gap-4 mt-6">
						<Button
							type="submit"
							class="flex-1"
						>
							Обновить
						</Button>
						<Button
							variant="outline"
							class="flex-1"
							@click="handleCancel"
						>
							Отмена
						</Button>
					</div>
				</form>
			</CardContent>
		</Card>
	</div>
</template>

<script setup lang="ts">
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
  FormMessage
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import type { Suppliers, UpdateSupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// Props
const props = defineProps<{
	initialData: Suppliers
}>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: UpdateSupplierDTO): void
	(e: 'onCancel'): void
}>()

// Define Zod schema
const schema = toTypedSchema(
	z.object({
		name: z.string().min(2, 'Название должно содержать минимум 2 символа'),
		contactEmail: z.string().email('Введите действительный адрес электронной почты'),
		contactPhone: z.string().min(7, 'Телефон должен содержать минимум 7 символов'),
		address: z.string().min(5, 'Адрес должен содержать минимум 5 символов')
	})
)

// Initialize form
const { handleSubmit } = useForm<UpdateSupplierDTO>({
	validationSchema: schema,
	initialValues: props.initialData,
})

// Submit form
const submitForm = handleSubmit((formValues) => {
	emit('onSubmit', formValues)
})

// Handle cancel
const handleCancel = () => {
	emit('onCancel')
}
</script>
