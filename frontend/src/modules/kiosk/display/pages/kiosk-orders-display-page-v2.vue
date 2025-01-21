<template>
	<div class="flex flex-col items-center h-screen overflow-hidden">
		<!-- Screen Content Wrapper -->
		<div class="flex items-start w-full h-full">
			<!-- In Progress Side -->
			<div class="flex flex-col flex-1 gap-8 bg-slate-800 px-12 p-6 h-full">
				<div class="text-left">
					<h2 class="mb-2 font-medium text-3xl text-slate-100 2xl:text-6xl =">В работе</h2>
				</div>

				<div class="flex-1">
					<transition-group
						name="fade"
						tag="ul"
						class="gap-4 grid"
						aria-live="polite"
					>
						<li
							v-for="(order, index) in currentInProgressPage"
							:key="`in-progress-${index}`"
							class="flex items-center gap-5 text-3xl lg:text-4xl 2xl:text-5xl"
						>
							<span
								class="bg-blue-100 p-4 rounded-md w-fit font-bold text-blue-900"
								>{{ order.number }}</span
							>
							<span class="font-medium text-slate-100">{{ order.name }}</span>
						</li>
					</transition-group>
				</div>

				<div class="flex justify-center mt-auto">
					<div class="flex space-x-2">
						<span
							v-for="i in totalInProgressPages"
							:key="`in-progress-indicator-${i}`"
							:class="[
								'inline-block w-3 h-3 2xl:w-5 2xl:h-5 rounded-full',
								i - 1 === readyPageIndex ? 'bg-blue-100' : 'bg-slate-700'
							]"
						></span>
					</div>
				</div>
			</div>

			<!-- Ready Side -->
			<div class="flex flex-col flex-1 gap-8 bg-slate-900 px-12 p-6 h-full">
				<div class="text-left">
					<h2 class="mb-2 font-medium text-3xl text-slate-100 2xl:text-6xl">В Готовы</h2>
				</div>

				<div class="flex-1">
					<transition-group
						name="fade"
						tag="ul"
						class="gap-4 grid"
						aria-live="polite"
					>
						<li
							v-for="(order, index) in currentReadyPage"
							:key="`in-progress-${index}`"
							class="flex items-center gap-5 text-3xl lg:text-4xl 2xl:text-5xl"
						>
							<span
								class="bg-emerald-100 p-4 rounded-md w-fit font-bold text-emerald-900"
								>{{ order.number }}</span
							>
							<span class="font-medium text-slate-100">{{ order.name }}</span>
						</li>
					</transition-group>
				</div>
				<div class="flex justify-center mt-auto">
					<div class="flex space-x-2">
						<span
							v-for="i in totalReadyPages"
							:key="`ready-indicator-${i}`"
							:class="[
								'inline-block w-3 h-3 2xl:w-5 2xl:h-5 rounded-full',
								i - 1 === readyPageIndex ? 'bg-emerald-100' : 'bg-slate-700'
							]"
						></span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue"

// Constants
const ORDERS_PER_PAGE = 6; // Fixed number of orders per page
const TRANSITION_INTERVAL = 8000; // Transition interval in milliseconds

// Mock Data
const inProgressOrders = ref([
	{ number: 737, name: "Екатерина" },
	{ number: 741, name: "Валерия" },
	{ number: 745, name: "Катрин" },
	{ number: 749, name: "Кристина" },
	{ number: 752, name: "Олеся" },
	{ number: 755, name: "Михаил" },
	{ number: 758, name: "Анна" },
	{ number: 760, name: "Иван" },
	{ number: 762, name: "Елена" },
	{ number: 764, name: "Андрей" },
]);

const readyOrders = ref([
	{ number: 734, name: "Александр" },
	{ number: 738, name: "Ольга" },
	{ number: 742, name: "Сергей" },
	{ number: 746, name: "Юлия" },
	{ number: 750, name: "Мария" },
]);

// Pagination State
const inProgressPageIndex = ref(0);
const readyPageIndex = ref(0);

// Computed Pages
const totalInProgressPages = computed(() => Math.ceil(inProgressOrders.value.length / ORDERS_PER_PAGE));
const totalReadyPages = computed(() => Math.ceil(readyOrders.value.length / ORDERS_PER_PAGE));

const currentInProgressPage = computed(() =>
	inProgressOrders.value.slice(
		inProgressPageIndex.value * ORDERS_PER_PAGE,
		(inProgressPageIndex.value + 1) * ORDERS_PER_PAGE
	)
);

const currentReadyPage = computed(() =>
	readyOrders.value.slice(
		readyPageIndex.value * ORDERS_PER_PAGE,
		(readyPageIndex.value + 1) * ORDERS_PER_PAGE
	)
);

// Auto Transition
function transitionPages() {
	setInterval(() => {
		// Update In Progress Page
		inProgressPageIndex.value = (inProgressPageIndex.value + 1) % totalInProgressPages.value;

		// Update Ready Page
		readyPageIndex.value = (readyPageIndex.value + 1) % totalReadyPages.value;
	}, TRANSITION_INTERVAL);
}

// Lifecycle Hook
onMounted(() => {
	if (inProgressOrders.value.length > ORDERS_PER_PAGE || readyOrders.value.length > ORDERS_PER_PAGE) {
		transitionPages();
	}
});
</script>

<style>
.fade-enter-active, .fade-leave-active {
	transition: opacity 0.5s ease-in-out;
}
.fade-enter-from, .fade-leave-to {
	opacity: 0;
}
</style>
