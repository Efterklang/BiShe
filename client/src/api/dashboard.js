import api from "./axios";

export const getDashboardStats = () => {
	return api.get("/api/dashboard/stats");
};

export const getFissionRanking = () => {
	return api.get("/api/fission/ranking");
};

export const getRevenueTrend = () => {
	return api.get("/api/dashboard/revenue-trend");
};

export const getServiceRanking = () => {
	return api.get("/api/dashboard/service-ranking");
};

export const getMonthlyStats = () => {
	return api.get("/api/dashboard/monthly-stats");
};
