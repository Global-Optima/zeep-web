import './styles.scss'

import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import { i18nConfig } from './core/config/locale.config'
import { router } from './router'

import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'

const app = createApp(App)
const pinia = createPinia().use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.use(i18nConfig)
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
