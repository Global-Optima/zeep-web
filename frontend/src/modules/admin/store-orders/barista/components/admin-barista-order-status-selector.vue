<template>
	<header
		class="top-0 left-0 z-20 sticky flex justify-between items-center bg-white p-4 border-b w-full overflow-x-auto no-scrollbar"
	>
		<!-- Back Button -->
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

		<!-- Reload Button -->

		<div class="flex items-center gap-2">
			<AdminBaristaQrSizesDialog />

			<Button
				size="icon"
				variant="outline"
				@click="onReloadClick"
			>
				<RefreshCcw class="size-5" />
			</Button>
		</div>
	</header>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import AdminBaristaQrSizesDialog from '@/modules/admin/store-orders/barista/components/admin-barista-qr-sizes-dialog.vue'
import { ChevronLeft, RefreshCcw } from 'lucide-vue-next'
import { toRefs } from 'vue'

interface StatusOption {
  label: string;
  count: number;
  // If needed, you could add: status?: OrderStatus
}

/**
 * Define our props.
 * - `statuses`: an array of statuses, each with a `label` and `count`.
 * - `selectedStatus`: the currently active/selected status filter.
 */
const props = defineProps<{
  statuses: StatusOption[];
  selectedStatus: StatusOption;
}>()

/**
 * Define our emits:
 * - `selectStatus`: user clicks a status button
 * - `back`: user clicks the Back button
 * - `reload`: user clicks the Reload button
 */
const emits = defineEmits<{
  (e: 'selectStatus', status: StatusOption): void;
  (e: 'back'): void;
  (e: 'reload'): void;
}>()

const { statuses, selectedStatus } = toRefs(props)

/** Emit an event when a status is chosen */
function selectStatus(status: StatusOption) {
  emits('selectStatus', status)
}

/** Emit the back event (to navigate away) */
function onBackClick() {
  emits('back')
}

/** Emit the reload event (to refresh data/page) */
function onReloadClick() {
  emits('reload')
}

/**
 * Helper to style each status button,
 * highlighting the selected one.
 */
function statusButtonClasses(status: StatusOption) {
  return [
    'flex items-center gap-2 px-5 py-2 rounded-xl text-base whitespace-nowrap',
    status.label === selectedStatus.value.label ? 'bg-primary text-primary-foreground' : ''
  ]
}

/**
 * Helper to style the count badge,
 * also highlighting it if selected.
 */
function statusCountClasses(status: StatusOption) {
  return [
    'bg-gray-100 px-2 py-1 rounded-sm text-black text-xs',
    status.label === selectedStatus.value.label ? 'bg-green-700 text-primary-foreground' : ''
  ]
}
</script>
