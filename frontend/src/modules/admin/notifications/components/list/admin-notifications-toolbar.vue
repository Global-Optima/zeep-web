<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/core/components/ui/select'
import type { GetNotificationsFilter } from '@/modules/admin/notifications/models/notifications.model'
import { useDebounce } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

const props = defineProps<{ filter: GetNotificationsFilter }>();

const emit = defineEmits<{
  (e: 'update:filter', filter: GetNotificationsFilter): void,
}>();

// Local copy of the filter to avoid direct mutations
const localFilter = ref({ ...props.filter });

// Debounced search term to minimize excessive updates
const searchTerm = ref(localFilter.value.search || '');
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500);

// Watch debounced search term and update filter
watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue;
  emit('update:filter', { ...localFilter.value });
});

// Update isRead filter
const updateIsReadFilter = (value: string) => {
  localFilter.value.isRead = value === 'all' ? undefined : value === 'true';
  emit('update:filter', { ...localFilter.value });
};
</script>

<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Search Input and Filter -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<!-- Search Input -->
			<Input
				v-model="searchTerm"
				placeholder="Поиск уведомлений"
				type="search"
				class="bg-white w-full md:w-64"
			/>

			<!-- Select Filter -->
			<Select
				:model-value="localFilter.isRead === undefined ? 'all' : localFilter.isRead.toString()"
				@update:model-value="updateIsReadFilter"
			>
				<SelectTrigger class="w-fit">
					<SelectValue placeholder="Фильтр" />
				</SelectTrigger>
				<SelectContent>
					<SelectItem value="all">Все уведомления</SelectItem>
					<SelectItem value="false">Непрочитанные</SelectItem>
					<SelectItem value="true">Прочитанные</SelectItem>
				</SelectContent>
			</Select>
		</div>

		<!-- Right Side: Export and Add Notification Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
				>Экспорт</Button
			>
		</div>
	</div>
</template>
