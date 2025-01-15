<template>
	<AdminStoreOrdersToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminStoreOrdersList from '@/modules/admin/store-orders/components/list/admin-store-orders-list.vue'
import AdminStoreOrdersToolbar from '@/modules/admin/store-orders/components/list/admin-store-orders-toolbar.vue'
import type { OrdersFilterQuery } from '@/modules/orders/models/orders.models'
import { orderService } from '@/modules/orders/services/orders.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<OrdersFilterQuery>({})


const { data: storeOrders } = useQuery({
  queryKey: computed(() => ['store-orders', filter.value]),
  queryFn: () => orderService.getAllOrders(filter.value),
  // enabled: computed(() => !!currentStoreId)
})

function updateFilter(updatedFilter: OrdersFilterQuery) {
  filter.value = {...filter.value, ...updatedFilter}
}

function updatePage(page: number) {
  updateFilter({ pageSize: DEFAULT_PAGINATION_META.pageSize, page: page})
}

function updatePageSize(pageSize: number) {
  updateFilter({ pageSize: pageSize, page: DEFAULT_PAGINATION_META.page})
}
</script>

<style scoped></style>
