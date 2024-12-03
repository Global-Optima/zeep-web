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

		<div class="flex items-center">
			<button
				v-for="status in statuses"
				:key="status.label"
				@click="selectStatus(status)"
				:class="[
          'flex items-center gap-2 px-5 py-2 rounded-xl text-lg whitespace-nowrap',
          status.label === selectedStatus.label ? 'bg-primary text-primary-foreground' : ''
        ]"
			>
				<p>{{ status.label }}</p>
				<p
					:class="[
            'bg-gray-100 px-2 py-1 rounded-sm text-black text-xs',
            status.label === selectedStatus.label ? 'bg-green-700 text-primary-foreground' : ''
          ]"
				>
					{{ status.count }}
				</p>
			</button>
		</div>

		<div></div>
	</header>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { ChevronLeft } from 'lucide-vue-next'

interface Status {
  label: string;
  count: number;
}

defineProps<{
  statuses: Status[];
  selectedStatus: Status;
}>();

const emits = defineEmits<{
  (e: 'selectStatus', status: Status): void;
  (e: 'back'): void;
}>();

const selectStatus = (status: Status) => {
  emits('selectStatus', status);
};

const onBackClick = () => {
  emits('back');
};
</script>
