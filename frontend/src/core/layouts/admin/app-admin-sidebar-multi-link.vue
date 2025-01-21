<template>
	<Collapsible v-model:open="isOpen">
		<CollapsibleTrigger class="p-0 w-full">
			<div
				class="flex justify-between items-center gap-4 px-3 py-2 rounded-lg"
				:class="{
					'bg-muted text-foreground': isActiveParent,
					'text-muted-foreground hover:text-foreground': !isActiveParent,
				}"
			>
				<div class="flex items-center gap-4">
					<component
						:is="item.icon"
						class="w-5 h-5"
					/>
					<span class="text-left">{{ item.label }}</span>
				</div>
				<CaretSortIcon class="w-4 h-4 text-gray-500" />
			</div>
		</CollapsibleTrigger>
		<CollapsibleContent>
			<div class="ml-4 py-1">
				<router-link
					v-for="child in item.items"
					:key="child.routeKey"
					:to="{ name: getRouteName(child.routeKey) }"
					class="flex items-center gap-4 px-3 py-2 rounded-lg"
					:class="{
						'bg-muted text-foreground': isActiveRoute(child.routeKey),
						'text-muted-foreground hover:text-foreground': !isActiveRoute(child.routeKey),
					}"
				>
					{{ child.name }}
				</router-link>
			</div>
		</CollapsibleContent>
	</Collapsible>
</template>

<script setup lang="ts">
import {
    Collapsible,
    CollapsibleContent,
    CollapsibleTrigger,
} from '@/core/components/ui/collapsible'
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import type { CollapsibleNavItem } from '@/core/layouts/admin/app-admin-sidebar-links.constants'
import { CaretSortIcon } from '@radix-icons/vue'
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute();
const { item } = defineProps<{ item: CollapsibleNavItem }>();

const isOpen = ref(false);

// Check if the parent is active
const isActiveParent = computed(() => {
	return item.items.some(child => isActiveRoute(child.routeKey));
});

const isActiveRoute = (routeKey: RouteKey) => {
	const routeName = getRouteName(routeKey);
	return (
		route.name === routeName || (route.name && route.name.toString().startsWith(routeName))
	);
};
</script>
