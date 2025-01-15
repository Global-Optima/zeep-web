<template>
	<p v-if="!storeAdditiveDetails">Товар не найден</p>

	<AdminStoreAdditiveDetailsForm
		v-else
		:initialAdditive="storeAdditiveDetails"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreAdditiveDetailsForm from '@/modules/admin/store-additives/components/details/admin-store-additive-details-form.vue'
import type { UpdateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const storeAdditiveId = route.params.id as string

const queryClient = useQueryClient()

const { data: storeAdditiveDetails } = useQuery({
  queryKey: computed(() => ['admin-store-additive-details', storeAdditiveId]),
	queryFn: () => storeAdditivesService.getStoreAdditiveById(Number(storeAdditiveId)),
  enabled: !isNaN(Number(storeAdditiveId)),
})

const updateMutation = useMutation({
	mutationFn: ({id, dto}:{id: number, dto: UpdateStoreAdditiveDTO}) => {
    return storeAdditivesService.updateStoreAdditive(id, dto)
  },
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-additives'] })
		queryClient.invalidateQueries({ queryKey: ['admin-store-additive-details', storeAdditiveId] })
		router.push({ name: getRouteName("ADMIN_STORE_ADDITIVES") })
	},
})

function handleUpdate(updatedData: UpdateStoreAdditiveDTO) {
  if (isNaN(Number(storeAdditiveId))) return router.back()

	updateMutation.mutate({id: Number(storeAdditiveId), dto: updatedData})
}

function handleCancel() {
	router.back()
}
</script>
