<template>
	<div class="container mx-auto p-6">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-3xl font-bold">用户管理</h1>
			<button class="btn btn-primary" @click="openCreateModal">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24"
					stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
				</svg>
				创建新用户
			</button>
		</div>

		<!-- Users Table -->
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<div v-if="isLoading" class="flex justify-center py-12">
					<span class="loading loading-spinner loading-lg"></span>
				</div>

				<div v-else-if="users.length === 0" class="text-center py-12">
					<p class="text-base-content/60">暂无用户数据</p>
				</div>

				<div v-else class="overflow-x-auto">
					<table class="table table-zebra w-full">
						<thead>
							<tr>
								<th>ID</th>
								<th>用户名</th>
								<th>角色</th>
								<th>状态</th>
								<th>创建时间</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="user in users" :key="user.id">
								<td>{{ user.id }}</td>
								<td class="font-medium">{{ user.username }}</td>
								<td>
									<div class="badge" :class="user.role === 'manager' ? 'badge-primary' : 'badge-secondary'
										">
										{{ user.role === "manager" ? "店长" : "操作员" }}
									</div>
								</td>
								<td>
									<div class="badge" :class="user.is_active ? 'badge-success' : 'badge-error'">
										{{ user.is_active ? "正常" : "已停用" }}
									</div>
								</td>
								<td>{{ formatDate(user.created_at) }}</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- Total Count -->
				<div class="text-sm text-base-content/60 mt-4">
					共 {{ users.length }} 个用户
				</div>
			</div>
		</div>

		<!-- Create User Modal -->
		<dialog ref="createModalRef" class="modal">
			<div class="modal-box">
				<h3 class="font-bold text-lg mb-4">创建新用户</h3>

				<form @submit.prevent="handleCreateUser">
					<!-- Username -->
					<div class="form-control">
						<label class="label">
							<span class="label-text">用户名</span>
						</label>
						<input v-model="newUser.username" type="text" placeholder="请输入用户名（3-64字符）"
							class="input input-bordered" :class="{ 'input-error': formErrors.username }" required />
						<label v-if="formErrors.username" class="label">
							<span class="label-text-alt text-error">{{
								formErrors.username
							}}</span>
						</label>
					</div>

					<!-- Password -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text">密码</span>
						</label>
						<input v-model="newUser.password" type="password" placeholder="至少8位，包含大小写字母和数字"
							class="input input-bordered" :class="{ 'input-error': formErrors.password }" required />
						<label v-if="formErrors.password" class="label">
							<span class="label-text-alt text-error">{{
								formErrors.password
							}}</span>
						</label>
					</div>

					<!-- Confirm Password -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text">确认密码</span>
						</label>
						<input v-model="newUser.confirmPassword" type="password" placeholder="请再次输入密码"
							class="input input-bordered" :class="{ 'input-error': formErrors.confirmPassword }"
							required />
						<label v-if="formErrors.confirmPassword" class="label">
							<span class="label-text-alt text-error">{{
								formErrors.confirmPassword
							}}</span>
						</label>
					</div>

					<!-- Role -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text">角色</span>
						</label>
						<select v-model="newUser.role" class="select select-bordered" required>
							<option value="operator">操作员</option>
							<option value="manager">店长</option>
						</select>
					</div>

					<!-- Error Alert -->
					<div v-if="createError" class="alert alert-error mt-4">
						<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
							viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						<span>{{ createError }}</span>
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