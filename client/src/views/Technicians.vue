<script setup>
import { ref, onMounted } from "vue";
import {
    getTechnicians,
    createTechnician,
    updateTechnician,
    deleteTechnician,
} from "../api/technicians";
import { getServices } from "../api/services";
import { getAppointments } from "../api/appointments";
import TechnicianSchedule from "../components/TechnicianSchedule.vue";
import Avatar from "../components/Avatar.vue";
import { usePermission } from "../composables/usePermission";
import {
    Plus,
    Users,
    Star,
    X,
    Calendar,
    Clock,
    Banknote,
    ChevronLeft,
    ChevronRight,
    UserPlus,
    User,
    Briefcase
} from 'lucide-vue-next';
import "cally";

const { canManageTechnicians } = usePermission();

const activeTab = ref("overview");
const technicians = ref([]);
const loading = ref(true);

// Modals refs
const createModalRef = ref(null);
const skillsModalRef = ref(null);
const appointmentModalRef = ref(null);

const submitting = ref(false);
const editingId = ref(null);

// Service list for skills selection
const services = ref([]);

// Appointment Modal
const appointmentModalLoading = ref(false);
const selectedAppointmentTech = ref(null);
const selectedAppointmentDate = ref(new Date().toISOString().split('T')[0]);
const calendarOpen = ref(false);
const technicianAppointments = ref([]);

const formData = ref({
    name: "",
    skills: [], // Changed to array of service IDs
    status: 0,
});

const fetchTechnicians = async () => {
    loading.value = true;
    try {
        const res = await getTechnicians();
        technicians.value = res || [];
    } catch (error) {
        console.error("Failed to load technicians:", error);
    } finally {
        loading.value = false;
    }
};

// Fetch services list for skills selection
const fetchServices = async () => {
    try {
        const res = await getServices();
        services.value = res || [];
    } catch (error) {
        console.error("Failed to fetch services:", error);
        services.value = [];
    }
};

onMounted(fetchTechnicians);

const openCreateModal = () => {
    editingId.value = null;
    formData.value = { name: "", skills: [], status: 0 };
    createModalRef.value?.showModal();
    fetchServices();
};

const closeCreateModal = () => {
    createModalRef.value?.close();
};

const handleEdit = (tech) => {
    editingId.value = tech.id;
    let skillsArray = [];

    if (tech.skills || tech.Skills) {
        try {
            const skillsData = JSON.parse(tech.skills || tech.Skills);
            if (Array.isArray(skillsData)) {
                // Check if array contains IDs (new format) or strings (old format)
                if (skillsData.length > 0 && typeof skillsData[0] === 'number') {
                    // New format: array of service IDs - use directly
                    skillsArray = skillsData;
                } else {
                    // Old format: array of service names - convert to IDs
                    skillsArray = skillsData.map(name => {
                        const service = services.value.find(s => s.name === name);
                        return service ? service.id : null;
                    }).filter(id => id !== null);
                }
            }
        } catch (e) {
            console.error("Failed to parse skills:", e);
            skillsArray = [];
        }
    }

    formData.value = {
        name: tech.name,
        skills: skillsArray, // Array of IDs
        status: tech.status,
    };
    createModalRef.value?.showModal();
    fetchServices();
};

const handleSchedule = (tech) => {
    selectedAppointmentTech.value = tech;
    selectedAppointmentDate.value = new Date().toISOString().split('T')[0];
    appointmentModalRef.value?.showModal();
    fetchTechnicianAppointments();
};

const closeAppointmentModal = () => {
    appointmentModalRef.value?.close();
};

const fetchTechnicianAppointments = async () => {
    if (!selectedAppointmentTech.value || !selectedAppointmentDate.value) return;

    appointmentModalLoading.value = true;
    try {
        const allAppts = await getAppointments();
        technicianAppointments.value = (allAppts || []).filter((app) =>
            app.tech_id === selectedAppointmentTech.value.id &&
            app.start_time.startsWith(selectedAppointmentDate.value),
        );
    } catch (error) {
        console.error("Failed to fetch appointments:", error);
        technicianAppointments.value = [];
    } finally {
        appointmentModalLoading.value = false;
    }
};

