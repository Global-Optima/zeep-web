<template>
	<AdminStockMaterialCategoryCreateForm
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminStockMaterialCategoryCreateForm from '@/modules/admin/stock-material-categories/components/create/admin-stock-material-category-create-form.vue'
import type { CreateStockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()

const createMutation = useMutation({
	mutationFn: (dto: CreateStockMaterialCategoryDTO) => stockMaterialCategoryService.create(dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories'] })
		router.back()
	},
})

function handleCreate(dto: CreateStockMaterialCategoryDTO) {
	createMutation.mutate(dto)
}

function handleCancel() {
	router.back()
}
</script>
