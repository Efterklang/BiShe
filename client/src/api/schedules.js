import api from "./axios";

/**
 * 获取排班列表
 * @param {Object} params - 查询参数
 * @param {string} params.start_date - 开始日期 YYYY-MM-DD
 * @param {string} params.end_date - 结束日期 YYYY-MM-DD
 * @param {Array} params.tech_ids - 技师ID数组
 */
export const getSchedules = (params) => {
	return api.get("/api/schedules", { params });
};

/**
 * 批量设置排班
 * @param {Object} payload
 * @param {Array<number>} payload.tech_ids - 技师ID数组
 * @param {Array<string>} payload.dates - 日期数组 ["2023-10-01"]
 * @param {boolean} payload.is_available - 是否在岗
 */
export const batchSetSchedule = (payload) => {
	// Adapter for legacy calls
	if (payload.tech_id && !payload.tech_ids) {
		payload = { ...payload, tech_ids: [payload.tech_id] };
	}
	return api.post("/api/schedules/batch", payload);
};

/**
 * 获取技师排班详情（包含预约信息）
 * @param {Object} params - 查询参数
 * @param {number} params.tech_id - 技师ID
 * @param {string} params.date - 日期 YYYY-MM-DD
 */
export const getTechnicianScheduleDetail = (params) => {
	return api.get("/api/schedules/detail", { params });
};

/**
 * 获取指定时间段的可用技师列表
 * @param {Object} params - 查询参数
 * @param {string} params.start_time - 开始时间 RFC3339格式
 * @param {number} params.service_id - 服务项目ID
 */
export const getAvailableTechnicians = (params) => {
	return api.get("/api/schedules/available-technicians", { params });
};

/**
 * 获取时间段可用性
 * @param {Object} params - 查询参数
 * @param {string} params.date - 日期 YYYY-MM-DD
 * @param {number} params.service_id - 服务项目ID
 */
export const getTimeSlots = (params) => {
	return api.get("/api/schedules/slots", { params });
};
