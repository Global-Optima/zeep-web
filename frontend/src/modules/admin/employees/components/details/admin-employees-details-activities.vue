<template>
	<div>
		<!-- Display Message if No Activities -->
		<p
			v-if="activities.length === 0"
			class="py-4 text-center text-muted-foreground"
		>
			Активности отсутствуют
		</p>

		<!-- Table of Activities -->
		<Table
			v-else
			class="bg-white rounded-xl"
		>
			<TableHeader>
				<TableRow>
					<TableHead>Дата</TableHead>
					<TableHead>Активность</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				<TableRow
					v-for="(activity, index) in activities"
					:key="index"
					class="h-12"
				>
					<!-- Date -->
					<TableCell>
						{{ activity.date }}
					</TableCell>
					<!-- Activity -->
					<TableCell>
						{{ activity.activity }}
					</TableCell>
				</TableRow>
			</TableBody>
		</Table>
	</div>
</template>

<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'

interface Activity {
  date: string;
  activity: string;
  status: string;
}

// Props for activities
defineProps<{
  activities: Activity[];
}>();
</script>

<style scoped>
/* Add any additional custom styles here if necessary */
</style>
