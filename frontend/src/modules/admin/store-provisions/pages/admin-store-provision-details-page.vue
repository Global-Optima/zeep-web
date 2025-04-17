<template>
	<p v-if="!storeProvision">Заготовка не найдена</p>

	<div
		v-else
		class="mx-auto max-w-6xl"
	>
		<!-- Header -->
		<div class="flex items-center gap-4">
    <Button
      variant="outline"
				size="icon"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="flex-1 sm:grow-0 font-semibold text-xl tracking-tight whitespace-nowrap shrink-0">
				Детали заготовки от {{ formattedCreatedAt }}
			</h1>
		</div>

		<!-- Main Grid -->
		<div class="gap-4 grid grid-cols-2 md:grid-cols-3 mt-4">
			<!-- Left: Ingredients / Technical Map -->
			<div class="col-span-2">
				<AdminStoreProvisionDetailsIngredients :storeProvision="storeProvision" />
			</div>

			<!-- Right: Detailed Info & Actions -->
			<div class="col-span-full md:col-span-1">
				<AdminStoreProvisionDetailsInfo
					:storeProvision="storeProvision"
					@on-complete="onCompleteProvision"
				/>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { Button } from '@/core/components/ui/button'
import { useToast } from '@/core/components/ui/toast'
import { useAxiosLocaleToast, type AxiosLocalizedError } from '@/core/hooks/use-axios-locale-toast.hooks'
import AdminStoreProvisionDetailsInfo from '@/modules/admin/store-provisions/components/details/admin-store-provision-details-info.vue'
import AdminStoreProvisionDetailsIngredients from '@/modules/admin/store-provisions/components/details/admin-store-provision-details-ingredients.vue'
import type { StoreProvisionDetailsDTO } from '@/modules/admin/store-provisions/models/store-provision.models'
import { storeProvisionsService } from '@/modules/admin/store-provisions/services/store-provision.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { ChevronLeft } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()
const queryClient = useQueryClient()
const {toastLocalizedError} = useAxiosLocaleToast()

const storeProvisionId = route.params.id as string

function onCancel() {
  router.back()
}

const { data: storeProvision } = useQuery<StoreProvisionDetailsDTO>({
  queryKey: computed(() => ['admin-store-provision', storeProvisionId]),
  queryFn: () => storeProvisionsService.getStoreProvisionById(Number(storeProvisionId)),
  enabled: !isNaN(Number(storeProvisionId)),
})

const {mutate: completeProvision} = useMutation({
	mutationFn: (id: number) => storeProvisionsService.completeStoreProvision(id),
	onMutate: () => {
		toast({
			title: 'Обновление...',
			description: 'Обновление данных заготовки. Пожалуйста, подождите.',
		})
	},
	onSuccess: () => {
		queryClient.invalidateQueries({ queryKey: ['admin-store-provisions'] })
		queryClient.invalidateQueries({ queryKey: ['admin-store-provision', storeProvisionId] })
		toast({
			title: 'Успех!',
      variant: 'success',
			description: 'Данные успешно обновлены.',
		})
	},
	onError: (err: AxiosLocalizedError) => {
    toastLocalizedError(err, 'Произошла ошибка при обновлении.')
	},
})

const onCompleteProvision = () => {
  if (isNaN(Number(storeProvisionId))) {
		toast({
			title: 'Ошибка',
			description: 'Неверный идентификатор.',
			variant: 'destructive',
	  })

    return router.back()
  }

  completeProvision(Number(storeProvisionId))
}

const formattedCreatedAt = computed(() => format(storeProvision.value?.createdAt ?? new Date(), 'dd MMM yyyy, HH:mm', { locale: ru }))
</script>
