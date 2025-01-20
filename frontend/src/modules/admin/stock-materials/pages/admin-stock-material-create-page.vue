<template>
	<AdminStockMaterialCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStockMaterialCreateForm from '@/modules/admin/stock-materials/components/create/admin-stock-material-create-form.vue'
import type { CreateStockMaterialDTO } from '@/modules/admin/stock-materials/models/stock-materials.model'
import { stockMaterialsService } from '@/modules/admin/stock-materials/services/stock-materials.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateStockMaterialDTO) => stockMaterialsService.createStockMaterial(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-materials'] })
		router.back()
	},
})

function handleCreate(dto: CreateStockMaterialDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
