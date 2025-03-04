import { watch, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/core/components/ui/button'
import { Input } from '@/core/components/ui/input'
import { getRouteName } from '@/core/config/routes.config'
import { useDebouncedSearch } from '@/core/hooks/use-debounced-search.hook'
import type { AdditiveFilterQuery } from '@/modules/admin/additives/models/additives.model'

const router = useRouter()

const props = defineProps<{ filter: AdditiveFilterQuery }>()
const emit = defineEmits<{
  (event: 'update:filter', value: AdditiveFilterQuery): void
}>()

// Ensure searchTerm is reactive and initialized from filter
const searchTerm = ref(props.filter.search || '')

// Ensure searchTerm updates when props.filter.search changes
watch(
  () => props.filter.search,
  (newSearch) => {
    if (newSearch !== searchTerm.value) {
      searchTerm.value = newSearch || ''
    }
  },
  { immediate: true } // Ensure it runs on mount
)

// Use debounced search hook
const { debouncedSearchTerm } = useDebouncedSearch(searchTerm, 500)

// Watch debounced search and update the filter
watch(debouncedSearchTerm, (newValue) => {
  emit('update:filter', { ...props.filter, search: newValue.trim() })
})

// Navigate to add store
const addStore = () => {
  router.push({ name: getRouteName('ADMIN_ADDITIVE_CREATE') })
}
