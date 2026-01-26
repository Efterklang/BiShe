<script setup>
import { ref, onMounted } from "vue";
import {
    getServices,
    createService,
    updateService,
    deleteService,
} from "../api/services";
import { usePermission } from "../composables/usePermission";
import {
    Plus,
    Image,
    Clock,
    Edit,
    Trash2,
    X,
    Package,
    PackagePlus,
    Archive
} from 'lucide-vue-next';

const { canManageServices } = usePermission();

const services = ref([]);
const loading = ref(true);
const createModalRef = ref(null);
const submitting = ref(false);
const editingId = ref(null);

const formData = ref({
    name: "",
    duration: 60,
    price: 0,
    is_active: true,
    image_url: "",
});

const fetchServices = async () => {
    loading.value = true;
    try {
        const res = await getServices();
        services.value = res || [];
    } catch (error) {
        console.error("Failed to load services:", error);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchServices);

const openCreateModal = () => {
    editingId.value = null;
    formData.value = { name: "", duration: 60, price: 0, is_active: true, image_url: "" };
    createModalRef.value?.showModal();
};

const closeCreateModal = () => {
    createModalRef.value?.close();
};

const handleEdit = (service) => {
    editingId.value = service.id;
    formData.value = {
        name: service.name,
        duration: service.duration,
        price: service.price,
        is_active: service.is_active || service.IsActive,
        image_url: service.image_url || "",
    };
    createModalRef.value?.showModal();
};

const handleDelete = async (id) => {
    if (!confirm("确定要删除该服务项目吗？")) return;
    try {
        await deleteService(id);
        await fetchServices();
    } catch (error) {
        alert("删除失败: " + (error.message || "未知错误"));
    }
};

const handleSubmit = async () => {
    submitting.value = true;
    try {
        const payload = {
            name: formData.value.name,
            duration: Number(formData.value.duration),
            price: Number(formData.value.price),
            is_active: formData.value.is_active,
            image_url: formData.value.image_url,
        };

        if (editingId.value) {
            await updateService(editingId.value, payload);
        } else {
            await createService(payload);
        }

        closeCreateModal();
        await fetchServices();
    } catch (error) {
        alert(
            (editingId.value ? "更新" : "创建") +
            "失败: " +
            (error.message || "未知错误"),
        );
    } finally {
        submitting.value = false;
    }
};
</script>

<template>
    <div>
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
            <div>
                <h1 class="text-2xl font-bold text-base-content">
                    服务项目
                </h1>
                <p class="mt-1 text-base-content/60">
                    配置店内服务菜单，设定时长与价格标准。
                </p>
            </div>
            <button v-if="canManageServices" @click="openCreateModal" class="btn btn-primary">
                <Plus class="w-4 h-4 mr-1" />
                新增项目
            </button>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div v-for="i in 6" :key="i" class="h-40 rounded-xl border border-base-300 bg-base-200 animate-pulse"></div>
        </div>

        <!-- Empty State -->
        <div v-else-if="services.length === 0"
            class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50 text-center">
            <div class="p-4 bg-base-300 rounded-full mb-4">
                <Package class="w-8 h-8 text-base-content/40" />
            </div>
            <h3 class="text-lg font-medium text-base-content">暂无服务项目</h3>
            <p class="text-base-content/60 mt-1">请添加您的第一个服务项目。</p>
        </div>

        <!-- Services Grid -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div v-for="service in services" :key="service.id"
                class="group relative flex flex-col bg-base-100 border border-base-300 rounded-xl overflow-hidden hover:border-primary/50 transition-all duration-200 hover:shadow-sm">
                <figure>
                    <div class="w-full h-48 bg-base-200/50 flex items-center justify-center overflow-hidden">
                        <img v-if="service.image_url" :src="service.image_url" :alt="service.name"
                            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" />
                        <Image v-else class="w-16 h-16 text-base-content/20" />
                    </div>
                </figure>
                <div class="p-5 flex flex-col grow">
                    <div class="flex justify-between items-start mb-3">
                        <div>
                            <h3 class="font-semibold text-lg line-clamp-1" :title="service.name"> {{ service.name }}
                            </h3>
                            <div class="flex items-center gap-1.5 mt-1 text-sm text-base-content/60">
                                <Clock class="w-3.5 h-3.5" />
                                <span>{{ service.duration }} 分钟</span>
                            </div>
                        </div>
                        <div class="text-right shrink-0">
                            <span class="block text-xl font-bold text-base-content font-mono">¥{{ service.price
                                }}</span>
                        </div>
                    </div>

                    <div class="mt-auto pt-4 border-t border-base-200 flex items-center justify-between">
                        <div class="flex items-center gap-2">
                            <span class="relative flex h-2 w-2">
                                <span v-if="service.is_active || service.IsActive"
                                    class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
                                <span class="relative inline-flex rounded-full h-2 w-2" :class="service.is_active || service.IsActive
                                    ? 'bg-success'
                                    : 'bg-base-300'
                                    "></span>
                            </span>
                            <span class="text-xs font-medium" :class="service.is_active || service.IsActive
                                ? 'text-success'
                                : 'text-base-content/60'
                                ">
                                {{
                                    service.is_active || service.IsActive
                                        ? "上架中"
                                        : "已下架"
                                }}
                            </span>
                        </div>

                        <div v-if="canManageServices" class="flex gap-1">
                            <button @click="handleEdit(service)"
                                class="btn btn-ghost btn-sm btn-square text-base-content/60 hover:text-primary hover:bg-primary/10">
                                <Edit class="w-4 h-4" />
                            </button>
                            <button @click="handleDelete(service.id)"
                                class="btn btn-ghost btn-sm btn-square text-base-content/60 hover:text-error hover:bg-error/10">
                                <Trash2 class="w-4 h-4" />
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Create Modal -->
        <dialog ref="createModalRef" class="modal">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <h3 class="font-semibold text-lg text-base-content flex items-center gap-2">
                        <PackagePlus v-if="!editingId" class="w-5 h-5 text-primary" />
                        <Package v-else class="w-5 h-5 text-primary" />
                        {{ editingId ? "编辑服务项目" : "新增服务项目" }}
                    </h3>
                    <button @click="closeCreateModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
                        <X class="w-5 h-5" />
                    </button>
                </div>

                <!-- Modal Body -->
                <div class="p-6">
                    <form @submit.prevent="handleSubmit" class="space-y-5">
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1">项目名称</label>
                            <input type="text" v-model="formData.name" placeholder="例如: 全身精油SPA"
                                class="input input-bordered w-full bg-base-100" required />
                        </div>

                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <label class="block text-sm font-medium text-base-content/80 mb-1">价格 (元)</label>
                                <div class="relative">
                                    <span class="absolute left-3 top-3 text-base-content/60">¥</span>
                                    <input type="number" v-model="formData.price" min="0"
                                        class="input input-bordered w-full pl-7 bg-base-100" required />
                                </div>
                            </div>
                            <div>
                                <label class="block text-sm font-medium text-base-content/80 mb-1">时长 (分钟)</label>
                                <input type="number" v-model="formData.duration" min="1"
                                    class="input input-bordered w-full bg-base-100" required />
                            </div>
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1">服务图片URL</label>
                            <div class="relative">
                                <input type="url" v-model="formData.image_url"
                                    placeholder="https://example.com/image.jpg"
                                    class="input input-bordered w-full pl-10 bg-base-100" />
                                <Image class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
                            </div>
                        </div>

                        <div class="flex items-center justify-between py-2 px-1">
                            <span class="text-sm font-medium text-base-content/80">立即上架</span>
                            <input type="checkbox" v-model="formData.is_active"
                                class="toggle toggle-success toggle-sm" />
                        </div>

                        <div class="pt-2">
                            <button type="submit" class="btn btn-primary w-full" :disabled="submitting">
                                <span v-if="submitting" class="loading loading-spinner loading-xs mr-2"></span>
                                {{
                                    submitting
                                        ? "保存中..."
                                        : editingId
                                            ? "确认修改"
                                            : "确认新增"
                                }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
                <button @click="closeCreateModal">close</button>
            </form>
        </dialog>
    </div>
</template>