<template>
	<div
		class="flex justify-between items-center gap-6 bg-blue-100 p-4 border border-blue-200 rounded-xl"
	>
		<div>
			<p class="font-semibold text-blue-900">Данные требуют синхронизации</p>
			<p class="text-blue-900 text-sm">
				Внесены изменения в продукты, добавки или технологические карты. Нажмите кнопку ниже для
				обновления.
			</p>
		</div>
		<Button
			class="flex items-center gap-2 bg-blue-500 disabled:opacity-50 px-4 py-2 rounded-lg font-medium text-white"
			:disabled="isLoading"
			@click="syncData"
		>
			<Loader
				v-if="isLoading"
				class="size-5 text-white animate-spin"
			/>
			<span v-if="!isLoading">Синхронизировать</span>
			<span v-else>Синхронизация...</span>
		</Button>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { useToast } from '@/core/components/ui/toast'
import { storeSyncService } from '@/modules/admin/stores/services/stores-sync.service'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { Loader } from 'lucide-vue-next'
import { ref } from 'vue'

const queryClient = useQueryClient()
const { toast } = useToast()
const isLoading = ref(false)

const { mutate: syncData } = useMutation({
  mutationFn: () => {
    isLoading.value = true
    return storeSyncService.syncStoreStocksAndAdditives()
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['admin-store-is-sync'] })
    toast({
      title: 'Синхронизация успешно завершена',
      description: 'Все изменения в продуктах, добавках и технологических картах обновлены.',
      variant: 'success',
    })
  },
  onError: () => {
    toast({
      title: 'Ошибка синхронизации',
      description: 'Не удалось обновить данные. Попробуйте еще раз.',
      variant: 'destructive',
    })
  },
  onSettled: () => {
    isLoading.value = false
  },
})
</script>

<style scoped></style>
