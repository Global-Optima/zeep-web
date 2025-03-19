<template>
	<Card class="border-none">
		<CardHeader class="p-0">
			<CardTitle class="text-lg sm:text-xl">Вход для сотрудника кафе</CardTitle>
			<CardDescription> Введите ваши учетные данные для входа в систему</CardDescription>
		</CardHeader>
		<CardContent class="mt-6 p-0">
			<form
				class="space-y-6 w-full"
				@submit="onSubmit"
			>
				<FormField
					v-slot="{ componentField }"
					name="selectedStoreId"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Заведение</FormLabel>
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

				<FormField
					v-slot="{ componentField }"
					name="selectedEmployeeEmail"
				>
					<FormItem v-if="values.selectedStoreId">
						<FormLabel class="text-sm sm:text-base">Выберите сотрудника</FormLabel>
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
											placeholder="Сотрудники"
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
								type="password"
								placeholder="Введите пароль сотрудника"
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
					class="w-full"
				>
					Войти
				</Button>
			</form>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import PasswordInput from '@/core/components/password-input/PasswordInput.vue'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/core/components/ui/select'
import { loginPasswordValidationSchema } from '@/core/validators/password.validator'
import type { EmployeeLoginDTO } from '@/modules/admin/employees/models/employees.models'
import { storesService } from "@/modules/admin/stores/services/stores.service"
import { authService } from '@/modules/auth/services/auth.service'
import { useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed } from 'vue'
import * as z from 'zod'

const emits = defineEmits<{
  'login': [payload: EmployeeLoginDTO]
}>()

const formSchema = toTypedSchema(
  z.object({
    selectedStoreId: z.coerce.number().min(1, {message: "Пожалуйста, выберите кафе"}),
    selectedEmployeeEmail: z.string().min(1, {message: "Пожалуйста, выберите сотрудника"}),
    password: loginPasswordValidationSchema
  })
)

const { values, isSubmitting, handleSubmit } = useForm({
  validationSchema: formSchema,
})


const { data: stores, isLoading: storesLoading, isError: storesError } = useQuery({
  queryKey: ['stores-all'],
  queryFn: () => storesService.getAll(),
})

const { data: employees, isLoading: employeesLoading, isError: employeesError } = useQuery({
  queryKey: computed(() => ['store-employees', values.selectedStoreId]),
  queryFn: () => authService.getStoreAccounts(values.selectedStoreId!),
  enabled: computed(() => Boolean(values.selectedStoreId)),
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

<style scoped></style>
