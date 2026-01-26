import { createRouter, createWebHistory } from "vue-router";
import { useAppStore } from "../stores/app";
import Layout from "../components/Layout.vue";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/login",
			name: "login",
			component: () => import("../views/Login.vue"),
			meta: { requiresAuth: false },
		},
		{
			path: "/",
			component: Layout,
			meta: { requiresAuth: true },
			children: [
				{
					path: "",
					name: "dashboard",
					component: () => import("../views/Dashboard.vue"),
				},
				{
					path: "appointments",
					name: "appointments",
					component: () => import("../views/Appointments.vue"),
				},
				{
					path: "technician-schedule",
					name: "technician-schedule",
					component: () => import("../views/TechnicianScheduleView.vue"),
				},
				{
					path: "technicians",
					name: "technicians",
					component: () => import("../views/Technicians.vue"),
				},
				{
					path: "services",
					name: "services",
					component: () => import("../views/Services.vue"),
				},
				{
					path: "members",
					name: "members",
					component: () => import("../views/Members.vue"),
				},
				{
					path: "history",
					name: "history",
					component: () => import("../views/History.vue"),
				},
				{
					path: "products",
					name: "products",
					component: () => import("../views/Products.vue"),
					meta: { requiresAuth: true },
				},
				{
					path: "users",
					name: "users",
					component: () => import("../views/UserManagement.vue"),
					meta: { requiresAuth: true, roles: ["manager"] },
				},
			],
		},
	],
});

// Navigation guard for authentication and authorization
router.beforeEach((to, from, next) => {
	const appStore = useAppStore();

	// Initialize auth from localStorage
	if (!appStore.user && localStorage.getItem("token")) {
		appStore.initAuth();
	}

	const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
	const isAuthenticated = appStore.isAuthenticated;

	// Check if route requires authentication
	if (requiresAuth && !isAuthenticated) {
		// Redirect to login page
		next({ name: "login", query: { redirect: to.fullPath } });
		return;
	}

	// Redirect authenticated users away from login page
	if (to.name === "login" && isAuthenticated) {
		next({ name: "dashboard" });
		return;
	}

	// Check role-based access
	const requiredRoles = to.meta.roles;
	if (requiredRoles && requiredRoles.length > 0) {
		const userRole = appStore.user?.role;
		if (!requiredRoles.includes(userRole)) {
			// User doesn't have required role, redirect to dashboard
			alert("您没有权限访问此页面");
			next({ name: "dashboard" });
			return;
		}
	}

	next();
});

export default router;
