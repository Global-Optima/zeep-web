<script setup lang="ts">
import PasswordInput from '@/core/components/password-input/PasswordInput.vue'
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/core/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/core/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/core/components/ui/select'
import { easyPasswordValidationSchema } from '@/core/validators/password.validator'
import type { EmployeeLoginDTO } from '@/modules/admin/employees/models/employees.models'
import { authService } from '@/modules/auth/services/auth.service'
import { useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

const emits = defineEmits<{
  'login': [payload: EmployeeLoginDTO]
}>()

// Updated validation schema with Latin-only password rule
const formSchema = toTypedSchema(
  z.object({
    selectedEmployeeEmail: z.string().min(1, { message: "Пожалуйста, выберите сотрудника" }),
    password: easyPasswordValidationSchema
  }
))

const { isSubmitting, handleSubmit } = useForm({
  validationSchema: formSchema,
})

const { data: employees, isLoading: employeesLoading, isError: employeesError } = useQuery({
  queryKey: ['admin-employees'],
  queryFn: () => authService.getAdminsAccounts(),
  initialData: [],
})

const onSubmit = handleSubmit((values) => {
  const dto: EmployeeLoginDTO = {
    email: values.selectedEmployeeEmail,
    password: values.password,
  }
  emits("login", dto)
})
</script>

<template>
	<Card class="border-none">
		<CardHeader class="p-0">
			<CardTitle class="text-lg sm:text-xl">Вход для администраторов</CardTitle>
			<CardDescription> Введите ваши учетные данные для входа в систему</CardDescription>
		</CardHeader>
		<CardContent class="mt-6 p-0">
			<form
				class="space-y-6 w-full"
				@submit="onSubmit"
			>
				<FormField
					v-slot="{ componentField }"
					name="selectedEmployeeEmail"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Выберите учетную запись</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger class="w-full">
									<template v-if="employeesLoading">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Загрузка учетных записей..."
										/>
									</template>
									<template v-else-if="employeesError">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Ошибка загрузки учетных записей"
										/>
									</template>
									<template v-else>
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Учетные записи"
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
				<FormField
					v-slot="{ componentField }"
					name="password"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Пароль</FormLabel>
						<FormControl>
							<PasswordInput
								placeholder="Введите ваш пароль"
								v-bind="componentField"
								class="text-sm sm:text-base"
								required
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
				<Button
					:disabled="isSubmitting"
					type="submit"
					size="lg"
					class="w-full"
				>
					Войти
				</Button>
			</form>
		</CardContent>
	</Card>
</template>
