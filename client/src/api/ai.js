import axios from "axios";

// Create a dedicated client for AI requests with a longer timeout
const aiClient = axios.create({
	baseURL: "/api", // Assumes Vite proxy is configured to forward /api to the backend
	timeout: 120000, // 120 seconds timeout for LLM generation
});

/**
 * Calls the backend to generate an AI business analysis report.
 * @returns {Promise<Object>} The report data containing markdown text and raw stats.
 */
export const generateAIReport = async () => {
	try {
		const response = await aiClient.get("/ai/report");
		// Backend response format: { code: 200, data: { report: "...", raw_data: {...} }, msg: "..." }
		return response.data.data;
	} catch (error) {
		console.error("Failed to generate AI report:", error);
		throw error;
	}
};
