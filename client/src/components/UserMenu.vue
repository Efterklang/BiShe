<template>
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar placeholder group">
			<div
				class="rounded-full w-10 h-10 bg-base-100 border border-base-300 flex items-center justify-center ring-2 ring-transparent group-hover:ring-primary/50 transition-all duration-300 shadow-sm group-hover:shadow-md">
				<User class="w-5 h-5 text-base-content/70" />
			</div>
		</div>
		<ul tabindex="0"
			class="dropdown-content z-50 menu p-2 shadow-2xl bg-base-100/95 backdrop-blur-sm rounded-2xl w-64 border border-base-200/50 mt-2 origin-top-right ring-1 ring-base-300/20">
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
						<Users class="h-4 w-4" />
					</div>
					<span class="font-medium">用户管理</span>
				</router-link>
			</li>

			<!-- Logout -->
			<li class="mt-1">
				<a @click="handleLogout"
					class="gap-3 py-3 rounded-xl text-error hover:bg-error/5 hover:text-error active:bg-error/10">
					<div class="p-1.5 bg-error/10 rounded-lg transition-colors">
						<LogOut class="h-4 w-4" />
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
import { Users, LogOut, User } from 'lucide-vue-next';
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