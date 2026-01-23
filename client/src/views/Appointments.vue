<script setup>
import { ref, onMounted, watch } from "vue";
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

// Payment Modal State
const showPaymentModal = ref(false);
const currentPaymentAppt = ref(null);
const paymentMethod = ref("balance"); // balance, cash, mixed
const paymentBalance = ref(0);
const paymentCash = ref(0);

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

const handleComplete = (appt) => {
    currentPaymentAppt.value = appt;
    const price = appt.actual_price || appt.ActualPrice || 0;
    const memberBalance = appt.member?.balance || appt.member?.Balance || 0;

    // Default selection logic
    if (memberBalance >= price) {
        paymentMethod.value = "balance";
        paymentBalance.value = price;
        paymentCash.value = 0;
    } else {
        paymentMethod.value = "mixed";
        if (memberBalance > 0) {
            paymentBalance.value = memberBalance;
            paymentCash.value = price - memberBalance;
        } else {
            paymentMethod.value = "cash";
            paymentBalance.value = 0;
            paymentCash.value = price;
        }
    }
    showPaymentModal.value = true;
};

// Watch for payment method change to reset amounts
watch(paymentMethod, (newMethod) => {
    if (!currentPaymentAppt.value) return;
    const price = currentPaymentAppt.value.actual_price || currentPaymentAppt.value.ActualPrice || 0;
    const memberBalance = currentPaymentAppt.value.member?.balance || currentPaymentAppt.value.member?.Balance || 0;

    if (newMethod === "balance") {
        paymentBalance.value = price;
        paymentCash.value = 0;
    } else if (newMethod === "cash") {
        paymentBalance.value = 0;
        paymentCash.value = price;
    } else if (newMethod === "mixed") {
        // Default split for mixed: use max available balance
        if (memberBalance >= price) {
            paymentBalance.value = price;
            paymentCash.value = 0;
        } else {
            paymentBalance.value = memberBalance;
            paymentCash.value = price - memberBalance;
        }
    }
});

const onBalanceInput = () => {
    if (paymentMethod.value !== "mixed") return;
    const price = currentPaymentAppt.value?.actual_price || currentPaymentAppt.value?.ActualPrice || 0;

    // Ensure balance doesn't exceed price or member balance is handled by validation,
    // but here we auto-calc cash
    if (paymentBalance.value > price) paymentBalance.value = price;
    if (paymentBalance.value < 0) paymentBalance.value = 0;

    // Auto calculate cash
    paymentCash.value = price - paymentBalance.value;
};

const onCashInput = () => {
    if (paymentMethod.value !== "mixed") return;
    const price = currentPaymentAppt.value?.actual_price || currentPaymentAppt.value?.ActualPrice || 0;

    if (paymentCash.value > price) paymentCash.value = price;
    if (paymentCash.value < 0) paymentCash.value = 0;

    // Auto calculate balance
    paymentBalance.value = price - paymentCash.value;
};

