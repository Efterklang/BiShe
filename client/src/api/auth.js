import api from "./axios";

/**
 * User login
 * @param {string} username
 * @param {string} password
 * @returns {Promise<{token: string, user: object}>}
 */
export const login = (username, password) => {
	return api.post("/api/auth/login", { username, password });
};

/**
 * Register a new user (manager only)
 * @param {string} username
 * @param {string} password
 * @param {string} role - "manager" or "operator"
 * @returns {Promise<object>}
 */
export const register = (username, password, role) => {
	return api.post("/api/auth/register", { username, password, role });
};

/**
 * Get current authenticated user
 * @returns {Promise<object>}
 */
export const getCurrentUser = () => {
	return api.get("/api/auth/me");
};

/**
 * Get list of all users (manager only)
 * @returns {Promise<{users: array, total: number}>}
 */
export const getUserList = () => {
	return api.get("/api/auth/users");
};

/**
 * Logout (client-side only, clears local storage)
 */
export const logout = () => {
	localStorage.removeItem("token");
	localStorage.removeItem("user");
};
