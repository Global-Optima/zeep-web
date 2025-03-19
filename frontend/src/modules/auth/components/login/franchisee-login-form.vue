<template>
	<Card class="border-none">
		<CardHeader class="p-0">
			<CardTitle class="text-lg sm:text-xl"> Вход для сотрудника франчайзинговой сети </CardTitle>
			<CardDescription> Пожалуйста, введите учетные данные для доступа к системе. </CardDescription>
		</CardHeader>
		<CardContent class="mt-6 p-0">
			<form
				class="space-y-6 w-full"
				@submit="onSubmit"
			>
				<!-- Филиал (заведение) выбора -->
				<FormField
					v-slot="{ componentField }"
					name="selectedStoreId"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Филиал</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger class="w-full">
									<template v-if="storesLoading">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Загрузка филиалов..."
										/>
									</template>
									<template v-else-if="storesError">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Ошибка загрузки филиалов"
										/>
									</template>
									<template v-else>
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Выберите филиал"
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

				<!-- Выбор сотрудника (активируется после выбора филиала) -->
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

				<!-- Поле для ввода пароля -->
				<FormField
					v-slot="{ componentField }"
					name="password"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Пароль</FormLabel>
						<FormControl>
							<PasswordInput
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

				<!-- Кнопка входа -->
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'

import PasswordInput from '@/core/components/password-input/PasswordInput.vue'
import { passwordValidationSchema } from '@/core/validators/password.validator'
import type { EmployeeLoginDTO } from '@/modules/admin/employees/models/employees.models'
import { franchiseeService } from '@/modules/admin/franchisees/services/franchisee.service'
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
      .min(1, { message: 'Пожалуйста, выберите филиал' }),
    selectedEmployeeEmail: z.string().min(1, {
      message: 'Пожалуйста, выберите сотрудника',
    }),
    password: passwordValidationSchema,
  })
)

// Настройка формы с VeeValidate
const { values, isSubmitting, handleSubmit } = useForm({
  validationSchema: formSchema,
})

// Получение списка филиалов (заведения)
const {
  data: stores,
  isLoading: storesLoading,
  isError: storesError,
} = useQuery({
  queryKey: ['franchisees-all'],
  queryFn: () => franchiseeService.getAll(),
})

// Получение списка сотрудников для выбранного филиала
const {
  data: employees,
  isLoading: employeesLoading,
  isError: employeesError,
} = useQuery({
  queryKey: computed(() => ['franchisees-employees', values.selectedStoreId]),
  queryFn: () => authService.getFranchiseeAccounts(values.selectedStoreId!),
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
/* Добавьте здесь специфичные стили для компонента, если необходимо */
</style>
