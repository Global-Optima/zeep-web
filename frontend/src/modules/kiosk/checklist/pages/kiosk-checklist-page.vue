<script setup lang="ts">
import { Button } from '@/core/components/ui/button'
import { getRouteName } from '@/core/config/routes.config'
import { getKaspiConfig } from '@/core/integrations/kaspi.service'
import { ArrowRight, CheckCircle, XCircle } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const checklistItems = ref([
  {
    key: 'kaspiConfig',
    title: 'Настройки Kaspi',
    description: 'Конфигурация платежной системы Kaspi.',
    checkFn: () => Boolean(getKaspiConfig()),
    route: getRouteName("ADMIN_STORE_SETTINGS"),
    checked: false
  },
]);

const allItemsChecked = computed(() => checklistItems.value.every(item => item.checkFn()));
const router = useRouter();

function goToPage(route: string) {
  router.push({ name: route });
}

onMounted(() => {
  checklistItems.value = checklistItems.value.map(item => ({
    ...item,
    checked: item.checkFn(),
  }));
});
</script>

<template>
	<div class="flex flex-col justify-center items-center bg-gray-50 px-6 min-h-screen">
		<div class="w-full max-w-2xl">
			<h1 class="font-semibold text-gray-900 text-3xl text-center">Проверка готовности киоска</h1>
			<p class="mt-3 text-gray-600 text-center">
				Перед началом работы убедитесь, что все пункты настроены.
			</p>

			<div class="flex flex-col gap-2 mt-6">
				<div
					v-for="item in checklistItems"
					:key="item.key"
					class="flex justify-between items-center bg-white shadow-sm px-6 py-4 border rounded-lg"
				>
					<div class="flex items-center space-x-4">
						<component
							:is="item.checked ? CheckCircle : XCircle"
							class="w-6 h-6"
							:class="item.checked ? 'text-green-500' : 'text-red-500'"
						/>
						<div>
							<h2 class="font-medium text-gray-900 text-lg">{{ item.title }}</h2>
							<p class="text-gray-600 text-sm">{{ item.description }}</p>
						</div>
					</div>
					<Button
						@click="goToPage(item.route)"
						variant="outline"
					>
						<ArrowRight class="w-4 h-4" />
					</Button>
				</div>
			</div>

			<Button
				v-if="allItemsChecked"
				@click="router.push({ name: getRouteName('KIOSK_HOME') })"
				class="mt-6 w-full"
			>
				Перейти в киоск
			</Button>
		</div>
	</div>
</template>
