<template>
	<Card class="border-none">
		<CardHeader class="p-0">
			<CardTitle class="text-lg sm:text-xl">Вход для сотрудников склада</CardTitle>
			<CardDescription> Введите ваши учетные данные для входа в систему</CardDescription>
		</CardHeader>
		<CardContent class="mt-6 p-0">
			<form
				class="space-y-6 w-full"
				@submit="onSubmit"
			>
				<FormField
					v-slot="{ componentField }"
					name="selectedWarehouseId"
				>
					<FormItem>
						<FormLabel class="text-sm sm:text-base">Склад</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger class="w-full">
									<template v-if="warehousesLoading">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Загрузка склада..."
										/>
									</template>
									<template v-else-if="warehousesError">
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Ошибка загрузки склада"
										/>
									</template>
									<template v-else>
										<SelectValue
											class="text-sm sm:text-base"
											placeholder="Выберите склад"
										/>
									</template>
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="warehouse in warehouses"
										:key="warehouse.id"
										:value="warehouse.id.toString()"
										class="text-sm sm:text-base"
									>
										{{ warehouse.name }}
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
					<FormItem v-if="values.selectedWarehouseId">
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
							<Input
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/core/components/ui/select'
import type { EmployeeLoginDTO } from '@/modules/admin/employees/models/employees.models'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'
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
    selectedWarehouseId: z.coerce.number().min(1, {message: "Пожалуйста, выберите склад"}),
    selectedEmployeeEmail: z.string().min(1, {message: "Пожалуйста, выберите сотрудника"}),
    password: z.string().min(2, "Пароль должен содержать не менее 2 символов"),
  })
)

const { values, isSubmitting, handleSubmit } = useForm({
  validationSchema: formSchema,
})

const { data: warehouses, isLoading: warehousesLoading, isError: warehousesError } = useQuery({
  queryKey: ['warehouses-all'],
  queryFn: () => warehouseService.getAll(),
})

const { data: employees, isLoading: employeesLoading, isError: employeesError } = useQuery({
  queryKey: computed(() => ['warehouse-employees', values.selectedWarehouseId]),
  queryFn: () => authService.getWarehouseAccounts(values.selectedWarehouseId!),
  enabled: computed(() => Boolean(values.selectedWarehouseId)),
  initialData: []
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
