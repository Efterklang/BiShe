<template>
	<div>
		<div class="flex justify-between items-center mb-6">
			<div>
				<h1 class="text-2xl font-bold tracking-tight">用户管理</h1>
				<p class="text-base-content/60 mt-1">管理系统用户、角色分配与账号状态</p>
			</div>
			<button class="btn btn-primary" @click="openCreateModal">
				<Plus class="w-4 h-4 mr-1" />
				创建新用户
			</button>
		</div>

		<!-- Users Table -->
		<div class="card bg-base-100 border border-base-300 shadow-sm rounded-xl overflow-hidden">
			<div class="card-body p-0">
				<div v-if="isLoading" class="flex justify-center py-12">
					<span class="loading loading-spinner loading-lg text-primary"></span>
				</div>

				<div v-else-if="users.length === 0" class="flex flex-col items-center justify-center py-16 text-center">
					<div class="w-16 h-16 bg-base-200 rounded-full flex items-center justify-center mb-4">
						<Users class="w-8 h-8 text-base-content/40" />
					</div>
					<h3 class="text-lg font-bold text-base-content">暂无用户数据</h3>
					<p class="text-base-content/60 mt-1 max-w-sm">
						当前系统中还没有任何用户。点击右上角的"创建新用户"按钮来添加。
					</p>
				</div>

				<div v-else class="overflow-x-auto">
					<table class="table w-full">
						<thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
							<tr>
								<th class="font-medium">ID</th>
								<th class="font-medium">用户名</th>
								<th class="font-medium">角色</th>
								<th class="font-medium">状态</th>
								<th class="font-medium">创建时间</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-base-200">
							<tr v-for="user in users" :key="user.id" class="hover:bg-base-50/50 transition-colors">
								<td class="font-mono text-xs opacity-60">#{{ user.id }}</td>
								<td>
									<div class="flex items-center gap-3">
										<div class="avatar placeholder">
											<div class="bg-neutral text-neutral-content rounded-full w-8 h-8">
												<span class="text-xs">{{ user.username.charAt(0).toUpperCase() }}</span>
											</div>
										</div>
										<span class="font-medium">{{ user.username }}</span>
									</div>
								</td>
								<td>
									<div class="badge gap-2"
										:class="user.role === 'manager' ? 'badge-primary badge-outline' : 'badge-ghost'">
										<ShieldCheck v-if="user.role === 'manager'" class="w-3 h-3" />
										<User v-else class="w-3 h-3" />
										{{ user.role === "manager" ? "店长" : "操作员" }}
									</div>
								</td>
								<td>
									<div class="badge gap-1.5"
										:class="user.is_active ? 'badge-success/10 text-success' : 'badge-error/10 text-error'">
										<span class="w-1.5 h-1.5 rounded-full"
											:class="user.is_active ? 'bg-success' : 'bg-error'"></span>
										{{ user.is_active ? "正常" : "已停用" }}
									</div>
								</td>
								<td class="text-base-content/60 text-sm font-mono">
									{{ formatDate(user.created_at) }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			<!-- Total Count Footer -->
			<div class="bg-base-50 px-6 py-3 border-t border-base-200 text-xs text-base-content/60 flex justify-between items-center"
				v-if="users.length > 0">
				<span>共 {{ users.length }} 个用户</span>
			</div>
		</div>

		<!-- Create User Modal -->
		<dialog ref="createModalRef" class="modal">
			<div class="modal-box">
				<h3 class="font-bold text-lg mb-4 flex items-center gap-2">
					<UserPlus class="w-5 h-5 text-primary" />
					创建新用户
				</h3>

				<form @submit.prevent="handleCreateUser">
					<!-- Username -->
					<div class="form-control">
						<label class="label">
							<span class="label-text">用户名</span>
						</label>
						<div class="relative">
							<input v-model="newUser.username" type="text" placeholder="请输入用户名（3-64字符）"
								class="input input-bordered w-full pl-10" :class="{ 'input-error': formErrors.username }"
								required />
							<User class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
						</div>
						<label v-if="formErrors.username" class="label">
							<span class="label-text-alt text-error flex items-center gap-1">
								<AlertCircle class="w-3 h-3" />
								{{ formErrors.username }}
							</span>
						</label>
					</div>

					<!-- Password -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text">密码</span>
						</label>
						<div class="relative">
							<input v-model="newUser.password" type="password" placeholder="至少8位，包含大小写字母和数字"
								class="input input-bordered w-full pl-10" :class="{ 'input-error': formErrors.password }"
								required />
							<Lock class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
						</div>
						<label v-if="formErrors.password" class="label">
							<span class="label-text-alt text-error flex items-center gap-1">
								<AlertCircle class="w-3 h-3" />
								{{ formErrors.password }}
							</span>
						</label>
					</div>

					<!-- Confirm Password -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text">确认密码</span>
						</label>
						<div class="relative">
							<input v-model="newUser.confirmPassword" type="password" placeholder="请再次输入密码"
								class="input input-bordered w-full pl-10"
								:class="{ 'input-error': formErrors.confirmPassword }" required />
							<Lock class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
						</div>
						<label v-if="formErrors.confirmPassword" class="label">
							<span class="label-text-alt text-error flex items-center gap-1">
								<AlertCircle class="w-3 h-3" />
								{{ formErrors.confirmPassword }}
							</span>
						</label>
					</div>

					<!-- Role -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text">角色</span>
						</label>
						<div class="relative">
							<select v-model="newUser.role" class="select select-bordered w-full pl-10" required>
								<option value="operator">操作员</option>
								<option value="manager">店长</option>
							</select>
							<Shield class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
						</div>
					</div>

					<!-- Error Alert -->
					<div v-if="createError" class="alert alert-error mt-6 shadow-sm">
						<AlertCircle class="stroke-current shrink-0 h-5 w-5" />
						<span class="text-sm">{{ createError }}</span>
					</div>

					<!-- Modal Actions -->
					<div class="modal-action">
						<button type="button" class="btn btn-ghost" @click="closeCreateModal" :disabled="isCreating">
							取消
						</button>
						<button type="submit" class="btn btn-primary" :class="{ loading: isCreating }"
							:disabled="isCreating">
							<span v-if="!isCreating">创建</span>
							<span v-else>创建中...</span>
						</button>
					</div>
				</form>
			</div>
			<form method="dialog" class="modal-backdrop">
				<button @click="closeCreateModal">close</button>
			</form>
		</dialog>
	</div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { getUserList, register } from "../api/auth";
import {
	Plus,
	Users,
	Shield,
	ShieldCheck,
	User,
	UserPlus,
	Lock,
	AlertCircle
} from 'lucide-vue-next';

const users = ref([]);
const isLoading = ref(false);
const isCreating = ref(false);
const createError = ref("");

const createModalRef = ref(null);

const newUser = reactive({
	username: "",
	password: "",
	confirmPassword: "",
	role: "operator",
});

const formErrors = reactive({
	username: "",
	password: "",
	confirmPassword: "",
});

// Fetch users list
const fetchUsers = async () => {
	isLoading.value = true;
	try {
		const data = await getUserList();
		users.value = data.users || [];
	} catch (error) {
		console.error("Failed to fetch users:", error);
	} finally {
		isLoading.value = false;
	}
};

// Validate password strength
const validatePassword = (password) => {
	if (password.length < 8) {
		return "密码至少8个字符";
	}

	const hasUpper = /[A-Z]/.test(password);
	const hasLower = /[a-z]/.test(password);
	const hasNumber = /[0-9]/.test(password);

	if (!hasUpper || !hasLower || !hasNumber) {
		return "密码必须包含大小写字母和数字";
	}

	return "";
};

// Validate form
const validateForm = () => {
	formErrors.username = "";
	formErrors.password = "";
	formErrors.confirmPassword = "";

	let isValid = true;

	if (newUser.username.length < 3 || newUser.username.length > 64) {
		formErrors.username = "用户名长度必须在3-64字符之间";
		isValid = false;
	}

	const passwordError = validatePassword(newUser.password);
	if (passwordError) {
		formErrors.password = passwordError;
		isValid = false;
	}

	if (newUser.password !== newUser.confirmPassword) {
		formErrors.confirmPassword = "两次输入的密码不一致";
		isValid = false;
	}

	return isValid;
};

// Open create modal
const openCreateModal = () => {
	resetForm();
	createModalRef.value?.showModal();
};

// Close create modal
const closeCreateModal = () => {
	createModalRef.value?.close();
	resetForm();
};

// Reset form
const resetForm = () => {
	newUser.username = "";
	newUser.password = "";
	newUser.confirmPassword = "";
	newUser.role = "operator";
	formErrors.username = "";
	formErrors.password = "";
	formErrors.confirmPassword = "";
	createError.value = "";
};

// Handle create user
const handleCreateUser = async () => {
	if (!validateForm()) {
		return;
	}

	isCreating.value = true;
	createError.value = "";

	try {
		await register(newUser.username, newUser.password, newUser.role);

		// Success - refresh list and close modal
		await fetchUsers();
		closeCreateModal();

		// Show success message (you can add a toast notification here)
		alert("用户创建成功！");
	} catch (error) {
		console.error("Failed to create user:", error);

		if (error.response?.data?.msg) {
			createError.value = error.response.data.msg;
		} else if (error.message) {
			createError.value = error.message;
		} else {
			createError.value = "创建用户失败，请重试";
		}
	} finally {
		isCreating.value = false;
	}
};

// Format date
const formatDate = (dateString) => {
	if (!dateString) return "-";
	const date = new Date(dateString);
	return date.toLocaleString("zh-CN", {
		year: "numeric",
		month: "2-digit",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit",
	});
};

// Load data on mount
onMounted(() => {
	fetchUsers();
});
</script>

<style scoped>
/* Add any custom styles here */
</style>