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
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import { useEmployeeAuthStore } from '@/modules/auth/store/employee-auth.store'
import { EmployeeRoles } from '@/modules/employees/models/employees.models'
import { Apple, ChartBar, ListPlus, Package, ShoppingCart, Store, Users, Warehouse, type LucideIcon } from 'lucide-vue-next'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

interface NavItem {
  name: string
  routeKey: RouteKey
  icon: LucideIcon
  badge?: number
  accessRoles: EmployeeRoles[]
}

const route = useRoute()
const currentEmployeeStore = useEmployeeAuthStore()
const { currentEmployee } = storeToRefs(currentEmployeeStore)

const isActiveRoute = (routeKey: RouteKey) => {
  const routeName = getRouteName(routeKey)
  return route.name === routeName || (route.name && route.name.toString().startsWith(routeName))
}

const canAccessRoute = (accessRoles: EmployeeRoles[]) => {
  const employeeRole = currentEmployee.value?.role
  return employeeRole ? accessRoles.includes(employeeRole) : false
}

const navItems: NavItem[] = [
  { name: 'Аналитика', routeKey: 'ADMIN_DASHBOARD', icon: ChartBar, accessRoles: [EmployeeRoles.ADMIN, EmployeeRoles.DIRECTOR] },
  { name: 'Аналитика магазина', routeKey: 'ADMIN_STORE_DASHBOARD', icon: ChartBar, accessRoles: [EmployeeRoles.MANAGER] },
  { name: 'Заказы магазина', routeKey: 'ADMIN_STORE_ORDERS', icon: ShoppingCart, accessRoles: [EmployeeRoles.MANAGER, EmployeeRoles.EMPLOYEE] },
  { name: 'Все товары', routeKey: 'ADMIN_PRODUCTS', icon: Package, accessRoles: [EmployeeRoles.ADMIN, EmployeeRoles.DIRECTOR] },
  { name: 'Все топпинги', routeKey: 'ADMIN_ADDITIVES', icon: ListPlus, accessRoles: [EmployeeRoles.ADMIN, EmployeeRoles.DIRECTOR] },
  { name: 'Товары магазина', routeKey: 'ADMIN_STORE_PRODUCTS', icon: Package, accessRoles: [EmployeeRoles.MANAGER, EmployeeRoles.EMPLOYEE] },
  { name: 'Склад магазина', routeKey: 'ADMIN_STORE_WAREHOUSE', icon: Warehouse, accessRoles: [EmployeeRoles.MANAGER, EmployeeRoles.EMPLOYEE] },
  { name: 'Сотрудники', routeKey: 'ADMIN_EMPLOYEES', icon: Users, accessRoles: [EmployeeRoles.MANAGER] },
  { name: 'Ингредиенты', routeKey: "ADMIN_INGREDIENTS", icon: Apple, accessRoles: [EmployeeRoles.ADMIN, EmployeeRoles.DIRECTOR]  },
  { name: 'Магазины', routeKey: 'ADMIN_STORES', icon: Store, accessRoles: [EmployeeRoles.ADMIN, EmployeeRoles.DIRECTOR] },
]


const accessibleNavItems = computed(() =>
  navItems.filter(item => canAccessRoute(item.accessRoles))
)
</script>
