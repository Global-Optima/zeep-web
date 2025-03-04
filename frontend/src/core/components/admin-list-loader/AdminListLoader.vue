<script lang="ts" setup>
import { Skeleton } from '@/core/components/ui/skeleton'
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between items-center gap-2">
			<Skeleton class="bg-slate-200 rounded-md w-[256px] h-10" />
			<Skeleton class="bg-slate-200 rounded-md w-[200px] h-10" />
		</div>

		<Skeleton class="bg-slate-200 rounded-xl w-full h-[400px]" />
	</div>
</template>
