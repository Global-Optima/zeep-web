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
import type { EmployeeDTO } from '@/modules/admin/employees/models/employees.models'
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
			size="icon"
			class="rounded-md"
			@click="open=true"
		>
			<Search class="size-[16px]" />
			<span class="sr-only">Toggle user menu</span>
		</Button>

		<CommandDialog v-model:open="open">
			<CommandInput
				placeholder="Поиск по платформе"
				class="py-7 text-lg"
			/>
			<CommandList class="p-1">
				<CommandEmpty class="text-gray-500 text-lg">Ничего не найдено</CommandEmpty>
				<CommandGroup
					heading="Страницы"
					class="!text-base"
				>
					<CommandItem
						v-for="page in pages"
						:key="page.routeKey"
						:value="page.name"
						@click="onRouteClick(page.routeKey)"
					>
						<div class="flex items-center gap-4 py-1.5 w-full text-base">
							{{ page.name }}
						</div>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</CommandDialog>
	</div>
</template>
