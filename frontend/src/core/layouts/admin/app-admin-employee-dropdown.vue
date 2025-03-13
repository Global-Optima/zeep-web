<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/core/components/ui/dropdown-menu'
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import { EMPLOYEE_ROLES_FORMATTED, type EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
import { authService } from '@/modules/auth/services/auth.service'
import { useMutation } from '@tanstack/vue-query'
import { User } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {currentEmployee} = defineProps<{currentEmployee: EmployeeDTO}>()
const router = useRouter()

const { toast } = useToast()

const {mutate: logoutEmployee} = useMutation({
		mutationFn: () => authService.logoutEmployee(),
		onSuccess: () => {
      localStorage.clear();
			router.push({name: getRouteName("LOGIN")})
		},
		onError: () => {
      toast({title: "Произошла ошибка при выходе"})
		},
})

const onLogoutClick = () => {
  logoutEmployee()
}

const onEmployeeClick = () => {
  router.push({name: getRouteName("ADMIN_EMPLOYEE_PROFILE")});

};
</script>

<template>
	<DropdownMenu>
		<DropdownMenuTrigger as-child>
			<Button
				variant="outline"
				size="icon"
				class="rounded-md"
			>
				<User class="size-[18px]" />
				<span class="sr-only">Toggle user menu</span>
			</Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent align="end">
			<DropdownMenuLabel class="hover:bg-gray-100 rounded-sm cursor-pointer">
				<div @click="onEmployeeClick">
					<p>{{ currentEmployee.firstName }} {{ currentEmployee.lastName }}</p>
					<p class="font-normal text-xs">{{ EMPLOYEE_ROLES_FORMATTED[currentEmployee.role] }}</p>
				</div>
			</DropdownMenuLabel>
			<DropdownMenuSeparator />
			<DropdownMenuItem @click="onLogoutClick">Выйти</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
