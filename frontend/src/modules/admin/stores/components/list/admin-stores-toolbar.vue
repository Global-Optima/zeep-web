<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Search Input and Filter Menu -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<!-- Search Input -->
			<Input
				v-model="searchQueryRef"
				placeholder="Поиск"
				class="bg-white w-full md:w-64"
				@input="onSearchInput"
			/>
			<!-- Filter Menu -->
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button
						variant="outline"
						class="whitespace-nowrap"
					>
						Фильтр
						<ChevronDown class="ml-2 w-4 h-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem @click="applyFilter('all')">Все</DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter('active')">Активные</DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter('inactive')">Неактивные</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="ghost"
				class="w-full md:w-auto"
				@click="exportData"
			>
				<Download class="mr-2 w-4 h-4" />
				Экспорт
			</Button>
			<Button
				class="w-full md:w-auto"
				@click="addStore"
			>
				<Plus class="mr-2 w-4 h-4" />
				Добавить
			</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { ChevronDown, Download, Plus } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter();



// Emits
const emits = defineEmits<{
  (e: 'update:searchQuery', value: string): void;
  (e: 'filterChanged', filter: string): void;
  (e: 'exportData'): void;
  (e: 'addStore'): void;
}>();

// State variables
const searchQueryRef = ref("");

// Methods
const onSearchInput = () => {
  emits('update:searchQuery', searchQueryRef.value);
};

const applyFilter = (filter: string) => {
  emits('filterChanged', filter);
};

const exportData = () => {
  emits('exportData');
};

const addStore = () => {
  router.push({ name: getRouteName("ADMIN_STORE_CREATE") });
};
</script>
