<template>
	<AdminStoreOrdersToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!storeOrders || storeOrders.data.length === 0"
					class="text-muted-foreground"
				>
					Заказы не найдены
				</p>

				<AdminStoreOrdersList
					v-else
					:orders="storeOrders.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="storeOrders"
					:meta="storeOrders.pagination"
					@update:page="updatePage"
					@update:pageSize="updatePageSize"
				/>
			</CardFooter>
		</Card>
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import {Card, CardContent} from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import {usePaginationFilter} from '@/core/hooks/use-pagination-filter.hook'
import AdminStoreOrdersList from '@/modules/admin/store-orders/components/list/admin-store-orders-list.vue'
import AdminStoreOrdersToolbar from '@/modules/admin/store-orders/components/list/admin-store-orders-toolbar.vue'
import type {OrdersFilterQuery} from '@/modules/admin/store-orders/models/orders.models'
import {ordersService} from '@/modules/admin/store-orders/services/orders.service'
import {useQuery} from '@tanstack/vue-query'
import {computed} from 'vue'
import {useHasRole} from "@/core/hooks/use-has-roles.hook";
import {EmployeeRole} from "@/modules/admin/employees/models/employees.models";

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<OrdersFilterQuery>({})
const isFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])
const isStore = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA])

const { data: storeOrders, isLoading } = useQuery({
  queryKey: computed(() => ['store-orders', filter.value]),
  queryFn: () => ordersService.getAllOrders(filter.value),
  enabled: computed(() =>
    isStore.value || (isFranchisee.value && Boolean(filter.value.storeId))
  )
})
</script>

<style scoped></style>
