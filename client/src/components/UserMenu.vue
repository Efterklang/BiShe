<template>
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar placeholder group">
			<div
				class="bg-neutral text-neutral-content rounded-full w-10 ring-2 ring-transparent group-hover:ring-primary/50 transition-all duration-300">
				<span class="text-xl font-medium">{{ userInitial }}</span>
			</div>
		</div>
		<ul tabindex="0"
			class="dropdown-content z-[50] menu p-2 shadow-2xl bg-base-100 rounded-2xl w-64 border border-base-200/50 mt-2 origin-top-right">
			<!-- User Info Header -->
			<li class="menu-title px-4 py-3">
				<div class="flex flex-col gap-2">
					<div class="flex items-center justify-between">
						<span class="font-bold text-lg text-base-content tracking-tight">{{ username }}</span>
						<span class="badge badge-sm border-0"
							:class="isManager ? 'bg-primary/10 text-primary' : 'bg-secondary/10 text-secondary'">
							{{ roleText }}
						</span>
					</div>
					<div class="h-px w-full bg-base-200"></div>
				</div>
			</li>

			<!-- User Management (Manager Only) -->
			<li v-if="isManager">
				<router-link to="/users" class="gap-3 py-3 rounded-xl hover:bg-base-200/50 active:bg-base-200">
					<div class="p-1.5 bg-base-200 rounded-lg group-hover:bg-base-300 transition-colors">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24"
							stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
						</svg>
					</div>
					<span class="font-medium">用户管理</span>
				</router-link>
			</li>

			<!-- Logout -->
			<li class="mt-1">
				<a @click="handleLogout"
					class="gap-3 py-3 rounded-xl text-error hover:bg-error/5 hover:text-error active:bg-error/10">
					<div class="p-1.5 bg-error/10 rounded-lg transition-colors">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24"
							stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
						</svg>
					</div>
					<span class="font-medium">退出登录</span>
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