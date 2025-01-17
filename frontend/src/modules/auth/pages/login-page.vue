<template>
	<div class="flex justify-center items-center bg-gray-200 w-full h-full">
		<Tabs default-value="cafe">
			<TabsList class="grid grid-cols-3 w-full">
				<TabsTrigger
					value="cafe"
					class="p-2 text-base"
				>
					Кафе
				</TabsTrigger>
				<TabsTrigger
					value="warehouse"
					class="p-2 text-base"
				>
					Склад
				</TabsTrigger>

				<TabsTrigger
					value="admin"
					class="p-2 text-base"
				>
					Администраторы
				</TabsTrigger>
			</TabsList>
			<TabsContent value="cafe">
				<StoreLoginForm @login="onLoginEmployee" />
			</TabsContent>
			<TabsContent value="warehouse">
				<WarehouseLoginForm @login="onLoginEmployee" />
			</TabsContent>
			<TabsContent value="admin">
				<AdminLoginForm @login="onLoginEmployee" />
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
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import AdminLoginForm from '@/modules/auth/components/login/admin-login-form.vue'
import StoreLoginForm from '@/modules/auth/components/login/store-login-form.vue'
import WarehouseLoginForm from '@/modules/auth/components/login/warehouse-login-form.vue'
import { authService } from '@/modules/auth/services/auth.service'
import type { EmployeeLoginDTO } from '@/modules/employees/models/employees.models'
import { useMutation } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const { toast } = useToast()

const {mutate: loginEmployee} = useMutation({
		mutationFn: (dto: EmployeeLoginDTO) => authService.loginEmployee(dto),
		onSuccess: () => {
			toast({title: "Вы вошли в систему"})
			router.push({name: getRouteName("ADMIN_DASHBOARD")})
		},
		onError: () => {
			toast({title: "Произошла ошибка при входе"})
		},
})

const onLoginEmployee = (dto: EmployeeLoginDTO) => {
  return loginEmployee(dto)
}
</script>
