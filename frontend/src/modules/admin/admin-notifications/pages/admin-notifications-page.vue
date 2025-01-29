<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { useToast } from '@/core/components/ui/toast'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminNotificationsList from '@/modules/admin/admin-notifications/components/list/admin-notifications-list.vue'
import AdminNotificationsToolbar from '@/modules/admin/admin-notifications/components/list/admin-notifications-toolbar.vue'
import type {
  GetNotificationsFilter,
  MarkNotificationAsReadDTO,
  MarkNotificationsAsReadDTO,
} from '@/modules/admin/admin-notifications/models/notifications.model'
import { notificationsService } from '@/modules/admin/admin-notifications/services/notifications.service'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<GetNotificationsFilter>({});

const selectedNotificationIds = ref<number[]>([]);

const { data: notificationsResponse } = useQuery({
  queryKey: computed(() => ['admin-notifications', filter.value]),
  queryFn: () => notificationsService.getNotifications(filter.value),
});

const { toast } = useToast();
const queryClient = useQueryClient();

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

const { mutate: readNotification } = useMutation({
  mutationFn: (dto: MarkNotificationAsReadDTO) =>
    notificationsService.markNotificationAsRead(dto),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['admin-notifications'] });
    toast({ title: 'Уведомление успешно прочитано' });
  },
  onError: () => {
    toast({ title: 'Произошла ошибка при чтении уведомления' });
  },
});

const updateFilter = (updatedFilter: Partial<GetNotificationsFilter>) => {
  filter.value = { ...filter.value, ...updatedFilter };
};

const toggleSelect = (id: number) => {
  const index = selectedNotificationIds.value.indexOf(id);
  if (index > -1) {
    selectedNotificationIds.value.splice(index, 1);
  } else {
    selectedNotificationIds.value.push(id);
  }
};

const toggleSelectAll = () => {
  if (selectedNotificationIds.value.length === notificationsResponse.value?.data.length) {
    selectedNotificationIds.value = [];
  } else {
    selectedNotificationIds.value =
      notificationsResponse.value?.data.map((n) => n.id) ?? [];
  }
};

const markAsRead = (id: number) => {
  readNotification({ notificationId: id });
};

const markSelectedAsRead = () => {
  if (selectedNotificationIds.value.length) {
    readMultipleNotifications({ notificationIds: selectedNotificationIds.value });
    selectedNotificationIds.value = [];
  }
};
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>Уведомления</CardTitle>
			<CardDescription>Управляйте своими уведомлениями</CardDescription>
		</CardHeader>

		<CardContent>
			<!-- Toolbar -->
			<AdminNotificationsToolbar
				:filter="filter"
				@update:filter="updateFilter"
			/>

			<!-- Notification List -->
			<AdminNotificationsList
				:notifications="notificationsResponse?.data || []"
				:selectedNotifications="selectedNotificationIds"
				@toggleSelect="toggleSelect"
				@toggleSelectAll="toggleSelectAll"
				@markAsRead="markAsRead"
			/>

			<!-- Bulk Actions -->
			<div
				v-if="selectedNotificationIds.length"
				class="flex justify-center items-center gap-5 mt-6"
			>
				<p class="text-gray-700 text-sm">Выбрано: {{ selectedNotificationIds.length }}</p>
				<Button
					size="sm"
					variant="outline"
					@click="markSelectedAsRead"
					>Прочитать</Button
				>
			</div>
		</CardContent>

		<!-- Pagination -->
		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="notificationsResponse"
				:meta="notificationsResponse.pagination"
				@update:page="updateFilter({ page: $event, pageSize: DEFAULT_PAGINATION_META.pageSize })"
				@update:pageSize="updateFilter({ pageSize: $event , page: DEFAULT_PAGINATION_META.page  })"
			/>
		</CardFooter>
	</Card>
</template>
