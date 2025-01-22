import {
	adminNavItems,
	isCollapsibleNavItem,
	type NavItem,
	type SidebarNavItem,
} from './app-admin-sidebar-links.constants'

function flattenSidebarNavItems(
	sidebarItems: SidebarNavItem[],
	parentLabel: string = '',
): NavItem[] {
	return sidebarItems.flatMap(item => {
		if (isCollapsibleNavItem(item)) {
			const updatedLabel = parentLabel ? `${parentLabel} ${item.label}` : item.label

			return flattenSidebarNavItems(item.items, updatedLabel)
		}

		return {
			...item,
			name: parentLabel ? `${parentLabel} ${item.name.toLowerCase()}` : item.name,
		}
	})
}

export const adminSearchItems: NavItem[] = flattenSidebarNavItems(adminNavItems)
