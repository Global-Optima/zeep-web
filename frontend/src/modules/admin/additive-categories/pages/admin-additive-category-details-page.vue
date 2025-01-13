<template>
	<p v-if="!categoryDetails">Категория не найдена</p>

	<AdminAdditiveCategoryDetailsForm
		v-else
		:category="categoryDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminAdditiveCategoryDetailsForm from '@/modules/admin/additive-categories/components/details/admin-additive-category-details-form.vue'
import type { UpdateAdditiveCategoryDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const categoryId = route.params.id as string

const queryClient = useQueryClient()

const { data: categoryDetails } = useQuery({
  queryKey: computed(() => ['admin-additive-categories-details', categoryId]),
	queryFn: () => additivesService.getAdditiveCategoryById(Number(categoryId)),
  enabled: !isNaN(Number(categoryId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateAdditiveCategoryDTO}) => {
    return additivesService.updateAdditiveCategory(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories'] })
		queryClient.invalidateQueries({ queryKey: ['admin-additive-categories-details', categoryId] })
		router.push({ name: getRouteName("ADMIN_ADDITIVE_CATEGORIES") })
	},
})

function handleUpdate(updatedData: UpdateAdditiveCategoryDTO) {
  if (isNaN(Number(categoryId))) return router.back()

	updateMutation.mutate({id: Number(categoryId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
