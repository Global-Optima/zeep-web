import pluginVitest from '@vitest/eslint-plugin'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import pluginVue from 'eslint-plugin-vue'

export default defineConfigWithVueTs(
	pluginVue.configs['flat/essential'],
	vueTsConfigs.recommended,
	{
		name: 'app/files-to-ignore',
		ignores: ['**/dist/**', '**/dev-dist/**', '**/dist-ssr/**', '**/coverage/**'],
	},
	{
		...pluginVitest.configs.recommended,
		files: ['src/**/__tests__/*'],
	},
	{
		rules: {
			'vue/multi-word-component-names': 'off',
			'@typescript-eslint/no-unused-vars': [
				'warn',
				{
					varsIgnorePattern: '^_',
					argsIgnorePattern: '^_',
				},
			],
		},
	},
)
