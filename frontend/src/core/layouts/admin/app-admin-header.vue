<template>
	<header class="flex items-center gap-4 bg-muted/40 px-4 lg:px-6 border-b h-14 lg:h-[60px]">
		<Sheet>
			<SheetTrigger as-child>
				<Button
					variant="outline"
					size="icon"
					class="md:hidden shrink-0"
				>
					<Menu class="w-5 h-5" />
					<span class="sr-only">Toggle navigation menu</span>
				</Button>
			</SheetTrigger>

			<SheetContent
				side="left"
				class="flex flex-col"
				:include-close-button="false"
			>
				<AppAdminSidebar />

				<div class="mt-auto">
					<Button
						variant="outline"
						@click="onBaristaClick"
						class="flex items-center gap-3 mb-1 py-6 w-full"
					>
						<Store class="w-5 h-5" />
						<p class="text-base">Бариста</p>
					</Button>

					<Button
						@click="onKioskClick"
						class="flex items-center gap-3 py-6 w-full"
					>
						<Store class="w-5 h-5" />
						<p class="text-base">Киоск</p>
					</Button>
				</div>
			</SheetContent>
		</Sheet>

		<div class="flex-1 w-full">
			<form>
				<div class="relative">
					<Search class="top-2.5 left-2.5 absolute w-4 h-4 text-muted-foreground" />
					<Input
						type="search"
						placeholder="Поиск"
						class="bg-background shadow-none pl-8 w-full md:w-2/3 lg:w-1/3 appearance-none"
					/>
				</div>
			</form>
		</div>

		<DropdownMenu v-if="currentEmployee">
			<DropdownMenuTrigger as-child>
				<Button
					variant="secondary"
					size="icon"
					class="rounded-full"
				>
					<CircleUser class="w-5 h-5" />
					<span class="sr-only">Toggle user menu</span>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end">
				<DropdownMenuLabel>
					{{ currentEmployee.firstName }} {{ currentEmployee.lastName }}
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuItem @click="onLogoutClick">Выйти</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</header>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import { Sheet, SheetContent, SheetTrigger } from '@/core/components/ui/sheet'
import { getRouteName } from '@/core/config/routes.config'
import { toastError, toastSuccess } from '@/core/config/toast.config'
import AppAdminSidebar from '@/core/layouts/admin/app-admin-sidebar.vue'
import { authService } from '@/modules/auth/services/auth.service'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { useMutation } from '@tanstack/vue-query'
import { CircleUser, Menu, Search, Store } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {currentEmployee, setCurrentEmployee} = useEmployeeAuthStore()

console.log(currentEmployee)

const {mutate: logoutEmployee} = useMutation({
		mutationFn: () => authService.logoutEmployee(),
		onSuccess: () => {
			toastSuccess("Вы вышли из системы")
      setCurrentEmployee(null)
			router.push({name: getRouteName("LOGIN")})
		},
		onError: () => {
			toastError("Произошла ошибка при выходе")
		},
})

const router = useRouter()

const onKioskClick = () => {
	router.push({name: getRouteName('KIOSK_HOME')})
}

const onBaristaClick = () => {
	router.push({name: getRouteName('KIOSK_ORDERS')})
}

const onLogoutClick = () => {
  logoutEmployee()
}
</script>

<style scoped></style>
