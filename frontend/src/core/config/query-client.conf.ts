import {
	QueryClient,
	type VueQueryPluginOptions,
	defaultShouldDehydrateQuery,
} from '@tanstack/vue-query'

const myClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 60 * 1000,
		},
		dehydrate: {
			shouldDehydrateQuery: query =>
				defaultShouldDehydrateQuery(query) || query.state.status === 'pending',
},
	},
})

export const vueQueryPluginOptions: VueQueryPluginOptions = {
	queryClient: myClient,
}
