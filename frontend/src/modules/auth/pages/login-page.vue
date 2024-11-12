<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader
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
import { getRouteName } from '@/core/config/routes.config'
import { CURRENT_STORE_COOKIES_CONFIG } from '@/modules/stores/constants/store-cookies.constant'
import { storesService } from "@/modules/stores/services/stores.service"
import { useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { watch } from 'vue'
import { useRouter } from "vue-router"
import * as z from 'zod'

const formSchema = toTypedSchema(
  z.object({
    selectedStoreId: z.string().min(1, {message: "Please select a store"}),
    selectedEmployeeId: z.string().min(1, {message: "Please select an employee"}),
    password: z.string().min(6, "Password must be at least 6 characters long"),
  })
)

const { values, isSubmitting, handleSubmit } = useForm({
  validationSchema: formSchema,
})

const { data: stores, isLoading: storesLoading, isError: storesError } = useQuery({
  queryKey: ['stores'],
  queryFn: storesService.getStores,
  initialData: [],
})

const { data: employees, refetch: refetchEmployees, isLoading: employeesLoading, isError: employeesError } = useQuery({
  queryKey: ['employees', values.selectedStoreId],
  queryFn: () => storesService.getStoreEmployees(Number(values.selectedStoreId)),
  initialData: [],
  enabled: false,
})

watch(() => values.selectedStoreId, (newStore) => {
  if (newStore) refetchEmployees()
})

const router = useRouter()
const onSubmit = handleSubmit((values) => {
  localStorage.setItem(CURRENT_STORE_COOKIES_CONFIG.key, values.selectedStoreId)
  router.push({ name: getRouteName("ADMIN_DASHBOARD") })
})
</script>

<template>
	<div class="w-full h-full flex items-center justify-center">
		<Card class="">
			<CardHeader>
				<CardTitle class="text-xl sm:text-2xl">Вход для сотрудника</CardTitle>
				<CardDescription>Введите данные для входа в портал сотрудника</CardDescription>
			</CardHeader>
			<CardContent>
				<form
					class="w-full space-y-6"
					@submit="onSubmit"
				>
					<FormField
						v-slot="{ componentField }"
						name="selectedStoreId"
					>
						<FormItem>
							<FormLabel class="text-sm sm:text-base">Заведение</FormLabel>
							<FormControl>
								<Select
									v-model="values.selectedStoreId"
									v-bind="componentField"
								>
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
						name="selectedEmployeeId"
					>
						<FormItem v-if="values.selectedStoreId">
							<FormLabel class="text-sm sm:text-base">Выберите сотрудника</FormLabel>
							<FormControl>
								<Select
									v-model="values.selectedEmployeeId"
									v-bind="componentField"
								>
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
											:key="employee.id"
											:value="employee.id.toString()"
											class="text-sm sm:text-base"
										>
											{{ employee.name }}
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
	</div>
</template>
