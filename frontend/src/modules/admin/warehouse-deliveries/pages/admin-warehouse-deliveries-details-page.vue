<template>
	<AdminWarehouseDeliveriesDetailsForm
		v-if="deliveryDetails"
		:delivery="deliveryDetails"
		@on-cancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminWarehouseDeliveriesDetailsForm from '@/modules/admin/warehouse-deliveries/components/details/admin-warehouse-deliveries-details-form.vue'
import { warehouseStocksService } from '@/modules/admin/warehouse-stocks/services/warehouse-stocks.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const warehouseDeliveryId = route.params.id as string

const { data: deliveryDetails } = useQuery({
  queryKey: computed(() => ['warehouse-delivery', warehouseDeliveryId]),
	queryFn: () => warehouseStocksService.getWarehouseDeliveryId(Number(warehouseDeliveryId)),
  enabled: !isNaN(Number(warehouseDeliveryId)),
})

function handleCancel() {
	router.back()
}
</script>
