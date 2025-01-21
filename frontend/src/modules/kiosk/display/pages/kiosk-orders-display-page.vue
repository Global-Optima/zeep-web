<template>
	<div class="relative flex flex-col items-center bg-white h-screen overflow-hidden">
		<!-- Screen Content Wrapper -->
		<div class="flex items-start w-full h-full">
			<!-- Preparing Side -->
			<div class="flex flex-col flex-1 bg-white p-6 h-full">
				<div class="mb-6 text-center">
					<h2 class="font-semibold text-3xl text-gray-800 2xl:text-4xl">Preparing</h2>
				</div>

				<div class="flex-1">
					<TransitionGroup
						name="fade"
						tag="div"
						class="gap-2 grid grid-cols-4"
						aria-live="polite"
						type="transition"
					>
						<div
							v-for="(order, index) in currentPreparingPage"
							:key="`preparing-${index}`"
							class="flex justify-center items-center bg-blue-100 py-6 rounded-md font-bold text-3xl text-blue-900 lg:text-4xl 2xl:text-6xl"
						>
							{{ order }}
						</div>
					</TransitionGroup>
				</div>

				<div class="flex justify-center mt-auto">
					<div class="flex space-x-2">
						<span
							v-for="i in totalPreparingPages"
							:key="`prep-indicator-${i}`"
							:class="[
								'inline-block w-5 h-5 rounded-full',
								i - 1 === preparingPageIndex ? 'bg-blue-900' : 'bg-gray-300'
							]"
						></span>
					</div>
				</div>
			</div>

			<!-- Middle Separator -->
			<div class="bg-gray-100 my-auto w-0.5 h-[95vh]"></div>

			<!-- Ready Side -->
			<div class="flex flex-col flex-1 bg-white p-6 h-full">
				<div class="mb-6 text-center">
					<h2 class="font-semibold text-3xl text-gray-800 2xl:text-4xl">Ready</h2>
				</div>
				<div class="flex-1">
					<TransitionGroup
						name="fade"
						tag="div"
						class="gap-2 grid grid-cols-4"
						aria-live="polite"
						type="transition"
					>
						<div
							v-for="(order, index) in currentReadyPage"
							:key="`ready-${index}`"
							class="flex justify-center items-center bg-emerald-100 py-6 rounded-md font-bold text-3xl text-emerald-900 lg:text-4xl 2xl:text-6xl"
						>
							{{ order }}
						</div>
					</TransitionGroup>
				</div>
				<div class="flex justify-center mt-auto">
					<div class="flex space-x-2">
						<span
							v-for="i in totalReadyPages"
							:key="`ready-indicator-${i}`"
							:class="[
								'inline-block w-5 h-5 rounded-full',
								i - 1 === readyPageIndex ? 'bg-emerald-900' : 'bg-gray-300'
							]"
						></span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

// Constants
const ORDERS_PER_PAGE = 20; // Fixed number of orders per page
const TRANSITION_INTERVAL = 5000; // Transition interval in milliseconds

// Mock Data
const preparingOrders = ref<number[]>([145, 111, 561, 165, 651, 145, 312, 452, 412, 541, 674, 189, 205, 621, 145, 111, 561, 165, 651, 145, 312, 452, 412, 541, 674, 189, 205, 621,]);
const readyOrders = ref<number[]>([123, 145, 712, 432, 615, 153, 143, 251, 411, 564]);

// Pagination State
const preparingPageIndex = ref(0);
const readyPageIndex = ref(0);

// Computed Pages
const totalPreparingPages = computed(() => Math.ceil(preparingOrders.value.length / ORDERS_PER_PAGE));
const totalReadyPages = computed(() => Math.ceil(readyOrders.value.length / ORDERS_PER_PAGE));

const currentPreparingPage = computed(() =>
	preparingOrders.value.slice(
		preparingPageIndex.value * ORDERS_PER_PAGE,
		(preparingPageIndex.value + 1) * ORDERS_PER_PAGE
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
		// Update Preparing Page
		preparingPageIndex.value = (preparingPageIndex.value + 1) % totalPreparingPages.value;

		// Update Ready Page
		readyPageIndex.value = (readyPageIndex.value + 1) % totalReadyPages.value;
	}, TRANSITION_INTERVAL);
}

// Lifecycle Hook
onMounted(() => {
	if (preparingOrders.value.length > ORDERS_PER_PAGE || readyOrders.value.length > ORDERS_PER_PAGE) {
		transitionPages();
	}
});
</script>

<style>
.fade-enter-active, .fade-leave-active {
	transition: opacity 0.3s ease-in;
}
.fade-enter-from, .fade-leave-to {
	opacity: 0;
}
</style>
