<template>
	<header
		class="top-0 left-0 z-20 sticky flex justify-between items-center bg-white p-4 border-b w-full overflow-x-auto no-scrollbar"
	>
		<Button
			size="icon"
			variant="outline"
			@click="onBackClick"
		>
			<ChevronLeft />
		</Button>

		<!-- Status Filter Buttons -->
		<div class="flex items-center">
			<button
				v-for="status in statuses"
				:key="status.label"
				@click="selectStatus(status)"
				:class="statusButtonClasses(status)"
			>
				<p>{{ status.label }}</p>
				<p :class="statusCountClasses(status)">
					{{ status.count }}
				</p>
			</button>
		</div>

		<Button
			size="icon"
			variant="outline"
			@click="onReloadClick"
		>
			<RefreshCcw class="size-5" />
		</Button>
	</header>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { ChevronLeft, RefreshCcw } from 'lucide-vue-next'
import { toRefs } from 'vue'

interface Status {
  label: string;
  count: number;
}

/**
 * Define props & emits using script setup
 */
const props = defineProps<{
  statuses: Status[];
  selectedStatus: Status;
}>()

const emits = defineEmits<{
  (e: 'selectStatus', status: Status): void;
  (e: 'back'): void;
  (e: 'reload'): void;
}>()

const { statuses, selectedStatus } = toRefs(props)

/**
 * Emit events
 */
function selectStatus(status: Status) {
  emits('selectStatus', status)
}

function onBackClick() {
  emits('back')
}

function onReloadClick() {
  emits('reload')
}

/**
 * Helper to style each status button:
 */
function statusButtonClasses(status: Status) {
  return [
    'flex items-center gap-2 px-5 py-2 rounded-xl text-base whitespace-nowrap',
    status.label === selectedStatus.value.label ? 'bg-primary text-primary-foreground' : ''
  ]
}

/**
 * Helper to style the count badge:
 */
function statusCountClasses(status: Status) {
  return [
    'bg-gray-100 px-2 py-1 rounded-sm text-black text-xs',
    status.label === selectedStatus.value.label ? 'bg-green-700 text-primary-foreground' : ''
  ]
}
</script>
