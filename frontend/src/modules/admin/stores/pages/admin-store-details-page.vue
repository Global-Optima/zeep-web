<template>
	<p v-if="!storeData">Кафе не найдено</p>

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
			<AdminStoreDetailsForm
				:store="storeData"
				@onSubmit="handleUpdate"
				@onCancel="handleCancel"
				:readonly="!canUpdate"
			/>
		</TabsContent>

		<TabsContent value="employees">
			<AdminStoreDetailsEmployees
				:store="storeData"
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
import AdminStoreDetailsEmployees from '@/modules/admin/stores/components/details/admin-store-details-employees.vue'
import AdminStoreDetailsForm from '@/modules/admin/stores/components/details/admin-store-details-form.vue'
import type { UpdateStoreDTO } from '@/modules/admin/stores/models/stores-dto.model'
import { storesService } from '@/modules/admin/stores/services/stores.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()

const storeId = route.params.id as string

const queryClient = useQueryClient()

const canUpdate = useHasRole(EmployeeRole.ADMIN)

const { data: storeData } = useQuery({
	queryKey: ['store', storeId],
	queryFn: () => storesService.getStore(Number(storeId)),
	enabled: computed(() => !!storeId),
})

const updateMutation = useMutation({
	mutationFn: (updatedData: UpdateStoreDTO) =>
		storesService.updateStore(Number(storeId), updatedData),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Пожалуйста, подождите, данные обновляются.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['stores'] })
		queryClient.invalidateQueries({ queryKey: ['store', storeId] })
		toast({
			title: 'Успех!',
			description: 'Данные кафе успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных кафе.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(updatedData: UpdateStoreDTO) {
	updateMutation.mutate(updatedData)
}

function handleCancel() {
	router.back()
}
</script>
