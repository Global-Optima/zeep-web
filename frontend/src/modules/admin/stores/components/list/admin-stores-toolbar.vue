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
			<!-- Filter Menu -->
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button
						variant="outline"
						class="whitespace-nowrap"
					>
						Фильтр
						<ChevronDown class="ml-2 w-4 h-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem @click="setFranchiseFilter(undefined)">Все</DropdownMenuItem>
					<DropdownMenuItem @click="setFranchiseFilter(true)">Франшиза</DropdownMenuItem>
					<DropdownMenuItem @click="setFranchiseFilter(false)">Не франшиза</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<!-- Right Side: Export and Add Store Buttons -->
		<div class="flex items-center space-x-2 w-full md:w-auto">
			<Button variant="outline"> Экспорт </Button>
			<Button @click="addStore"> Добавить </Button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/core/components/ui/dropdown-menu'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import type { StoresFilter } from '@/modules/stores/models/stores-dto.model'
import { useDebounce } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

// Props
const props = defineProps<{ filter: StoresFilter }>()
const emit = defineEmits(['update:filter'])

// Local copy of the filter to avoid direct mutation
const localFilter = ref({ ...props.filter })

// Search term with debouncing
const searchTerm = ref(localFilter.value.searchTerm || '')
const debouncedSearchTerm = useDebounce(computed(() => searchTerm.value), 500)

// Watch debounced search term and emit updates
watch(debouncedSearchTerm, (newValue) => {
	localFilter.value.searchTerm = newValue
	emit('update:filter', { searchTerm: newValue.trim() })
})

// Update franchise filter
function setFranchiseFilter(isFranchise: boolean | undefined) {
	localFilter.value.isFranchise = isFranchise
	emit('update:filter', { isFranchise })
}

// Navigate to add store page
const router = useRouter()
const addStore = () => {
	router.push({ name: getRouteName('ADMIN_STORE_CREATE') })
}
</script>
