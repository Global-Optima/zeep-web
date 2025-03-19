<template>
	<p v-if="!stockRequest">Запрос на склад не найден</p>

	<div
		v-else
		class="mx-auto max-w-6xl"
	>
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Детали запроса на склад #{{ stockRequest.requestId }}
			</h1>

			<div class="hidden md:flex items-center gap-2 md:ml-auto"></div>
		</div>

		<div class="gap-4 grid grid-cols-2 md:grid-cols-3 mt-4">
			<div class="col-span-2">
				<AdminStockRequestsDetailsMaterialsTable :stockRequest="stockRequest" />
			</div>

			<div class="col-span-full md:col-span-1">
				<AdminStockRequestsDetailsInfo :request="stockRequest" />
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { Button } from '@/core/components/ui/button'
import AdminStockRequestsDetailsInfo from '@/modules/admin/stock-requests/components/details/admin-stock-requests-details-info.vue'
import AdminStockRequestsDetailsMaterialsTable from '@/modules/admin/stock-requests/components/details/admin-stock-requests-details-materials-table.vue'
import { stockRequestsService } from '@/modules/admin/stock-requests/services/stock-requests.service'
import { useQuery } from '@tanstack/vue-query'
import { ChevronLeft } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const storeStockRequestId = route.params.id as string

function onCancel() {
	router.back()
}

const {
  data: stockRequest,
} = useQuery({
  queryKey: computed(() => ['stock-request', Number(storeStockRequestId)]),
  queryFn: () => stockRequestsService.getStockRequestById(Number(storeStockRequestId)),
  enabled: !isNaN(Number(storeStockRequestId)),
})
</script>
