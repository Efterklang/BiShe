import api from "./axios";

export const listOrders = (params = {}) => {
	return api.get("/api/orders", { params });
};
