import api from "./axios";

export const getMembers = () => {
	return api.get("/api/members");
};

export const createMember = (data) => {
	return api.post("/api/members", data);
};

export const updateMemberBalance = (id, balance) => {
	return api.put(`/api/members/${id}/balance`, { balance });
};

export const deleteMember = (id) => {
	return api.delete(`/api/members/${id}`);
};
