<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/core/components/ui/table'
import type { NotificationDTO } from '@/modules/admin/admin-notifications/models/notifications.model'
import { formatDistanceToNow } from 'date-fns'
import { ru } from 'date-fns/locale'
import { Check } from 'lucide-vue-next'

const {notifications} = defineProps<{
  notifications: NotificationDTO[];
  selectedNotifications: number[];
}>();

const emit = defineEmits(['toggleSelect', 'toggleSelectAll', 'markAsRead']);

const PRIORITY_COLORS = {
  HIGH: 'bg-yellow-100 text-yellow-800',
  MEDIUM: 'bg-green-100 text-green-800',
  LOW: 'bg-gray-100 text-gray-800',
};

const PRIORITY_LABELS = {
  HIGH: 'Высокий',
  MEDIUM: 'Средний',
  LOW: 'Низкий',
};

const formatTimeAgo = (date: string) => formatDistanceToNow(new Date(date), { locale: ru });
const formatMessage = (s: string) => s.replace(/\*(.*?)\*/g, '<span class="font-medium">$1</span>');
</script>

<template>
	<Table class="bg-white shadow-md mt-6 rounded-lg">
		<TableHeader>
			<TableRow>
				<TableHead class="w-10"></TableHead>
				<TableHead class="w-10">
					<input
						type="checkbox"
						:checked="selectedNotifications.length === notifications.length"
						@change="emit('toggleSelectAll')"
						class="size-4"
					/>
				</TableHead>
				<TableHead>Сообщение</TableHead>
				<TableHead>Дата</TableHead>
				<TableHead>Приоритет</TableHead>
				<TableHead></TableHead>
			</TableRow>
		</TableHeader>
		<TableBody>
			<TableRow
				v-for="notification in notifications"
				:key="notification.id"
				:class="{ 'bg-emerald-50 bg-opacity-30': !notification.isRead }"
			>
				<TableCell class="w-8">
					<span
						v-if="!notification.isRead"
						class="block bg-emerald-500 mx-auto mb-1 rounded-full size-2.5"
					/>
				</TableCell>
				<TableCell>
					<input
						type="checkbox"
						:checked="selectedNotifications.includes(notification.id)"
						@change="emit('toggleSelect', notification.id)"
						class="size-4"
					/>
				</TableCell>
				<TableCell class="py-4">
					<span v-html="formatMessage(notification.messages.ru)"></span>
				</TableCell>
				<TableCell>{{ formatTimeAgo(notification.createdAt) }}</TableCell>
				<TableCell>
					<span
						:class="PRIORITY_COLORS[notification.priority]"
						class="px-2 py-1 rounded-md text-xs"
					>
						{{ PRIORITY_LABELS[notification.priority] }}
					</span>
				</TableCell>
				<TableCell class="flex justify-center">
					<Button
						size="xs"
						variant="ghost"
						@click="emit('markAsRead', notification.id)"
						v-if="!notification.isRead"
						class=""
					>
						<Check class="text-primary size-5" />
					</Button>
				</TableCell>
			</TableRow>
		</TableBody>
	</Table>
</template>
