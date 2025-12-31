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
		return Promise.reject(error);
	},
);

export default api;
