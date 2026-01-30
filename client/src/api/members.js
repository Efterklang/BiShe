import api from "./axios";

export const getMembers = () => {
	return api.get("/api/members");
};

export const createMember = (data) => {
	return api.post("/api/members", data);
};
