import api from "./axios";

export const getDashboardStats = () => {
	return api.get("/api/dashboard/stats");
};

export const getFissionRanking = () => {
	return api.get("/api/fission/ranking");
};

export const getRevenueTrend = (params) => {
	return api.get("/api/dashboard/revenue-trend", { params });
};

export const getServiceRanking = () => {
	return api.get("/api/dashboard/service-ranking");
};

export const getMonthlyStats = () => {
	return api.get("/api/dashboard/monthly-stats");
};

export const getProductSales = (params) => {
	return api.get("/api/dashboard/product-sales", { params });
};
