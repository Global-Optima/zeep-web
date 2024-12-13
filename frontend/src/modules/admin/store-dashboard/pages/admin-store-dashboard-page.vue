<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import AdminDashboardDateRangePicker from '@/modules/admin/dashboard/components/admin-dashboard-date-range-picker.vue'
import AdminDashboardOverview from '@/modules/admin/dashboard/components/admin-dashboard-overview.vue'
import AdminDashboardRecentSales from '@/modules/admin/dashboard/components/admin-dashboard-recent-sales.vue'

// Import Lucide icons
import { Activity, DollarSign, ShoppingCart, Users } from 'lucide-vue-next'

// Constant array for statistic cards
const statisticCards = [
  {
    title: 'Общая выручка',
    value: '$45,231.89',
    description: '+20.1% с прошлого месяца',
    icon: DollarSign,
  },
  {
    title: 'Подписки',
    value: '+2350',
    description: '+180.1% с прошлого месяца',
    icon: Users,
  },
  {
    title: 'Продажи',
    value: '+12,234',
    description: '+19% с прошлого месяца',
    icon: ShoppingCart,
  },
  {
    title: 'Активные сейчас',
    value: '+573',
    description: '+201 за последний час',
    icon: Activity,
  },
]
</script>

<template>
	<div class="flex-1 space-y-4">
		<!-- Top Bar -->
		<div class="flex items-center space-x-2">
			<AdminDashboardDateRangePicker />
			<Button>Скачать</Button>
		</div>

		<!-- Statistic Cards -->
		<div class="space-y-4">
			<div class="gap-4 grid md:grid-cols-2 lg:grid-cols-4">
				<Card
					v-for="(card, index) in statisticCards"
					:key="index"
				>
					<CardHeader class="flex flex-row justify-between items-center space-y-0 pb-2">
						<CardTitle class="font-medium text-sm">
							{{ card.title }}
						</CardTitle>
						<component
							:is="card.icon"
							class="w-4 h-4 text-muted-foreground"
						/>
					</CardHeader>
					<CardContent>
						<div class="font-bold text-2xl">
							{{ card.value }}
						</div>
						<p class="text-muted-foreground text-xs">
							{{ card.description }}
						</p>
					</CardContent>
				</Card>
			</div>

			<!-- Overview and Recent Sales -->
			<div class="gap-4 grid md:grid-cols-2 lg:grid-cols-7">
				<Card class="col-span-4">
					<CardHeader>
						<CardTitle>Обзор</CardTitle>
					</CardHeader>
					<CardContent class="pl-2">
						<AdminDashboardOverview />
					</CardContent>
				</Card>
				<Card class="col-span-3">
					<CardHeader>
						<CardTitle>Последние продажи</CardTitle>
						<CardDescription> Вы сделали 265 продаж в этом месяце. </CardDescription>
					</CardHeader>
					<CardContent>
						<AdminDashboardRecentSales />
					</CardContent>
				</Card>
			</div>
		</div>
	</div>
</template>
