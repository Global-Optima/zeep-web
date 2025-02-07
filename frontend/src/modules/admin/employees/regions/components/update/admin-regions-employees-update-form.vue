<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="handleCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Обновить {{ getEmployeeShortName(employee) }}
			</h1>

			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="handleCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					@click="submitForm"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Обновить сотрудника</CardTitle>
				<CardDescription> Заполните форму ниже, чтобы обновить данные сотрудника. </CardDescription>
			</CardHeader>
			<CardContent>
				<form
					@submit="submitForm"
					class="gap-6 grid"
				>
					<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
						<FormField
							name="firstName"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Имя</FormLabel>
								<FormControl>
									<Input
										v-bind="componentField"
										placeholder="Введите имя"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
						<FormField
							name="lastName"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Фамилия</FormLabel>
								<FormControl>
									<Input
										v-bind="componentField"
										placeholder="Введите фамилию"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
						<FormField
							name="email"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Электронная почта</FormLabel>
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

						<FormField
							name="phone"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Телефон</FormLabel>
								<FormControl>
									<Input
										type="tel"
										v-bind="componentField"
										placeholder="+7 (___) ___-__-__"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<FormField
						name="role"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Должность</FormLabel>
							<FormControl>
								<Select v-bind="componentField">
									<SelectTrigger id="role">
										<SelectValue placeholder="Выберите должность" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem
											v-for="role in roles"
											:key="role.value"
											:value="role.value"
										>
											{{ role.label }}
										</SelectItem>
									</SelectContent>
								</Select>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-slot="{ value, handleChange }"
						name="isActive"
					>
						<FormItem
							class="flex flex-row justify-between items-center gap-12 p-4 border rounded-lg"
						>
							<div class="flex flex-col space-y-0.5">
								<FormLabel class="font-medium text-base"> Активировать сотрудника </FormLabel>
								<FormDescription class="text-sm">
									Вы можете деактивировать его учетную запись, если он больше не работает.
								</FormDescription>
							</div>

							<FormControl>
								<Switch
									:checked="value"
									@update:checked="handleChange"
								/>
							</FormControl>
						</FormItem>
					</FormField>
				</form>
			</CardContent>
		</Card>

		<!-- Footer -->
		<div class="flex justify-center items-center gap-2 md:hidden">
			<Button
				variant="outline"
				@click="handleCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="submitForm"
			>
				Сохранить
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

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
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import { Switch } from '@/core/components/ui/switch'
import { getEmployeeShortName } from '@/core/utils/user-formatting.utils'
import { EmployeeRole, type EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import type { UpdateRegionEmployeeDTO } from '@/modules/admin/employees/regions/models/region-employees.model'
import { ChevronLeft } from 'lucide-vue-next'

const {employee} = defineProps<{
	employee: EmployeeDTO
}>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: UpdateRegionEmployeeDTO): void
	(e: 'onCancel'): void
}>()

const roles = ref([
	{ value: EmployeeRole.REGION_WAREHOUSE_MANAGER, label: 'Менеджер' },
]);

const schema = toTypedSchema(
	z.object({
		firstName: z.string().min(2, 'Имя должно содержать минимум 2 символа').max(50, 'Имя должно содержать не более 50 символов'),
		lastName: z.string().min(2, 'Фамилия должна содержать минимум 2 символа').max(50, 'Фамилия должна содержать не более 50 символов'),
		email: z.string().email('Введите действительный адрес электронной почты'),
		phone: z.string().min(7, 'Телефон должен содержать минимум 7 символов').max(15, 'Телефон должен содержать не более 15 символов'),
		role: z.nativeEnum(EmployeeRole),
    isActive: z.boolean(),
	})
);

const { handleSubmit } = useForm({
	validationSchema: schema,
	initialValues: employee,
})

const submitForm = handleSubmit((formValues) => {
	emit('onSubmit', formValues)
})

const handleCancel = () => {
	emit('onCancel')
}
</script>
