<script setup lang="ts">
import { getRouteName } from '@/core/config/routes.config'
import { Home, LineChart, Package, ShoppingCart, Users, Warehouse, type LucideIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

interface NavItem {
  name: string
  routeName: string
  icon: LucideIcon
  badge?: number
}

const route = useRoute()
const currentRouteName = computed(() => route.name)

const isActiveRoute = (routeName: string) => {
  return currentRouteName.value?.toString().includes(routeName) ?? false
}

const navItems: NavItem[] = [
  {
    name: 'Главная',
    routeName: getRouteName("ADMIN_DASHBOARD"),
    icon: Home,
  },
  {
    name: 'Заказы',
    routeName: getRouteName("ADMIN_ORDERS"),
    icon: ShoppingCart,
    badge: 6,
  },
  {
    name: 'Товары',
    routeName: getRouteName("ADMIN_PRODUCTS"),
    icon: Package,
  },
  {
    name: 'Склад',
    routeName: getRouteName("ADMIN_WAREHOUSE"),
    icon: Warehouse,
  },
  {
    name: 'Аналитика',
    routeName: getRouteName("ADMIN_ANALYTICS"),
    icon: LineChart,
  },
  {
    name: 'Сотрудники',
    routeName: getRouteName("ADMIN_EMPLOYEES"),
    icon: Users,
  },
]
</script>

<template>
	<nav class="items-start gap-2 grid sm:px-2 lg:px-4 font-medium text-lg sm:text-sm">
		<div
			v-for="item in navItems"
			:key="item.routeName"
		>
			<router-link
				:to="{name: item.routeName}"
				class="flex items-center gap-4 mx-[-0.65rem] px-3 py-2 rounded-lg hover:text-foreground"
				:class="{
          'bg-muted text-foreground': isActiveRoute(item.routeName),
          'text-muted-foreground': !isActiveRoute(item.routeName),
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

<style scoped></style>
