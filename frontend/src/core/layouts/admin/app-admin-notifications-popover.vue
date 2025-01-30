<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/core/components/ui/popover'
import { useToast } from '@/core/components/ui/toast'
import { getRouteName } from '@/core/config/routes.config'
import { formatLocalizedMessage } from '@/core/utils/format-localized-messages.utils'
import type {
  GetNotificationsFilter,
  MarkNotificationsAsReadDTO,
} from '@/modules/admin/admin-notifications/models/notifications.model'
import { notificationsService } from '@/modules/admin/admin-notifications/services/notifications.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { formatDistanceToNow } from 'date-fns'
import { ru } from 'date-fns/locale'
import { BellIcon, CheckCheck } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const { toast } = useToast();
const router = useRouter()
const queryClient = useQueryClient();

const open = ref<boolean>(false)

const filter = ref<GetNotificationsFilter>({
  page: 1,
  pageSize: 5,
  isRead: false,
});

const { data: notificationsResponse } = useQuery({
  queryKey: computed(() => ['admin-notifications', filter.value]),
  queryFn: () => notificationsService.getNotifications(filter.value),
  refetchInterval: 5_000,
});

const { mutate: readMultipleNotifications } = useMutation({
  mutationFn: (dto: MarkNotificationsAsReadDTO) =>
    notificationsService.markNotificationsAsRead(dto),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['admin-notifications'] });
    toast({ title: 'Уведомления успешно прочитаны' });
  },
  onError: () => {
    toast({ title: 'Произошла ошибка при чтении уведомлений' });
  },
});

const notificationIds = computed(() => notificationsResponse.value?.data.map((n) => n.id) || [])

// Mark all unread notifications as read
const readAll = () => {
  if (notificationIds.value.length > 0) {
    readMultipleNotifications({ notificationIds: notificationIds.value });
  }
};

const onSeeAllClick = () => {
  open.value = false
  router.push({name: getRouteName("ADMIN_NOTIFICATIONS")})
}

// Format time into "time ago" format
const formatTimeAgo = (date: string) => formatDistanceToNow(new Date(date), { locale: ru });
</script>

<template>
	<Popover
		:open="open"
		@update:open="open = !open"
	>
		<PopoverTrigger>
			<div class="relative">
				<Button
					variant="outline"
					size="icon"
					class="rounded-md"
				>
					<BellIcon class="text-gray-800 size-4" />
					<span class="sr-only">Переключить уведомления</span>
				</Button>

				<span
					v-if="notificationIds.length > 0"
					class="top-1.5 right-1.5 absolute flex justify-center items-center bg-emerald-500 rounded-full text-white text-xs -translate-y-1/2 translate-x-1/2 size-4"
				>
					{{ notificationIds.length }}
				</span>
			</div>
		</PopoverTrigger>

		<PopoverContent
			align="end"
			class="p-0 w-full max-w-sm"
		>
			<!-- Header -->
			<div class="flex justify-between items-center gap-10 px-4 pt-4">
				<h3 class="font-semibold text-gray-800 text-lg">Уведомления</h3>

				<Button
					:disabled="notificationIds.length === 0"
					@click="readAll"
					variant="link"
					size="sm"
					class="gap-2 p-0 text-primary"
				>
					<CheckCheck class="size-4" />
					Прочитать все
				</Button>
			</div>

			<!-- Notification List -->
			<div class="space-y-4 px-4 py-5">
				<div
					v-for="notification in notificationsResponse?.data ?? []"
					:key="notification.id"
					class="flex items-start space-x-5"
				>
					<!-- Icon -->
					<span class="flex justify-center items-center bg-gray-100 p-2 rounded-full text-gray-500">
						<BellIcon class="size-5" />
					</span>

					<!-- Message and Time -->
					<div class="flex-1">
						<p class="text-sm">
							<span
								class="line-clamp-2 font-medium text-gray-700"
								v-html="formatLocalizedMessage(notification.messages.ru)"
							></span>
						</p>

						<p class="mt-1 text-gray-500 text-xs">
							{{ formatTimeAgo(notification.createdAt) }}
						</p>
					</div>

					<!-- Unread Dot -->
					<span
						v-if="!notification.isRead"
						class="bg-primary rounded-full size-2.5"
					/>
				</div>

				<!-- No Unread Notifications -->
				<p
					v-if="!notificationsResponse?.data?.length"
					class="text-center text-gray-500 text-sm"
				>
					Нет новых уведомлений
				</p>
			</div>

			<!-- Footer -->
			<div class="flex justify-center items-center px-4 pb-4">
				<Button
					size="sm"
					variant="link"
					@click="onSeeAllClick"
				>
					Посмотреть все
				</Button>
			</div>
		</PopoverContent>
	</Popover>
</template>

<style scoped></style>
