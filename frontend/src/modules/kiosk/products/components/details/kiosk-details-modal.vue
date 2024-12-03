<script lang="ts" setup>
import KioskDetailsModalContent from '@/modules/kiosk/products/components/details/kiosk-details-modal-content.vue'
import { useCurrentProductStore } from '@/modules/kiosk/products/components/hooks/use-current-product.hook'
import { SwipeModal } from "@takuma-ru/vue-swipe-modal"
import { storeToRefs } from 'pinia'

const currentProductStore = useCurrentProductStore();
const { isModalOpen, currentProductId } = storeToRefs(currentProductStore);

const handleModalClose = () => {
  currentProductStore.closeModal();
};
</script>

<template>
	<div class="flex justify-center">
		<SwipeModal
			:model-value="isModalOpen"
			@update:model-value="currentProductStore.closeModal"
			snapPoint="95dvh"
		>
			<div class="rounded-md w-full h-full no-scrollbar">
				<KioskDetailsModalContent
					v-if="currentProductId"
					:productId="currentProductId"
					@close="handleModalClose"
				/>
			</div>
		</SwipeModal>
	</div>
</template>

<style lang="scss" scoped>
.modal-style {
  max-width: calc(100% - 112px) !important;
  margin: 0 auto !important;
  background-color: #F5F5F7 !important;
  border-radius: 3rem 3rem 0 0 !important;
  overflow: clip !important;;
}

:deep(.swipe-modal-content>.panel) {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

:deep(.swipe-modal-content>.panel)::-webkit-scrollbar {
		display: none;
}

:deep(.swipe-modal-drag-handle-wrapper) {
  background-color: #F5F5F7;
  border: none;
}
</style>
