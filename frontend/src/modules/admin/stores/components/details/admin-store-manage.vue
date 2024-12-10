<template>
	<div class="flex flex-col gap-6 mx-auto w-full md:w-2/3">
		<Card>
			<CardHeader>
				<CardTitle>{{ isEditing ? 'Обновить магазин' : 'Создать магазин' }}</CardTitle>
				<CardDescription>
					Заполните форму ниже, чтобы {{ isEditing ? 'обновить' : 'создать' }} магазин.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					@submit="submitForm"
					class="gap-6 grid"
				>
					<!-- Store Name and Is Franchise -->
					<div class="flex items-end gap-4">
						<div class="flex-grow">
							<FormField
								name="name"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Название магазина</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											placeholder="Введите название магазина"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<div>
							<FormField
								name="isFranchise"
								v-slot="{ field }"
							>
								<FormItem class="flex items-center gap-3 mb-3">
									<Switch
										id="is-franchise"
										v-model="field.value"
									/>
									<Label
										class="!m-0"
										for="is-franchise"
										>Франшиза</Label
									>
								</FormItem>
							</FormField>
						</div>
					</div>
					<!-- Facility Address -->
					<FormField
						name="facilityAddress.address"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Адрес магазина</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="Введите адрес"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Contact Phone and Email -->
					<div class="flex gap-4">
						<div class="w-1/2">
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
						</div>
						<div class="w-1/2">
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
						</div>
					</div>
					<!-- Store Hours -->
					<FormField
						name="storeHours"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Часы работы</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="Введите часы работы (например, 9:00-18:00)"
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
							{{ isEditing ? 'Обновить' : 'Создать' }}
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
import { Label } from '@/core/components/ui/label'
import { Switch } from '@/core/components/ui/switch'
import type { Store } from '@/modules/stores/models/stores.models'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

// Props
const props = defineProps<{
	isEditing: boolean
	initialData: Partial<Store>
}>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: Store): void
	(e: 'onCancel'): void
}>()

// Define Zod schema
const schema = toTypedSchema(
	z.object({
		name: z.string().min(2, 'Название должно содержать минимум 2 символа'),
		isFranchise: z.boolean(),
		facilityAddress: z.object({
			address: z.string().min(5, 'Адрес должен содержать минимум 5 символов'),
			longitude: z.number(),
			latitude: z.number(),
		}),
		contactPhone: z.string().min(7, 'Телефон должен содержать минимум 7 символов'),
		contactEmail: z.string().email('Введите действительный адрес электронной почты'),
		storeHours: z.string().min(5, 'Часы работы должны быть указаны'),
	}),
)

// Initialize form
const { handleSubmit } = useForm<Store>({
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