const confirmPayment = async () => {
    if (!currentPaymentAppt.value) return;

    const price = currentPaymentAppt.value.actual_price || currentPaymentAppt.value.ActualPrice || 0;
    const memberBalance = currentPaymentAppt.value.member?.balance || currentPaymentAppt.value.member?.Balance || 0;

    // Validation
    if (Math.abs(Number(paymentBalance.value) + Number(paymentCash.value) - price) > 0.01) {
        alert(`支付总额必须等于订单金额 (¥${price})`);
        return;
    }

    if (Number(paymentBalance.value) > memberBalance) {
        alert(`余额不足，当前余额: ¥${memberBalance}`);
        return;
    }

    try {
        await completeAppointment(currentPaymentAppt.value.id, {
            payment_method: paymentMethod.value,
            balance_amount: Number(paymentBalance.value),
            cash_amount: Number(paymentCash.value)
        });
        alert("订单已完成");
        showPaymentModal.value = false;
        await fetchData();
    } catch (error) {
        console.error(error);
        alert("操作失败: " + (error.response?.data?.msg || error.message));
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
                <select v-model="filterStatus" @change="fetchData" class="select select-bordered w-36 shrink-0">
                    <option selected value="">所有状态</option>
                    <option value="pending">待服务</option>
                    <option value="waiting">候补中</option>
                    <option value="completed">已完成</option>
                    <option value="cancelled">已取消</option>
                </select>
                <button @click="showModal = true" class="btn btn-primary">
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
                <table class="table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>会员</th>
                            <th>技师</th>
                            <th>服务项目</th>
                            <th>时间段</th>
                            <th>状态</th>
                            <th>价格</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
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
                            <td class="font-medium text-base-content">
                                {{ getMemberName(appt.member) }}
                            </td>
                            <td>
                                <div class="flex items-center gap-2">
                                    <Avatar :name="getTechName(appt.technician)" size="xs" />
                                    <span class="text-base-content/80">{{
                                        getTechName(appt.technician)
                                    }}</span>
                                </div>
                            </td>
                            <td>
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
                            <td>
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
                            <td>
                                ¥{{ appt.actual_price || appt.ActualPrice }}
                            </td>
                            <td>
                                <button @click="showDetails(appt)" class="btn btn-ghost btn-xs">
                                    详情
                                </button>
                                <button v-if="
                                    (appt.status || appt.Status) ===
                                    'pending'
                                " @click="handleComplete(appt)" class="btn btn-success btn-outline btn-xs">
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

        <!-- Payment Modal -->
        <div v-if="showPaymentModal" class="modal modal-open">
            <div class="modal-box">
                <h3 class="font-bold text-lg">订单结算</h3>
                <div class="py-4 space-y-4" v-if="currentPaymentAppt">
                    <!-- Info -->
                    <div class="flex justify-between items-center bg-base-200 p-3 rounded-lg">
                        <span>订单金额:</span>
                        <span class="text-xl font-bold">¥{{ currentPaymentAppt.actual_price ||
                            currentPaymentAppt.ActualPrice }}</span>
                    </div>
                    <div class="text-sm text-base-content/70">
                        会员余额: <span class="font-bold text-primary">¥{{ currentPaymentAppt.member?.balance ||
                            currentPaymentAppt.member?.Balance || 0 }}</span>
                    </div>

                    <!-- Payment Method -->
                    <div class="form-control">
                        <label class="label"><span class="label-text font-medium">支付方式</span></label>
                        <div class="flex gap-4">
                            <label class="label cursor-pointer gap-2">
                                <input type="radio" name="payment" class="radio radio-primary" value="balance"
                                    v-model="paymentMethod" />
                                <span class="label-text">余额支付</span>
                            </label>
                            <label class="label cursor-pointer gap-2">
                                <input type="radio" name="payment" class="radio radio-primary" value="cash"
                                    v-model="paymentMethod" />
                                <span class="label-text">现金/其他</span>
                            </label>
                            <label class="label cursor-pointer gap-2">
                                <input type="radio" name="payment" class="radio radio-primary" value="mixed"
                                    v-model="paymentMethod" />
                                <span class="label-text">组合支付</span>
                            </label>
                        </div>
                    </div>

                    <!-- Amount Inputs -->
                    <div class="grid grid-cols-2 gap-4">
                        <div class="form-control">
                            <label class="label"><span class="label-text">余额扣除</span></label>
                            <input type="number" v-model.number="paymentBalance" @input="onBalanceInput"
                                :disabled="paymentMethod === 'cash'" class="input input-bordered w-full" step="0.01"
                                min="0" />
                        </div>
                        <div class="form-control">
                            <label class="label"><span class="label-text">现金支付</span></label>
                            <input type="number" v-model.number="paymentCash" @input="onCashInput"
                                :disabled="paymentMethod === 'balance'" class="input input-bordered w-full" step="0.01"
                                min="0" />
                        </div>
                    </div>

                    <div class="text-xs text-warning"
                        v-if="(currentPaymentAppt.member?.balance || currentPaymentAppt.member?.Balance || 0) < (currentPaymentAppt.actual_price || currentPaymentAppt.ActualPrice) && paymentMethod === 'balance'">
                        警告：余额不足以全额支付
                    </div>
                </div>
                <div class="modal-action">
                    <button class="btn" @click="showPaymentModal = false">取消</button>
                    <button class="btn btn-primary" @click="confirmPayment">确认支付</button>
                </div>
            </div>
        </div>
    </div>
</template>