<script setup lang="ts">

import PageLoader from '@/core/components/page-loader/PageLoader.vue'
import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import AppAdminHeader from '@/core/layouts/admin/app-admin-header.vue'
import AppAdminSidebar from '@/core/layouts/admin/app-admin-sidebar.vue'
import { EmployeeRole, EmployeeType } from '@/modules/admin/employees/models/employees.models'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { Kanban, Store, TvMinimal } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const {currentEmployee} = useEmployeeAuthStore()

const onDisplayClick = () => {
	router.push({name: getRouteName("EVENT_ORDERS_DISPLAY")})
}

const onKioskClick = () => {
	router.replace({name: getRouteName("KIOSK_LANDING")})
}

const onBaristaClick = () => {
	router.push({name: getRouteName("EVENT_ORDERS_BARISTA")})
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
	<div class="flex bg-white pt-safe w-full min-h-screen">
		<div class="hidden md:block relative bg-muted/40 border-r">
			<div class="top-0 sticky flex flex-col gap-2 h-full max-h-screen">
				<div class="flex items-center px-4 lg:px-6 border-b h-14 lg:h-[60px]">
					<router-link
						:to='{name: getRouteName("ADMIN_DASHBOARD")}'
						class="flex items-center gap-3 font-semibold"
					>
						<div class="flex items-center gap-2">
							<span class="font-bold text-xl">Zeep</span>
							<span class="font-normal text-base">|</span>
							<span class="font-normal">
								{{ layoutLabel }}
							</span>
						</div>
					</router-link>
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
						<Kanban class="w-5 h-5" />
						<p class="text-base">Заказы</p>
					</Button>

					<Button
						variant="outline"
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
				<RouterView v-slot="{ Component, route }">
					<template v-if="Component">
						<Transition
							name="fade"
							mode="out-in"
						>
							<Suspense>
								<div :key="route.matched[0]?.path">
									<component :is="Component" />
								</div>
								<template #fallback> <PageLoader /> </template>
							</Suspense>
						</Transition>
					</template>
				</RouterView>
			</main>
		</div>
	</div>
</template>
