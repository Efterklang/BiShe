<template>
	<div>
		<!-- Header Section -->
		<div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
			<div>
				<h1 class="text-2xl font-bold tracking-tight text-base-content">
					实体商品管理
				</h1>
				<p class="mt-1 text-base-content/60">
					管理店内实体商品，控制库存和价格信息。
				</p>
			</div>
			<div class="flex gap-2">
				<button v-if="canManageProducts" @click="openCreateModal" class="btn btn-primary">
					<Plus class="w-4 h-4 mr-1" />
					添加商品
				</button>
			</div>
		</div>

		<!-- Stats Cards -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
			<div class="stats shadow-sm border border-base-200 bg-base-100 rounded-xl">
				<div class="stat">
					<div class="stat-title text-base-content/60">商品总数</div>
					<div class="stat-value text-primary">{{ stats.total_products }}</div>
				</div>
			</div>
			<div class="stats shadow-sm border border-base-200 bg-base-100 rounded-xl">
				<div class="stat">
					<div class="stat-title text-base-content/60">上架商品</div>
					<div class="stat-value text-success">{{ stats.active_products }}</div>
				</div>
			</div>
			<div class="stats shadow-sm border border-base-200 bg-base-100 rounded-xl">
				<div class="stat">
					<div class="stat-title text-base-content/60">库存总值</div>
					<div class="stat-value text-info text-2xl">
						¥{{ stats.total_value?.toFixed(2) || 0 }}
					</div>
				</div>
			</div>
			<div class="stats shadow-sm border border-base-200 bg-base-100 rounded-xl">
				<div class="stat">
					<div class="stat-title text-base-content/60">低库存</div>
					<div class="stat-value text-warning">{{ stats.low_stock_count }}</div>
					<div class="stat-desc text-error font-medium">缺货: {{ stats.out_of_stock_count }}</div>
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
		<div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			<div v-for="i in 6" :key="i" class="h-64 rounded-xl border border-base-300 bg-base-200 animate-pulse"></div>
		</div>

		<!-- Empty State -->
		<div v-else-if="products.length === 0"
			class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50 text-center">
			<div class="p-4 bg-base-300 rounded-full mb-4">
				<Package class="w-8 h-8 text-base-content/40" />
			</div>
			<h3 class="text-lg font-medium text-base-content">暂无商品</h3>
			<p class="text-base-content/60 mt-1">点击右上角按钮添加第一个商品。</p>
		</div>

		<!-- Products Grid -->
		<div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			<div v-for="product in products" :key="product.id"
				class="card bg-base-100 border border-base-300 hover:border-primary/50 transition-all duration-200 shadow-sm hover:shadow-md">
				<figure>
					<div class="w-full h-48 bg-base-200/50 flex items-center justify-center overflow-hidden">
						<img v-if="product.image_url" :src="product.image_url" :alt="product.name"
							class="w-full h-full object-cover transition-transform duration-500 hover:scale-105" />
						<Image v-else class="w-16 h-16 text-base-content/20" />
					</div>
				</figure>
				<div class="card-body p-5">
					<div class="flex items-start justify-between">
						<h2 class="card-title text-base-content text-lg line-clamp-1" :title="product.name">{{ product.name }}</h2>
						<div class="badge badge-sm" :class="product.is_active ? 'badge-success' : 'badge-ghost'">
							{{ product.is_active ? "上架" : "下架" }}
						</div>
					</div>

					<p v-if="product.description" class="text-sm text-base-content/60 line-clamp-2 min-h-[2.5em]">
						{{ product.description }}
					</p>
					<p v-else class="text-sm text-base-content/40 italic min-h-[2.5em]">暂无描述</p>

					<div class="grid grid-cols-2 gap-4 mt-2">
						<div>
							<div class="text-xs text-base-content/60 mb-0.5">零售价</div>
							<div class="text-lg font-bold text-primary font-mono">
								¥{{ product.retail_price }}
							</div>
						</div>
						<div>
							<div class="text-xs text-base-content/60 mb-0.5">进货价</div>
							<div class="text-lg font-semibold text-base-content/70 font-mono">
								¥{{ product.cost_price }}
							</div>
						</div>
					</div>

					<div class="mt-4 bg-base-200/50 rounded-lg p-3">
						<div class="flex items-center justify-between">
							<span class="text-sm text-base-content/60 font-medium">当前库存</span>
							<div class="flex items-center gap-3">
								<span class="font-bold text-lg font-mono" :class="{
									'text-error': product.stock === 0,
									'text-warning': product.stock > 0 && product.stock < 10,
									'text-success': product.stock >= 10,
								}">
									{{ product.stock }}
								</span>
								<button @click="openInventoryModal(product)" class="btn btn-xs btn-ghost btn-square"
									title="查看库存记录">
									<History class="w-4 h-4 text-base-content/60" />
								</button>
							</div>
						</div>
					</div>

					<div class="card-actions justify-end mt-4 pt-4 border-t border-base-200">
						<button @click="openStockChangeModal(product, 'restock')"
							class="btn btn-sm btn-success btn-outline gap-1" title="入库">
							<ArrowDownToLine class="w-3.5 h-3.5" />
							入库
						</button>
						<button @click="openStockChangeModal(product, 'sale')"
							class="btn btn-sm btn-warning btn-outline gap-1" title="出库">
							<ArrowUpFromLine class="w-3.5 h-3.5" />
							出库
						</button>
						<div v-if="canManageProducts" class="flex gap-1 ml-auto">
							<button @click="handleEdit(product)" class="btn btn-sm btn-ghost btn-square text-base-content/60 hover:text-primary">
								<Edit class="w-4 h-4" />
							</button>
							<button @click="handleDelete(product.id)"
								class="btn btn-sm btn-ghost btn-square text-base-content/60 hover:text-error">
								<Trash2 class="w-4 h-4" />
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Create/Edit Modal -->
		<dialog ref="createModalRef" class="modal">
			<div class="modal-box max-w-2xl bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden">
				<!-- Modal Header -->
				<div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
					<h3 class="font-bold text-lg flex items-center gap-2">
						<Package class="w-5 h-5 text-primary" />
						{{ editingId ? "编辑商品" : "添加新商品" }}
					</h3>
					<button @click="closeCreateModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
						<X class="w-5 h-5" />
					</button>
				</div>

				<div class="p-6">
					<form @submit.prevent="handleSubmit" class="space-y-4">
						<div class="grid grid-cols-2 gap-4">
							<div class="col-span-2">
								<label class="label">
									<span class="label-text font-medium">商品名称</span>
								</label>
								<input v-model="formData.name" type="text" placeholder="请输入商品名称"
									class="input input-bordered w-full" required />
							</div>

							<div>
								<label class="label">
									<span class="label-text font-medium">零售价 (元)</span>
								</label>
								<div class="relative">
									<span class="absolute left-3 top-1/2 -translate-y-1/2 text-base-content/60">¥</span>
									<input v-model.number="formData.retail_price" type="number" step="0.01" min="0"
										placeholder="0.00" class="input input-bordered w-full pl-7" required />
								</div>
							</div>

							<div>
								<label class="label">
									<span class="label-text font-medium">进货价 (元)</span>
								</label>
								<div class="relative">
									<span class="absolute left-3 top-1/2 -translate-y-1/2 text-base-content/60">¥</span>
									<input v-model.number="formData.cost_price" type="number" step="0.01" min="0"
										placeholder="0.00" class="input input-bordered w-full pl-7" required />
								</div>
							</div>

							<div v-if="!editingId">
								<label class="label">
									<span class="label-text font-medium">初始库存</span>
								</label>
								<input v-model.number="formData.stock" type="number" min="0" placeholder="0"
									class="input input-bordered w-full" required />
							</div>

							<div :class="editingId ? 'col-span-2' : ''">
								<label class="label">
									<span class="label-text font-medium">商品图片URL</span>
								</label>
								<div class="relative">
									<input v-model="formData.image_url" type="url" placeholder="https://example.com/image.jpg"
										class="input input-bordered w-full pl-10" />
									<Image class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
								</div>
							</div>

							<div class="col-span-2">
								<label class="label">
									<span class="label-text font-medium">商品描述</span>
								</label>
								<textarea v-model="formData.description" placeholder="请输入商品描述"
									class="textarea textarea-bordered w-full" rows="3"></textarea>
							</div>

							<div class="col-span-2">
								<label class="label cursor-pointer justify-start gap-3 p-0">
									<span class="label-text font-medium">立即上架</span>
									<input v-model="formData.is_active" type="checkbox" class="toggle toggle-success toggle-sm" />
								</label>
							</div>
						</div>

						<div class="pt-2">
							<button type="submit" class="btn btn-primary w-full" :class="{ loading: submitting }"
								:disabled="submitting">
								{{ submitting ? "保存中..." : editingId ? "保存修改" : "添加商品" }}
							</button>
						</div>
					</form>
				</div>
			</div>
			<form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
				<button @click="closeCreateModal">close</button>
			</form>
		</dialog>

		<!-- Stock Change Modal -->
		<dialog ref="stockModalRef" class="modal">
			<div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
				<div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
					<h3 class="font-bold text-lg flex items-center gap-2">
						<component :is="stockChangeType === 'restock' ? ArrowDownToLine : ArrowUpFromLine" class="w-5 h-5" :class="stockChangeType === 'restock' ? 'text-success' : 'text-warning'" />
						{{ stockChangeType === "restock" ? "商品入库" : "商品出库" }}
					</h3>
					<button @click="closeStockModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
						<X class="w-5 h-5" />
					</button>
				</div>

				<div class="p-6">
					<div v-if="selectedProduct" class="mb-6 p-4 bg-base-200/50 rounded-lg border border-base-200">
						<div class="font-semibold text-base">{{ selectedProduct.name }}</div>
						<div class="text-sm text-base-content/60 mt-1">
							当前库存: <span class="font-mono font-bold text-base-content">{{ selectedProduct.stock }}</span>
						</div>
					</div>

					<form @submit.prevent="handleStockChange" class="space-y-4">
						<div>
							<label class="label">
								<span class="label-text font-medium">
									{{ stockChangeType === "restock" ? "入库数量" : "出库数量" }}
								</span>
							</label>
							<input v-model.number="stockChangeAmount" type="number"
								:min="stockChangeType === 'restock' ? 1 : -selectedProduct?.stock || 0"
								:max="stockChangeType === 'sale' ? 0 : undefined" placeholder="请输入数量"
								class="input input-bordered w-full font-mono" required />
						</div>

						<div>
							<label class="label">
								<span class="label-text font-medium">备注</span>
							</label>
							<textarea v-model="stockChangeRemark" placeholder="请输入备注信息（可选）"
								class="textarea textarea-bordered w-full" rows="2"></textarea>
						</div>

						<div class="pt-2">
							<button type="submit" class="btn w-full"
								:class="stockChangeType === 'restock' ? 'btn-success text-white' : 'btn-warning text-white'"
								:disabled="submitting">
								{{ submitting ? "处理中..." : "确认" }}
							</button>
						</div>
					</form>
				</div>
			</div>
			<form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
				<button @click="closeStockModal">close</button>
			</form>
		</dialog>

		<!-- Inventory History Modal -->
		<dialog ref="inventoryModalRef" class="modal">
			<div class="modal-box max-w-4xl bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden h-[80vh] flex flex-col">
				<div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
					<h3 class="font-bold text-lg flex items-center gap-2">
						<History class="w-5 h-5 text-primary" />
						库存变动记录
					</h3>
					<button @click="closeInventoryModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
						<X class="w-5 h-5" />
					</button>
				</div>

				<div class="p-6 overflow-y-auto flex-1">
					<div v-if="selectedProduct" class="mb-6 p-4 bg-base-200/50 rounded-lg border border-base-200 flex justify-between items-center">
						<div>
							<div class="font-semibold text-lg">{{ selectedProduct.name }}</div>
							<div class="text-sm text-base-content/60 mt-1">
								当前库存: <span class="font-mono font-bold text-base-content">{{ selectedProduct.stock }}</span>
							</div>
						</div>
					</div>

					<div v-if="loadingInventory" class="flex justify-center py-12">
						<span class="loading loading-spinner loading-lg text-primary"></span>
					</div>

					<div v-else-if="inventoryLogs.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
						<History class="w-12 h-12 text-base-content/20 mb-3" />
						<p class="text-base-content/60">暂无库存变动记录</p>
					</div>

					<div v-else class="overflow-x-auto border border-base-200 rounded-lg">
						<table class="table table-zebra w-full">
							<thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
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
							<tbody class="text-sm">
								<tr v-for="log in inventoryLogs" :key="log.id">
									<td class="font-mono text-xs">{{ formatDate(log.created_at) }}</td>
									<td>
										<div class="badge badge-sm gap-1" :class="{
											'badge-success badge-outline': log.action_type === 'restock',
											'badge-warning badge-outline': log.action_type === 'sale',
											'badge-info badge-outline': log.action_type === 'adjustment',
										}">
											<ArrowDownToLine v-if="log.action_type === 'restock'" class="w-3 h-3" />
											<ArrowUpFromLine v-else-if="log.action_type === 'sale'" class="w-3 h-3" />
											{{
												log.action_type === "restock"
													? "入库"
													: log.action_type === "sale"
														? "销售"
														: "纠错"
											}}
										</div>
									</td>
									<td class="font-mono font-medium" :class="{
										'text-success': log.change_amount > 0,
										'text-error': log.change_amount < 0,
									}">
										{{ log.change_amount > 0 ? "+" : "" }}{{ log.change_amount }}
									</td>
									<td class="font-mono text-base-content/60">{{ log.before_stock }}</td>
									<td class="font-mono text-base-content">{{ log.after_stock }}</td>
									<td>{{ log.operator?.username || "-" }}</td>
									<td class="text-base-content/60 max-w-xs truncate" :title="log.remark">
										{{ log.remark || "-" }}
									</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
			<form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
				<button @click="closeInventoryModal">close</button>
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
import { 
    Plus, 
    Package, 
    Image, 
    History, 
    ArrowDownToLine, 
    ArrowUpFromLine, 
    Edit, 
    Trash2, 
    X 
} from 'lucide-vue-next';

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

// Modal Refs
const createModalRef = ref(null);
const stockModalRef = ref(null);
const inventoryModalRef = ref(null);

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
	createModalRef.value?.showModal();
};

const closeCreateModal = () => {
	createModalRef.value?.close();
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
	createModalRef.value?.showModal();
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

		closeCreateModal();
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
	stockModalRef.value?.showModal();
};

const closeStockModal = () => {
	stockModalRef.value?.close();
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

		closeStockModal();
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
	inventoryModalRef.value?.showModal();
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

const closeInventoryModal = () => {
	inventoryModalRef.value?.close();
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
.line-clamp-1 {
	display: -webkit-box;
	-webkit-line-clamp: 1;
	-webkit-box-orient: vertical;
	overflow: hidden;
}
.line-clamp-2 {
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
}
</style>