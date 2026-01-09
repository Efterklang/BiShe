import { createRouter, createWebHistory } from "vue-router";
import { useAppStore } from "../stores/app";
import Layout from "../components/Layout.vue";
import Appointments from "../views/Appointments.vue";
import Dashboard from "../views/Dashboard.vue";
import History from "../views/History.vue";
import Members from "../views/Members.vue";
import Services from "../views/Services.vue";
import Technicians from "../views/Technicians.vue";
import Login from "../views/Login.vue";
import UserManagement from "../views/UserManagement.vue";
import Products from "../views/Products.vue";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/login",
			name: "login",
			component: Login,
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
					component: Dashboard,
				},
				{
					path: "appointments",
					name: "appointments",
					component: Appointments,
				},
				{
					path: "technicians",
					name: "technicians",
					component: Technicians,
				},
				{
					path: "services",
					name: "services",
					component: Services,
				},
				{
					path: "members",
					name: "members",
					component: Members,
				},
				{
					path: "history",
					name: "history",
					component: History,
				},
				{
					path: "products",
					name: "products",
					component: Products,
					meta: { requiresAuth: true },
				},
				{
					path: "users",
					name: "users",
					component: UserManagement,
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
