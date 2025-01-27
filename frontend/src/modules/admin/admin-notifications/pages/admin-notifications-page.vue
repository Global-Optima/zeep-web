<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { BellIcon, CheckCheck, SearchIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'

const notifications = ref([
  { id: 1, message: 'Новый комментарий к вашему посту', date: '2025-01-27', priority: 'High', isRead: false, isSystem: false },
  { id: 2, message: 'Новый подписчик', date: '2025-01-26', priority: 'Medium', isRead: true, isSystem: false },
  { id: 3, message: 'Заказ #12345 был отправлен', date: '2025-01-25', priority: 'Low', isRead: false, isSystem: true },
  { id: 4, message: 'Системное уведомление: Обновление завершено', date: '2025-01-24', priority: 'Critical', isRead: true, isSystem: true },
])

const searchQuery = ref('')
const filter = ref<string>('all') // 'all' or 'unread'
const selectedNotifications = ref<number[]>([])

// Filtered notifications based on search and filter
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

// Mark all selected notifications as read
const markSelectedAsRead = () => {
  notifications.value.forEach(notification => {
    if (selectedNotifications.value.includes(notification.id)) {
      notification.isRead = true
    }
  })
  selectedNotifications.value = [] // Clear selection after marking as read
}

// Mark a single notification as read
const markAsRead = (id: number) => {
  const notification = notifications.value.find(notification => notification.id === id)
  if (notification) {
    notification.isRead = true
  }
}

// Toggle selection of notifications
const toggleSelect = (id: number) => {
  if (selectedNotifications.value.includes(id)) {
    selectedNotifications.value = selectedNotifications.value.filter(selectedId => selectedId !== id)
  } else {
    selectedNotifications.value.push(id)
  }
}
</script>

<template>
	<div class="space-y-6 p-6">
		<!-- Header -->
		<div class="flex justify-between items-center">
			<h1 class="font-semibold text-2xl text-gray-800">Уведомления</h1>
			<div class="flex space-x-4">
				<Button
					size="sm"
					variant="outline"
					@click="markSelectedAsRead"
					:disabled="selectedNotifications.length === 0"
				>
					<CheckCheck class="size-4" />
					Отметить выбранные как прочитанные
				</Button>
				<Button
					size="sm"
					variant="outline"
					@click="() => (filter = 'all')"
					:class="{ 'bg-blue-100 text-blue-700': filter === 'all' }"
				>
					Все
				</Button>
				<Button
					size="sm"
					variant="outline"
					@click="() => (filter = 'unread')"
					:class="{ 'bg-blue-100 text-blue-700': filter === 'unread' }"
				>
					Непрочитанные
				</Button>
			</div>
		</div>

		<!-- Search Bar -->
		<div class="relative">
			<Input
				v-model="searchQuery"
				placeholder="Поиск уведомлений"
				class="pl-10"
			/>
			<SearchIcon
				class="top-1/2 left-3 absolute w-5 h-5 text-gray-500 transform -translate-y-1/2"
			/>
		</div>

		<!-- Notifications List -->
		<div class="space-y-4">
			<div
				v-for="notification in filteredNotifications"
				:key="notification.id"
				class="flex items-center space-x-4 hover:bg-gray-100 shadow-sm p-4 border rounded-lg"
			>
				<!-- Icon -->
				<div
					class="flex justify-center items-center rounded-full w-10 h-10"
					:class="notification.isSystem ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-700'"
				>
					<BellIcon class="w-6 h-6" />
				</div>

				<!-- Content -->
				<div class="flex-1">
					<p class="font-medium text-gray-800 text-sm">{{ notification.message }}</p>
					<p class="text-gray-500 text-xs">{{ notification.date }}</p>
					<p class="text-gray-500 text-xs">Приоритет: {{ notification.priority }}</p>
				</div>

				<!-- Unread Dot -->
				<div
					v-if="!notification.isRead"
					class="bg-primary rounded-full w-2.5 h-2.5"
				></div>

				<!-- Actions -->
				<div class="flex items-center space-x-2">
					<input
						type="checkbox"
						:checked="selectedNotifications.includes(notification.id)"
						@change="toggleSelect(notification.id)"
						class="border-gray-300 rounded"
					/>
					<Button
						size="sm"
						variant="ghost"
						@click="markAsRead(notification.id)"
					>
						<CheckCheck class="w-4 h-4 text-primary" />
					</Button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
