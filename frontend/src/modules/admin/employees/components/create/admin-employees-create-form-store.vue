<template>
	<Card>
		<CardHeader>
			<CardTitle>Создать сотрудника</CardTitle>
			<CardDescription> Выберите ниже, какого сотрудника вы хотите создать. </CardDescription>
		</CardHeader>

		<CardContent>
			<form
				@submit="onSubmit"
				class="gap-6 grid"
			>
				<!-- Name & Surname -->
				<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
					<FormField
						name="name"
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
						name="surname"
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

				<!-- Email & Phone -->
				<div class="gap-4 grid grid-cols-1 sm:grid-cols-2">
					<FormField
						name="email"
						v-slot="{ componentField }"
					>
						<FormItem>
							<FormLabel>Email</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="example@store.com"
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
									v-bind="componentField"
									placeholder="+7 (___) ___-__-__"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<!-- Password -->
				<FormField
					name="password"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormLabel>Пароль</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="password"
								placeholder="Введите пароль"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Role -->
				<FormField
					name="role"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormLabel>Должность</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger>
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

				<!-- Store Selection -->
				<FormField
					name="storeDetails.storeId"
					v-slot="{ componentField }"
				>
					<FormItem>
						<FormLabel>Выберите магазин</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger class="w-full">
									<template v-if="storesLoading">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Загрузка кафе..."
										/>
									</template>
									<template v-else-if="storesError">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Ошибка загрузки кафе"
										/>
									</template>
									<template v-else>
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Выберите кафе"
										/>
									</template>
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
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Action Buttons -->
				<div class="flex gap-4 mt-6">
					<Button
						type="submit"
						class="flex-1"
						>Создать</Button
					>
					<Button
						variant="outline"
						class="flex-1"
						@click.prevent="emitCancel"
					>
						Отмена
					</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import * as z from 'zod'

import { Button } from '@/core/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/core/components/ui/form'
import { Input } from '@/core/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'

// Query fetching
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { EmployeeRole, type CreateEmployeeDto } from '@/modules/employees/models/employees.models'
import { storesService } from '@/modules/stores/services/stores.service'
import { useQuery } from '@tanstack/vue-query'

// Emit definition
const emit = defineEmits<{
  (e: 'onSubmit', p: CreateEmployeeDto): void
  (e: 'onCancel'): void
}>()

// Roles for store employees (example)
const roles = ref([
  { value: EmployeeRole.MANAGER, label: 'Менеджер' },
  { value: EmployeeRole.BARISTA, label: 'Бариста' },
])

// Fetch stores
const {
  data: stores,
  isLoading: storesLoading,
  isError: storesError,
} = useQuery({
  queryKey: ['stores'],
  queryFn: () => storesService.getStores(),
  initialData: [],
})

// Validation schema specific to store employees
const storeSchema = toTypedSchema(
  z
    .object({
      name: z.string().min(2, 'Имя должно содержать минимум 2 символа'),
      surname: z.string().min(2, 'Фамилия должна содержать минимум 2 символа'),
      email: z.string().email('Введите корректный email'),
      phone: z.string().min(7, 'Телефон должен содержать минимум 7 символов'),
      password: z.string().min(6, 'Пароль должен содержать минимум 6 символов'),
      role: z.string().min(1, 'Выберите должность'),
      storeDetails: z.object({
        storeId: z.coerce.number().positive('Необходимо выбрать магазин'),
      }),
    })
)

// UseForm
const { handleSubmit } = useForm<CreateEmployeeDto>({
  validationSchema: storeSchema,
  initialValues: {
    storeDetails: { storeId: 0 },
  },
})

// Handles form submission
const onSubmit = handleSubmit(async (values) => {
  // Usually you'd call an API here, e.g.: await createStoreEmployee(values)
  // For demonstration, we just emit success
  emit('onSubmit', values)
})

// Cancel
const emitCancel = () => {
  emit('onCancel')
}
</script>
