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

import { Input } from '@/core/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'
import { formatDistanceToNow } from 'date-fns'
import { ru } from 'date-fns/locale'
import { CheckCheck, SearchIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'

const notifications = ref([
  { id: 1, message: 'Новый комментарий к вашему посту', date: '2025-01-27T09:00:00Z', priority: 'High', isRead: false, isSystem: false },
  { id: 2, message: 'Новый подписчик', date: '2025-01-26T14:00:00Z', priority: 'Medium', isRead: true, isSystem: false },
  { id: 3, message: 'Заказ #12345 был отправлен', date: '2025-01-25T18:00:00Z', priority: 'Low', isRead: false, isSystem: true },
  { id: 4, message: 'Системное уведомление: Обновление завершено', date: '2025-01-24T20:00:00Z', priority: 'Critical', isRead: true, isSystem: true },
])

const PRIORITY_COLORS: Record<string, string> = {
  Critical: 'bg-red-100 text-red-800',
  High: 'bg-yellow-100 text-yellow-800',
  Medium: 'bg-green-100 text-green-800',
  Low: 'bg-gray-100 text-gray-800',
}

const PRIORITY_LABELS: Record<string, string> = {
  Critical: 'Критичный',
  High: 'Высокий',
  Medium: 'Средний',
  Low: 'Низкий',
}

const searchQuery = ref('')
const filter = ref<string>('all')
const selectedNotifications = ref<number[]>([])

const filteredNotifications = computed(() => {
  let result = notifications.value
  if (filter.value === 'unread') {
    result = result.filter(notification => !notification.isRead)
  }
  if (searchQuery.value) {
    result = result.filter(notification =>
      notification.message.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  return result
})

const markSelectedAsRead = () => {
  notifications.value.forEach(notification => {
    if (selectedNotifications.value.includes(notification.id)) {
      notification.isRead = true
    }
  })
  selectedNotifications.value = []
}

const markAsRead = (id: number) => {
  const notification = notifications.value.find(notification => notification.id === id)
  if (notification) {
    notification.isRead = true
  }
}

const toggleSelect = (id: number) => {
  if (selectedNotifications.value.includes(id)) {
    selectedNotifications.value = selectedNotifications.value.filter(selectedId => selectedId !== id)
  } else {
    selectedNotifications.value.push(id)
  }
}

const toggleSelectAll = () => {
  if (selectedNotifications.value.length === filteredNotifications.value.length) {
    selectedNotifications.value = []
  } else {
    selectedNotifications.value = filteredNotifications.value.map(notification => notification.id)
  }
}

const formatTimeAgo = (date: string) => {
  return formatDistanceToNow(new Date(date), { locale: ru })
}
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>Уведомления</CardTitle>
			<CardDescription>Управляйте своими уведомлениями</CardDescription>
		</CardHeader>

		<CardContent>
			<div class="space-y-6">
				<!-- Header -->
				<div class="flex sm:flex-row flex-col sm:justify-between sm:items-center gap-4">
					<div class="relative w-full sm:w-1/2">
						<Input
							v-model="searchQuery"
							placeholder="Поиск уведомлений"
							class="pl-10 max-w-sm"
						/>
						<SearchIcon
							class="top-1/2 left-3 absolute text-gray-500 transform -translate-y-1/2 size-4"
						/>
					</div>
					<div>
						<div class="relative">
							<Select v-model="filter">
								<SelectTrigger>
									<SelectValue placeholder="Фильтр" />
								</SelectTrigger>

								<SelectContent>
									<SelectItem value="all">Все уведомления</SelectItem>
									<SelectItem value="unread">Непрочитанные</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				</div>

				<!-- Notifications Table -->
				<Table class="bg-white shadow-md rounded-lg">
					<TableHeader>
						<TableRow>
							<TableHead class="w-8"></TableHead>
							<TableHead class="w-8">
								<input
									type="checkbox"
									:checked="selectedNotifications.length === filteredNotifications.length"
									@change="toggleSelectAll"
									class="border-gray-300 rounded size-4"
								/>
							</TableHead>
							<TableHead>Сообщение</TableHead>
							<TableHead>Дата</TableHead>
							<TableHead>Приоритет</TableHead>
							<TableHead>Действия</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow
							v-for="notification in filteredNotifications"
							:key="notification.id"
							:class="{ 'bg-blue-50 bg-opacity-30': !notification.isRead }"
						>
							<!-- Dot -->
							<TableCell class="w-8">
								<span
									v-if="!notification.isRead"
									class="block bg-blue-500 mx-auto mb-1 rounded-full w-2 h-2"
								/>
							</TableCell>

							<!-- Select Checkbox -->
							<TableCell class="w-8">
								<input
									type="checkbox"
									:checked="selectedNotifications.includes(notification.id)"
									@change="toggleSelect(notification.id)"
									class="border-gray-300 rounded size-4"
								/>
							</TableCell>

							<!-- Message -->
							<TableCell>
								<p class="font-medium text-gray-800 text-sm">{{ notification.message }}</p>
							</TableCell>

							<!-- Date -->
							<TableCell class="py-4 text-gray-500 text-sm">
								{{ formatTimeAgo(notification.date) }}
							</TableCell>

							<!-- Priority -->
							<TableCell>
								<span
									:class="PRIORITY_COLORS[notification.priority]"
									class="inline-flex items-center px-2.5 py-1 rounded-md w-fit text-xs"
								>
									{{ PRIORITY_LABELS[notification.priority] }}
								</span>
							</TableCell>

							<!-- Actions -->
							<TableCell>
								<Button
									size="xs"
									variant="ghost"
									@click="markAsRead(notification.id)"
									v-if="!notification.isRead"
								>
									<CheckCheck class="text-primary size-5" />
								</Button>
							</TableCell>
						</TableRow>
					</TableBody>
				</Table>

				<!-- Selected Actions -->
				<div
					v-if="selectedNotifications.length"
					class="flex justify-center items-center gap-4"
				>
					<p class="text-gray-700 text-sm">
						{{ `Выбрано: ${selectedNotifications.length}` }}
					</p>
					<Button
						size="sm"
						variant="outline"
						@click="markSelectedAsRead"
						class="gap-2"
					>
						<CheckCheck class="size-4" />
						Прочитать
					</Button>
				</div>
			</div>
		</CardContent>

		<CardFooter class="flex justify-end">
			<PaginationWithMeta
				v-if="categoriesResponse"
				:meta="categoriesResponse.pagination"
				@update:page="updatePage"
				@update:pageSize="updatePageSize"
			/>
		</CardFooter>
	</Card>
</template>
