<template>
  <div class="gap-4 grid grid-cols-2 sm:grid-cols-3">
    <Skeleton
      v-for="n in skeletonCount"
      :key="n"
      class="bg-slate-200 bg-opacity-80 rounded-[38px] w-full h-[404px]"
    />
  </div>
</template>

<script setup lang="ts">
import { Skeleton } from '@/core/components/ui/skeleton'

const props = defineProps<{ skeletonCount?: number }>()
const skeletonCount = props.skeletonCount || 6
</script>
