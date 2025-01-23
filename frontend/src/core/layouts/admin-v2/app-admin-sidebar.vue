<template>
	<nav class="items-start gap-2 grid px-2 font-medium text-lg sm:text-sm">
		<div
			v-for="item in accessibleNavItems"
			:key="isCollapsibleNavItem(item) ? item.label : item.routeKey"
		>
			<AppAdminSidebarMultiLink
				v-if="isCollapsibleNavItem(item)"
				:item="item"
			/>
			<router-link
				v-else
				:to="{ name: getRouteName(item.routeKey) }"
				class="flex items-center gap-4 px-3 py-2 rounded-lg hover:text-foreground"
				:class="{
					'bg-muted text-foreground': isActiveRoute(item.routeKey),
					'text-muted-foreground': !isActiveRoute(item.routeKey),
				}"
			>
				<component
					:is="item.icon"
					class="w-5 h-5"
				/>
				{{ item.name }}
			</router-link>
		</div>
	</nav>
</template>

<script setup lang="ts">
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import { adminNavItems, isCollapsibleNavItem } from '@/core/layouts/admin-v2/app-admin-sidebar-links.constants'
import AppAdminSidebarMultiLink from '@/core/layouts/admin-v2/app-admin-sidebar-multi-link.vue'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute();
const { currentEmployee } = useEmployeeAuthStore();

const isActiveRoute = (routeKey: RouteKey) => {
	const routeName = getRouteName(routeKey);
	return (
		route.name === routeName || (route.name && route.name.toString().startsWith(routeName))
	);
};

const canAccessRoute = (accessRoles: string[]) => {
	if (!currentEmployee) return false;
	return currentEmployee.role ? accessRoles.includes(currentEmployee.role) : false;
};

const accessibleNavItems = computed(() => {
	if (!currentEmployee) return [];
	return adminNavItems.filter(item => {
		if (isCollapsibleNavItem(item)) {
			return canAccessRoute(item.accessRoles);
		}
		return canAccessRoute(item.accessRoles);
	});
});
</script>
