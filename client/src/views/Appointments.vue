<script setup>
import { ref, onMounted } from "vue";
import {
    getAppointments,
    cancelAppointment,
    completeAppointment,
} from "../api/appointments";
import AppointmentWizard from "../components/AppointmentWizard.vue";
import Avatar from "../components/Avatar.vue";

const appointments = ref([]);
const loading = ref(true);
const showModal = ref(false);
const filterStatus = ref("");

// Fetch appointments data
const fetchData = async () => {
    loading.value = true;
    try {
        const data = await getAppointments({
            status: filterStatus.value || undefined,
        });
        appointments.value = data || [];
    } catch (error) {
        console.error("Error fetching data:", error);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchData);

// Helper to format date
const formatDate = (dateStr) => {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleString("zh-CN", {
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
    });
};

// Helper for status badge
const getStatusBadge = (status) => {
    const map = {
        待服务: "badge badge-info badge-outline",
        完成: "badge badge-success badge-outline",
        候补: "badge badge-warning badge-outline",
        取消: "badge badge-neutral badge-outline",
    };

    // Handle numeric or different string inputs if backend changes
    if (status === 0 || status === "0" || status === "pending")
        return map["待服务"];
    if (status === 1 || status === "1" || status === "completed")
        return map["完成"];
    if (status === 2 || status === "2" || status === "waiting")
        return map["候补"];
    if (status === 3 || status === "3" || status === "cancelled")
        return map["取消"];

    return map[status] || "badge badge-ghost";
};

const getStatusText = (status) => {
    const map = {
        0: "待服务",
        1: "完成",
        2: "候补",
        3: "取消",
        pending: "待服务",
        completed: "完成",
        waiting: "候补",
        cancelled: "取消",
    };
    return map[status] || status;
};

// Handle appointment created
const handleAppointmentCreated = async () => {
    await fetchData();
};

const handleCancel = async (id) => {
    if (!confirm("确定要取消该预约吗？")) return;
    try {
        await cancelAppointment(id);
        alert("预约已取消");
        await fetchData();
    } catch (error) {
        alert("取消失败: " + (error.message || "未知错误"));
    }
};

const handleComplete = async (id) => {
    if (!confirm("确定要完成该订单吗？这将结算费用并计算佣金。")) return;
    try {
        await completeAppointment(id);
        alert("订单已完成");
        await fetchData();
    } catch (error) {
        alert("操作失败: " + (error.message || "未知错误"));
    }
};

// Get names for display (handling potential missing data)
const getTechName = (tech) => tech?.name || `技师#${tech?.id || "未知"}`;
const getServiceName = (service) =>
    service?.name || `项目#${service?.id || "未知"}`;
const getMemberName = (member) =>
    member?.name || `会员#${member?.id || "未知"}`;

const showDetails = (appt) => {
    const info = [
        `订单 ID: ${appt.id}`,
        `会员: ${getMemberName(appt.member)}`,
        `技师: ${getTechName(appt.technician)}`,
        `项目: ${getServiceName(appt.service_item)}`,
        `状态: ${getStatusText(appt.status || appt.Status)}`,
        `实付: ¥${appt.actual_price || appt.ActualPrice}`,
    ].join("\n");
    alert(info);
};
</script>

<template>
    <div class="max-w-7xl mx-auto">
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between mb-10 gap-4">
            <div>
                <h1 class="text-3xl font-bold tracking-tight text-base-content">
                    预约管理
                </h1>
                <p class="mt-2 text-base-content/60">
                    查看所有预约订单，处理候补队列与调度冲突。
                </p>
            </div>
            <div class="flex items-center gap-3">
                <select v-model="filterStatus" @change="fetchData"
                    class="select select-bordered select-sm w-full max-w-xs">
                    <option value="">所有状态</option>
                    <option value="pending">待服务</option>
                    <option value="waiting">候补中</option>
                    <option value="completed">已完成</option>
                    <option value="cancelled">已取消</option>
                </select>
                <button @click="showModal = true" class="btn btn-primary btn-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                        stroke="currentColor" class="w-4 h-4 mr-2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                    </svg>
                    新建预约
                </button>
            </div>
        </div>

        <!-- Appointments Table -->
        <div class="bg-base-100 rounded-box border border-base-200 shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
                <table class="table table-zebra w-full">
                    <thead class="bg-base-200 text-base-content/70 uppercase text-xs">
                        <tr>
                            <th class="px-6 py-3 font-medium">ID</th>
                            <th class="px-6 py-3 font-medium">会员</th>
                            <th class="px-6 py-3 font-medium">技师</th>
                            <th class="px-6 py-3 font-medium">服务项目</th>
                            <th class="px-6 py-3 font-medium">时间段</th>
                            <th class="px-6 py-3 font-medium">状态</th>
                            <th class="px-6 py-3 font-medium">价格</th>
                            <th class="px-6 py-3 font-medium text-right">
                                操作
                            </th>
                        </tr>
                    </thead>
                    <tbody class="text-sm">
                        <tr v-if="loading">
                            <td colspan="8" class="px-6 py-12 text-center">
                                <span class="loading loading-spinner loading-lg text-base-content/30"></span>
                            </td>
                        </tr>
                        <tr v-else-if="appointments.length === 0">
                            <td colspan="8" class="px-6 py-12 text-center text-base-content/50">
                                暂无预约记录
                            </td>
                        </tr>
                        <tr v-else v-for="appt in appointments" :key="appt.id"
                            class="hover:bg-base-200/50 transition-colors">
                            <td class="px-6 py-4 text-base-content/50 font-mono text-xs">
                                #{{ appt.id }}
                            </td>
                            <td class="px-6 py-4 font-medium text-base-content">
                                {{ getMemberName(appt.member) }}
                            </td>
                            <td class="px-6 py-4">
                                <div class="flex items-center gap-2">
                                    <Avatar :name="getTechName(appt.technician)" size="xs" />
                                    <span class="text-base-content/80">{{
                                        getTechName(appt.technician)
                                        }}</span>
                                </div>
                            </td>
                            <td class="px-6 py-4 text-base-content/80">
                                {{ getServiceName(appt.service_item) }}
                            </td>
                            <td class="px-6 py-4 text-base-content/60 text-xs">
                                <div>
                                    {{
                                        formatDate(
                                            appt.start_time || appt.StartTime,
                                        )
                                    }}
                                </div>
                                <div class="text-base-content/40 mt-0.5">
                                    至
                                    {{
                                        formatDate(
                                            appt.end_time || appt.EndTime,
                                        ).split(" ")[1]
                                    }}
                                </div>
                            </td>
                            <td class="px-6 py-4">
                                <span :class="getStatusBadge(
                                    appt.status || appt.Status,
                                )
                                    ">
                                    {{
                                        getStatusText(
                                            appt.status || appt.Status,
                                        )
                                    }}
                                </span>
                            </td>
                            <td class="px-6 py-4 font-medium text-base-content">
                                ¥{{ appt.actual_price || appt.ActualPrice }}
                            </td>
                            <td class="px-6 py-4 text-right flex justify-end gap-2">
                                <button @click="showDetails(appt)" class="btn btn-ghost btn-xs">
                                    详情
                                </button>
                                <button v-if="
                                    (appt.status || appt.Status) ===
                                    'pending'
                                " @click="handleComplete(appt.id)" class="btn btn-success btn-outline btn-xs">
                                    完成
                                </button>
                                <button v-if="
                                    ['pending', 'waiting'].includes(
                                        appt.status || appt.Status,
                                    )
                                " @click="handleCancel(appt.id)" class="btn btn-error btn-outline btn-xs">
                                    取消
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Appointment Wizard Modal -->
        <AppointmentWizard :show="showModal" @close="showModal = false" @success="handleAppointmentCreated" />
    </div>
</template>