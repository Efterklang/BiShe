import axios from "axios";

const api = axios.create({
	baseURL: "http://localhost:8080",
	timeout: 10000,
	headers: {
		"Content-Type": "application/json",
	},
});

// Request interceptor
api.interceptors.request.use(
	(config) => {
		// Add JWT token to Authorization header if available
		const token = localStorage.getItem("token");
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);

// Response interceptor
api.interceptors.response.use(
	(response) => {
		const res = response.data;
		// Standard response format: {"code": 200, "data": {}, "msg": ""}
		if (res.code !== 200) {
			console.error("API Error:", res.msg);
			return Promise.reject(new Error(res.msg || "Error"));
		} else {
			return res.data;
		}
	},
	(error) => {
		console.error("Network Error:", error);

		// Handle 401 Unauthorized - redirect to login
		if (error.response && error.response.status === 401) {
			// Clear authentication data
			localStorage.removeItem("token");
			localStorage.removeItem("user");

			// Redirect to login page
			window.location.href = "/login";
		}

		return Promise.reject(error);
	},
);

export default api;
