import api from "./axios";

/**
 * Get list of all physical products
 * @param {object} params - Query parameters
 * @param {boolean} params.is_active - Filter by active status
 * @returns {Promise<{products: array, total: number}>}
 */
export const getProducts = (params = {}) => {
	return api.get("/api/products", { params });
};

/**
 * Get a single product by ID
 * @param {number} id - Product ID
 * @returns {Promise<object>}
 */
export const getProduct = (id) => {
	return api.get(`/api/products/${id}`);
};

/**
 * Create a new physical product
 * @param {object} data - Product data
 * @param {string} data.name - Product name
 * @param {number} data.stock - Initial stock
 * @param {number} data.retail_price - Retail price
 * @param {number} data.cost_price - Cost price
 * @param {string} data.description - Product description
 * @param {string} data.image_url - Product image URL
 * @param {boolean} data.is_active - Whether product is active
 * @returns {Promise<object>}
 */
export const createProduct = (data) => {
	return api.post("/api/products", data);
};

/**
 * Update an existing product
 * @param {number} id - Product ID
 * @param {object} data - Product data to update
 * @returns {Promise<object>}
 */
export const updateProduct = (id, data) => {
	return api.put(`/api/products/${id}`, data);
};

/**
 * Delete a product
 * @param {number} id - Product ID
 * @returns {Promise<void>}
 */
export const deleteProduct = (id) => {
	return api.delete(`/api/products/${id}`);
};

/**
 * Get product statistics
 * @returns {Promise<object>}
 */
export const getProductStats = () => {
	return api.get("/api/products/stats");
};

/**
 * Get list of inventory logs
 * @param {object} params - Query parameters
 * @param {number} params.product_id - Filter by product ID
 * @param {string} params.action_type - Filter by action type (restock/sale/adjustment)
 * @param {number} params.page - Page number
 * @param {number} params.page_size - Page size
 * @returns {Promise<{logs: array, total: number, page: number, page_size: number}>}
 */
export const getInventoryLogs = (params = {}) => {
	return api.get("/api/inventory/logs", { params });
};

/**
 * Get inventory logs for a specific product
 * @param {number} productId - Product ID
 * @returns {Promise<{logs: array, total: number}>}
 */
export const getProductInventoryLogs = (productId) => {
	return api.get(`/products/${productId}/inventory-logs`).then((res) => res.data.data);
};

// Get all inventory sale logs for history page
export const getAllInventoryLogs = (params = {}) => {
	return api.get('/inventory-logs', { params }).then((res) => res.data.data);
};

/**
 * Create an inventory change record
 * @param {object} data - Inventory change data
 * @param {number} data.product_id - Product ID
 * @param {number} data.change_amount - Change amount (positive for restock, negative for sale)
 * @param {string} data.action_type - Action type: restock, sale, or adjustment
 * @param {string} data.remark - Remark/note
 * @returns {Promise<object>}
 */
export const createInventoryChange = (data) => {
	return api.post("/api/inventory/change", data);
};

/**
 * Batch restock multiple products
 * @param {array} items - Array of restock items
 * @param {number} items[].product_id - Product ID
 * @param {number} items[].quantity - Quantity to restock
 * @param {string} items[].remark - Remark
 * @returns {Promise<object>}
 */
export const batchRestock = (items) => {
	return api.post("/api/inventory/batch-restock", { items });
};

/**
 * Get inventory statistics
 * @returns {Promise<object>}
 */
export const getInventoryStats = () => {
	return api.get("/api/inventory/stats");
};
