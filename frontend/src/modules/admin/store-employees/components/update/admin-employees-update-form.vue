<template>
	<div class="flex md:flex-row flex-col gap-6">
		<!-- Left Side: Employee Details Form -->
		<div class="w-full md:w-2/3">
			<Card>
				<CardHeader>
					<CardTitle>Обновить сотрудника</CardTitle>
					<CardDescription>
						Заполните форму ниже, чтобы обновить данные сотрудника.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form
						@submit="submitForm"
						class="gap-6 grid"
					>
						<!-- Name and Surname -->
						<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
							<!-- Name -->
							<FormField
								name="name"
								v-slot="{ field, errorMessage }"
							>
								<FormItem>
									<FormLabel>Имя</FormLabel>
									<FormControl>
										<Input
											v-model="field.value"
											placeholder="Введите имя"
										/>
									</FormControl>
									<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
								</FormItem>
							</FormField>
							<!-- Surname -->
							<FormField
								name="surname"
								v-slot="{ field, errorMessage }"
							>
								<FormItem>
									<FormLabel>Фамилия</FormLabel>
									<FormControl>
										<Input
											v-model="field.value"
											placeholder="Введите фамилию"
										/>
									</FormControl>
									<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
								</FormItem>
							</FormField>
						</div>
						<!-- Email and Phone -->
						<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
							<!-- Email -->
							<FormField
								name="email"
								v-slot="{ field, errorMessage }"
							>
								<FormItem>
									<FormLabel>Электронная почта</FormLabel>
									<FormControl>
										<Input
											type="email"
											v-model="field.value"
											placeholder="example@example.com"
										/>
									</FormControl>
									<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
								</FormItem>
							</FormField>
							<!-- Phone -->
							<FormField
								name="phone"
								v-slot="{ field, errorMessage }"
							>
								<FormItem>
									<FormLabel>Телефон</FormLabel>
									<FormControl>
										<Input
											v-model="field.value"
											placeholder="+7 (___) ___-__-__"
										/>
									</FormControl>
									<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
								</FormItem>
							</FormField>
						</div>

						<!-- Roles Selector -->
						<FormField
							name="role"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Должность</FormLabel>
								<FormControl>
									<Select v-model="field.value">
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
								<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
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

		<!-- Right Side: Image Upload and Working Hours -->
		<div class="flex flex-col gap-6 w-full md:w-1/3">
			<!-- Image Upload -->
			<Card>
				<CardHeader>
					<CardTitle>Изображение профиля</CardTitle>
					<CardDescription>Загрузите изображение (необязательно)</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminEmployeesCreateImage v-model="employee.image" />
				</CardContent>
			</Card>
			<!-- Working Hours -->
			<Card>
				<CardHeader>
					<CardTitle>Рабочие часы</CardTitle>
					<CardDescription>Установите рабочие часы для каждого дня</CardDescription>
				</CardHeader>
				<CardContent>
					<AdminEmployeesCreateWorkHours v-model="employee.workingHours" />
				</CardContent>
			</Card>
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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import AdminEmployeesCreateImage from '@/modules/admin/employees/components/create/admin-employees-create-image.vue'
import AdminEmployeesCreateWorkHours from '@/modules/admin/employees/components/create/admin-employees-create-work-hours.vue'
import { EmployeeRole, type Employee, type UpdateEmployeeDto } from '@/modules/admin/employees/models/employees.models'

// Props
const props = defineProps<{
	initialData: Employee
}>()

const emit = defineEmits<{
	(e: 'onSubmit', formValues: UpdateEmployeeDto): void
	(e: 'onCancel'): void
}>()

// Roles
const roles = ref([
	{ value: EmployeeRole.MANAGER, label: 'Менеджер' },
	{ value: EmployeeRole.BARISTA, label: 'Бариста' },
]);

// Define the Zod schema for form validation
const schema = toTypedSchema(
	z.object({
		name: z.string().min(2, 'Имя должно содержать минимум 2 символа').max(50, 'Имя должно содержать не более 50 символов'),
		surname: z.string().min(2, 'Фамилия должна содержать минимум 2 символа').max(50, 'Фамилия должна содержать не более 50 символов'),
		email: z.string().email('Введите действительный адрес электронной почты'),
		phone: z.string().min(7, 'Телефон должен содержать минимум 7 символов').max(15, 'Телефон должен содержать не более 15 символов'),
		role: z.string().min(1, 'Выберите должность'),
	})
);

// Initialize the form with initial data
const { handleSubmit } = useForm<UpdateEmployeeDto>({
	validationSchema: schema,
	initialValues: props.initialData,
})

// Employee image and working hours state
const employee = ref({
	image: null,
	workingHours: {},
})

const submitForm = handleSubmit((formValues) => {
	emit('onSubmit', formValues)
})

const handleCancel = () => {
	emit('onCancel')
}
</script>
