import api from "./axios";

export const getDashboardStats = () => {
	return api.get("/api/dashboard/stats");
};

export const getFissionRanking = () => {
	return api.get("/api/fission/ranking");
};
