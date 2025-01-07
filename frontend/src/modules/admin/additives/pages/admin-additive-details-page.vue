<template>
	<p v-if="!additiveDetails">Топпинг не найден</p>

	<AdminAdditiveDetailsForm
		v-else
		:additive="additiveDetails"
		@onSubmit="handleCreate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import AdminAdditiveDetailsForm, { type UpdateAdditiveFormSchema } from '@/modules/admin/additives/components/details/admin-additive-details-form.vue'
import type { UpdateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute() 

const additiveId = route.params.id as string

const { data: additiveDetails } = useQuery({
  queryKey: ['admin-additive-details', additiveId],
	queryFn: () => additivesService.getAdditiveById(Number(additiveId)),
  enabled: !isNaN(Number(additiveId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateAdditiveDTO}) => additivesService.updateAdditive(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-additives'] })
    queryClient.invalidateQueries({ queryKey: ['admin-additive-details', additiveId] })
		router.back()
	},
})

function handleCreate(data: UpdateAdditiveFormSchema) {
  if (isNaN(Number(additiveId))) return router.back()

  const dto: UpdateAdditiveDTO = {
    name: data.name,
    description: data.description,
    price: data.price,
    imageUrl: data.imageUrl,
    size: data.size,
    additiveCategoryId: data.additiveCategoryId,
  }

	updateMutation.mutate({id: Number(additiveId), dto})
}

function handleCancel() {
	router.back()
}
</script>
