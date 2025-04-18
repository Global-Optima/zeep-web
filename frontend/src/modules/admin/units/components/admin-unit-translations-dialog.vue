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
import type { UnitTranslationsDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { Globe } from 'lucide-vue-next'
import { computed, defineAsyncComponent, ref, toRefs } from 'vue'

// Async dialog component
const AdminTranslationsDialog = defineAsyncComponent(
  () => import('@/core/components/admin-translations-dialog/AdminTranslationsDialog.vue')
)

const props = defineProps<{ unitId: number }>()
const {unitId} = toRefs(props)

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
  queryKey: computed(() => ['unit-translations', unitId.value]),
  queryFn: () => unitsService.getUnitTranslations(unitId.value),
  enabled: computed(() => !!unitId.value),
})

const fields = computed<TranslationFieldLocale[]>(() => [
  {
    field: 'name',
    label: 'Название',
    locales: translations.value?.name ?? {},
  },
])

// Mutation to upsert translations
const { mutate: saveTranslations, isPending: isSaving } = useMutation({
  mutationFn: (dto: UnitTranslationsDTO) => unitsService.upsertUnitTranslations(unitId.value, dto),
  onMutate: () => {
    toast({
      title: 'Сохранение переводов...',
    })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['unit-translations', unitId.value] })
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
  const dto: UnitTranslationsDTO = {}

  payload.forEach((item) => {
    const [field, locales] = Object.entries(item)[0]
    dto[field as keyof UnitTranslationsDTO] = locales
  })

  saveTranslations(dto)
}
</script>

<style scoped>
/* Add any component-specific styles if needed */
</style>
