<template>
	<div class="min-h-screen flex items-center justify-center bg-base-200">
		<div class="card w-full max-w-md bg-base-100 shadow-xl">
			<div class="card-body">
				<!-- Logo and Title -->
				<div class="text-center mb-6">
					<h1 class="text-3xl font-bold text-primary">养生店管理系统</h1>
				</div>

				<!-- Login Form -->
				<form @submit.prevent="handleLogin">
					<!-- Username Input -->
					<div class="form-control">
						<label class="label">
							<span class="label-text text-lg">用户名</span>
						</label>
						<input v-model="formData.username" type="text" placeholder="请输入用户名"
							class="input input-bordered w-full" :class="{ 'input-error': errors.username }" required
							autocomplete="username" />
						<label v-if="errors.username" class="label">
							<span class="label-text-alt text-error">{{ errors.username }}</span>
						</label>
					</div>

					<!-- Password Input -->
					<div class="form-control mt-4">
						<label class="label">
							<span class="label-text text-lg">密码</span>
						</label>
						<div class="relative">
							<input v-model="formData.password" :type="showPassword ? 'text' : 'password'"
								placeholder="请输入密码" class="input input-bordered w-full pr-10"
								:class="{ 'input-error': errors.password }" required autocomplete="current-password" />
							<button type="button"
								class="absolute right-3 top-1/2 -translate-y-1/2 text-base-content/50 hover:text-base-content"
								@click="showPassword = !showPassword">
								<Eye v-if="showPassword" class="h-5 w-5" />
								<EyeOff v-else class="h-5 w-5" />
							</button>
						</div>
						<label v-if="errors.password" class="label">
							<span class="label-text-alt text-error">{{ errors.password }}</span>
						</label>
					</div>

					<!-- Error Message -->
					<div v-if="loginError" class="alert alert-error mt-4">
						<AlertCircle class="stroke-current shrink-0 h-6 w-6" />
						<span>{{ loginError }}</span>
					</div>

					<!-- Submit Button -->
					<div class="form-control mt-6">
						<button type="submit" class="btn btn-primary" :class="{ loading: isLoading }"
							:disabled="isLoading">
							<span v-if="!isLoading">登录</span>
							<span v-else>登录中...</span>
						</button>
					</div>
				</form>

				<!-- Default Account Info -->
				<div class="divider">默认账号信息</div>
				<div class="bg-base-200 p-4 rounded-lg text-sm">
					<p class="font-semibold mb-2">管理员账号：</p>
					<p>用户名: <code class="bg-base-300 px-2 py-1 rounded">admin</code></p>
					<p>密码: <code class="bg-base-300 px-2 py-1 rounded">Admin123!</code></p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useAppStore } from "../stores/app";
import { Eye, EyeOff, AlertCircle } from 'lucide-vue-next';

const router = useRouter();
const appStore = useAppStore();

const formData = reactive({
	username: "",
	password: "",
});

const errors = reactive({
	username: "",
	password: "",
});

const loginError = ref("");
const isLoading = ref(false);
const showPassword = ref(false);

// Validate form
const validateForm = () => {
	errors.username = "";
	errors.password = "";
	loginError.value = "";

	let isValid = true;

	if (!formData.username || formData.username.trim().length < 3) {
		errors.username = "用户名至少3个字符";
		isValid = false;
	}

	if (!formData.password || formData.password.length < 6) {
		errors.password = "密码至少6个字符";
		isValid = false;
	}

	return isValid;
};

// Handle login
const handleLogin = async () => {
	if (!validateForm()) {
		return;
	}

	isLoading.value = true;
	loginError.value = "";

	try {
		await appStore.login({
			username: formData.username.trim(),
			password: formData.password,
		});

		// Login successful, redirect to dashboard
		router.push("/");
	} catch (error) {
		console.error("Login error:", error);

		// Display error message
		if (error.response?.data?.msg) {
			loginError.value = error.response.data.msg;
		} else if (error.message) {
			loginError.value = error.message;
		} else {
			loginError.value = "登录失败，请检查用户名和密码";
		}
	} finally {
		isLoading.value = false;
	}
};
</script>

<style scoped>
code {
	font-family: "Courier New", Courier, monospace;
}
</style>
