<script setup lang="ts">
import { Separator } from '@/core/components/ui/separator'
import { ArrowLeft, ChevronRight } from 'lucide-vue-next'
import { ref } from 'vue'

import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import type { EmployeeLoginDTO } from '@/modules/admin/employees/models/employees.models'
import AdminLoginForm from '@/modules/auth/components/login/admin-login-form.vue'
import FranchiseeLoginForm from '@/modules/auth/components/login/franchisee-login-form.vue'
import RegionLoginForm from '@/modules/auth/components/login/region-login-form.vue'
import StoreLoginForm from '@/modules/auth/components/login/store-login-form.vue'
import WarehouseLoginForm from '@/modules/auth/components/login/warehouse-login-form.vue'
import { authService } from '@/modules/auth/services/auth.service'
import { useMutation } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

// Role selection state
const selectedRole = ref<string | null>(null)

// Lazy load role-specific login components
const roleComponents: Record<string, unknown> = {
  Store: StoreLoginForm,
  Warehouse: WarehouseLoginForm,
  Region: RegionLoginForm,
  Franchise: FranchiseeLoginForm,
  Admin: AdminLoginForm
}

// Role options in Russian
const roles = [
  { key: 'Store', label: 'Кафе' },
  { key: 'Warehouse', label: 'Центральный склад' },
  { key: 'Region', label: 'Регионы' },
  { key: 'Franchise', label: 'Франчайзи' },
  { key: 'Admin', label: 'Администратор' }
]


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

<template>
	<!-- Outer container centers the card on the page -->
	<div
		class="flex flex-col justify-center items-center gap-6 bg-gradient-to-r from-slate-900 to-slate-700 p-4 min-h-screen"
	>
		<div class="border-gray-100 bg-white shadow-lg p-6 md:p-8 rounded-xl w-full max-w-md h-full">
			<!-- Wrap the conditional blocks in a Transition for smooth animation -->
			<transition
				name="fade-slide"
				mode="out-in"
			>
				<!-- Role Selection Block -->
				<template v-if="!selectedRole">
					<div class="flex flex-col gap-5">
						<h1 class="font-semibold text-gray-800 text-xl md:text-2xl">
							Выберите способ авторизации
						</h1>
						<div class="flex flex-col gap-2">
							<div
								v-for="(role, index) in roles"
								:key="role.key"
								class="w-full"
							>
								<button
									@click="selectedRole = role.key"
									class="flex justify-between items-center py-2 md:py-3 w-full text-lg md:text-xl hover:text-primary"
								>
									{{ role.label }}
									<ChevronRight class="w-5 h-5" />
								</button>
								<Separator
									v-if="index !== roles.length - 1"
									class="bg-gray-100"
								/>
							</div>
						</div>
					</div>
				</template>

				<!-- Dynamic Login Form Block -->
				<template v-else>
					<div>
						<!-- Back Button -->
						<button
							@click="selectedRole = null"
							class="flex items-center gap-2 mb-4 text-gray-600 hover:text-gray-900 transition-all duration-200"
						>
							<ArrowLeft class="w-4 h-4" />
							Назад
						</button>
						<Suspense>
							<template #default>
								<component
									:is="roleComponents[selectedRole]"
									@login="onLoginEmployee"
								/>
							</template>
						</Suspense>
					</div>
				</template>
			</transition>
		</div>
	</div>
</template>

<style scoped>
/* Transition classes for a smooth fade and slide effect */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}
.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
.fade-slide-enter-to,
.fade-slide-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
