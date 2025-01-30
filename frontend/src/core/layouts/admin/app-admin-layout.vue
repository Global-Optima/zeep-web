<script setup lang="ts">

import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import AppAdminHeader from '@/core/layouts/admin/app-admin-header.vue'
import AppAdminSidebar from '@/core/layouts/admin/app-admin-sidebar.vue'
import { EmployeeRole, EmployeeType } from '@/modules/admin/store-employees/models/employees.models'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { Coffee, Store, TvMinimal } from 'lucide-vue-next'
import { computed } from 'vue'
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

const showBottomButtons = computed(() => currentEmployee?.role ? [EmployeeRole.BARISTA, EmployeeRole.STORE_MANAGER].includes(currentEmployee.role) : false )

const layoutLabel = computed(() => {
  if (!currentEmployee?.type) return "ZEEP"

  switch (currentEmployee.type) {
    case EmployeeType.STORE:
      return "Кафе"
    case EmployeeType.ADMIN:
      return "Администрация"
    case EmployeeType.FRANCHISEE:
      return "Франшизы"
    case EmployeeType.REGION:
      return "Регионы"
    case EmployeeType.WAREHOUSE:
      return "Склад"
    default:
      return "ZEEP"
  }
})
</script>

<template>
	<div class="flex bg-white w-full min-h-screen">
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
							<span class="font-normal">
								{{ layoutLabel }}
							</span>
						</div>
					</a>
				</div>

				<div class="flex-1">
					<AppAdminSidebar />
				</div>

				<div
					class="mt-auto p-4"
					v-if="showBottomButtons"
				>
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
			<AppAdminHeader />

			<main class="flex-1 bg-slate-50 p-4 lg:p-6">
				<RouterView />
			</main>
		</div>
	</div>
</template>
