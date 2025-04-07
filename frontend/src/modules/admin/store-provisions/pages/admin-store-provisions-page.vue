<template>
  <AdminProvisionsToolbar
    :filter="filter"
    @update:filter="updateFilter"
  />

  <AdminListLoader v-if="isPending" />

  <div v-else>
    <Card>
      <CardContent class="mt-4">
        <!-- No Data Message -->
        <p
          v-if="!provisionsResponse || provisionsResponse.data.length === 0"
          class="text-muted-foreground"
        >
          Заготовки не найдены
        </p>

        <!-- Ingredients List -->
        <AdminProvisionsList
          v-else
          :provisions="provisionsResponse.data"
        />
      </CardContent>
      <CardFooter class="flex justify-end">
        <PaginationWithMeta
          v-if="provisionsResponse"
          :meta="provisionsResponse.pagination"
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
import { usePaginationFilter } from "@/core/hooks/use-pagination-filter.hook"
import AdminProvisionsList from "@/modules/admin/provisions/components/list/admin-provisions-list.vue"
import AdminProvisionsToolbar from "@/modules/admin/provisions/components/list/admin-provisions-toolbar.vue"
import type { ProvisionFilter } from "@/modules/admin/provisions/models/provision.models"
import { provisionsService } from "@/modules/admin/provisions/services/provisions.service"
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { filter, updateFilter, updatePage, updatePageSize } = usePaginationFilter<ProvisionFilter>({})

// Fetch data using Vue Query
const { data: provisionsResponse, isPending } = useQuery({
  queryKey: computed(() => ['admin-store-provisions', filter.value]),
  queryFn: () => provisionsService.getProvisions(filter.value),
})
</script>

<style scoped></style>