const handleDateChange = (event) => {
    const newDate = event.target?.value || event.detail?.value || event.detail;
    if (newDate) {
        selectedAppointmentDate.value = newDate;
        calendarOpen.value = false;
        fetchTechnicianAppointments();
    }
};

const handleDateSelect = (event) => {
    handleDateChange(event);
};

const formatDisplayDate = (dateStr) => {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'short'
    });
};

const handleDelete = async (tech) => {
    if (
        !confirm(
            `确定要删除技师 ${tech.name} 吗？如果该技师有待服务的订单，订单将被移至候补中。`,
        )
    ) {
        return;
    }

    try {
        await deleteTechnician(tech.id);
        alert("删除成功");
        await fetchTechnicians();
    } catch (error) {
        alert("删除失败: " + (error.response?.data?.msg || error.message));
    }
};

const handleSubmit = async () => {
    submitting.value = true;
    try {
        const payload = {
            name: formData.value.name,
            skills: formData.value.skills, // 直接发送ID数组
            status: Number(formData.value.status),
        };

        if (editingId.value) {
            await updateTechnician(editingId.value, payload);
        } else {
            await createTechnician(payload);
        }

        closeCreateModal();
        await fetchTechnicians();
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

const getStatusInfo = (status) => {
    switch (Number(status)) {
        case 0:
            return { text: "空闲", class: "badge-success" };
        case 1:
            return { text: "忙碌", class: "badge-warning" };
        case 2:
            return { text: "请假", class: "badge-ghost" };
        default:
            return { text: "未知", class: "badge-ghost" };
    }
};

// Toggle skill selection
const toggleSkill = (serviceId) => {
    const index = formData.value.skills.indexOf(serviceId);
    if (index === -1) {
        formData.value.skills.push(serviceId);
    } else {
        formData.value.skills.splice(index, 1);
    }
};

// Check if skill is selected
const isSkillSelected = (serviceId) => {
    return formData.value.skills.includes(serviceId);
};

// Get service name by ID
const getSkillName = (serviceId) => {
    const service = services.value.find(s => s.id === serviceId);
    return service ? service.name : '未知服务';
};

const openSkillsModal = () => {
    skillsModalRef.value?.showModal();
};

const closeSkillsModal = () => {
    skillsModalRef.value?.close();
};
</script>

<template>
    <div>
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
            <div>
                <h1 class="text-2xl font-bold tracking-tight text-base-content">
                    技师管理
                </h1>
                <p class="mt-1 text-base-content/60">
                    管理店内技师团队，查看实时状态与技能分布。
                </p>
            </div>
            <button v-if="activeTab === 'overview' && canManageTechnicians" @click="openCreateModal"
                class="btn btn-primary">
                <Plus class="w-4 h-4 mr-1" />
                添加技师
            </button>
        </div>

        <!-- Tabs -->
        <div role="tablist" class="tabs tabs-bordered mb-6">
            <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'overview' }"
                @click="activeTab = 'overview'">技师总览</a>
            <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'schedule' }"
                @click="activeTab = 'schedule'">排班管理</a>
        </div>

        <div v-if="activeTab === 'overview'">
            <!-- Loading State -->
            <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                <div v-for="i in 4" :key="i" class="h-64 rounded-xl border border-base-300 bg-base-200 animate-pulse">
                </div>
            </div>

            <!-- Empty State -->
            <div v-else-if="technicians.length === 0"
                class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50 text-center">
                <div class="p-4 bg-base-300 rounded-full mb-4">
                    <Users class="w-8 h-8 text-base-content/40" />
                </div>
                <h3 class="text-lg font-medium text-base-content">暂无技师</h3>
                <p class="text-base-content/60 mt-1">
                    点击右上角按钮添加第一位技师。
                </p>
            </div>

            <!-- Technicians Grid -->
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                <div v-for="tech in technicians" :key="tech.id"
                    class="group relative flex flex-col bg-base-100 border border-base-300 rounded-xl p-6 hover:border-primary/50 transition-all duration-200 hover:shadow-sm">
                    <!-- Status Badge -->
                    <div class="absolute top-4 right-4">
                        <span class="badge badge-sm" :class="getStatusInfo(tech.status).class">
                            {{ getStatusInfo(tech.status).text }}
                        </span>
                    </div>

                    <!-- Avatar & Info -->
                    <div class="flex flex-col items-center text-center mb-4">
                        <div class="mb-4 ring-4 ring-base-100 shadow-sm rounded-full">
                            <Avatar :name="tech.name" size="xl" />
                        </div>
                        <h3 class="text-lg font-semibold text-base-content">
                            {{ tech.name }}
                        </h3>

                        <!-- Rating -->
                        <div class="flex items-center gap-1 mt-1 text-warning text-sm font-medium">
                            <Star class="w-4 h-4 fill-current" />
                            <span>{{
                                tech.average_rating || tech.AverageRating || 5.0
                                }}</span>
                        </div>
                    </div>

                    <!-- Skills -->
                    <div class="flex-1">
                        <div class="flex flex-wrap gap-2 justify-center">
                            <span v-for="(skill, idx) in tech.skill_names" :key="idx"
                                class="badge badge-outline text-xs">
                                {{ skill }}
                            </span>
                        </div>
                    </div>

                    <!-- Actions -->
                    <div class="mt-6 pt-4 border-t border-base-200 flex gap-2">
                        <button @click="handleSchedule(tech)" class="btn btn-outline flex-1 btn-sm">
                            查看预约
                        </button>
                        <div class="dropdown dropdown-top dropdown-end">
                            <div tabindex="0" role="button" class="btn btn-square btn-outline btn-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                    stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                                    <path stroke-linecap="round" stroke-linejoin="round"
                                        d="M12 6.75a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5ZM12 12.75a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5ZM12 18.75a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5Z" />
                                </svg>
                            </div>
                            <ul tabindex="0"
                                class="dropdown-content z-1 menu p-2 shadow bg-base-100 rounded-box w-32 border border-base-200">
                                <li v-if="canManageTechnicians"><a @click="handleEdit(tech)">编辑</a></li>
                                <li v-if="canManageTechnicians"><a @click="handleDelete(tech)" class="text-error">删除</a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-else-if="activeTab === 'schedule'">
            <TechnicianSchedule :selected-technician="selectedTechnician" />
        </div>

        <!-- Create/Edit Modal -->
        <dialog ref="createModalRef" class="modal">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <h3 class="font-semibold text-lg text-base-content flex items-center gap-2">
                        <UserPlus v-if="!editingId" class="w-5 h-5 text-primary" />
                        <User v-else class="w-5 h-5 text-primary" />
                        {{ editingId ? "编辑技师" : "添加新技师" }}
                    </h3>
                    <button @click="closeCreateModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
                        <X class="w-5 h-5" />
                    </button>
                </div>

                <!-- Modal Body -->
                <div class="p-6">
                    <form @submit.prevent="handleSubmit" class="space-y-5">
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1">姓名</label>
                            <input type="text" v-model="formData.name" placeholder="请输入技师姓名"
                                class="input input-bordered w-full bg-base-100" required />
                        </div>

                        <!-- Skills Selection -->
                        <div class="form-control">
                            <label class="label">
                                <span class="label-text font-medium">
                                    掌握技能
                                    <span class="text-base-content/40 font-normal">(可多选)</span>
                                </span>
                            </label>

                            <!-- Selected Skills Tags -->
                            <div class="flex flex-wrap gap-2 mb-2" v-if="formData.skills.length > 0">
                                <div v-for="serviceId in formData.skills" :key="serviceId"
                                    class="badge badge-primary gap-1 cursor-pointer hover:opacity-80">
                                    {{ getSkillName(serviceId) }}
                                    <button type="button" @click="toggleSkill(serviceId)"
                                        class="btn btn-xs btn-circle btn-ghost">
                                        ✕
                                    </button>
                                </div>
                            </div>

                            <!-- Add Skills Button -->
                            <button type="button" @click="openSkillsModal" class="btn btn-outline w-full">
                                <Plus class="w-4 h-4" />
                                选择技能
                                <span class="text-xs text-base-content/60 ml-2" v-if="formData.skills.length > 0">
                                    (已选 {{ formData.skills.length }} 项)
                                </span>
                            </button>

                            <label class="label" v-if="services.length === 0">
                                <span class="label-text-alt text-base-content/60">
                                    请先在"服务管理"中添加服务项目
                                </span>
                            </label>
                        </div>

                        <div v-if="editingId">
                            <label class="block text-sm font-medium text-base-content/80 mb-1">状态</label>
                            <select v-model="formData.status" class="select select-bordered w-full bg-base-100">
                                <option :value="0">空闲</option>
                                <option :value="1">忙碌</option>
                                <option :value="2">请假</option>
                            </select>
                        </div>

                        <div class="pt-2">
                            <button type="submit" class="btn btn-primary w-full" :disabled="submitting">
                                <span v-if="submitting" class="loading loading-spinner loading-xs mr-2"></span>
                                {{
                                    submitting
                                        ? "保存中..."
                                        : editingId
                                            ? "确认修改"
                                            : "确认添加"
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

        <!-- Skills Selection Modal -->
        <dialog ref="skillsModalRef" class="modal">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-4xl">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <div>
                        <h3 class="font-semibold text-lg text-base-content flex items-center gap-2">
                            <Briefcase class="w-5 h-5 text-primary" />
                            选择技师技能
                        </h3>
                        <p class="text-sm text-base-content/60 mt-1">
                            已选择 {{ formData.skills.length }} 项服务
                        </p>
                    </div>
                    <button @click="closeSkillsModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
                        <X class="w-5 h-5" />
                    </button>
                </div>

                <!-- Modal Body -->
                <div class="p-6">
                    <!-- Empty State -->
                    <div v-if="services.length === 0" class="text-center py-12 text-base-content/60">
                        <Briefcase class="w-16 h-16 mx-auto mb-4 opacity-30" />
                        <p class="text-lg font-medium">暂无服务项目</p>
                        <p class="text-sm mt-1">请先在"服务管理"中添加服务项目</p>
                    </div>

                    <!-- Skills Grid -->
                    <div v-else
                        class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 max-h-96 overflow-y-auto pr-2">
                        <div v-for="service in services" :key="service.id"
                            class="card bg-base-100 border-2 border-base-300 hover:border-primary transition-all cursor-pointer"
                            :class="{ 'border-primary bg-primary/5': isSkillSelected(service.id) }"
                            @click="toggleSkill(service.id)">
                            <div class="card-body p-4">
                                <div class="flex items-start gap-3">
                                    <input type="checkbox" :checked="isSkillSelected(service.id)"
                                        @click.stop="toggleSkill(service.id)"
                                        class="checkbox checkbox-primary mt-1 checkbox-sm" />
                                    <div class="flex-1">
                                        <div class="font-semibold text-base-content text-sm">
                                            {{ service.name }}
                                        </div>
                                        <div class="text-xs text-base-content/60 mt-1">
                                            ¥{{ service.price }} · {{ service.duration }}分钟
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Modal Footer -->
                <div class="modal-action px-6 py-4 border-t border-base-200 bg-base-200/30 flex justify-between">
                    <button @click="closeSkillsModal" class="btn btn-ghost">
                        取消
                    </button>
                    <button @click="closeSkillsModal" class="btn btn-primary">
                        确定 (已选 {{ formData.skills.length }} 项)
                    </button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
                <button @click="closeSkillsModal">close</button>
            </form>
        </dialog>

        <!-- Appointment Modal -->
        <dialog ref="appointmentModalRef" class="modal">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-2xl">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <div>
                        <h3 class="font-semibold text-lg text-base-content flex items-center gap-2">
                            <Calendar class="w-5 h-5 text-primary" />
                            {{ selectedAppointmentTech?.name }} 的预约安排
                        </h3>
                        <p class="text-sm text-base-content/60 mt-1">
                            查看指定日期的预约情况
                        </p>
                    </div>
                    <button @click="closeAppointmentModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
                        <X class="w-5 h-5" />
                    </button>
                </div>

                <!-- Modal Body -->
                <div class="p-6">
                    <!-- Date Picker -->
                    <div class="mb-6">
                        <label class="block text-sm font-medium text-base-content/80 mb-2">选择日期</label>
                        <div class="relative">
                            <input type="text" :value="formatDisplayDate(selectedAppointmentDate)" readonly
                                class="input input-bordered w-full bg-base-100 cursor-pointer pl-10"
                                placeholder="点击选择日期" @click="calendarOpen = !calendarOpen" />
                            <Calendar class="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />

                            <calendar-date v-if="calendarOpen"
                                class="cally absolute top-full mt-2 z-10 bg-base-100 border border-base-300 shadow-lg rounded-box"
                                :value="selectedAppointmentDate" @select="handleDateSelect" @change="handleDateChange"
                                locale="zh-CN">
                                <ChevronLeft slot="previous" class="w-4 h-4" />
                                <ChevronRight slot="next" class="w-4 h-4" />
                                <calendar-month></calendar-month>
                            </calendar-date>
                        </div>
                    </div>

                    <!-- Loading State -->
                    <div v-if="appointmentModalLoading" class="flex justify-center py-12">
                        <span class="loading loading-spinner loading-lg text-primary"></span>
                    </div>

                    <!-- Appointments List -->
                    <div v-else>
                        <div v-if="technicianAppointments.length === 0" class="text-center py-12 text-base-content/60">
                            <Calendar class="w-16 h-16 mx-auto mb-4 opacity-30" />
                            <p class="text-lg font-medium">暂无预约</p>
                            <p class="text-sm mt-1">{{ selectedAppointmentDate }} 当天没有预约</p>
                        </div>

                        <div v-else class="space-y-4">
                            <h4
                                class="font-semibold flex items-center gap-2 mb-4 text-sm uppercase tracking-wider text-base-content/60">
                                预约列表 ({{ technicianAppointments.length }})
                            </h4>

                            <div class="space-y-3 max-h-96 overflow-y-auto">
                                <div v-for="appt in technicianAppointments" :key="appt.id"
                                    class="p-4 border border-base-200 rounded-lg hover:border-primary/50 transition-colors bg-base-50/50">
                                    <div class="flex justify-between items-start mb-3">
                                        <div>
                                            <div class="font-semibold text-base flex items-center gap-2">
                                                {{ appt.member?.name || "未知客户" }}
                                            </div>
                                            <div class="text-sm text-base-content/60 mt-0.5">
                                                {{ appt.service_item?.name || "未知服务" }}
                                            </div>
                                        </div>
                                        <span class="badge badge-sm" :class="{
                                            'badge-warning': appt.status === 'pending',
                                            'badge-success': appt.status === 'completed',
                                            'badge-info': appt.status === 'waitlist' || appt.status === 'waiting',
                                            'badge-error': appt.status === 'cancelled',
                                        }">
                                            {{
                                                appt.status === "pending" ? "待服务" :
                                                    appt.status === "completed" ? "已完成" :
                                                        appt.status === "waitlist" || appt.status === "waiting" ? "候补中" :
                                                            appt.status === "cancelled" ? "已取消" : appt.status
                                            }}
                                        </span>
                                    </div>

                                    <div class="flex items-center gap-4 text-sm text-base-content/70">
                                        <div class="flex items-center gap-1.5">
                                            <Clock class="w-3.5 h-3.5" />
                                            {{ appt.start_time.substring(11, 16) }} - {{ appt.end_time.substring(11, 16)
                                            }}
                                        </div>
                                        <div class="flex items-center gap-1.5">
                                            <Banknote class="w-3.5 h-3.5" />
                                            ¥{{ appt.actual_price }}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Modal Footer -->
                <div class="modal-action px-6 py-4 border-t border-base-200 bg-base-200/30">
                    <button @click="closeAppointmentModal" class="btn btn-neutral">
                        关闭
                    </button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
                <button @click="closeAppointmentModal">close</button>
            </form>
        </dialog>
    </div>
</template>