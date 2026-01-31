import api from "./axios";

/**
 * Calls the backend to generate an AI business analysis report.
 * @returns {Promise<Object>} The report data containing markdown text and raw stats.
 */
export const generateAIReport = async () => {
	try {
		return await api.get("/api/ai/report", { timeout: 120000 });
	} catch (error) {
		console.error("Failed to generate AI report:", error);
		throw error;
	}
};

export const generateMemberProfile = async (memberId) => {
	return await api.get(`/api/members/${memberId}/ai-profile`, { timeout: 120000 });
};
