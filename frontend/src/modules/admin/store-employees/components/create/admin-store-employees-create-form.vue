<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useFieldArray, useForm } from 'vee-validate'
import * as z from 'zod'

// UI Components
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { Checkbox } from '@/core/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import { Select, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import { EmployeeRole, type CreateStoreEmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'

// Enum for Workdays
enum Weekdays {
  MONDAY = 'MONDAY',
  TUESDAY = 'TUESDAY',
  WEDNESDAY = 'WEDNESDAY',
  THURSDAY = 'THURSDAY',
  FRIDAY = 'FRIDAY',
  SATURDAY = 'SATURDAY',
  SUNDAY = 'SUNDAY',
}

// Role Options (Restricted to MANAGER and BARISTA)
const roleOptions = [
  { label: 'Менеджер', value: EmployeeRole.MANAGER },
  { label: 'Бариста', value: EmployeeRole.BARISTA },
];

// Emit Events
const emits = defineEmits<{
  (e: 'submit', value: CreateStoreEmployeeDTO): void;
  (e: 'cancel'): void;
}>();

// Validation Schema
const createEmployeeSchema = toTypedSchema(
  z.object({
    firstName: z.string().min(1, 'Введите имя'),
    lastName: z.string().min(1, 'Введите фамилию'),
    phone: z.string().optional(),
    email: z.string().email('Неверный email').min(1, 'Введите email'),
    role: z.enum([EmployeeRole.MANAGER, EmployeeRole.BARISTA]),
    password: z.string().min(8, 'Пароль должен быть минимум 8 символов'),
    isActive: z.boolean(),
    workdays: z
      .array(
        z.object({
          day: z.enum([Weekdays.MONDAY, Weekdays.TUESDAY, Weekdays.WEDNESDAY, Weekdays.THURSDAY, Weekdays.FRIDAY, Weekdays.SATURDAY, Weekdays.SUNDAY]),
          startAt: z.string().regex(/^(2[0-3]|[0-1]\d):[0-5]\d$/, 'Неверный формат времени (HH:mm)'),
          endAt: z.string().regex(/^(2[0-3]|[0-1]\d):[0-5]\d$/, 'Неверный формат времени (HH:mm)'),
        })
      )
      .min(1, 'Добавьте хотя бы один рабочий день'),
  })
);

// Form Setup
const { handleSubmit, resetForm } = useForm({
  validationSchema: createEmployeeSchema,
});

// Workdays Array
const { fields: workdays, push, remove } = useFieldArray('workdays');

// Handlers
const onSubmit = handleSubmit((formValues) => {
  const dto: CreateStoreEmployeeDTO = {
    storeId: 1,
    firstName: formValues.firstName,
    lastName: formValues.lastName,
    email: formValues.email,
    role: formValues.role,
    password: formValues.password,
    isActive: true,
    workdays: formValues.workdays,
  };

  emits('submit', dto);
});

const onCancel = () => {
  resetForm();
  emits('cancel');
};
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-6xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<h1 class="flex-1 font-semibold text-xl tracking-tight">Создать Пользователя</h1>
			<div class="md:flex items-center gap-2 hidden">
				<Button
					variant="outline"
					@click="onCancel"
					>Отменить</Button
				>
				<Button @click="onSubmit">Сохранить</Button>
			</div>
		</div>

		<!-- Main Form -->
		<div class="gap-4 grid grid-cols-1 md:grid-cols-[2fr_1fr]">
			<!-- Base Info -->
			<div class="gap-4 grid auto-rows-max">
				<Card>
					<CardHeader>
						<CardTitle>Основная информация</CardTitle>
						<CardDescription>Заполните базовую информацию о пользователе.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-6 grid">
							<FormField
								name="firstName"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Имя</FormLabel>
									<FormControl>
										<Input
											type="text"
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
											type="text"
											v-bind="componentField"
											placeholder="Введите фамилию"
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
											type="text"
											v-bind="componentField"
											placeholder="Введите телефон"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="email"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Email</FormLabel>
									<FormControl>
										<Input
											type="email"
											v-bind="componentField"
											placeholder="Введите email"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								name="password"
								v-slot="{ componentField }"
							>
								<FormItem>
									<FormLabel>Пароль</FormLabel>
									<FormControl>
										<Input
											type="password"
											v-bind="componentField"
											placeholder="Введите пароль"
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>

				<!-- Workdays -->
				<Card>
					<CardHeader>
						<CardTitle>Рабочие дни</CardTitle>
						<CardDescription>Добавьте рабочие дни для пользователя.</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="gap-4 grid">
							<div
								v-for="(workday, index) in workdays"
								:key="index"
								class="flex items-center gap-4"
							>
								<FormField
									name="workdays[index].day"
									v-slot="{ componentField }"
								>
									<FormItem>
										<FormLabel>День</FormLabel>
										<FormControl>
											<Select v-bind="componentField">
												<SelectTrigger><SelectValue /></SelectTrigger>
												<SelectItem
													v-for="(label, value) in Weekdays"
													:key="value"
													:value="value"
												>
													{{ label }}
												</SelectItem>
											</Select>
										</FormControl>
									</FormItem>
								</FormField>

								<FormField
									name="workdays[index].startAt"
									v-slot="{ componentField }"
								>
									<FormItem>
										<FormLabel>Начало</FormLabel>
										<FormControl>
											<Input
												type="time"
												v-bind="componentField"
											/>
										</FormControl>
									</FormItem>
								</FormField>

								<FormField
									name="workdays[index].endAt"
									v-slot="{ componentField }"
								>
									<FormItem>
										<FormLabel>Конец</FormLabel>
										<FormControl>
											<Input
												type="time"
												v-bind="componentField"
											/>
										</FormControl>
									</FormItem>
								</FormField>

								<Button
									variant="destructive"
									@click="remove(index)"
									>Удалить</Button
								>
							</div>

							<Button
								variant="secondary"
								@click="push({ day: 'MONDAY', startAt: '', endAt: '' })"
							>
								Добавить рабочий день
							</Button>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Roles -->
			<div>
				<Card>
					<CardHeader>
						<CardTitle>Роли</CardTitle>
						<CardDescription>Выберите роль для пользователя.</CardDescription>
					</CardHeader>
					<CardContent>
						<FormField
							name="role"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormControl>
									<Select v-bind="componentField">
										<SelectTrigger><SelectValue /></SelectTrigger>
										<SelectItem
											v-for="option in roleOptions"
											:key="option.value"
											:value="option.value"
										>
											{{ option.label }}
										</SelectItem>
									</Select>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							name="isActive"
							v-slot="{ componentField }"
						>
							<FormItem>
								<FormLabel>Активен</FormLabel>
								<FormControl>
									<Checkbox v-bind="componentField" />
								</FormControl>
							</FormItem>
						</FormField>
					</CardContent>
				</Card>
			</div>
		</div>

		<!-- Footer -->
		<div class="flex justify-center gap-2 mt-4">
			<Button
				variant="outline"
				@click="onCancel"
				>Отменить</Button
			>
			<Button @click="onSubmit">Сохранить</Button>
		</div>
	</div>
</template>
