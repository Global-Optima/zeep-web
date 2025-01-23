<script setup lang="ts">

import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import AppAdminHeader from '@/core/layouts/admin/app-admin-header.vue'
import AppAdminSidebar from '@/core/layouts/admin/app-admin-sidebar.vue'
import { EmployeeType } from '@/modules/admin/store-employees/models/employees.models'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { Coffee, Store, TvMinimal } from 'lucide-vue-next'
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
</script>

<template>
	<div class="grid md:grid-cols-[200px_1fr] lg:grid-cols-[220px_1fr] bg-white w-full min-h-screen">
		<div class="md:block relative hidden bg-muted/40 border-r">
			<div class="top-0 sticky flex flex-col gap-2 h-full max-h-screen">
				<div class="flex items-center px-4 lg:px-6 border-b h-14 lg:h-[60px]">
					<a
						href="/admin"
						class="flex items-center gap-3 font-semibold"
					>
						<div class="flex items-center gap-2">
							<span class="font-bold text-xl">Zeep</span>
							<span class="font-normal text-base">|</span>
							<span
								class="font-normal"
								>{{ currentEmployee?.type === EmployeeType.STORE ? "Кафе" : "Склад" }}</span
							>
						</div>
					</a>
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

		<div class="flex flex-col overflow-y-auto">
			<AppAdminHeader />

			<main class="flex-1 bg-slate-50 p-4 lg:p-6">
				<RouterView />
			</main>
		</div>
	</div>
</template>
