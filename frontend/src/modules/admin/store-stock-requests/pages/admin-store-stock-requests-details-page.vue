<template>
	<div class="mx-auto max-w-6xl">
		<div class="gap-4 grid grid-cols-2 md:grid-cols-3">
			<div class="col-span-2">
				<AdminStoreStockRequestsDetailsMaterialsTable
					v-if="stockRequest"
					:stockRequest="stockRequest"
				/>
			</div>

			<div class="col-span-full md:col-span-1">
				<AdminStoreStockRequestsDetailsInfo
					v-if="stockRequest"
					:request="stockRequest"
				/>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import AdminStoreStockRequestsDetailsInfo from '@/modules/admin/store-stock-requests/components/details/admin-store-stock-requests-details-info.vue'
import AdminStoreStockRequestsDetailsMaterialsTable from '@/modules/admin/store-stock-requests/components/details/admin-store-stock-requests-details-materials-table.vue'
import { stockRequestsService } from '@/modules/admin/store-stock-requests/services/stock-requests.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const storeStockRequestId = route.params.id as string

const {
  data: stockRequest,
} = useQuery({
  queryKey: computed(() => ['stock-request', Number(storeStockRequestId)]),
  queryFn: () => stockRequestsService.getStockRequestById(Number(storeStockRequestId)),
  enabled: !isNaN(Number(storeStockRequestId)),
})
</script>
