<template>
	<p v-if="!regionDetails">Склад не найден</p>

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
				value="employees"
				>Сотрудники</TabsTrigger
			>
		</TabsList>
		<TabsContent value="details">
			<AdminWarehouseDetailsForm
				:warehouse="regionDetails"
				@onSubmit="handleUpdate"
				@onCancel="handleCancel"
				:readonly="!canUpdate"
			/>
		</TabsContent>

		<TabsContent value="employees">
			<AdminWarehouseEmployees
				:warehouse="regionDetails"
				@on-cancel="handleCancel"
				:readonly="!canUpdate"
			/>
		</TabsContent>
	</Tabs>
</template>

<script lang="ts" setup>
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from '@/core/components/ui/tabs'
import { useToast } from '@/core/components/ui/toast/use-toast'
import { useHasRole } from '@/core/hooks/use-has-roles.hook'
import { EmployeeRole } from '@/modules/admin/employees/models/employees.models'
import AdminWarehouseDetailsForm from '@/modules/admin/warehouses/components/details/admin-warehouse-details-form.vue'
import AdminWarehouseEmployees from '@/modules/admin/warehouses/components/details/admin-warehouse-employees.vue'
import type { UpdateWarehouseDTO } from '@/modules/admin/warehouses/models/warehouse.model'
import { warehouseService } from '@/modules/admin/warehouses/services/warehouse.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const warehouseId = route.params.id as string

const canUpdate = useHasRole([EmployeeRole.ADMIN])

const { data: regionDetails } = useQuery({
	queryKey: ['admin-warehouse-details', warehouseId],
	queryFn: () => warehouseService.getById(Number(warehouseId)),
	enabled: !isNaN(Number(warehouseId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateWarehouseDTO }) => warehouseService.update(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных склада. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-warehouses'] })
		queryClient.invalidateQueries({ queryKey: ['admin-warehouse-details', warehouseId] })
		toast({
			title: 'Успех!',
variant: 'success',
			description: 'Данные склада успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных склада.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(data: UpdateWarehouseDTO) {
	if (isNaN(Number(warehouseId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор склада.',
			variant: 'destructive',
		})
		return router.back()
	}

  console.log("DTOOO", data)

	updateMutation.mutate({ id: Number(warehouseId), dto: data })
}

function handleCancel() {
	router.back()
}
</script>
