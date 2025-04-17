<template>
	<div>
		<Button
			@click="openDialog"
			variant="outline"
			size="icon"
			type="button"
			:disabled="isLoadingTranslations || isSaving"
		>
			<Globe
				class="size-4"
				:stroke-width="1.4"
			/>
		</Button>

		<AdminTranslationsDialog
			v-model:open="isDialogOpen"
			:fields="fields"
			:loading="isLoadingTranslations || isSaving"
			@submit="handleSubmit"
		/>
	</div>
</template>

<script setup lang="ts">
import type { TranslationFieldLocale, TranslationsLanguage } from '@/core/components/admin-translations-dialog'
import { Button } from '@/core/components/ui/button'
import { useToast } from '@/core/components/ui/toast'
import type { AdditiveTranslationsDTO } from '@/modules/admin/additives/models/additives.model'
import { additivesService } from '@/modules/admin/additives/services/additives.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { Globe } from 'lucide-vue-next'
import { computed, defineAsyncComponent, defineProps, ref, toRefs } from 'vue'

// Async dialog component
const AdminTranslationsDialog = defineAsyncComponent(
  () => import('@/core/components/admin-translations-dialog/AdminTranslationsDialog.vue')
)

// Props: productId must be passed in
const props = defineProps<{ additiveId: number }>()
const {additiveId} = toRefs(props)

// Dialog state
const isDialogOpen = ref(false)
function openDialog() {
  isDialogOpen.value = true
}

// Toast & query client
const {toast} = useToast()
const queryClient = useQueryClient()

// Fetch existing translations
const { data: translations, isLoading: isLoadingTranslations } = useQuery({
  queryKey: computed(() => ['additive-translations', additiveId.value]),
  queryFn: () => additivesService.getAdditiveTranslations(additiveId.value),
  enabled: computed(() => !!additiveId.value),
})

const fields = computed<TranslationFieldLocale[]>(() => [
  {
    field: 'name',
    label: 'Название',
    locales: translations.value?.name ?? {},
  },
  {
    field: 'description',
    label: 'Описание',
    locales: translations.value?.description ?? {},
  },
])

// Mutation to upsert translations
const { mutate: saveTranslations, isPending: isSaving } = useMutation({
  mutationFn: (dto: AdditiveTranslationsDTO) => additivesService.upsertAdditiveTranslations(additiveId.value, dto),
  onMutate: () => {
    toast({
      title: 'Сохранение переводов...',
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['additive-translations', additiveId.value] })
    toast({
      title: 'Успех!',
      description: 'Переводы успешно сохранены.',
      variant: 'success',
    })
    isDialogOpen.value = false
  },
  onError: () => {
    toast({
      title: 'Ошибка',
      description: 'Не удалось сохранить переводы.',
      variant: 'destructive',
    })
  },
})

// Handle dialog submit
function handleSubmit(
  payload: Record<string, Partial<Record<TranslationsLanguage, string>>>[]
) {
  // Convert array payload to DTO object
  const dto: AdditiveTranslationsDTO = {}

  payload.forEach((item) => {
    const [field, locales] = Object.entries(item)[0]
    dto[field as keyof AdditiveTranslationsDTO] = locales
  })

  saveTranslations(dto)
}
</script>

<style scoped>
/* Add any component-specific styles if needed */
</style>
