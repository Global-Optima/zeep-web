<template>
	<Card class="border-none">
		<CardHeader class="p-0">
			<CardTitle class="text-lg sm:text-xl"> Вход для регионального сотрудника</CardTitle>
			<CardDescription> Введите ваши учетные данные для входа в систему </CardDescription>
		</CardHeader>
		<CardContent class="mt-6 p-0">
			<form
				class="space-y-6 w-full"
				@submit="onSubmit"
			>
				<!-- Region Selection -->
				<FormField
					v-slot="{ componentField }"
					name="selectedStoreId"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Регион</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger class="w-full">
									<template v-if="storesLoading">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Загрузка регионов..."
										/>
									</template>
									<template v-else-if="storesError">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Ошибка загрузки регионов"
										/>
									</template>
									<template v-else>
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Выберите регион"
										/>
									</template>
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="store in stores"
										:key="store.id"
										:value="store.id.toString()"
										class="text-sm sm:text-base"
									>
										{{ store.name }}
									</SelectItem>
								</SelectContent>
							</Select>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Employee Selection (dependent on the region) -->
				<FormField
					v-slot="{ componentField }"
					name="selectedEmployeeEmail"
				>
					<FormItem v-if="values.selectedStoreId">
						<FormLabel class="text-sm sm:text-base">Сотрудник</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger class="w-full">
									<template v-if="employeesLoading">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Загрузка сотрудников..."
										/>
									</template>
									<template v-else-if="employeesError">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Ошибка загрузки сотрудников"
										/>
									</template>
									<template v-else>
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Выберите сотрудника"
										/>
									</template>
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="employee in employees"
										:key="employee.email"
										:value="employee.email"
										class="text-sm sm:text-base"
									>
										{{ employee.firstName }} {{ employee.lastName }}
									</SelectItem>
								</SelectContent>
							</Select>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Password Field -->
				<FormField
					v-slot="{ componentField }"
					name="password"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Пароль</FormLabel>
						<FormControl>
							<Input
								type="password"
								placeholder="Введите пароль"
								v-bind="componentField"
								class="text-sm sm:text-base"
								required
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Submit Button -->
				<Button
					:disabled="isSubmitting"
					type="submit"
					class="w-full"
				>
					Войти
				</Button>
			</form>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed } from 'vue'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'

import type { EmployeeLoginDTO } from '@/modules/admin/employees/models/employees.models'
import { regionsService } from '@/modules/admin/regions/services/regions.service'
import { authService } from '@/modules/auth/services/auth.service'

// Объявление событий компонента
const emits = defineEmits<{
  login: [payload: EmployeeLoginDTO]
}>()

// Определение схемы валидации с использованием Zod
const formSchema = toTypedSchema(
  z.object({
    selectedStoreId: z.coerce
      .number()
      .min(1, { message: 'Пожалуйста, выберите регион' }),
    selectedEmployeeEmail: z.string().min(1, {
      message: 'Пожалуйста, выберите сотрудника',
    }),
    password: z.string().min(2, 'Пароль должен содержать не менее 2 символов'),
  })
)

// Настройка формы с VeeValidate
const { values, isSubmitting, handleSubmit } = useForm({
  validationSchema: formSchema,
})

// Получение списка регионов
const {
  data: stores,
  isLoading: storesLoading,
  isError: storesError,
} = useQuery({
  queryKey: ['regions-all'],
  queryFn: () => regionsService.getAllRegions(),
})

// Получение списка сотрудников для выбранного региона
const {
  data: employees,
  isLoading: employeesLoading,
  isError: employeesError,
} = useQuery({
  queryKey: computed(() => ['region-employees', values.selectedStoreId]),
  queryFn: () => authService.getRegionAccounts(values.selectedStoreId!),
  enabled: computed(() => Boolean(values.selectedStoreId)),
  initialData: [],
})

// Обработка отправки формы
const onSubmit = handleSubmit((formValues) => {
  const dto: EmployeeLoginDTO = {
    email: formValues.selectedEmployeeEmail,
    password: formValues.password,
  }
  emits('login', dto)
})
</script>

<style scoped>
/* Добавьте любые специфические стили для компонента */
</style>
