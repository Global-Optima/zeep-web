<template>
	<div
		v-if=" !isStoreSyncPending && isStoreSyncResponse && !isStoreSyncResponse.isSync"
		class="mb-4"
	>
		<AdminStoresSyncCard />
	</div>

	<AdminStoreAdditivesToolbar
		:filter="filter"
		@update:filter="updateFilter"
	/>

	<AdminListLoader v-if="isLoading" />

	<div v-else>
		<Card>
			<CardContent class="mt-4">
				<p
					v-if="!storeAdditivesResponse || storeAdditivesResponse.data.length === 0"
					class="text-muted-foreground"
				>
					Модификаторы кафе не найдены
				</p>

				<AdminStoreAdditivesList
					v-else
					:additives="storeAdditivesResponse.data"
				/>
			</CardContent>
			<CardFooter class="flex justify-end">
				<PaginationWithMeta
					v-if="storeAdditivesResponse"
					:meta="storeAdditivesResponse.pagination"
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
import { Card, CardContent } from '@/core/components/ui/card'
import CardFooter from '@/core/components/ui/card/CardFooter.vue'
import { useHasRole } from "@/core/hooks/use-has-roles.hook"
import { usePaginationFilter } from '@/core/hooks/use-pagination-filter.hook'
import type { AdditiveFilterQuery } from '@/modules/admin/additives/models/additives.model'
import { EmployeeRole } from "@/modules/admin/employees/models/employees.models"
import AdminStoreAdditivesList from '@/modules/admin/store-additives/components/list/admin-store-additives-list.vue'
import AdminStoreAdditivesToolbar from '@/modules/admin/store-additives/components/list/admin-store-additives-toolbar.vue'
import { storeAdditivesService } from '@/modules/admin/store-additives/services/store-additives.service'
import AdminStoresSyncCard from '@/modules/admin/store-products/components/stores-sync/admin-stores-sync-card.vue'
import { storeSyncService } from '@/modules/admin/stores/services/stores-sync.service'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<AdditiveFilterQuery>({})
const isFranchisee = useHasRole([EmployeeRole.FRANCHISEE_MANAGER, EmployeeRole.FRANCHISEE_OWNER])
const isStore = useHasRole([EmployeeRole.STORE_MANAGER, EmployeeRole.BARISTA])

const { data: storeAdditivesResponse, isLoading } = useQuery({
  queryKey: computed(() => ['admin-store-additives', filter.value]),
  queryFn: () => storeAdditivesService.getStoreAdditives(filter.value),
  enabled: computed(() =>
    isStore.value || (isFranchisee.value && Boolean(filter.value.storeId))
  )
})

const { data: isStoreSyncResponse, isLoading: isStoreSyncPending } = useQuery({
  queryKey: ['admin-store-is-sync'],
  queryFn: () => storeSyncService.isStoreSynchronized(),
  enabled: computed(() =>
    isStore.value || (isFranchisee.value && Boolean(filter.value.storeId))
  )
})
</script>

<style scoped></style>
