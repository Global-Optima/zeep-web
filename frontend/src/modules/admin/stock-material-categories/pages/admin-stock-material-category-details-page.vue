<template>
	<p v-if="!categoryDetails">Категория не найдена</p>

	<AdminStockMaterialCategoryDetailsForm
		v-else
		:category="categoryDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStockMaterialCategoryDetailsForm from '@/modules/admin/stock-material-categories/components/details/admin-stock-material-category-details-form.vue'
import type { UpdateStockMaterialCategoryDTO } from '@/modules/admin/stock-material-categories/models/stock-material-categories.model'
import { stockMaterialCategoryService } from '@/modules/admin/stock-material-categories/services/stock-materials.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const categoryId = route.params.id as string

const queryClient = useQueryClient()

const { data: categoryDetails } = useQuery({
  queryKey: computed(() => ['admin-stock-material-categories-details', categoryId]),
	queryFn: () => stockMaterialCategoryService.getById(Number(categoryId)),
  enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateStockMaterialCategoryDTO}) => {
    return stockMaterialCategoryService.update(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-stock-material-categories-details', categoryId] })
		router.push({ name: getRouteName("ADMIN_STOCK_MATERIAL_CATEGORIES") })
	},
})

function handleUpdate(updatedData: UpdateStockMaterialCategoryDTO) {
  if (isNaN(Number(categoryId))) return router.back()

	updateMutation.mutate({id: Number(categoryId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
