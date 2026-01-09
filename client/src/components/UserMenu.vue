<template>
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar placeholder">
			<div class="bg-neutral text-neutral-content rounded-full w-10">
				<span class="text-xl">{{ userInitial }}</span>
			</div>
		</div>
		<ul
			tabindex="0"
			class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-56"
		>
			<!-- User Info -->
			<li class="menu-title">
				<div class="flex flex-col gap-1">
					<span class="font-semibold">{{ username }}</span>
					<span
						class="badge badge-sm"
						:class="isManager ? 'badge-primary' : 'badge-secondary'"
					>
						{{ roleText }}
					</span>
				</div>
			</li>

			<li class="divider"></li>

			<!-- User Management (Manager Only) -->
			<li v-if="isManager">
				<router-link to="/users">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
						/>
					</svg>
					用户管理
				</router-link>
			</li>

			<!-- Logout -->
			<li>
				<a @click="handleLogout" class="text-error">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
						/>
					</svg>
					退出登录
				</a>
			</li>
		</ul>
	</div>
</template>

<script setup>
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useAppStore } from "../stores/app";

const router = useRouter();
const appStore = useAppStore();

// Computed properties
const username = computed(() => appStore.user?.username || "User");
const isManager = computed(() => appStore.isManager);
const roleText = computed(() => (isManager.value ? "店长" : "操作员"));

// Get user initial for avatar
const userInitial = computed(() => {
	const name = username.value;
	return name ? name.charAt(0).toUpperCase() : "U";
});

// Handle logout
const handleLogout = () => {
	if (confirm("确定要退出登录吗？")) {
		appStore.logout();
		router.push("/login");
	}
};
</script>

<style scoped>
.menu-title {
	padding: 0.75rem 1rem;
}

.divider {
	margin: 0.5rem 0;
	height: 1px;
	background-color: hsl(var(--bc) / 0.1);
}
</style>
