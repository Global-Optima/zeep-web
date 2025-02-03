import type { App, Directive } from 'vue'
import { roleDirective } from './roles.directive'

const directives: { [key: string]: Directive } = {
	appRole: roleDirective,
}

export default {
	install(app: App) {
		Object.keys(directives).forEach(key => {
			app.directive(key, directives[key])
		})
	},
}
