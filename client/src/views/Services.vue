<script setup>
import { ref, onMounted } from "vue";
import {
    getServices,
    createService,
    updateService,
    deleteService,
} from "../api/services";
import { usePermission } from "../composables/usePermission";

const { canManageServices } = usePermission();

const services = ref([]);
const loading = ref(true);
const showModal = ref(false);
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
    showModal.value = true;
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
    showModal.value = true;
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

        showModal.value = false;
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
    <div class="max-w-7xl mx-auto">
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between mb-10 gap-4">
            <div>
                <h1 class="text-3xl font-bold tracking-tight text-base-content">
                    服务项目
                </h1>
                <p class="mt-2 text-base-content/60">
                    配置店内服务菜单，设定时长与价格标准。
                </p>
            </div>
            <button v-if="canManageServices" @click="openCreateModal" class="btn btn-neutral">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                    stroke="currentColor" class="w-4 h-4 mr-2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
                新增项目
            </button>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div v-for="i in 6" :key="i" class="h-40 rounded-xl border border-base-300 bg-base-200 animate-pulse"></div>
        </div>

        <!-- Empty State -->
        <div v-else-if="services.length === 0"
            class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50">
            <div class="p-4 bg-base-300 rounded-full mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-8 h-8 text-base-content/40">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25zM6.75 12h.008v.008H6.75V12zm0 3h.008v.008H6.75V15zm0 3h.008v.008H6.75V18z" />
                </svg>
            </div>
            <h3 class="text-lg font-medium text-base-content">暂无服务项目</h3>
            <p class="text-base-content/60 mt-1">请添加您的第一个服务项目。</p>
        </div>

        <!-- Services Grid -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div v-for="service in services" :key="service.id"
                class="group relative flex flex-col bg-base-100 border border-base-300 rounded-xl overflow-hidden hover:border-base-content/20 transition-all duration-200 hover:shadow-sm">
                <figure class="px-6 pt-6">
                    <div class="w-full h-40 bg-base-200 rounded-xl flex items-center justify-center">
                        <img v-if="service.image_url" :src="service.image_url" :alt="service.name"
                            class="w-full h-full object-cover rounded-xl" />
                        <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1"
                            stroke="currentColor" class="w-16 h-16 text-base-content/20">
                            <path stroke-linecap="round" stroke-linejoin="round"
                                d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.25 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25zM6.75 12h.008v.008H6.75V12zm0 3h.008v.008H6.75V15zm0 3h.008v.008H6.75V18z" />
                        </svg>
                    </div>
                </figure>
                <div class="p-6 flex flex-col flex-grow">
                    <div class="flex justify-between items-start mb-4">
                    <div>
                        <h3 class="text-lg font-semibold text-base-content group-hover:text-primary transition-colors">
                            {{ service.name }}
                        </h3>
                        <div class="flex items-center gap-2 mt-1 text-sm text-base-content/60">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-4 h-4">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            <span>{{ service.duration }} 分钟</span>
                        </div>
                    </div>
                    <div class="text-right">
                        <span class="block text-xl font-bold text-base-content">¥{{ service.price }}</span>
                    </div>
                </div>

                <div class="mt-auto pt-4 border-t border-base-200 flex items-center justify-between">
                    <div class="flex items-center gap-2">
                        <span class="relative flex h-2.5 w-2.5">
                            <span v-if="service.is_active || service.IsActive"
                                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
                            <span class="relative inline-flex rounded-full h-2.5 w-2.5" :class="service.is_active || service.IsActive
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

                    <div v-if="canManageServices" class="flex gap-2">
                        <button @click="handleEdit(service)"
                            class="btn btn-ghost btn-sm btn-square text-base-content/60 hover:text-base-content">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-4 h-4">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                            </svg>
                        </button>
                        <button @click="handleDelete(service.id)"
                            class="btn btn-ghost btn-sm btn-square text-base-content/60 hover:text-error hover:bg-error/10">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-4 h-4">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                            </svg>
                        </button>
                    </div>
                </div>
                </div>
            </div>
        </div>

        <!-- Create Modal -->
        <dialog class="modal" :class="{ 'modal-open': showModal }">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <h3 class="font-semibold text-lg text-base-content">
                        {{ editingId ? "编辑服务项目" : "新增服务项目" }}
                    </h3>
                    <button @click="showModal = false" class="btn btn-ghost btn-sm btn-square text-base-content/60">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                            stroke="currentColor" class="w-5 h-5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
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
                            <input type="url" v-model="formData.image_url" placeholder="https://example.com/image.jpg"
                                class="input input-bordered w-full bg-base-100" />
                        </div>

                        <div class="flex items-center justify-between py-2">
                            <span class="text-sm font-medium text-base-content/80">立即上架</span>
                            <input type="checkbox" v-model="formData.is_active"
                                class="toggle toggle-success toggle-sm" />
                        </div>

                        <div class="pt-2">
                            <button type="submit" class="btn btn-neutral w-full" :disabled="submitting">
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
                <button @click="showModal = false">close</button>
            </form>
        </dialog>
    </div>
</template>