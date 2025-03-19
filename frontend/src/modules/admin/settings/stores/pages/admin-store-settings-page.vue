// src/components/StoreSettings.vue
<template>
	<Card class="mx-auto max-w-7xl">
		<CardContent>
			<div class="flex py-4 h-full">
				<aside class="pr-8 border-gray-200 border-r">
					<h3 class="font-medium text-2xl">Настройки</h3>

					<ul class="space-y-4 mt-4">
						<li
							v-for="tab in tabs"
							:key="tab.id"
							@click="activeTab = tab.id"
							class="text-gray-500 transition cursor-pointer"
							:class="{'text-primary': activeTab === tab.id}"
						>
							{{ tab.name }}
						</li>
					</ul>
				</aside>

				<div class="flex-1 pl-8 w-full">
					<component :is="activeComponent" />
				</div>
			</div>
		</CardContent>
	</Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/core/components/ui/card'
import { computed, defineAsyncComponent, ref, shallowRef } from 'vue'

const tabs = shallowRef([
  { id: 'kaspi', name: 'Каспи Интеграция', component: defineAsyncComponent(() => import('../components/admin-store-settings-kaspi.vue')) },
  { id: 'kaspiTEst', name: 'Каспи Тест', component: defineAsyncComponent(() => import('../components/admin-store-settings-kaspi-test.vue')) },
]);

const activeTab = ref(tabs.value[0].id);
const activeComponent = computed(() => tabs.value.find(tab => tab.id === activeTab.value)?.component || null);
</script>

<style scoped></style>
