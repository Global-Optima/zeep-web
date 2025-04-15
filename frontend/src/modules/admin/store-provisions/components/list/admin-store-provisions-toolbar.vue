<template>
	<div
		class="flex md:flex-row flex-col justify-between items-start md:items-center space-y-4 md:space-y-0 mb-4"
	>
		<!-- Left Side: Search Input and Filter Menu -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<!-- Search Input -->
			<Input
				v-model="searchTerm"
				placeholder="Поиск"
				type="search"
				class="bg-white w-full md:w-64"
			/>

      <MultiSelectFilter
        title="Статусы"
        :options="statusOptions"
        v-model="selectedStatuses"
      />
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button
				variant="outline"
				disabled
			>
				Экспорт
			</Button>
			<Button @click="addStore"> Добавить </Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import {Button} from '@/core/components/ui/button'
import {Input} from '@/core/components/ui/input'
import {getRouteName} from '@/core/config/routes.config'
import {DEFAULT_PAGINATION_META} from '@/core/utils/pagination.utils'
import {
  type StoreProvisionFilter,
  StoreProvisionStatus
} from '@/modules/admin/store-provisions/models/store-provision.models'
import {useDebounce} from '@vueuse/core'
import {computed, ref, watch} from 'vue'
import {useRouter} from 'vue-router'
import MultiSelectFilter, {type MultiSelectOption} from "@/core/components/multi-select-filter/MultiSelectFilter.vue";

const props = defineProps<{ filter: StoreProvisionFilter }>()
const emit = defineEmits<{
  (event: 'update:filter', value: StoreProvisionFilter): void
}>()

const router = useRouter()

const localFilter = ref({ ...props.filter })

const searchTerm = ref(localFilter.value.search || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

watch(debouncedSearchTerm, (newValue) => {
  localFilter.value.search = newValue
  emit('update:filter', {
    ...localFilter.value,
    search: newValue.trim(),
    page: DEFAULT_PAGINATION_META.page,       // Reset to default page (e.g., 1)
    pageSize: DEFAULT_PAGINATION_META.pageSize  // Reset to default page size
  })
})

const statusOptions: MultiSelectOption[] = [
  { label: 'Готовится', value: StoreProvisionStatus.PREPARING },
  { label: 'Приготовлен', value: StoreProvisionStatus.COMPLETED },
  { label: 'Просрочен', value: StoreProvisionStatus.EXPIRED },
  { label: 'Пустой', value: StoreProvisionStatus.EMPTY },
];

const selectedStatuses = ref<StoreProvisionStatus[]>(props.filter?.statuses ?? [])

watch(selectedStatuses, (newStatuses) => {
  emit('update:filter', {
    ...localFilter.value,
    statuses: newStatuses.length ? newStatuses : undefined,
    page: DEFAULT_PAGINATION_META.page,
    pageSize: DEFAULT_PAGINATION_META.pageSize
  })
})

const addStore = () => {
	router.push({ name: getRouteName('ADMIN_STORE_PROVISION_CREATE') })
}
</script>
