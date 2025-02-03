<template>
	<p v-if="!regionDetails">Топпинг не найден</p>

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
			<AdminRegionDetailsForm
				:region="regionDetails"
				@onSubmit="handleUpdate"
				@onCancel="handleCancel"
			/>
		</TabsContent>

		<TabsContent value="variants">
			<AdminRegionWarehouses
				:region="regionDetails"
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
import type { UpdateAdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import AdminRegionDetailsForm from '@/modules/admin/regions/components/details/admin-region-details-form.vue'
import AdminRegionWarehouses from '@/modules/admin/regions/components/details/admin-region-warehouses.vue'
import type { UpdateRegionDTO } from '@/modules/admin/regions/models/regions.model'
import { regionsService } from '@/modules/admin/regions/services/regions.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const queryClient = useQueryClient()
const route = useRoute()
const { toast } = useToast()

const additiveId = route.params.id as string

const { data: regionDetails } = useQuery({
	queryKey: ['admin-region-details', additiveId],
	queryFn: () => regionsService.getById(Number(additiveId)),
	enabled: !isNaN(Number(additiveId)),
})

const updateMutation = useMutation({
	mutationFn: ({ id, dto }: { id: number; dto: UpdateRegionDTO }) =>
  regionsService.update(id, dto),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных региона. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-regions'] })
		queryClient.invalidateQueries({ queryKey: ['admin-region-details', additiveId] })
		toast({
			title: 'Успех!',
			description: 'Данные региона успешно обновлены.',
		})
	},
	onError: () => {
		toast({
			title: 'Ошибка',
			description: 'Произошла ошибка при обновлении данных региона.',
			variant: 'destructive',
		})
	},
})

function handleUpdate(data: UpdateAdditiveDTO) {
	if (isNaN(Number(additiveId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор региона.',
			variant: 'destructive',
		})
		return router.back()
	}

	updateMutation.mutate({ id: Number(additiveId), dto: data })
}

function handleCancel() {
	router.back()
}
</script>
