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
					>Детали</TabsTrigger
				>
				<TabsTrigger
					class="py-2"
					value="variants"
					>Товары</TabsTrigger
				>
			</TabsList>
			<TabsContent value="details">
				<AdminSupplierDetailsForm
					:supplier="supplier"
					@on-submit="onUpdateSupplier"
					@on-cancel="onCancel"
					:readonly="!canUpdate"
				/>
			</TabsContent>

			<TabsContent value="variants">
				<AdminSupplierMaterialsForm
					:supplier="supplier"
					:stock-materials="stockMaterials"
					@on-submit="onUpdateSupplierMaterials"
					@on-cancel="onCancel"
					:readonly="!canUpdate"
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
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminSupplierDetailsForm from '@/modules/admin/suppliers/components/details/admin-supplier-details-form.vue'
import AdminSupplierMaterialsForm from '@/modules/admin/suppliers/components/details/admin-supplier-materials-form.vue'
import type { UpdateSupplierDTO, UpdateSupplierMaterialDTO, UpsertSupplierMaterialsDTO } from '@/modules/admin/suppliers/models/suppliers.model'
import { suppliersService } from '@/modules/admin/suppliers/services/suppliers.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()
const { toast } = useToast()

const canUpdate = useHasRole([EmployeeRole.ADMIN])

const supplierId = route.params.id as string

const { data: supplier } = useQuery({
	queryKey: ['admin-supplier-details', supplierId],
	queryFn: () => suppliersService.getSupplierByID(Number(supplierId)),
	enabled: !isNaN(Number(supplierId)),
})

const { mutate: updateSupplier } = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateSupplierDTO }) =>
		suppliersService.updateSupplier(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных поставщика.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
		queryClient.invalidateQueries({ queryKey: ['admin-supplier-details', supplierId] })
		queryClient.invalidateQueries({ queryKey: ['admin-supplier-materials', Number(supplierId)] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные поставщика успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных поставщика.',
			variant: 'destructive',
		})
	},
})

const { mutate: updateSupplierMaterials } = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpsertSupplierMaterialsDTO }) =>
		suppliersService.updateSupplierMaterials(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление материалов поставщика.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-suppliers'] })
		queryClient.invalidateQueries({ queryKey: ['admin-supplier-details', supplierId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Материалы поставщика успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении материалов поставщика.',
			variant: 'destructive',
		})
	},
})

const { data: stockMaterials } = useQuery({
	queryKey: ['admin-supplier-materials', supplierId],
	queryFn: () => suppliersService.getMaterialsBySupplier(Number(supplierId)),
	initialData: [],
})

function onUpdateSupplier(dto: UpdateSupplierDTO) {
	if (isNaN(Number(supplierId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор поставщика.',
			variant: 'destructive',
		})
		return router.back()
	}
	updateSupplier({ id: Number(supplierId), dto })
}

function onUpdateSupplierMaterials(data: UpdateSupplierMaterialDTO[]) {
	if (isNaN(Number(supplierId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор поставщика.',
			variant: 'destructive',
		})
		return router.back()
	}
	updateSupplierMaterials({ id: Number(supplierId), dto: { materials: data } })
}

function onCancel() {
	router.back()
}
</script>
