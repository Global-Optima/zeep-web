<script setup lang="ts">
import { computed, ref, type Ref } from 'vue'

// UI Components (shadcn or custom)
import { Button } from '@/core/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/core/components/ui/table'

// Icons
import type { AdditiveDTO } from '@/modules/admin/additives/models/additives.model'
import AdminSelectAvailableToAddAdditiveDialog from '@/modules/admin/store-additives/components/admin-select-available-to-add-additive-dialog.vue'
import type { CreateStoreAdditiveDTO } from '@/modules/admin/store-additives/models/store-additves.model'
import { ChevronLeft, Trash } from 'lucide-vue-next'


// Props & Emits
const emits = defineEmits<{
  onSubmit: [dto: CreateStoreAdditiveDTO[]]
  onCancel: []
}>()

// Local reactive state to manage selected additives
const additivesList: Ref<(AdditiveDTO & { storePrice: number })[]> = ref([])

// State to control the additive selection dialog
const openAdditiveDialog = ref(false)

// ----- Computed Properties -----
/** Check if there are any selected additives */
const hasAdditives = computed(() => additivesList.value.length > 0)

// ----- Methods -----
/** Submit the additives as `CreateStoreAdditiveDTO[]` */
function onSubmit() {
  const createDTO: CreateStoreAdditiveDTO[] = additivesList.value.map((additive) => ({
    additiveId: additive.id,
    storePrice: additive.storePrice,
  }))
  emits('onSubmit', createDTO)
}

/** Cancel the process and emit `onCancel` */
function onCancel() {
  additivesList.value = []
  emits('onCancel')
}

/** Called when an additive is selected in the dialog */
function selectAdditive(additive: AdditiveDTO) {
  // Check if the additive is already added
  const exists = additivesList.value.some((a) => a.id === additive.id)
  if (exists) {
    // Optional: Notify the user (e.g., toast) about the duplicate entry
    openAdditiveDialog.value = false
    return
  }

  // Add the selected additive with the default storePrice as the basePrice
  additivesList.value.push({ ...additive, storePrice: additive.basePrice })
  openAdditiveDialog.value = false
}

/** Remove an additive from the list */
function removeAdditive(index: number) {
  additivesList.value.splice(index, 1)
}
</script>

<template>
	<div class="flex-1 gap-4 grid auto-rows-max mx-auto max-w-4xl">
		<!-- Header -->
		<div class="flex items-center gap-4">
			<Button
				variant="outline"
				size="icon"
				type="button"
				@click="onCancel"
			>
				<ChevronLeft class="w-5 h-5" />
				<span class="sr-only">Назад</span>
			</Button>
			<h1 class="font-semibold text-xl tracking-tight shrink-0">Добавить добавки в кафе</h1>
			<div class="md:flex items-center gap-2 hidden md:ml-auto">
				<Button
					variant="outline"
					type="button"
					@click="onCancel"
				>
					Отменить
				</Button>
				<Button
					type="submit"
					@click="onSubmit"
				>
					Сохранить
				</Button>
			</div>
		</div>

		<!-- Table of selected additives -->
		<Card>
			<CardHeader>
				<div class="flex justify-between">
					<CardTitle>Добавки</CardTitle>
					<Button
						type="button"
						variant="outline"
						@click="openAdditiveDialog = true"
					>
						Добавить добавку
					</Button>
				</div>
				<CardDescription class="mt-1"> Список добавленных добавок и их цены </CardDescription>
			</CardHeader>

			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Название</TableHead>
							<TableHead>Базовая цена</TableHead>
							<TableHead>Цена в кафе</TableHead>
							<TableHead>Удалить</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<template v-if="hasAdditives">
							<TableRow
								v-for="(additive, index) in additivesList"
								:key="additive.id"
							>
								<TableCell>{{ additive.name }}</TableCell>
								<TableCell>{{ additive.basePrice }}</TableCell>
								<TableCell>
									<Input
										type="number"
										class="w-24"
										placeholder="Цена"
										v-model.number="additive.storePrice"
									/>
								</TableCell>
								<TableCell>
									<Button
										variant="ghost"
										size="icon"
										@click="removeAdditive(index)"
									>
										<Trash class="w-5 h-5 text-red-500 hover:text-red-700" />
									</Button>
								</TableCell>
							</TableRow>
						</template>
						<template v-else>
							<TableRow>
								<TableCell
									colspan="4"
									class="text-center text-gray-500"
								>
									Нет добавленных добавок
								</TableCell>
							</TableRow>
						</template>
					</TableBody>
				</Table>
			</CardContent>
		</Card>

		<!-- Footer (mobile only) -->
		<div class="flex justify-center items-center gap-2 md:hidden mt-4">
			<Button
				variant="outline"
				type="button"
				@click="onCancel"
			>
				Отменить
			</Button>
			<Button
				type="submit"
				@click="onSubmit"
			>
				Сохранить
			</Button>
		</div>

		<!-- Dialog for selecting additives -->
		<AdminSelectAvailableToAddAdditiveDialog
			:open="openAdditiveDialog"
			@close="openAdditiveDialog = false"
			@select="selectAdditive"
		/>
	</div>
</template>
