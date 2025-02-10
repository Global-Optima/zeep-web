<template>
	<p v-if="!regionDetails">Франчайзи не найден</p>

	<Tabs
		v-else
		default-value="details"
	>
		<TabsList class="grid grid-cols-3 mx-auto mb-6 w-full max-w-6xl">
			<TabsTrigger
				class="py-2"
				value="details"
				>Детали</TabsTrigger
			>
			<TabsTrigger
				class="py-2"
				value="variants"
				>Кафе</TabsTrigger
			>
			<TabsTrigger
				class="py-2"
				value="employees"
				>Сотрудники</TabsTrigger
			>
		</TabsList>
		<TabsContent value="details">
			<AdminFranchiseeDetailsForm
				:franchisee="regionDetails"
				@onSubmit="handleUpdate"
				@onCancel="handleCancel"
			/>
		</TabsContent>

		<TabsContent value="variants">
			<AdminFranchiseeStores
				:franchisee="regionDetails"
				@on-cancel="handleCancel"
			/>
		</TabsContent>

		<TabsContent value="employees">
			<AdminFranchiseeEmployees
				:franchisee="regionDetails"
				@on-cancel="handleCancel"
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
import AdminFranchiseeDetailsForm from '@/modules/admin/franchisees/components/details/admin-franchisee-details-form.vue'
import AdminFranchiseeEmployees from '@/modules/admin/franchisees/components/details/admin-franchisee-employees.vue'
import AdminFranchiseeStores from '@/modules/admin/franchisees/components/details/admin-franchisee-stores.vue'
import type { UpdateFranchiseeDTO } from '@/modules/admin/franchisees/models/franchisee.model'
import { franchiseeService } from '@/modules/admin/franchisees/services/franchisee.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const franchiseeId = route.params.id as string

const { data: regionDetails } = useQuery({
	queryKey: ['admin-franchisee-details', franchiseeId],
	queryFn: () => franchiseeService.getById(Number(franchiseeId)),
	enabled: !isNaN(Number(franchiseeId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateFranchiseeDTO }) =>
  franchiseeService.update(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных франчайзи. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-franchisees'] })
		queryClient.invalidateQueries({ queryKey: ['admin-franchisee-details', franchiseeId] })
		toast({
			title: 'Успех!',
			description: 'Данные франчайзи успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных франчайзи.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(data: UpdateFranchiseeDTO) {
	if (isNaN(Number(franchiseeId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор франчайзи.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(franchiseeId), dto: data })
}

function handleCancel() {
	router.back()
}
</script>
