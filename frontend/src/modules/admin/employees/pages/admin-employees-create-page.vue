<template>
	<div class="flex md:flex-row flex-col gap-6">
		<!-- Left Side: Employee Details Form -->
		<div class="w-full md:w-2/3">
			<Card>
				<CardHeader>
					<CardTitle>{{ isEditing ? 'Обновить сотрудника' : 'Создать сотрудника' }}</CardTitle>
					<CardDescription>
						Заполните форму ниже, чтобы {{ isEditing ? 'обновить' : 'создать' }} сотрудника.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form
						@submit.prevent="onSubmit"
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
						<!-- Password -->
						<FormField
							name="password"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Пароль</FormLabel>
								<FormControl>
									<Input
										type="password"
										v-model="field.value"
										placeholder="Введите пароль"
										:disabled="isEditing"
									/>
								</FormControl>
								<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
							</FormItem>
						</FormField>
						<!-- Store Selector -->
						<FormField
							name="storeId"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Магазин</FormLabel>
								<FormControl>
									<Select v-model="field.value">
										<SelectTrigger id="store">
											<SelectValue placeholder="Выберите магазин" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem
												v-for="store in stores"
												:key="store.id"
												:value="store.id.toString()"
											>
												{{ store.name }}
											</SelectItem>
										</SelectContent>
									</Select>
								</FormControl>
								<FormMessage v-if="errorMessage">{{ errorMessage }}</FormMessage>
							</FormItem>
						</FormField>
						<!-- Position Selector -->
						<FormField
							name="position"
							v-slot="{ field, errorMessage }"
						>
							<FormItem>
								<FormLabel>Должность</FormLabel>
								<FormControl>
									<Select v-model="field.value">
										<SelectTrigger id="position">
											<SelectValue placeholder="Выберите должность" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem
												v-for="position in positions"
												:key="position.value"
												:value="position.value"
											>
												{{ position.label }}
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
import { useRouter } from 'vue-router'
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

const router = useRouter();
const isEditing = ref(false); // Set to true if editing an existing employee

// Mock data for stores and positions
const stores = ref([
  { id: 1, name: 'Магазин А' },
  { id: 2, name: 'Магазин Б' },
  { id: 3, name: 'Магазин В' },
]);

const positions = ref([
  { value: 'manager', label: 'Менеджер' },
  { value: 'cashier', label: 'Кассир' },
  { value: 'barista', label: 'Бариста' },
]);

// Define the Zod schema for form validation
const schema = toTypedSchema(
  z.object({
    name: z.string().min(2, 'Имя должно содержать минимум 2 символа').max(50, 'Имя должно содержать не более 50 символов'),
    surname: z.string().min(2, 'Фамилия должна содержать минимум 2 символа').max(50, 'Фамилия должна содержать не более 50 символов'),
    email: z.string().email('Введите действительный адрес электронной почты'),
    phone: z.string().min(7, 'Телефон должен содержать минимум 7 символов').max(15, 'Телефон должен содержать не более 15 символов'),
    password: isEditing.value
      ? z.string().optional()
      : z.string().min(6, 'Пароль должен содержать минимум 6 символов'),
    storeId: z.string().nonempty('Выберите магазин'),
    position: z.string().nonempty('Выберите должность'),
  })
);

// Initialize the form with Vee-Validate
const { handleSubmit, errors, values } = useForm({
  validationSchema: schema,
  initialValues: {
    name: '',
    surname: '',
    email: '',
    phone: '',
    password: '',
    storeId: '',
    position: '',
  },
});

// Employee data (including image and working hours)
const employee = ref({
  image: null as File | null,
  workingHours: {},
});

// Handle form submission
const onSubmit = handleSubmit((formValues) => {
  const employeeData = {
    ...formValues,
    image: employee.value.image,
    workingHours: employee.value.workingHours,
  };
  if (isEditing.value) {
    // Update employee logic
    console.log('Employee updated:', employeeData);
  } else {
    // Create employee logic
    console.log('Employee created:', employeeData);
  }
  router.push('/employees'); // Redirect after saving
});

// Handle cancel action
const handleCancel = () => {
  router.push('/employees');
};
</script>
