<script setup>
import { ref, onMounted } from "vue";
import {
    getTechnicians,
    createTechnician,
    updateTechnician,
    deleteTechnician,
} from "../api/technicians";
import { getAppointments } from "../api/appointments";
import TechnicianSchedule from "../components/TechnicianSchedule.vue";
import Avatar from "../components/Avatar.vue";
import { usePermission } from "../composables/usePermission";
import "cally";

const { canManageTechnicians } = usePermission();

const activeTab = ref("overview");
const technicians = ref([]);
const loading = ref(true);
const showModal = ref(false);
const submitting = ref(false);
const editingId = ref(null);

// Appointment Modal
const showAppointmentModal = ref(false);
const appointmentModalLoading = ref(false);
const selectedAppointmentTech = ref(null);
const selectedAppointmentDate = ref(new Date().toISOString().split('T')[0]);
const calendarOpen = ref(false);
const technicianAppointments = ref([]);

const formData = ref({
    name: "",
    skills: "", // Comma separated string for input
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

onMounted(fetchTechnicians);

const openCreateModal = () => {
    editingId.value = null;
    formData.value = { name: "", skills: "", status: 0 };
    showModal.value = true;
};

const handleEdit = (tech) => {
    editingId.value = tech.id;
    const skills = parseSkills(tech.skills || tech.Skills);
    formData.value = {
        name: tech.name,
        skills: Array.isArray(skills) ? skills.join(", ") : "",
        status: tech.status,
    };
    showModal.value = true;
};

const handleSchedule = (tech) => {
    selectedAppointmentTech.value = tech;
    selectedAppointmentDate.value = new Date().toISOString().split('T')[0];
    showAppointmentModal.value = true;
    fetchTechnicianAppointments();
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
        calendarOpen.value = false; // Close calendar after selection
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
        // Convert comma-separated skills string to array
        const skillsArray = formData.value.skills
            .split(/[,，]/) // Split by comma (English or Chinese)
            .map((s) => s.trim())
            .filter((s) => s);

        const payload = {
            name: formData.value.name,
            skills: JSON.stringify(skillsArray),
            status: Number(formData.value.status),
        };

        if (editingId.value) {
            await updateTechnician(editingId.value, payload);
        } else {
            await createTechnician(payload);
        }

        showModal.value = false;
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

// Helper to parse skills
const parseSkills = (skills) => {
    if (Array.isArray(skills)) return skills;
    if (typeof skills === "string") {
        try {
            return JSON.parse(skills);
        } catch (e) {
            return [skills];
        }
    }
    return [];
};
</script>

<template>
    <div class="max-w-7xl mx-auto">
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between mb-10 gap-4">
            <div>
                <h1 class="text-3xl font-bold tracking-tight text-base-content">
                    技师管理
                </h1>
                <p class="mt-2 text-base-content/60">
                    管理店内技师团队，查看实时状态与技能分布。
                </p>
            </div>
            <button v-if="activeTab === 'overview' && canManageTechnicians" @click="openCreateModal"
                class="btn btn-neutral">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                    stroke="currentColor" class="w-4 h-4 mr-2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
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
                class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50">
                <div class="p-4 bg-base-300 rounded-full mb-4">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                        stroke="currentColor" class="w-8 h-8 text-base-content/40">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
                    </svg>
                </div>
                <h3 class="text-lg font-medium text-base-content">暂无技师</h3>
                <p class="text-base-content/60 mt-1">
                    点击右上角按钮添加第一位技师。
                </p>
            </div>

            <!-- Technicians Grid -->
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                <div v-for="tech in technicians" :key="tech.id"
                    class="group relative flex flex-col bg-base-100 border border-base-300 rounded-xl p-6 hover:border-base-content/20 transition-all duration-200 hover:shadow-sm">
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
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
                                class="w-4 h-4">
                                <path fill-rule="evenodd"
                                    d="M10.788 3.21c.448-1.077 1.976-1.077 2.424 0l2.082 5.007 5.404.433c1.164.093 1.636 1.545.749 2.305l-4.117 3.527 1.257 5.273c.271 1.136-.964 2.033-1.96 1.425L12 18.354 7.373 21.18c-.996.608-2.231-.29-1.96-1.425l1.257-5.273-4.117-3.527c-.887-.76-.415-2.212.749-2.305l5.404-.433 2.082-5.006z"
                                    clip-rule="evenodd" />
                            </svg>
                            <span>{{
                                tech.average_rating || tech.AverageRating || 5.0
                                }}</span>
                        </div>
                    </div>

                    <!-- Skills -->
                    <div class="flex-1">
                        <div class="flex flex-wrap gap-2 justify-center">
                            <span v-for="(skill, idx) in parseSkills(
                                tech.skills || tech.Skills,
                            )" :key="idx" class="badge badge-outline text-xs">
                                {{ skill }}
                            </span>
                        </div>
                    </div>

                    <!-- Actions -->
                    <div class="mt-6 pt-4 border-t border-base-200 flex gap-2">
                        <button @click="handleSchedule(tech)" class="btn btn-outline">
                            查看预约
                        </button>
                        <button v-if="canManageTechnicians" @click="handleEdit(tech)" class="btn btn-outline">
                            编辑
                        </button>
                        <button v-if="canManageTechnicians" @click="handleDelete(tech)"
                            class="btn btn-error btn-outline">
                            删除
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div v-else-if="activeTab === 'schedule'">
            <TechnicianSchedule :selected-technician="selectedTechnician" />
        </div>

        <!-- Create Modal -->
        <dialog class="modal" :class="{ 'modal-open': showModal }">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <h3 class="font-semibold text-lg text-base-content">
                        {{ editingId ? "编辑技师" : "添加新技师" }}
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
                            <label class="block text-sm font-medium text-base-content/80 mb-1">姓名</label>
                            <input type="text" v-model="formData.name" placeholder="请输入技师姓名"
                                class="input input-bordered w-full bg-base-100" required />
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1">
                                技能标签
                                <span class="text-base-content/40 font-normal">(用逗号分隔)</span>
                            </label>
                            <textarea v-model="formData.skills" placeholder="例如: 精油SPA, 足疗, 中式推拿" rows="3"
                                class="textarea textarea-bordered w-full bg-base-100 resize-none"></textarea>
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
                            <button type="submit" class="btn btn-neutral w-full" :disabled="submitting">
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
                <button @click="showModal = false">close</button>
            </form>
        </dialog>

        <!-- Appointment Modal -->
        <dialog class="modal" :class="{ 'modal-open': showAppointmentModal }">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-2xl">
                <!-- Modal Header -->
                <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                    <div>
                        <h3 class="font-semibold text-lg text-base-content">
                            {{ selectedAppointmentTech?.name }} 的预约安排
                        </h3>
                        <p class="text-sm text-base-content/60 mt-1">
                            查看指定日期的预约情况
                        </p>
                    </div>
                    <button @click="showAppointmentModal = false"
                        class="btn btn-ghost btn-sm btn-square text-base-content/60">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                            stroke="currentColor" class="w-5 h-5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Modal Body -->
                <div class="p-6">
                    <!-- Date Picker -->
                    <div class="mb-6">
                        <label class="block text-sm font-medium text-base-content/80 mb-2">选择日期</label>
                        <div class="relative">
                            <input type="text" :value="formatDisplayDate(selectedAppointmentDate)" readonly
                                class="input input-bordered w-full bg-base-100 cursor-pointer" placeholder="点击选择日期"
                                @click="calendarOpen = !calendarOpen" />
                            <calendar-date v-if="calendarOpen"
                                class="cally absolute top-full mt-2 z-10 bg-base-100 border border-base-300 shadow-lg rounded-box"
                                :value="selectedAppointmentDate" @select="handleDateSelect" @change="handleDateChange"
                                locale="zh-CN">
                                <svg aria-label="Previous" class="fill-current size-4" slot="previous"
                                    xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                                    <path fill="currentColor" d="M15.75 19.5 8.25 12l7.5-7.5"></path>
                                </svg>
                                <svg aria-label="Next" class="fill-current size-4" slot="next"
                                    xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                                    <path fill="currentColor" d="m8.25 4.5 7.5 7.5-7.5 7.5"></path>
                                </svg>
                                <calendar-month></calendar-month>
                            </calendar-date>
                        </div>
                    </div>

                    <!-- Loading State -->
                    <div v-if="appointmentModalLoading" class="flex justify-center py-12">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>

                    <!-- Appointments List -->
                    <div v-else>
                        <div v-if="technicianAppointments.length === 0" class="text-center py-12 text-base-content/60">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-16 h-16 mx-auto mb-4 opacity-30">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5" />
                            </svg>
                            <p class="text-lg font-medium">暂无预约</p>
                            <p class="text-sm mt-1">{{ selectedAppointmentDate }} 当天没有预约</p>
                        </div>

                        <div v-else class="space-y-4">
                            <h4 class="font-semibold flex items-center gap-2 mb-4">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                    stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                                    <path stroke-linecap="round" stroke-linejoin="round"
                                        d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5" />
                                </svg>
                                预约列表 ({{ technicianAppointments.length }})
                            </h4>

                            <div class="space-y-3 max-h-96 overflow-y-auto">
                                <div v-for="appt in technicianAppointments" :key="appt.id"
                                    class="p-4 border border-base-200 rounded-lg hover:border-primary/50 transition-colors">
                                    <div class="flex justify-between items-start mb-3">
                                        <div>
                                            <div class="font-semibold text-base">
                                                {{ appt.member?.name || "未知客户" }}
                                            </div>
                                            <div class="text-sm text-base-content/60">
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
                                        <div class="flex items-center gap-1">
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                                stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                                                <path stroke-linecap="round" stroke-linejoin="round"
                                                    d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                                            </svg>
                                            {{ appt.start_time.substring(11, 16) }} - {{ appt.end_time.substring(11, 16)
                                            }}
                                        </div>
                                        <div class="flex items-center gap-1">
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                                stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                                                <path stroke-linecap="round" stroke-linejoin="round"
                                                    d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                            </svg>
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
                    <button @click="showAppointmentModal = false" class="btn btn-neutral">
                        关闭
                    </button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
                <button @click="showAppointmentModal = false">close</button>
            </form>
        </dialog>
    </div>
</template>