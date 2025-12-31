import { createRouter, createWebHistory } from "vue-router";
import Layout from "../components/Layout.vue";
import Appointments from "../views/Appointments.vue";
import Dashboard from "../views/Dashboard.vue";
import History from "../views/History.vue";
import Members from "../views/Members.vue";
import Services from "../views/Services.vue";
import Technicians from "../views/Technicians.vue";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/",
			component: Layout,
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
			],
		},
	],
});

export default router;
