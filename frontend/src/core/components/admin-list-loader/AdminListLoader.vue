<script lang="ts" setup>
import { Skeleton } from '@/core/components/ui/skeleton'
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between items-center gap-2">
			<Skeleton class="rounded-md w-[256px] h-9" />
			<Skeleton class="rounded-md w-[200px] h-9" />
		</div>

		<Skeleton class="rounded-xl w-full h-[300px]" />
	</div>
</template>
