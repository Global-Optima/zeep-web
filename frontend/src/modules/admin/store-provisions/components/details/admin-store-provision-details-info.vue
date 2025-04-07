<template>
	<Card>
		<CardHeader>
			<div>
				<CardTitle>Информация о заготовке</CardTitle>
				<CardDescription class="mt-2">
					Подробная информация о созданной заготовке.
				</CardDescription>
			</div>
		</CardHeader>

		<CardContent>
			<div class="space-y-4">
				<div
					v-for="detail in provisionDetails"
					:key="detail.label"
				>
					<p class="mb-1 text-muted-foreground text-sm">{{ detail.label }}</p>
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
			</div>
		</CardContent>

		<CardFooter v-if="!readonly">
			<div class="flex gap-2 w-full">
				<Button
					variant="outline"
					class="flex-1"
					@click="onUpdateClick"
				>
					Редактировать
				</Button>
				<Button
					variant="outline"
					class="flex-1"
					v-if="isEditable"
					@click="onCompleteClick"
				>
					Завершить провизию
				</Button>
			</div>
		</CardFooter>
	</Card>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import {
  StoreProvisionStatus,
  type StoreProvisionDetailsDTO,
} from '@/modules/admin/store-provisions/models/store-provision.models'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
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

function formatDate(date?: Date | string): string {
  if (!date) return '-'
  return format(new Date(date), 'dd.MM.yyyy HH:mm', { locale: ru })
}

const STORE_PROVISION_STATUS_FORMATTED: Record<StoreProvisionStatus, string> = {
  [StoreProvisionStatus.COMPLETED]: 'Приготовлен',
  [StoreProvisionStatus.PREPARING]: 'Готовится',
  [StoreProvisionStatus.EXPIRED]: 'Просрочен',
}

const provisionDetails = computed(() => [
  { label: 'Номер провизии', value: storeProvision.id },
  { label: 'Заготовка', value: storeProvision.provision.name },
  { label: 'Объем', value: storeProvision.volume },
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
