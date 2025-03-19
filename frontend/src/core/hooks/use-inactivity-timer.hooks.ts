// file: src/modules/kiosk/composables/use-inactivity-timer.ts

import type { Ref } from 'vue'
import { onBeforeUnmount, onMounted } from 'vue'

export function useInactivityTimer(
	inactivityDurationMs: number,
	onTimeout: () => void,
	isEnabled: Ref<boolean> | boolean = true,
) {
	let inactivityTimeout: ReturnType<typeof setTimeout> | null = null

	const activityEvents = ['mousemove', 'keydown', 'click', 'touchstart']

	function resetInactivityTimer() {
		if (!isEnabled) return
		if (inactivityTimeout) {
			clearTimeout(inactivityTimeout)
		}
		inactivityTimeout = setTimeout(() => {
			onTimeout()
		}, inactivityDurationMs)
	}

	onMounted(() => {
		// Only enable if the kiosk is truly active
		if (!isEnabled) return
		resetInactivityTimer()
		activityEvents.forEach(event => {
			window.addEventListener(event, resetInactivityTimer, { passive: true })
		})
	})

	onBeforeUnmount(() => {
		if (inactivityTimeout) {
			clearTimeout(inactivityTimeout)
		}
		activityEvents.forEach(event => {
			window.removeEventListener(event, resetInactivityTimer)
		})
	})

	return {
		resetInactivityTimer,
	}
}
