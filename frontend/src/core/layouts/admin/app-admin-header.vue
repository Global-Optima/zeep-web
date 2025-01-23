<template>
	<header class="flex items-center gap-4 bg-muted/40 px-4 lg:px-6 border-b w-full h-14 lg:h-[60px]">
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

				<div class="hidden mt-auto">
					<Button
						variant="outline"
						@click="onDisplayClick"
						class="flex items-center gap-3 mb-1 py-6 w-full"
					>
						<TvMinimal class="w-5 h-5" />
						<p class="text-base">Экран</p>
					</Button>

					<Button
						variant="outline"
						@click="onBaristaClick"
						class="flex items-center gap-3 mb-1 py-6 w-full"
					>
						<Store class="w-5 h-5" />
						<p class="text-base">Бариста</p>
					</Button>

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

		<div
			v-if="currentEmployee"
			class="flex justify-end items-center gap-4 w-full"
		>
			<AppAdminSearch :current-employee="currentEmployee" />
			<AppAdminEmployeeDropdown :current-employee="currentEmployee" />
		</div>
	</header>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Sheet, SheetContent, SheetTrigger } from '@/core/components/ui/sheet'
import { getRouteName } from '@/core/config/routes.config'
import AppAdminEmployeeDropdown from '@/core/layouts/admin/app-admin-employee-dropdown.vue'
import AppAdminSearch from '@/core/layouts/admin/app-admin-search.vue'
import AppAdminSidebar from '@/core/layouts/admin/app-admin-sidebar.vue'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { Menu, Store, TvMinimal } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const {currentEmployee} = useEmployeeAuthStore()

const router = useRouter()

const onKioskClick = () => {
	router.push({name: getRouteName('KIOSK_HOME')})
}

const onBaristaClick = () => {
	router.push({name: getRouteName('KIOSK_ORDERS')})
}

const onDisplayClick = () => {
	router.push({name: getRouteName('KIOSK_ORDERS_DISPLAY')})
}
</script>

<style scoped></style>
