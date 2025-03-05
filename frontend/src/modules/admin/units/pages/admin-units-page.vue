<template>
	<AdminUnitsToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<Card>
		<CardContent class="mt-4">
			<p
				v-if="!unitsResponse || unitsResponse.data.length === 0"
				class="text-muted-foreground"
			>
				Категории топпингов не найдены
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
</template>

<script setup lang="ts">
import PaginationWithMeta from '@/core/components/ui/app-pagination/PaginationWithMeta.vue'
import { Card, CardContent, CardFooter } from '@/core/components/ui/card'
import { DEFAULT_PAGINATION_META } from '@/core/utils/pagination.utils'
import AdminUnitsList from '@/modules/admin/units/components/list/admin-units-list.vue'
import AdminUnitsToolbar from '@/modules/admin/units/components/list/admin-units-toolbar.vue'
import type { UnitsFilterDTO } from '@/modules/admin/units/models/units.model'
import { unitsService } from '@/modules/admin/units/services/units.service'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const filter = ref<UnitsFilterDTO>({})

const { data: unitsResponse } = useQuery({
  queryKey: computed(() => ['admin-units', filter.value]),
  queryFn: () => unitsService.getAllUnits(filter.value),
})

function updateFilter(updatedFilter: UnitsFilterDTO) {
  filter.value = {...filter.value, ...updatedFilter}
}

function updatePage(page: number) {
  updateFilter({ pageSize: DEFAULT_PAGINATION_META.pageSize, page: page})

}

function updatePageSize(pageSize: number) {
  updateFilter({ pageSize: pageSize, page: DEFAULT_PAGINATION_META.page})
}
</script>

<style scoped></style>
