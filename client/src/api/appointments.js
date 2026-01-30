import api from "./axios";

export const getAppointments = (params) => {
	return api.get("/api/appointments", { params });
};

export const createAppointment = (data) => {
	return api.post("/api/appointments", data);
};

export const cancelAppointment = (id) => {
	return api.put(`/api/appointments/${id}/cancel`);
};

export const completeAppointment = (id, data) => {
	return api.put(`/api/appointments/${id}/complete`, data);
};
