<script setup lang="ts">
import { adminSidebarMenu } from "@/core/config/admin-sidebar.config";
import { getRoute, getRouteName } from '@/core/config/routes.config';
import { Icon } from '@iconify/vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
const currentRouteName = computed(() => route.name)

const isActiveRoute = (routeName: string) => {
	return currentRouteName.value?.toString().includes(routeName) ?? false
}
</script>

<template>
	<div class="h-full flex flex-col overflow-hidden">
		<button class="w-full flex items-center  gap-3 border-b border-border p-6">
			<p class="text-xl">Zeep</p>
		</button>

		<div class="h-full flex flex-col gap-3 mt-8 overflow-y-auto no-scrollbar p-6 pt-0">
			<template
				v-for="sidebarItem in adminSidebarMenu"
				:key="sidebarItem.routeName"
			>
				<router-link
					v-slot="{ href, navigate }"
					:to="{name: sidebarItem.routeName}"
					custom
				>
					<a
						class="flex cursor-pointer "
						:class='{
								"text-primary": isActiveRoute(sidebarItem.routeName),
								"text-foreground": !isActiveRoute(sidebarItem.routeName)
							}'
						:href="href"
						@click="navigate"
					>
						<div class="flex items-center gap-4">
							<Icon
								:icon="sidebarItem.icon"
								class="text-xl"
							/>

							<p class="text-sm">{{ getRoute(sidebarItem.routeName)?.meta.title }}</p>
						</div>
					</a>
				</router-link>
			</template>
		</div>

		<div class="flex-1 p-6">
			<router-link
				:to="{name: getRouteName('KIOSK_HOME')}"
				class="w-full flex gap-4 cursor-pointer bg-card items-center justify-center px-4 py-2 rounded-lg text-card-foreground"
			>
				<Icon
					icon="mingcute:pad-line"
					class="text-lg pi pi-angle-left"
				/>

				<p class="text-sm">Kiosk</p>
			</router-link>
		</div>
	</div>
</template>

<style lang="scss" scoped></style>
