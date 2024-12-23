<template>
	<div class="flex items-center space-x-6 lg:space-x-8">
		<!-- Rows per page dropdown -->
		<div class="flex items-center space-x-2">
			<p class="font-medium text-sm">Cтрок на странице</p>
			<select
				class="border-gray-300 border rounded-md w-[70px] h-8 text-center text-sm"
				v-model="selectedPageSize"
			>
				<option
					v-for="size in pageSizeOptions"
					:key="size"
					:value="size"
				>
					{{ size }}
				</option>
			</select>
		</div>

		<!-- Page count display -->
		<div class="flex justify-center items-center w-[100px] font-medium text-nowrap text-sm">
			<span v-if="totalPages > 0">Страница {{ currentPage }} из {{ totalPages }}</span>
			<span v-else>Нет данных</span>
		</div>

		<!-- Navigation buttons -->
		<div class="flex items-center space-x-2">
			<!-- First page button -->
			<Button
				variant="outline"
				class="lg:flex hidden p-0 w-8 h-8"
				:disabled="totalPages === 0 || currentPage === 1"
				@click="handlePageChange(1)"
			>
				<span class="sr-only">Go to first page</span>
				<DoubleArrowLeftIcon class="w-4 h-4" />
			</Button>

			<!-- Previous page button -->
			<Button
				variant="outline"
				class="p-0 w-8 h-8"
				:disabled="totalPages === 0 || currentPage === 1"
				@click="handlePageChange(currentPage - 1)"
			>
				<span class="sr-only">Go to previous page</span>
				<ChevronLeftIcon class="w-4 h-4" />
			</Button>

			<!-- Next page button -->
			<Button
				variant="outline"
				class="p-0 w-8 h-8"
				:disabled="totalPages === 0 || currentPage === totalPages"
				@click="handlePageChange(currentPage + 1)"
			>
				<span class="sr-only">Go to next page</span>
				<ChevronRightIcon class="w-4 h-4" />
			</Button>

			<!-- Last page button -->
			<Button
				variant="outline"
				class="lg:flex hidden p-0 w-8 h-8"
				:disabled="totalPages === 0 || currentPage === totalPages"
				@click="handlePageChange(totalPages)"
			>
				<span class="sr-only">Go to last page</span>
				<DoubleArrowRightIcon class="w-4 h-4" />
			</Button>
		</div>
	</div>
</template>
<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import type { PaginationMeta } from '@/core/utils/pagination.utils'
import { DoubleArrowLeftIcon, DoubleArrowRightIcon } from '@radix-icons/vue'
import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

// Props: Accepts PaginationMeta
const props = defineProps<{
  meta: PaginationMeta,
}>();

const emit = defineEmits<{
  (e: 'update:page', page: number): void,
  (e: 'update:pageSize', pageSize: number): void
}>();

const pageSizeOptions = [10, 20, 30, 40, 50];
const selectedPageSize = ref(props.meta.pageSize);

const totalPages = computed(() => props.meta.totalPages || 0); // Ensure totalPages is at least 0
const currentPage = computed(() => Math.min(props.meta.page, totalPages.value)); // Prevent out-of-range pages

watch(selectedPageSize, (newPageSize) => {
  emit('update:pageSize', Number(newPageSize));
});

const handlePageChange = (newPage: number) => {
  if (newPage >= 1 && newPage <= totalPages.value) {
    emit('update:page', newPage);
  }
};
</script>
