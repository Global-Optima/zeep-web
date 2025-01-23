// projects.types.ts

/**
 * Enum for project statuses.
 */
export enum ProjectStatus {
	PLANNED = 'Запланирован',
	IN_PROGRESS = 'В процессе',
	COMPLETED = 'Завершен',
}

/**
 * Type for a single project.
 */
export interface Project {
	id: number
	name: string
	startDate: string // ISO format (e.g., '2024-01-01')
	endDate: string // ISO format
	status: ProjectStatus
	responsible: string // Full name of the responsible person
}

/**
 * Interface for project filters.
 */
export interface ProjectFilter {
	search?: string
	status?: ProjectStatus | null
}
