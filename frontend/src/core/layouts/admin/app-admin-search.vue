<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '@/core/components/ui/command'
import { getRouteName, type RouteKey } from '@/core/config/routes.config'
import { adminSearchItems } from '@/core/layouts/admin/app-admin-search-links.constants'
import type { EmployeeDTO } from '@/modules/admin/store-employees/models/employees.models'
import { Search } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const {currentEmployee} = defineProps<{currentEmployee: EmployeeDTO}>()

const router = useRouter()

const pages = computed(() => adminSearchItems.filter(i => i.accessRoles.includes(currentEmployee.role)))
const open = ref(false)

const onRouteClick = (routeKey: RouteKey) => {
  router.push({name: getRouteName(routeKey)})
  open.value = false

}
</script>

<template>
	<div>
		<Button
			variant="outline"
			size="sm"
			@click="open = true"
			class="gap-2"
		>
			<Search class="w-3.5 h-3.5" />
			Поиск
		</Button>

		<CommandDialog v-model:open="open">
			<CommandInput placeholder="Type a command or search..." />
			<CommandList>
				<CommandEmpty>No results found.</CommandEmpty>
				<CommandGroup heading="Страницы">
					<CommandItem
						v-for="page in pages"
						:key="page.routeKey"
						:value="page.name"
					>
						<div
							class="flex items-center gap-3 w-full"
							@click="onRouteClick(page.routeKey)"
						>
							<component
								:is="page.icon"
								class="w-5 h-5 text-gray-600"
							/>
							{{ page.name }}
						</div>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</CommandDialog>
	</div>
</template>
