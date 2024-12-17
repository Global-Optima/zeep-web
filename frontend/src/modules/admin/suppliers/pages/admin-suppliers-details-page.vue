<template>
	<AdminSuppliersDetailsForm
		v-if="supplierData"
		:initialData="supplierData"
		@onSubmit="handleUpdate"
		@onCancel="handleCancel"
	/>
</template>

<script lang="ts" setup>
import { getRouteName } from '@/core/config/routes.config'
import AdminSuppliersDetailsForm from '@/modules/admin/suppliers/components/details/admin-suppliers-details-form.vue'
import type { UpdateSupplierDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const supplierId = route.params.id as string

const queryClient = useQueryClient()

const { data: supplierData } = useQuery({
	queryKey: ['supplier', supplierId],
	queryFn: () => suppliersService.getSupplier(Number(supplierId)),
	enabled: computed(() => !!supplierId),
})

const updateMutation = useMutation({
  mutationFn: (updatedData: UpdateSupplierDTO) => suppliersService.updateSupplier(Number(supplierId), updatedData),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['suppliers'] })
		queryClient.invalidateQueries({queryKey: ['supplier', supplierId]})
		router.push({name: getRouteName("ADMIN_SUPPLIERS")})
	},
})

function handleUpdate(updatedData: UpdateSupplierDTO) {
	updateMutation.mutate(updatedData)
}

function handleCancel() {
	router.back()
}
</script>
