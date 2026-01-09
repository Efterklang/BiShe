<template>
	<div class="max-w-7xl mx-auto">
		<!-- Header Section -->
		<div
			class="flex flex-col md:flex-row md:items-center justify-between mb-10 gap-4"
		>
			<div>
				<h1 class="text-3xl font-bold tracking-tight text-base-content">
					实体商品管理
				</h1>
				<p class="mt-2 text-base-content/60">
					管理店内实体商品，控制库存和价格信息。
				</p>
			</div>
			<div class="flex gap-2">
				<button
					v-if="canManageProducts"
					@click="openCreateModal"
					class="btn btn-primary"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5 mr-2"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M12 4.5v15m7.5-7.5h-15"
						/>
					</svg>
					添加商品
				</button>
			</div>
		</div>

		<!-- Stats Cards -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
			<div class="stats shadow bg-base-100">
				<div class="stat">
					<div class="stat-title">商品总数</div>
					<div class="stat-value text-primary">{{ stats.total_products }}</div>
				</div>
			</div>
			<div class="stats shadow bg-base-100">
				<div class="stat">
					<div class="stat-title">上架商品</div>
					<div class="stat-value text-success">{{ stats.active_products }}</div>
				</div>
			</div>
			<div class="stats shadow bg-base-100">
				<div class="stat">
					<div class="stat-title">库存总值</div>
					<div class="stat-value text-info">
						¥{{ stats.total_value?.toFixed(2) || 0 }}
					</div>
				</div>
			</div>
			<div class="stats shadow bg-base-100">
				<div class="stat">
					<div class="stat-title">低库存</div>
					<div class="stat-value text-warning">{{ stats.low_stock_count }}</div>
					<div class="stat-desc">缺货: {{ stats.out_of_stock_count }}</div>
				</div>
			</div>
		</div>

		<!-- Filter -->
		<div class="mb-6 flex gap-2">
			<select v-model="filterStatus" class="select select-bordered" @change="fetchProducts">
				<option value="">全部商品</option>
				<option value="true">已上架</option>
				<option value="false">已下架</option>
			</select>
		</div>

		<!-- Loading State -->
		<div
			v-if="loading"
			class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
		>
			<div
				v-for="i in 6"
				:key="i"
				class="h-64 rounded-xl border border-base-300 bg-base-200 animate-pulse"
			></div>
		</div>

		<!-- Empty State -->
		<div
			v-else-if="products.length === 0"
			class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50"
		>
			<div class="p-4 bg-base-300 rounded-full mb-4">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-8 h-8 text-base-content/40"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"
					/>
				</svg>
			</div>
			<h3 class="text-lg font-medium text-base-content">暂无商品</h3>
			<p class="text-base-content/60 mt-1">点击右上角按钮添加第一个商品。</p>
		</div>

		<!-- Products Grid -->
		<div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			<div
				v-for="product in products"
				:key="product.id"
				class="card bg-base-100 border border-base-300 hover:border-primary/50 transition-all duration-200 shadow-sm hover:shadow-md"
			>
				<figure class="px-6 pt-6">
					<div
						class="w-full h-48 bg-base-200 rounded-xl flex items-center justify-center"
					>
						<img
							v-if="product.image_url"
							:src="product.image_url"
							:alt="product.name"
							class="w-full h-full object-cover rounded-xl"
						/>
						<svg
							v-else
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1"
							stroke="currentColor"
							class="w-16 h-16 text-base-content/20"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5m8.25 3v6.75m0 0l-3-3m3 3l3-3M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"
							/>
						</svg>
					</div>
				</figure>
				<div class="card-body">
					<div class="flex items-start justify-between">
						<h2 class="card-title text-base-content">{{ product.name }}</h2>
						<div
							class="badge"
							:class="product.is_active ? 'badge-success' : 'badge-ghost'"
						>
							{{ product.is_active ? "上架" : "下架" }}
						</div>
					</div>

					<p v-if="product.description" class="text-sm text-base-content/60 line-clamp-2">
						{{ product.description }}
					</p>

					<div class="grid grid-cols-2 gap-4 mt-4">
						<div>
							<div class="text-xs text-base-content/60">零售价</div>
							<div class="text-lg font-bold text-primary">
								¥{{ product.retail_price }}
							</div>
						</div>
						<div>
							<div class="text-xs text-base-content/60">进货价</div>
							<div class="text-lg font-semibold text-base-content/70">
								¥{{ product.cost_price }}
							</div>
						</div>
					</div>

					<div class="mt-4">
						<div class="flex items-center justify-between">
							<span class="text-sm text-base-content/60">库存</span>
							<div class="flex items-center gap-2">
								<span
									class="font-bold text-lg"
									:class="{
										'text-error': product.stock === 0,
										'text-warning': product.stock > 0 && product.stock < 10,
										'text-success': product.stock >= 10,
									}"
								>
									{{ product.stock }}
								</span>
								<button
									@click="openInventoryModal(product)"
									class="btn btn-xs btn-ghost"
									title="查看库存记录"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="1.5"
										stroke="currentColor"
										class="w-4 h-4"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 12h.007v.008H3.75V12zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm-.375 5.25h.007v.008H3.75v-.008zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"
										/>
									</svg>
								</button>
							</div>
						</div>
					</div>

					<div class="card-actions justify-end mt-4 pt-4 border-t border-base-200">
						<button
							@click="openStockChangeModal(product, 'restock')"
							class="btn btn-sm btn-success btn-outline"
							title="入库"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M12 4.5v15m7.5-7.5h-15"
								/>
							</svg>
							入库
						</button>
						<button
							@click="openStockChangeModal(product, 'sale')"
							class="btn btn-sm btn-warning btn-outline"
							title="出库"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M15 12H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
							出库
						</button>
						<button
							v-if="canManageProducts"
							@click="handleEdit(product)"
							class="btn btn-sm btn-ghost"
						>
							编辑
						</button>
						<button
							v-if="canManageProducts"
							@click="handleDelete(product.id)"
							class="btn btn-sm btn-ghost text-error"
						>
							删除
						</button>
					</div>
				</div>
			</div>
		</div>

		<!-- Create/Edit Modal -->
		<dialog class="modal" :class="{ 'modal-open': showModal }">
			<div class="modal-box max-w-2xl">
				<h3 class="font-bold text-lg mb-4">
					{{ editingId ? "编辑商品" : "添加新商品" }}
				</h3>

				<form @submit.prevent="handleSubmit" class="space-y-4">
					<div class="grid grid-cols-2 gap-4">
						<div class="col-span-2">
							<label class="label">
								<span class="label-text">商品名称</span>
							</label>
							<input
								v-model="formData.name"
								type="text"
								placeholder="请输入商品名称"
								class="input input-bordered w-full"
								required
							/>
						</div>

						<div>
							<label class="label">
								<span class="label-text">零售价 (元)</span>
							</label>
							<input
								v-model.number="formData.retail_price"
								type="number"
								step="0.01"
								min="0"
								placeholder="0.00"
								class="input input-bordered w-full"
								required
							/>
						</div>

						<div>
							<label class="label">
								<span class="label-text">进货价 (元)</span>
							</label>
							<input
								v-model.number="formData.cost_price"
								type="number"
								step="0.01"
								min="0"
								placeholder="0.00"
								class="input input-bordered w-full"
								required
							/>
						</div>

						<div v-if="!editingId">
							<label class="label">
								<span class="label-text">初始库存</span>
							</label>
							<input
								v-model.number="formData.stock"
								type="number"
								min="0"
								placeholder="0"
								class="input input-bordered w-full"
								required
							/>
						</div>

						<div :class="editingId ? 'col-span-2' : ''">
							<label class="label">
								<span class="label-text">商品图片URL</span>
							</label>
							<input
								v-model="formData.image_url"
								type="url"
								placeholder="https://example.com/image.jpg"
								class="input input-bordered w-full"
							/>
						</div>

						<div class="col-span-2">
							<label class="label">
								<span class="label-text">商品描述</span>
							</label>
							<textarea
								v-model="formData.description"
								placeholder="请输入商品描述"
								class="textarea textarea-bordered w-full"
								rows="3"
							></textarea>
						</div>

						<div class="col-span-2">
							<label class="label cursor-pointer justify-start gap-2">
								<input
									v-model="formData.is_active"
									type="checkbox"
									class="checkbox checkbox-primary"
								/>
								<span class="label-text">立即上架</span>
							</label>
						</div>
					</div>

					<div class="modal-action">
						<button
							type="button"
							class="btn btn-ghost"
							@click="showModal = false"
							:disabled="submitting"
						>
							取消
						</button>
						<button
							type="submit"
							class="btn btn-primary"
							:class="{ loading: submitting }"
							:disabled="submitting"
						>
							{{ submitting ? "保存中..." : editingId ? "保存修改" : "添加商品" }}
						</button>
					</div>
				</form>
			</div>
			<form method="dialog" class="modal-backdrop">
				<button @click="showModal = false">close</button>
			</form>
		</dialog>

		<!-- Stock Change Modal -->
		<dialog class="modal" :class="{ 'modal-open': showStockModal }">
			<div class="modal-box">
				<h3 class="font-bold text-lg mb-4">
					{{ stockChangeType === "restock" ? "商品入库" : "商品出库" }}
				</h3>

				<div v-if="selectedProduct" class="mb-4 p-4 bg-base-200 rounded-lg">
					<div class="font-semibold">{{ selectedProduct.name }}</div>
					<div class="text-sm text-base-content/60">
						当前库存: {{ selectedProduct.stock }}
					</div>
				</div>

				<form @submit.prevent="handleStockChange" class="space-y-4">
					<div>
						<label class="label">
							<span class="label-text">
								{{ stockChangeType === "restock" ? "入库数量" : "出库数量" }}
							</span>
						</label>
						<input
							v-model.number="stockChangeAmount"
							type="number"
							:min="stockChangeType === 'restock' ? 1 : -selectedProduct?.stock || 0"
							:max="stockChangeType === 'sale' ? 0 : undefined"
							placeholder="请输入数量"
							class="input input-bordered w-full"
							required
						/>
					</div>

					<div>
						<label class="label">
							<span class="label-text">备注</span>
						</label>
						<textarea
							v-model="stockChangeRemark"
							placeholder="请输入备注信息（可选）"
							class="textarea textarea-bordered w-full"
							rows="2"
						></textarea>
					</div>

					<div class="modal-action">
						<button
							type="button"
							class="btn btn-ghost"
							@click="showStockModal = false"
							:disabled="submitting"
						>
							取消
						</button>
						<button
							type="submit"
							class="btn"
							:class="stockChangeType === 'restock' ? 'btn-success' : 'btn-warning'"
							:disabled="submitting"
						>
							{{ submitting ? "处理中..." : "确认" }}
						</button>
					</div>
				</form>
			</div>
			<form method="dialog" class="modal-backdrop">
				<button @click="showStockModal = false">close</button>
			</form>
		</dialog>

		<!-- Inventory History Modal -->
		<dialog class="modal" :class="{ 'modal-open': showInventoryModal }">
			<div class="modal-box max-w-4xl">
				<h3 class="font-bold text-lg mb-4">库存变动记录</h3>

				<div v-if="selectedProduct" class="mb-4 p-4 bg-base-200 rounded-lg">
					<div class="font-semibold">{{ selectedProduct.name }}</div>
					<div class="text-sm text-base-content/60">
						当前库存: {{ selectedProduct.stock }}
					</div>
				</div>

				<div v-if="loadingInventory" class="flex justify-center py-8">
					<span class="loading loading-spinner loading-lg"></span>
				</div>

				<div v-else-if="inventoryLogs.length === 0" class="text-center py-8 text-base-content/60">
					暂无库存变动记录
				</div>

				<div v-else class="overflow-x-auto">
					<table class="table table-zebra">
						<thead>
							<tr>
								<th>时间</th>
								<th>类型</th>
								<th>变动量</th>
								<th>变动前</th>
								<th>变动后</th>
								<th>操作员</th>
								<th>备注</th>
							</tr>
						</thead>
						<tbody>
							<tr v-for="log in inventoryLogs" :key="log.id">
								<td class="text-sm">{{ formatDate(log.created_at) }}</td>
								<td>
									<div
										class="badge badge-sm"
										:class="{
											'badge-success': log.action_type === 'restock',
											'badge-warning': log.action_type === 'sale',
											'badge-info': log.action_type === 'adjustment',
										}"
									>
										{{
											log.action_type === "restock"
												? "入库"
												: log.action_type === "sale"
													? "销售"
													: "纠错"
										}}
									</div>
								</td>
								<td
									:class="{
										'text-success': log.change_amount > 0,
										'text-error': log.change_amount < 0,
									}"
								>
									{{ log.change_amount > 0 ? "+" : "" }}{{ log.change_amount }}
								</td>
								<td>{{ log.before_stock }}</td>
								<td>{{ log.after_stock }}</td>
								<td class="text-sm">{{ log.operator?.username || "-" }}</td>
								<td class="text-sm text-base-content/60">
									{{ log.remark || "-" }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<div class="modal-action">
					<button class="btn" @click="showInventoryModal = false">关闭</button>
				</div>
			</div>
			<form method="dialog" class="modal-backdrop">
				<button @click="showInventoryModal = false">close</button>
			</form>
		</dialog>
	</div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import {
	getProducts,
	getProductStats,
	createProduct,
	updateProduct,
	deleteProduct,
	createInventoryChange,
	getProductInventoryLogs,
} from "../api/products";
import { usePermission } from "../composables/usePermission";

const { canManageProducts } = usePermission();

const products = ref([]);
const stats = ref({
	total_products: 0,
	active_products: 0,
	total_value: 0,
	low_stock_count: 0,
	out_of_stock_count: 0,
});
const loading = ref(true);
const showModal = ref(false);
const showStockModal = ref(false);
const showInventoryModal = ref(false);
const submitting = ref(false);
const editingId = ref(null);
const filterStatus = ref("");

const selectedProduct = ref(null);
const stockChangeType = ref("restock");
const stockChangeAmount = ref(0);
const stockChangeRemark = ref("");
const inventoryLogs = ref([]);
const loadingInventory = ref(false);

const formData = reactive({
	name: "",
	stock: 0,
	retail_price: 0,
	cost_price: 0,
	description: "",
	image_url: "",
	is_active: true,
});

const fetchProducts = async () => {
	loading.value = true;
	try {
		const params = {};
		if (filterStatus.value) {
			params.is_active = filterStatus.value;
		}
		const res = await getProducts(params);
		products.value = res.products || [];
	} catch (error) {
		console.error("Failed to load products:", error);
	} finally {
		loading.value = false;
	}
};

const fetchStats = async () => {
	try {
		const res = await getProductStats();
		stats.value = res;
	} catch (error) {
		console.error("Failed to load stats:", error);
	}
};

onMounted(() => {
	fetchProducts();
	fetchStats();
});

const openCreateModal = () => {
	editingId.value = null;
	Object.assign(formData, {
		name: "",
		stock: 0,
		retail_price: 0,
		cost_price: 0,
		description: "",
		image_url: "",
		is_active: true,
	});
	showModal.value = true;
};

const handleEdit = (product) => {
	editingId.value = product.id;
	Object.assign(formData, {
		name: product.name,
		retail_price: product.retail_price,
		cost_price: product.cost_price,
		description: product.description || "",
		image_url: product.image_url || "",
		is_active: product.is_active,
	});
	showModal.value = true;
};

const handleSubmit = async () => {
	submitting.value = true;
	try {
		if (editingId.value) {
			await updateProduct(editingId.value, {
				name: formData.name,
				retail_price: formData.retail_price,
				cost_price: formData.cost_price,
				description: formData.description,
				image_url: formData.image_url,
				is_active: formData.is_active,
			});
		} else {
			await createProduct({
				name: formData.name,
				stock: formData.stock,
				retail_price: formData.retail_price,
				cost_price: formData.cost_price,
				description: formData.description,
				image_url: formData.image_url,
				is_active: formData.is_active,
			});
		}

		showModal.value = false;
		await fetchProducts();
		await fetchStats();
	} catch (error) {
		alert((editingId.value ? "更新" : "创建") + "失败: " + (error.message || "未知错误"));
	} finally {
		submitting.value = false;
	}
};

const handleDelete = async (id) => {
	if (!confirm("确定要删除该商品吗？")) return;
	try {
		await deleteProduct(id);
		await fetchProducts();
		await fetchStats();
	} catch (error) {
		alert("删除失败: " + (error.message || "未知错误"));
	}
};

const openStockChangeModal = (product, type) => {
	selectedProduct.value = product;
	stockChangeType.value = type;
	stockChangeAmount.value = type === "restock" ? 1 : -1;
	stockChangeRemark.value = "";
	showStockModal.value = true;
};

const handleStockChange = async () => {
	if (!selectedProduct.value) return;

	submitting.value = true;
	try {
		const changeAmount =
			stockChangeType.value === "restock"
				? Math.abs(stockChangeAmount.value)
				: -Math.abs(stockChangeAmount.value);

		await createInventoryChange({
			product_id: selectedProduct.value.id,
			change_amount: changeAmount,
			action_type: stockChangeType.value,
			remark: stockChangeRemark.value,
		});

		showStockModal.value = false;
		await fetchProducts();
		await fetchStats();
	} catch (error) {
		alert("操作失败: " + (error.message || "未知错误"));
	} finally {
		submitting.value = false;
	}
};

const openInventoryModal = async (product) => {
	selectedProduct.value = product;
	showInventoryModal.value = true;
	loadingInventory.value = true;

	try {
		const res = await getProductInventoryLogs(product.id);
		inventoryLogs.value = res.logs || [];
	} catch (error) {
		console.error("Failed to load inventory logs:", error);
	} finally {
		loadingInventory.value = false;
	}
};

const formatDate = (dateString) => {
	if (!dateString) return "-";
	const date = new Date(dateString);
	return date.toLocaleString("zh-CN", {
		year: "numeric",
		month: "2-digit",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit",
	});
};
</script>

<style scoped>
.line-clamp-2 {
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
}
</style>
