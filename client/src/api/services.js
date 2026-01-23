import api from "./axios";

export const getServices = (params) => {
	return api.get("/api/services", { params });
};

export const createService = (data) => {
	return api.post("/api/services", data);
};

export const updateService = (id, data) => {
	return api.put(`/api/services/${id}`, data);
};

export const deleteService = (id) => {
	return api.delete(`/api/services/${id}`);
};
