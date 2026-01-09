import { defineStore } from "pinia";
import { ref, computed } from "vue";
import * as authAPI from "../api/auth";

export const useAppStore = defineStore("app", () => {
	const isLoading = ref(false);
	const theme = ref("emerald");
	const token = ref(localStorage.getItem("token") || "");
	const user = ref(null);

	// Computed properties
	const isAuthenticated = computed(() => !!token.value && !!user.value);
	const isManager = computed(() => user.value?.role === "manager");
	const isOperator = computed(() => user.value?.role === "operator");

	function setLoading(status) {
		isLoading.value = status;
	}

	function setTheme(newTheme) {
		theme.value = newTheme;
		document.documentElement.setAttribute("data-theme", newTheme);
	}

	// Initialize auth state from localStorage
	function initAuth() {
		const storedUser = localStorage.getItem("user");
		if (storedUser && token.value) {
			try {
				user.value = JSON.parse(storedUser);
			} catch (e) {
				console.error("Failed to parse stored user:", e);
				logout();
			}
		}
	}

	// Login function
	async function login(credentials) {
		try {
			const response = await authAPI.login(
				credentials.username,
				credentials.password,
			);

			// Store token and user info
			token.value = response.token;
			user.value = response.user;

			localStorage.setItem("token", response.token);
			localStorage.setItem("user", JSON.stringify(response.user));

			return response;
		} catch (error) {
			console.error("Login failed:", error);
			throw error;
		}
	}

	// Logout function
	function logout() {
		token.value = "";
		user.value = null;
		localStorage.removeItem("token");
		localStorage.removeItem("user");
	}

	// Fetch current user info
	async function fetchCurrentUser() {
		try {
			const userData = await authAPI.getCurrentUser();
			user.value = userData;
			localStorage.setItem("user", JSON.stringify(userData));
			return userData;
		} catch (error) {
			console.error("Failed to fetch current user:", error);
			logout();
			throw error;
		}
	}

	return {
		isLoading,
		theme,
		token,
		user,
		isAuthenticated,
		isManager,
		isOperator,
		setLoading,
		setTheme,
		initAuth,
		login,
		logout,
		fetchCurrentUser,
	};
});
