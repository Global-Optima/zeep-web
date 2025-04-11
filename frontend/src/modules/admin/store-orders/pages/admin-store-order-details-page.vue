<template>
	<p v-if="!orderDetails">Заказ не найден</p>

	<div v-else>
		<AdminOrderDetails :order-details="orderDetails" />
	</div>
</template>

<script lang="ts" setup>
import AdminOrderDetails from '@/modules/admin/store-orders/components/details/admin-order-details.vue'
import { ordersService } from '@/modules/admin/store-orders/services/orders.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const orderId = route.params.id as string

const { data: orderDetails } = useQuery({
  queryKey: computed(() => ['admin-store-order', orderId]),
  queryFn: () => ordersService.getOrderById(Number(orderId)),
  enabled: !isNaN(Number(orderId)),
})
</script>
