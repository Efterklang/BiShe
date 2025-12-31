import api from "./axios";

export const getTechnicians = () => {
	return api.get("/api/technicians");
};

export const createTechnician = (data) => {
	return api.post("/api/technicians", data);
};

export const updateTechnician = (id, data) => {
	return api.put(`/api/technicians/${id}`, data);
};

export const deleteTechnician = (id) => {
	return api.delete(`/api/technicians/${id}`);
};
