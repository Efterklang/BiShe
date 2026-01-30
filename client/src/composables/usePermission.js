import { computed } from "vue";
import { useAppStore } from "../stores/app";

/**
 * Composable for checking user permissions based on role
 */
export function usePermission() {
	const appStore = useAppStore();

	/**
	 * Check if user has a specific role
	 * @param {string} role - The role to check ("manager" or "operator")
	 * @returns {boolean}
	 */
	const hasRole = (role) => {
		return appStore.user?.role === role;
	};

	/**
	 * Check if user has any of the specified roles
	 * @param {string[]} roles - Array of roles to check
	 * @returns {boolean}
	 */
	const hasAnyRole = (roles) => {
		return roles.includes(appStore.user?.role);
	};

	/**
	 * Check if user can delete resources (manager only)
	 */
	const canDelete = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can edit resources (manager only)
	 */
	const canEdit = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can create resources (manager only for restricted resources)
	 */
	const canCreate = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can view AI reports (manager only)
	 */
	const canViewAI = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can manage users (manager only)
	 */
	const canManageUsers = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can manage services (manager only)
	 */
	const canManageServices = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can manage technicians (manager only)
	 */
	const canManageTechnicians = computed(() => {
		return appStore.isManager;
	});

	/**
	 * Check if user can manage products (manager only)
	 */
	const canManageProducts = computed(() => {
		return appStore.isManager;
	});

	return {
		hasRole,
		hasAnyRole,
		canDelete,
		canEdit,
		canCreate,
		canViewAI,
		canManageUsers,
		canManageServices,
		canManageTechnicians,
		canManageProducts,
	};
}
