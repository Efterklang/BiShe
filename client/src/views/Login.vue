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
								<svg v-if="showPassword" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
									viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
										d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
								</svg>
								<svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
									viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
										d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
										d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
								</svg>
							</button>
						</div>
						<label v-if="errors.password" class="label">
							<span class="label-text-alt text-error">{{ errors.password }}</span>
						</label>
					</div>

					<!-- Error Message -->
					<div v-if="loginError" class="alert alert-error mt-4">
						<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
							viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
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
