<template>
	<div class="flex justify-center items-center bg-gray-200 w-full h-full">
		<Tabs
			default-value="cafe"
			class="w-[400px]"
		>
			<TabsList class="grid grid-cols-2 w-full">
				<TabsTrigger
					value="cafe"
					class="text-base"
				>
					Кафе
				</TabsTrigger>
				<TabsTrigger
					value="warehouse"
					class="text-base"
				>
					Склад
				</TabsTrigger>
			</TabsList>
			<TabsContent value="cafe">
				<StoreLoginForm @login="onLoginEmployee" />
			</TabsContent>
			<TabsContent value="warehouse">
				<WarehouseLoginForm @login="onLoginEmployee" />
			</TabsContent>
		</Tabs>
	</div>
</template>

<script setup lang="ts">
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/core/components/ui/tabs'
import StoreLoginForm from '@/modules/auth/components/login/store-login-form.vue'
import WarehouseLoginForm from '@/modules/auth/components/login/warehouse-login-form.vue'

import { getRouteName } from '@/core/config/routes.config'
import { toastError, toastSuccess } from '@/core/config/toast.config'
import { authService } from '@/modules/auth/services/auth.service'
import type { EmployeeLoginDTO } from '@/modules/employees/models/employees.models'
import { useMutation } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()

const {mutate: loginEmployee} = useMutation({
		mutationFn: (dto: EmployeeLoginDTO) => authService.loginEmployee(dto),
		onSuccess: () => {
			toastSuccess("Вы вошли в систему")
			router.push({name: getRouteName("ADMIN_DASHBOARD")})
		},
		onError: () => {
			toastError("Произошла ошибка при входе")
		},
})

const onLoginEmployee = (dto: EmployeeLoginDTO) => {
  return loginEmployee(dto)
}
</script>
