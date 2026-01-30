<script setup>
import { ref, onMounted, watch } from "vue";
import {
    getAppointments,
    cancelAppointment,
    completeAppointment,
} from "../api/appointments";
import AppointmentWizard from "../components/AppointmentWizard.vue";
import Avatar from "../components/Avatar.vue";
import {
    Plus,
    Calendar,
    Clock,
    CheckCircle2,
    XCircle,
    AlertCircle,
    Wallet,
    CircleCheckBig,
    CircleX
} from 'lucide-vue-next';

const appointments = ref([]);
const loading = ref(true);
const showModal = ref(false);
const filterStatus = ref("");

// Payment Modal State
const paymentModalRef = ref(null);
const currentPaymentAppt = ref(null);
const paymentMethod = ref("balance"); // balance, cash, mixed
const paymentBalance = ref(0);
const paymentCash = ref(0);

const paymentOptions = [
    { value: 'balance', label: '余额支付' },
    { value: 'cash', label: '现金/其他' },
    { value: 'mixed', label: '组合支付' }
];

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
        待服务: "badge badge-info badge-outline gap-1",
        完成: "badge badge-success badge-outline gap-1",
        候补: "badge badge-warning badge-outline gap-1",
        取消: "badge badge-neutral badge-outline gap-1",
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
    paymentModalRef.value?.showModal();
};

const closePaymentModal = () => {
    paymentModalRef.value?.close();
    currentPaymentAppt.value = null;
}

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

// Round to 2 decimal places to avoid floating point errors
const roundMoney = (val) => Math.round(val * 100) / 100;

const onBalanceInput = () => {
    if (paymentMethod.value !== "mixed") return;
    const price = currentPaymentAppt.value?.actual_price || currentPaymentAppt.value?.ActualPrice || 0;

    if (paymentBalance.value > price) paymentBalance.value = price;
    if (paymentBalance.value < 0) paymentBalance.value = 0;

    paymentCash.value = roundMoney(price - paymentBalance.value);
};

const onCashInput = () => {
    if (paymentMethod.value !== "mixed") return;
    const price = currentPaymentAppt.value?.actual_price || currentPaymentAppt.value?.ActualPrice || 0;

    if (paymentCash.value > price) paymentCash.value = price;
    if (paymentCash.value < 0) paymentCash.value = 0;

    paymentBalance.value = roundMoney(price - paymentCash.value);
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
        closePaymentModal();
        await fetchData();
    } catch (error) {
        console.error(error);
        alert("操作失败: " + (error.response?.data?.msg || error.message));
    }
};

// Get names for display (handling potential missing data)
const getServiceName = (service) =>
    service?.name || `项目#${service?.id || "未知"}`;
const getMemberName = (member) =>
    member?.name || `会员#${member?.id || "未知"}`;

</script>

