<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/core/components/ui/popover'
import { BellIcon, CheckCheck } from 'lucide-vue-next'
import { ref } from 'vue'

const notifications = ref([
  { id: 1, message: 'Новый комментарий к вашему посту', icon: '📩', timeAgo: '2 минуты назад', isRead: false },
  { id: 2, message: 'Новый подписчик', icon: '👤', timeAgo: '5 минут назад', isRead: false },
  { id: 3, message: 'Заказ #12345 был отправлен', icon: '📦', timeAgo: '10 минут назад', isRead: true },
])

const readAll = () => {
  console.log('Отметить все уведомления как прочитанные')
  notifications.value.forEach(notification => {
    notification.isRead = true
  })
}
</script>

<template>
	<Popover>
		<PopoverTrigger>
			<Button
				variant="outline"
				size="icon"
				class="rounded-full"
			>
				<BellIcon class="text-gray-800 size-4" />
				<span class="sr-only">Переключить уведомления</span>
			</Button>
		</PopoverTrigger>

		<PopoverContent
			align="end"
			class="p-0 w-fit"
		>
			<!-- Header -->
			<div class="flex justify-between items-center px-4 pt-4">
				<h3 class="font-semibold text-gray-800 text-lg">Уведомления</h3>

				<Button
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
					v-for="notification in notifications"
					:key="notification.id"
					class="flex items-start space-x-4"
				>
					<!-- Icon -->
					<span class="flex justify-center items-center bg-gray-100 p-1 rounded-full text-2xl">
						{{ notification.icon }}
					</span>

					<!-- Message and Time -->
					<div class="flex-1">
						<p class="line-clamp-2 font-medium text-gray-700 text-sm">{{ notification.message }}</p>
						<p class="text-gray-500 text-xs">{{ notification.timeAgo }}</p>
					</div>

					<!-- Unread Dot -->
					<span
						v-if="!notification.isRead"
						class="bg-primary rounded-full size-2.5"
					/>
				</div>
			</div>

			<!-- Footer -->
			<div class="flex justify-center items-center px-4 pb-4">
				<Button
					size="sm"
					variant="link"
					@click='$router.push("/admin/notifications")'
				>
					Посмотреть все
				</Button>
			</div>
		</PopoverContent>
	</Popover>
</template>

<style scoped></style>
