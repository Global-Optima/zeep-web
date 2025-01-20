<template>
	<div>
		<p v-if="!supplier">Поставщик не найден</p>

		<Tabs
			v-else
			default-value="details"
		>
			<TabsList class="grid grid-cols-2 mx-auto mb-6 w-full max-w-6xl">
				<TabsTrigger
					class="py-2"
					value="details"
				>
					Детали
				</TabsTrigger>
				<TabsTrigger
					class="py-2"
					value="variants"
				>
					Товары
				</TabsTrigger>
			</TabsList>
			<TabsContent value="details">
				<AdminSupplierDetailsForm
					:supplier="supplier"
					@on-submit="onUpdateSupplier"
					@on-cancel="onCancel"
				/>
			</TabsContent>

			<TabsContent value="variants">
				<AdminSupplierMaterialsForm
					:supplier="supplier"
					@on-submit="onUpdateSupplierMaterials"
					@on-cancel="onCancel"
				/>
			</TabsContent>
		</Tabs>
	</div>
</template>

<script setup lang="ts">
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/core/components/ui/tabs'
import AdminSupplierDetailsForm from '@/modules/admin/suppliers/components/details/admin-supplier-details-form.vue'
import AdminSupplierMaterialsForm from '@/modules/admin/suppliers/components/details/admin-supplier-materials-form.vue'
import type { UpdateSupplierDTO, UpdateSupplierMaterialDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()

const supplierId = route.params.id as string

const { data: supplier } = useQuery({
  queryKey: ['admin-supplier-details', supplierId],
	queryFn: () => suppliersService.getSupplierByID(Number(supplierId)),
  enabled: !isNaN(Number(supplierId)),
})

const {mutate: updateSupplier} = useMutation({
	mutationFn: ({id, dto} : {id: number, dto: UpdateSupplierDTO}) => suppliersService.updateSupplier(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
    queryClient.invalidateQueries({ queryKey: ['admin-supplier-details', supplierId] })
    queryClient.invalidateQueries({ queryKey: ['admin-supplier-materials', Number(supplierId)] })

	},
})

const {mutate: updateSupplierMaterials} = useMutation({
	mutationFn: ({id, dto} : {id: number, dto: UpdateSupplierMaterialDTO[]}) => suppliersService.updateSupplierMaterials(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
    queryClient.invalidateQueries({ queryKey: ['admin-supplier-details', supplierId] })
	},
})

function onUpdateSupplier(dto: UpdateSupplierDTO) {
  if (isNaN(Number(supplierId))) return router.back()
	updateSupplier({id: Number(supplierId), dto})
}

function onUpdateSupplierMaterials(dto: UpdateSupplierMaterialDTO[]) {
  if (isNaN(Number(supplierId))) return router.back()
	updateSupplierMaterials({id: Number(supplierId), dto})
}

function onCancel() {
	router.back()
}
</script>