<template>
    <div>
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
            <div>
                <h1 class="text-2xl font-bold text-base-content">
                    预约管理
                </h1>
                <p class="mt-1 text-base-content/60">
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
                    <Plus class="w-4 h-4 mr-1" />
                    新建预约
                </button>
            </div>
        </div>

        <!-- Appointments Table -->
        <div class="card bg-base-100 border border-base-300 shadow-sm rounded-xl overflow-hidden">
            <div class="card-body p-0">
                <div v-if="loading" class="flex justify-center py-16">
                    <span class="loading loading-spinner loading-lg text-primary"></span>
                </div>

                <div v-else-if="appointments.length === 0"
                    class="flex flex-col items-center justify-center py-16 text-center">
                    <div class="w-16 h-16 bg-base-200 rounded-full flex items-center justify-center mb-4">
                        <Calendar class="w-8 h-8 text-base-content/40" />
                    </div>
                    <h3 class="text-lg font-bold text-base-content">暂无预约记录</h3>
                    <p class="text-base-content/60 mt-1 max-w-sm">
                        当前系统中还没有任何预约。点击右上角的"新建预约"按钮来添加。
                    </p>
                </div>

                <div v-else class="overflow-x-auto">
                    <table class="table w-full">
                        <thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
                            <tr>
                                <th class="font-medium">ID</th>
                                <th class="font-medium">会员</th>
                                <th class="font-medium">技师</th>
                                <th class="font-medium">服务项目</th>
                                <th class="font-medium">时间段</th>
                                <th class="font-medium">状态</th>
                                <th class="font-medium">价格</th>
                                <th class="font-medium text-right pr-6">操作</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-base-200">
                            <tr v-for="appt in appointments" :key="appt.id"
                                class="hover:bg-base-50/50 transition-colors">
                                <td class="px-6 py-4 text-base-content/50 font-mono text-xs">
                                    #{{ appt.id }}
                                </td>
                                <td class="font-medium text-base-content">
                                    {{ getMemberName(appt.member) }}
                                </td>
                                <td>
                                    <div class="flex items-center gap-2">
                                        <Avatar :name="appt.technician.name" :src="appt.technician.avatar_url"
                                            class="w-2 text-xl" />
                                        <span class="text-base-content/80 text-sm">{{ appt.technician.name }}</span>
                                    </div>
                                </td>
                                <td>
                                    <div class="badge badge-ghost text-xs">
                                        {{ getServiceName(appt.service_item) }}
                                    </div>
                                </td>
                                <td class="px-6 py-4 text-base-content/60 text-xs">
                                    <div class="flex items-center gap-1.5 font-medium text-base-content/80">
                                        <Calendar class="w-3 h-3" />
                                        {{ formatDate(appt.start_time || appt.StartTime).split(" ")[0] }}
                                    </div>
                                    <div class="flex items-center gap-1.5 text-base-content/50 mt-1">
                                        <Clock class="w-3 h-3" />
                                        {{ formatDate(appt.start_time || appt.StartTime).split(" ")[1] }}
                                        -
                                        {{ formatDate(appt.end_time || appt.EndTime).split(" ")[1] }}
                                    </div>
                                </td>
                                <td>
                                    <span :class="getStatusBadge(appt.status || appt.Status)">
                                        <CheckCircle2 v-if="(appt.status || appt.Status) === 'completed'"
                                            class="w-3 h-3" />
                                        <Clock v-else-if="(appt.status || appt.Status) === 'pending'" class="w-3 h-3" />
                                        <AlertCircle v-else-if="(appt.status || appt.Status) === 'waiting'"
                                            class="w-3 h-3" />
                                        <XCircle v-else-if="(appt.status || appt.Status) === 'cancelled'"
                                            class="w-3 h-3" />
                                        {{ getStatusText(appt.status || appt.Status) }}
                                    </span>
                                </td>
                                <td class="font-mono text-sm">
                                    ¥{{ appt.actual_price || appt.ActualPrice }}
                                </td>
                                <td class="text-right pr-6">
                                    <div class="flex items-center justify-end gap-2">
                                        <button v-if="(appt.status || appt.Status) === 'pending'"
                                            @click="handleComplete(appt)" title="完成"
                                            class="btn btn-ghost btn-xs p-1 hover:bg-success/10">
                                            <CircleCheckBig class="w-4 h-4 text-success" />
                                        </button>
                                        <button v-if="['pending', 'waiting'].includes(appt.status || appt.Status)"
                                            @click="handleCancel(appt.id)" title="取消"
                                            class="btn btn-ghost btn-xs p-1 hover:bg-error/10">
                                            <CircleX class="w-4 h-4 text-error" />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <!-- Pagination or count could go here -->
                <div class="bg-base-50 px-6 py-3 border-t border-base-200 text-xs text-base-content/60 flex justify-between items-center"
                    v-if="appointments.length > 0">
                    <span>共 {{ appointments.length }} 条记录</span>
                </div>
            </div>
        </div>

        <!-- Appointment Wizard Modal -->
        <AppointmentWizard :show="showModal" @close="showModal = false" @success="handleAppointmentCreated" />

        <!-- Payment Modal -->
        <dialog ref="paymentModalRef" class="modal">
            <div class="modal-box">
                <h3 class="font-bold text-lg flex items-center gap-2">
                    <Wallet class="w-5 h-5 text-primary" />
                    订单结算
                </h3>
                <div class="py-6 space-y-6" v-if="currentPaymentAppt">
                    <!-- Info -->
                    <div class="flex justify-between items-center bg-base-200/50 p-4 rounded-xl border border-base-200">
                        <span class="text-base-content/70">订单金额</span>
                        <span class="text-2xl font-bold font-mono">¥{{ currentPaymentAppt.actual_price ||
                            currentPaymentAppt.ActualPrice }}</span>
                    </div>
                    <div class="text-sm flex justify-between px-1">
                        <span class="text-base-content/60">会员当前余额</span>
                        <span class="font-bold text-primary font-mono">¥{{ currentPaymentAppt.member?.balance ||
                            currentPaymentAppt.member?.Balance || 0 }}</span>
                    </div>
                    <!-- Payment Method -->
                    <div class="form-control">
                        <div class="join grid grid-cols-3 w-full">
                            <input v-for="opt in paymentOptions" :key="opt.value" type="radio" name="payment"
                                :value="opt.value" v-model="paymentMethod" class="join-item btn btn-outline"
                                :aria-label="opt.label" />
                        </div>
                    </div>

                    <!-- Amount Inputs -->
                    <div class="grid grid-cols-2 gap-4">
                        <label class="form-control">
                            <div class="label"><span class="label-text text-xs uppercase font-bold text-base-content/50">余额扣除</span></div>
                            <div class="input input-bordered flex items-center gap-2 px-3" :class="{ 'input-disabled': paymentMethod === 'cash' }">
                                <span class="text-base-content/40">¥</span>
                                <input type="number" v-model.number="paymentBalance" @input="onBalanceInput"
                                    :disabled="paymentMethod === 'cash'"
                                    class="w-full bg-transparent font-mono outline-none" step="0.01" min="0" />
                            </div>
                        </label>
                        <label class="form-control">
                            <div class="label"><span class="label-text text-xs uppercase font-bold text-base-content/50">现金支付</span></div>
                            <div class="input input-bordered flex items-center gap-2 px-3" :class="{ 'input-disabled': paymentMethod === 'balance' }">
                                <span class="text-base-content/40">¥</span>
                                <input type="number" v-model.number="paymentCash" @input="onCashInput"
                                    :disabled="paymentMethod === 'balance'"
                                    class="w-full bg-transparent font-mono outline-none" step="0.01" min="0" />
                            </div>
                        </label>
                    </div>

                    <div class="alert alert-warning text-xs py-2 shadow-sm"
                        v-if="(currentPaymentAppt.member?.balance || currentPaymentAppt.member?.Balance || 0) < (currentPaymentAppt.actual_price || currentPaymentAppt.ActualPrice) && paymentMethod === 'balance'">
                        <AlertCircle class="w-4 h-4" />
                        <span>余额不足以全额支付，请选择组合支付或现金支付</span>
                    </div>
                </div>
                <div class="modal-action">
                    <button class="btn" @click="closePaymentModal">取消</button>
                    <button class="btn btn-primary" @click="confirmPayment">确认支付</button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="closePaymentModal">close</button>
            </form>
        </dialog>
    </div>
</template>
