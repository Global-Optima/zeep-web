<template>
	<AdminStoreDetailsForm
		v-if="storeData"
		:initialData="storeData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
	<div v-else>Loading...</div>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminStoreDetailsForm from '@/modules/admin/stores/components/details/admin-store-details-form.vue'
import type { UpdateStoreDTO } from '@/modules/stores/models/stores-dto.model'
import { storesService } from '@/modules/stores/services/stores.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const storeId = route.params.id as string

const queryClient = useQueryClient()

const { data: storeData } = useQuery({
	queryKey: ['store', storeId],
	queryFn: () => storesService.getStore(Number(storeId)),
	enabled: computed(() => !!storeId),
})

const updateMutation = useMutation({
  mutationFn: (updatedData: UpdateStoreDTO) => storesService.updateStore(Number(storeId), updatedData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stores'] })
		queryClient.invalidateQueries({queryKey: ['store', storeId]})
		router.push({name: getRouteName("ADMIN_STORES")})
	},
})

function handleUpdate(updatedData: UpdateStoreDTO) {
  console.log("HERE")
	updateMutation.mutate(updatedData)
}

function handleCancel() {
	router.back()
}
</script>
