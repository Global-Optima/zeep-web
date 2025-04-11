<template>
	<Card>
		<CardHeader>
			<div>
				<CardTitle>Информация о заготовке</CardTitle>
			</div>
		</CardHeader>

		<CardContent>
			<div class="relative space-y-4">
				<template
					v-for="detail in provisionDetails"
					:key="detail.label"
				>
					<div v-if="detail.value">
						<p class="text-muted-foreground text-sm">{{ detail.label }}</p>
						<span
							v-if="detail.value"
							v-html="detail.value"
							class="text-sm"
						></span>
						<p
							v-else
							class="text-sm"
						>
							Отсутствует
						</p>
					</div>
				</template>
			</div>
		</CardContent>

		<CardFooter
			v-if="!readonly"
			class="flex flex-col gap-2"
		>
			<Button
				v-if="isEditable"
				@click="onCompleteClick"
				class="rounded-lg w-full h-10"
			>
				Завершить приготовление
			</Button>

			<Button
				v-if="isEditable"
				variant="outline"
				class="gap-3 rounded-lg w-full h-10"
				@click="onUpdateClick"
			>
				<Pencil class="size-4" />
				Редактировать
			</Button>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/core/components/ui/card'
import {
  StoreProvisionStatus,
  type StoreProvisionDetailsDTO,
} from '@/modules/admin/store-provisions/models/store-provision.models'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { Pencil } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const {storeProvision, readonly} = defineProps<{
  storeProvision: StoreProvisionDetailsDTO
  readonly?: boolean
}>()

const emits = defineEmits<{
  onComplete: []
}>()


const router = useRouter()

function formatDate(date?: Date | string): string | null {
  if (!date) return null
  return format(new Date(date), 'dd.MM.yyyy HH:mm', { locale: ru })
}

const STORE_PROVISION_STATUS_FORMATTED: Record<StoreProvisionStatus, string> = {
  [StoreProvisionStatus.COMPLETED]: 'Приготовлен',
  [StoreProvisionStatus.PREPARING]: 'Готовится',
  [StoreProvisionStatus.EXPIRED]: 'Просрочен',
  [StoreProvisionStatus.EMPTY]: 'Пустой',
}

const provisionDetails = computed(() => [
  { label: 'Номер провизии', value: storeProvision.id },
  { label: 'Заготовка', value: storeProvision.provision.name },
  { label: 'Объем', value: `${storeProvision.volume} ${storeProvision.provision.unit.name.toLowerCase()}` },
  { label: 'Начальный объем', value: `${storeProvision.initialVolume} ${storeProvision.provision.unit.name.toLowerCase()}` },
  { label: 'Срок годности (минут)', value: storeProvision.expirationInMinutes },
  { label: 'Статус', value:  STORE_PROVISION_STATUS_FORMATTED[storeProvision.status] },
  { label: 'Дата создания', value: formatDate(storeProvision.createdAt) },
  { label: 'Дата завершения', value: formatDate(storeProvision.completedAt) },
  { label: 'Дата истечения', value: formatDate(storeProvision.expiresAt) },
])


const isEditable = computed(
  () => storeProvision.status === StoreProvisionStatus.PREPARING
)

const onUpdateClick = () => {
  router.push(`/admin/store-provisions/${storeProvision.id}/update`)
}

const onCompleteClick = () => {
  emits("onComplete")
}
</script>
