export const DEFAULT_PAGINATION_META = {
	page: 1,
	pageSize: 10,
}

export interface PaginationMeta {
	page: number
	pageSize: number
	totalCount: number
	totalPages: number
}

export interface PaginatedResponse<T> {
	data: T
	pagination: PaginationMeta
}

export interface PaginationParams {
	page?: number
	pageSize?: number
}

export function hasMorePages(meta: PaginationMeta): boolean {
	return meta.page < meta.totalPages
}
