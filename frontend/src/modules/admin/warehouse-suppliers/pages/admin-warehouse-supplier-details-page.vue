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
				<AdminWarehouseSupplierDetailsForm
					:supplier="supplier"
					@on-submit="onUpdateSupplier"
					@on-cancel="onCancel"
				/>
			</TabsContent>

			<TabsContent value="variants">
				<AdminWarehouseSupplierMaterialsForm
					:supplier="supplier"
					:stock-materials="stockMaterials"
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
import type { UpdateSupplierDTO, UpdateSupplierMaterialDTO, UpsertSupplierMaterialsDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import AdminWarehouseSupplierDetailsForm from '@/modules/admin/warehouse-suppliers/components/details/admin-warehouse-supplier-details-form.vue'
import AdminWarehouseSupplierMaterialsForm from '@/modules/admin/warehouse-suppliers/components/details/admin-warehouse-supplier-materials-form.vue'
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
	mutationFn: ({id, dto} : {id: number, dto: UpsertSupplierMaterialsDTO}) => suppliersService.updateSupplierMaterials(id, dto),
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
    queryClient.invalidateQueries({ queryKey: ['admin-supplier-details', supplierId] })
	},
})

const { data: stockMaterials } = useQuery({
  queryKey:   ['admin-supplier-materials', supplierId],
  queryFn:  () => suppliersService.getMaterialsBySupplier(Number(supplierId)),
  initialData: [],
})

function onUpdateSupplier(dto: UpdateSupplierDTO) {
  if (isNaN(Number(supplierId))) return router.back()
	updateSupplier({id: Number(supplierId), dto})
}

function onUpdateSupplierMaterials(data: UpdateSupplierMaterialDTO[]) {
  if (isNaN(Number(supplierId))) return router.back()
	updateSupplierMaterials({id: Number(supplierId), dto: {materials: data}})
}

function onCancel() {
	router.back()
}
</script>
