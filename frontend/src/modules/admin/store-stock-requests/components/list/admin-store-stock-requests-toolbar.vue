<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Filter Menu -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
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
					<DropdownMenuItem @click="applyFilter(undefined)"> Все </DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter(StoreStockRequestStatus.CREATED)">
						Созданные
					</DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter(StoreStockRequestStatus.IN_DELIVERY)">
						В доставке
					</DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter(StoreStockRequestStatus.COMPLETED)">
						Завершённые
					</DropdownMenuItem>
					<DropdownMenuItem @click="applyFilter(StoreStockRequestStatus.REJECTED)">
						Отклонённые
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button variant="outline">Экспорт</Button>
			<Button @click="addStore">Добавить</Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/core/components/ui/dropdown-menu'
import { getRouteName } from '@/core/config/routes.config'
import { type GetStoreStockRequestsFilter, StoreStockRequestStatus } from '@/modules/admin/store-stock-requests/models/store-stock-request.model'
import { ChevronDown } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

// Props and Emit
const props = defineProps<{ filter: GetStoreStockRequestsFilter }>();
const emit = defineEmits(['update:filter']);

// Filter Selection
const applyFilter = (status: StoreStockRequestStatus | undefined) => {
  emit('update:filter', { ...props.filter, status });
};

// Add Store Navigation
const router = useRouter();
const addStore = () => {
  router.push({ name: getRouteName('ADMIN_STORE_STOCK_REQUESTS_CREATE') });
};
</script>
