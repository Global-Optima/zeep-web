<template>
	<AdminUnitsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isPending" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!unitsResponse || unitsResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Единицы измерения не найдены
				</p>
				<AdminUnitsList
					v-else
					:units="unitsResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="unitsResponse"
					:meta="unitsResponse.pagination"
					@update:page="updatePage"
					@update:pageSize="updatePageSize"
				/>
			</CardFooter>
		</Card>
	</div>
</template>

<script setup lang="ts">
import AdminListLoader from '@/core/components/admin-list-loader/AdminListLoader.vue'
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import AdminUnitsList from '@/modules/admin/units/components/list/admin-units-list.vue'
import AdminUnitsToolbar from '@/modules/admin/units/components/list/admin-units-toolbar.vue'
import type { UnitsFilterDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<UnitsFilterDTO>({});

const { data: unitsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-units', filter.value]),
  queryFn: () => unitsService.getAllUnits(filter.value),
})
</script>

<style scoped></style>
