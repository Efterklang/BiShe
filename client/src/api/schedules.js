import api from './axios';

/**
 * 获取排班列表
 * @param {Object} params - 查询参数
 * @param {string} params.start_date - 开始日期 YYYY-MM-DD
 * @param {string} params.end_date - 结束日期 YYYY-MM-DD
 * @param {Array} params.tech_ids - 技师ID数组
 */
export const getSchedules = (params) => {
  return api.get('/api/schedules', { params });
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
  return api.post('/api/schedules/batch', payload);
};
