<template>
	<Popover>
		<PopoverTrigger as-child>
			<Button
				variant="outline"
				size="sm"
				class="h-9"
			>
				<PlusCircle class="mr-2 w-4 h-4" />
				{{ title }}
				<template v-if="selectedValues.size > 0">
					<Separator
						orientation="vertical"
						class="mx-2 h-4"
					/>
					<Badge
						variant="secondary"
						class="lg:hidden px-1 rounded-sm font-normal"
					>
						{{ selectedValues.size }}
					</Badge>
					<div class="hidden lg:flex space-x-1">
						<Badge
							v-if="selectedValues.size > 2"
							variant="secondary"
							class="px-1 rounded-sm font-normal"
						>
							{{ selectedValues.size }} выбрано
						</Badge>
						<template v-else>
							<Badge
								v-for="option in options.filter((opt) => selectedValues.has(opt.value))"
								:key="option.value"
								variant="secondary"
								class="px-1 rounded-sm font-normal"
							>
								{{ option.label }}
							</Badge>
						</template>
					</div>
				</template>
			</Button>
		</PopoverTrigger>

		<PopoverContent
			class="p-0 w-fit"
			align="start"
		>
			<Command :filter-function="filterFunction">
				<CommandInput
					:placeholder="title"
					v-model="searchTerm"
				/>
				<CommandList>
					<CommandEmpty>Нет результатов</CommandEmpty>
					<CommandGroup>
						<CommandItem
							v-for="option in filteredOptions"
							:key="option.value"
							:value="option.value"
							@select="toggleOption(option.value)"
						>
							<div
								:class="cn(
                  'mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary',
                  selectedValues.has(option.value)
                    ? 'bg-primary text-primary-foreground'
                    : 'opacity-50 [&_svg]:invisible'
                )"
							>
								<Check class="w-4 h-4" />
							</div>
							<component
								v-if="option.icon"
								:is="option.icon"
								class="mr-2 w-4 h-4 text-muted-foreground"
							/>
							<span>{{ option.label }}</span>
						</CommandItem>
					</CommandGroup>

					<template v-if="selectedValues.size > 0">
						<CommandSeparator />
						<CommandGroup>
							<CommandItem
								value="clear-filters"
								class="justify-center text-center"
								@select="clearSelection"
							>
								Сбросить фильтры
							</CommandItem>
						</CommandGroup>
					</template>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>

<script setup lang="ts">
import { Badge } from '@/core/components/ui/badge'
import { Button } from '@/core/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '@/core/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/core/components/ui/popover'
import { Separator } from '@/core/components/ui/separator'
import { cn } from '@/core/utils/tailwind.utils'
import { Check, PlusCircle, type LucideIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'

interface MultiSelectOption {
  label: string
  value: string
  icon?: LucideIcon
}

const props = defineProps<{
  title: string
  options: MultiSelectOption[]
  modelValue: string[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void
}>()

const searchTerm = ref('')

const selectedValues = computed(() => new Set(props.modelValue))

const filteredOptions = computed(() => {
  const term = searchTerm.value.toLowerCase()
  if (!term) return props.options
  return props.options.filter((opt) => opt.label.toLowerCase().includes(term))
})

function toggleOption(value: string) {
  const updated = new Set(selectedValues.value)
  if (updated.has(value)) {
    updated.delete(value)
  } else {
    updated.add(value)
  }
  emit('update:modelValue', Array.from(updated))
}

function clearSelection() {
  emit('update:modelValue', [])
  searchTerm.value = ''
}

function filterFunction(
  _: unknown[],
  term: string
) {
  const filtered = props.options.filter((opt) =>
    opt.label.toLowerCase().includes(term.toLowerCase())
  )

  return filtered.map(item => item.value)
}
</script>

<style scoped></style>
