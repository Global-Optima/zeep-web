import './styles.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import { router } from './router'
import { i18nConfig } from './core/config/locale.config'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18nConfig)

app.mount('#app')
