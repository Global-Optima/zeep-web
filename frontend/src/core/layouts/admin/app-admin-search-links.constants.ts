import {
	adminNavItems,
	isCollapsibleNavItem,
	type NavItem,
	type SidebarNavItem,
} from './app-admin-sidebar-links.constants'

function flattenSidebarNavItems(sidebarItems: SidebarNavItem[]): NavItem[] {
	return sidebarItems.flatMap(item => {
		if (isCollapsibleNavItem(item)) {
			return flattenSidebarNavItems(item.items)
		}
    
		return item
	})
}

export const adminSearchItems: NavItem[] = flattenSidebarNavItems(adminNavItems)
