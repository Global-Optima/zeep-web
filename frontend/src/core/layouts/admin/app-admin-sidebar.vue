<template>
	<nav class="items-start gap-2 grid sm:px-2 lg:px-4 font-medium text-lg sm:text-sm">
		<div
			v-for="item in accessibleNavItems"
			:key="item.routeKey"
		>
			<router-link
				:to="{ name: getRouteName(item.routeKey) }"
				class="flex items-center gap-4 mx-[-0.65rem] px-3 py-2 rounded-lg hover:text-foreground"
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

				<Badge
					v-if="item.badge"
					class="flex justify-center items-center ml-auto rounded-full w-6 h-6 shrink-0"
				>
					{{ item.badge }}
				</Badge>
			</router-link>
		</div>
	</nav>
</template>

<script setup lang="ts">
import { Badge } from '@/core/components/ui/badge'
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { EmployeeRole } from '@/modules/employees/models/employees.models'
import {
  Apple,
  ChartBar,
  FileBox,
  LayoutList,
  ListPlus,
  Package,
  ShoppingCart,
  Store,
  Users,
  Warehouse,
  type LucideIcon,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

interface NavItem {
  name: string;
  routeKey: RouteKey;
  icon: LucideIcon;
  badge?: number;
  accessRoles: EmployeeRole[];
}

const route = useRoute();
const { currentEmployee } = useEmployeeAuthStore();

const isActiveRoute = (routeKey: RouteKey) => {
  const routeName = getRouteName(routeKey);
  return (
    route.name === routeName || (route.name && route.name.toString().startsWith(routeName))
  );
};

const canAccessRoute = (accessRoles: EmployeeRole[]) => {
  if (!currentEmployee) return false;
  return currentEmployee.role ? accessRoles.includes(currentEmployee.role) : false;
};

// Navigation items
const navItems: NavItem[] = [
  { name: 'Аналитика', routeKey: 'ADMIN_DASHBOARD', icon: ChartBar, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
  { name: 'Аналитика магазина', routeKey: 'ADMIN_STORE_DASHBOARD', icon: ChartBar, accessRoles: [EmployeeRole.MANAGER] },
  { name: 'Заказы магазина', routeKey: 'ADMIN_STORE_ORDERS', icon: ShoppingCart, accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA] },
  { name: 'Товары', routeKey: 'ADMIN_PRODUCTS', icon: Package, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
  { name: 'Топпинги', routeKey: 'ADMIN_ADDITIVES', icon: ListPlus, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
  { name: 'Товары магазина', routeKey: 'ADMIN_STORE_PRODUCTS', icon: Package, accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA] },
  { name: 'Склад магазина', routeKey: 'ADMIN_STORE_STOCKS', icon: Warehouse, accessRoles: [EmployeeRole.MANAGER, EmployeeRole.BARISTA] },
  { name: 'Сотрудники', routeKey: 'ADMIN_EMPLOYEES', icon: Users, accessRoles: [EmployeeRole.MANAGER] },
  { name: 'Ингредиенты', routeKey: 'ADMIN_INGREDIENTS', icon: Apple, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
  { name: 'Магазины', routeKey: 'ADMIN_STORES', icon: Store, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
  { name: 'Поставщики', routeKey: 'ADMIN_SUPPLIERS', icon: FileBox, accessRoles: [EmployeeRole.WAREHOUSE, EmployeeRole.DIRECTOR] },
  { name: 'Категории товаров', routeKey: 'ADMIN_PRODUCT_CATEGORIES', icon: LayoutList, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
  { name: 'Категории топпингов', routeKey: 'ADMIN_ADDITIVE_CATEGORIES', icon: LayoutList, accessRoles: [EmployeeRole.ADMIN, EmployeeRole.DIRECTOR] },
];

// Computed list of accessible navigation items
const accessibleNavItems = computed(() => {
  if (!currentEmployee) return [];
  return navItems.filter(item => canAccessRoute(item.accessRoles));
});
</script>
