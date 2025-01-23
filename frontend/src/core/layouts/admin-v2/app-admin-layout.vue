<script setup lang="ts">

import { Avatar, AvatarFallback } from '@/core/components/ui/avatar'
import { Button } from '@/core/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/core/components/ui/dropdown-menu'
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import AppAdminHeader from '@/core/layouts/admin-v2/app-admin-header.vue'
import AppAdminSidebar from '@/core/layouts/admin-v2/app-admin-sidebar.vue'
import { EMPLOYEE_ROLES_FORMATTED } from '@/modules/admin/store-employees/models/employees.models'
import { authService } from '@/modules/auth/services/auth.service'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { useMutation } from '@tanstack/vue-query'
import { ChevronsUpDown, Coffee, Store, TvMinimal } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()

const {currentEmployee} = useEmployeeAuthStore()

const onDisplayClick = () => {
	router.push({name: getRouteName('KIOSK_ORDERS_DISPLAY')})
}

const onKioskClick = () => {
	router.push({name: getRouteName('KIOSK_HOME')})
}

const onBaristaClick = () => {
	router.push({name: getRouteName('KIOSK_ORDERS')})
}

const { toast } = useToast()

const {mutate: logoutEmployee} = useMutation({
		mutationFn: () => authService.logoutEmployee(),
		onSuccess: () => {
			toast({title: "Вы вышли из системы"})
			router.push({name: getRouteName("LOGIN")})
		},
		onError: () => {
      toast({title: "Произошла ошибка при выходе"})
		},
})

const onLogoutClick = () => {
  logoutEmployee()
}
</script>

<template>
	<div class="flex bg-white w-full min-h-screen">
		<div class="md:block relative hidden bg-muted/40 border-r">
			<div class="top-0 sticky flex flex-col gap-2 h-full max-h-screen">
				<div
					class="p-2 border-b"
					v-if="currentEmployee"
				>
					<DropdownMenu>
						<DropdownMenuTrigger as-child>
							<div>
								<div
									class="flex justify-between items-center gap-8 hover:bg-muted p-2 rounded-md cursor-pointer"
								>
									<div class="flex items-center gap-4">
										<Avatar class="bg-gray-200 rounded-md w-11 h-11">
											<AvatarFallback>
												{{ `${currentEmployee.firstName.charAt(0)}${currentEmployee.lastName.charAt(0)}` }}
											</AvatarFallback>
										</Avatar>

										<div>
											<p class="font-medium text-sm">
												{{ currentEmployee.firstName }} {{ currentEmployee.lastName }}
											</p>
											<p class="mt-0.5 text-gray-500 text-xs">
												{{ EMPLOYEE_ROLES_FORMATTED[currentEmployee.role] }}
											</p>
										</div>
									</div>

									<ChevronsUpDown class="w-4 h-4 text-gray-400" />
								</div>
							</div>
						</DropdownMenuTrigger>
						<DropdownMenuContent
							class="rounded-lg min-w-56"
							side="bottom"
							:side-offset="4"
						>
							<DropdownMenuItem
								class="w-full"
								@click="onLogoutClick"
								>Выйти</DropdownMenuItem
							>
						</DropdownMenuContent>
					</DropdownMenu>
				</div>

				<div class="flex-1">
					<AppAdminSidebar />
				</div>

				<div class="hidden mt-auto p-4">
					<Button
						variant="outline"
						@click="onDisplayClick"
						class="flex items-center gap-3 mb-2 py-5 w-full"
					>
						<TvMinimal class="w-4 h-4" />
						<p class="text-base">Экран</p>
					</Button>

					<Button
						variant="outline"
						@click="onBaristaClick"
						class="flex items-center gap-3 mb-2 py-5 w-full"
					>
						<Coffee class="w-5 h-5" />
						<p class="text-base">Бариста</p>
					</Button>

					<Button
						@click="onKioskClick"
						class="flex items-center gap-3 py-5 w-full"
					>
						<Store class="w-5 h-5" />
						<p class="text-base">Киоск</p>
					</Button>
				</div>
			</div>
		</div>

		<div class="flex flex-col flex-1">
			<div class="block md:hidden">
				<AppAdminHeader />
			</div>

			<main class="flex-1 bg-slate-50 p-4 lg:p-6">
				<RouterView />
			</main>
		</div>
	</div>
</template>
