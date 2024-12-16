import './styles.scss'

import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import { i18nConfig } from './core/config/locale.config'
import { router } from './router'

import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import { toastConfig } from '@/core/config/toast.config'
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import Vue3Toastify from 'vue3-toastify'

const app = createApp(App)
const pinia = createPinia()

pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.use(i18nConfig)
app.use(Vue3Toastify, toastConfig)
app.use(VueQueryPlugin, {
	queryClient: new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 0,
			},
		},
	}),
})

app.mount('#app')
