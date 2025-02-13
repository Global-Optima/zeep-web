<template>
	<div :class="['flex flex-col flex-1 gap-8 px-12 p-6 h-full', classNames]">
		<div class="text-left">
			<h2 class="font-medium text-3xl text-slate-100 2xl:text-6xl">{{ title }}</h2>
		</div>

		<div class="flex-1 mt-2 overflow-auto">
			<!-- Transition for the orders -->
			<transition
				name="fade"
				mode="out-in"
				type="transition"
			>
				<ul
					v-if="pagedOrders.length"
					:key="currentPageIndex"
					class="gap-4 grid"
				>
					<li
						v-for="(item, index) in pagedOrders"
						:key="item.id || index"
						class="flex items-center gap-5 2xl:gap-10 text-3xl lg:text-4xl 2xl:text-5xl"
					>
						<span
							:class="cn('flex justify-center items-center p-4 rounded-md w-24 2xl:w-32 font-bold', type === 'PREPARING' && 'bg-blue-100 text-blue-900', type === 'COMPLETED' && 'bg-emerald-100 text-emerald-900')"
						>
							{{ item.displayNumber }}
						</span>
						<span class="font-medium text-slate-100">{{ item.customerName }}</span>
					</li>
				</ul>
			</transition>
		</div>

		<div
			v-if="totalPages > 1"
			class="flex justify-center mt-auto"
		>
			<div class="flex space-x-2">
				<span
					v-for="i in totalPages"
					:key="`indicator-${i}`"
					:class="[
						'inline-block w-3 h-3 rounded-full',
						i - 1 === currentPageIndex ? 'bg-blue-100' : 'bg-slate-700'
					]"
					style="cursor: pointer;"
					@click="$emit('pageChange', i - 1)"
				></span>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { cn } from '@/core/utils/tailwind.utils'
import type { OrderDTO } from '@/modules/admin/store-orders/models/orders.models'
import { computed } from 'vue'

// Props
const { orders, currentPageIndex, class: classNames } = defineProps<{
	title: string
	orders: OrderDTO[]
	currentPageIndex: number
	totalPages: number
	class: string
	type: 'PREPARING' | 'COMPLETED'
}>()

// Pagination Constants
const ORDERS_PER_PAGE = 6

// Paginated Orders
const pagedOrders = computed(() =>
	orders.slice(
		currentPageIndex * ORDERS_PER_PAGE,
		(currentPageIndex + 1) * ORDERS_PER_PAGE
	)
)
</script>

<style scoped>
/* Fade transition styles */
.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.5s ease;
}
.fade-enter-from,
.fade-leave-to {
	opacity: 0;
}
</style>
