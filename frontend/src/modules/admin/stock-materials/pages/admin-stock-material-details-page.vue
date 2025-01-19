<template>
	<p v-if="!additiveDetails">Товар не найден</p>

	<AdminStockMaterialDetailsForm
		v-else
		:stockMaterial="additiveDetails"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStockMaterialDetailsForm from '@/modules/admin/stock-materials/components/details/admin-stock-material-details-form.vue'
import type { UpdateStockMaterialDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()

const stockMaterialId = route.params.id as string

const { data: additiveDetails } = useQuery({
  queryKey: ['admin-stock-material-details', stockMaterialId],
	queryFn: () => stockMaterialsService.getStockMaterialById(Number(stockMaterialId)),
  enabled: !isNaN(Number(stockMaterialId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateStockMaterialDTO}) => stockMaterialsService.updateStockMaterial(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-materials'] })
    queryClient.invalidateQueries({ queryKey: ['admin-stock-material-details', stockMaterialId] })
		router.back()
	},
})

function handleCreate(data: UpdateStockMaterialDTO) {
  if (isNaN(Number(stockMaterialId))) return router.back()

	updateMutation.mutate({id: Number(stockMaterialId), dto: data})
}

function handleCancel() {
	router.back()
}
</script>
