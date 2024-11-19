<script setup lang="ts">
import { Home, LineChart, Package, ShoppingCart, Users } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const currentRouteName = computed(() => route.name)

const isActiveRoute = (routeName: string) => {
	return currentRouteName.value?.toString().includes(routeName) ?? false
}
</script>

<template>
	<nav class="items-start gap-2 grid sm:px-2 lg:px-4 font-medium text-lg sm:text-sm">
		<a
			href="#"
			class="flex items-center gap-2 font-semibold text-lg"
		>
			<Package2 class="w-6 h-6" />
			<span class="sr-only">Zeep</span>
		</a>
		<a
			href="#"
			class="flex items-center gap-4 mx-[-0.65rem] px-3 py-2 rounded-xl text-muted-foreground hover:text-foreground"
		>
			<Home class="w-5 h-5" />
			Главная
		</a>
		<a
			href="#"
			class="flex items-center gap-4 bg-muted mx-[-0.65rem] px-3 py-2 rounded-xl text-foreground hover:text-foreground"
		>
			<ShoppingCart class="w-5 h-5" />
			Заказы
			<Badge class="flex justify-center items-center ml-auto rounded-full w-6 h-6 shrink-0">
				6
			</Badge>
		</a>
		<a
			href="#"
			class="flex items-center gap-4 mx-[-0.65rem] px-3 py-2 rounded-xl text-muted-foreground hover:text-foreground"
		>
			<Package class="w-5 h-5" />
			Товары
		</a>
		<a
			href="#"
			class="flex items-center gap-4 mx-[-0.65rem] px-3 py-2 rounded-xl text-muted-foreground hover:text-foreground"
		>
			<Users class="w-5 h-5" />
			Покупатели
		</a>
		<a
			href="#"
			class="flex items-center gap-4 mx-[-0.65rem] px-3 py-2 rounded-xl text-muted-foreground hover:text-foreground"
		>
			<LineChart class="w-5 h-5" />
			Аналитика
		</a>
	</nav>
</template>

<style lang="scss" scoped></style>
