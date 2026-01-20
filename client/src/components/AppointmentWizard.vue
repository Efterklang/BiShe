<script setup>
import { ref, computed, watch } from "vue";
import { getServices } from "../api/services";
import { getMembers } from "../api/members";
import { getAvailableTechnicians } from "../api/schedules";
import { createAppointment } from "../api/appointments";

const props = defineProps({
    show: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(["close", "success"]);

// --- State ---
const currentStep = ref(1); // 1: 服务+时间, 2: 选择技师, 3: 确认
const loading = ref(false);
const submitting = ref(false);

const services = ref([]);
const members = ref([]);
const availableTechs = ref([]);
const unavailableTechs = ref([]);
const selectedServiceInfo = ref(null); // Stores service info from API response

// Form data
const formData = ref({
    member_id: "",
    service_id: "",
    start_time: "",
    tech_id: "",
    allow_waitlist: false,
});

// --- Computed ---
const selectedService = computed(() => {
    // Use API response first, fallback to local services list
    if (selectedServiceInfo.value) {
        return selectedServiceInfo.value;
    }
    return services.value.find(
        (s) => s.id === Number(formData.value.service_id),
    );
});

const selectedMember = computed(() => {
    return members.value.find((m) => m.id === Number(formData.value.member_id));
});

const selectedTech = computed(() => {
    return [...availableTechs.value, ...unavailableTechs.value].find(
        (t) => t.id === Number(formData.value.tech_id),
    );
});

const canProceedStep1 = computed(() => {
    return (
        formData.value.member_id &&
        formData.value.service_id &&
        formData.value.start_time
    );
});

const canProceedStep2 = computed(() => {
    return formData.value.tech_id;
});

const endTime = computed(() => {
    if (!formData.value.start_time || !selectedService.value) return "";
    const start = new Date(formData.value.start_time);
    const end = new Date(
        start.getTime() + selectedService.value.duration * 60000,
    );
    return end.toISOString();
});

const formatTime = (isoString) => {
    if (!isoString) return "";
    const date = new Date(isoString);
    return date.toLocaleString("zh-CN", {
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
    });
};

// --- Lifecycle ---
const fetchInitialData = async () => {
    loading.value = true;
    try {
        const [servicesData, membersData] = await Promise.all([
            getServices(),
            getMembers(),
        ]);
        services.value = servicesData || [];
        members.value = membersData || [];
    } catch (error) {
        console.error("Failed to fetch initial data:", error);
        alert("加载数据失败，请刷新页面重试");
    } finally {
        loading.value = false;
    }
};

// --- Actions ---
const fetchAvailableTechnicians = async () => {
    if (!formData.value.start_time || !formData.value.service_id) return;

    loading.value = true;
    try {
        const startTime = new Date(formData.value.start_time).toISOString();
        const data = await getAvailableTechnicians({
            start_time: startTime,
            service_id: formData.value.service_id,
        });

        availableTechs.value = data.available || [];
        unavailableTechs.value = data.unavailable || [];

        // Save service info from API response for skill matching display
        if (data.service) {
            selectedServiceInfo.value = data.service;
        }

        if (availableTechs.value.length === 0) {
            // 没有可用技师，询问用户
            const choice = confirm(
                `该时间段没有可用技师。\n\n点击"确定"加入候补队列\n点击"取消"返回重新选择时间`,
            );
            if (choice) {
                formData.value.allow_waitlist = true;
                // 显示所有技师供用户选择
            } else {
                // 返回第一步
                currentStep.value = 1;
                return;
            }
        }
    } catch (error) {
        console.error("Failed to fetch available technicians:", error);
        alert("查询可用技师失败: " + (error.message || "未知错误"));
    } finally {
        loading.value = false;
    }
};

const goToStep2 = async () => {
    if (!canProceedStep1.value) return;
    await fetchAvailableTechnicians();
    currentStep.value = 2;
};

const goToStep3 = () => {
    if (!canProceedStep2.value) return;
    currentStep.value = 3;
};

const goBack = () => {
    if (currentStep.value > 1) {
        currentStep.value--;
        if (currentStep.value === 1) {
            // 重置技师选择
            formData.value.tech_id = "";
            formData.value.allow_waitlist = false;
            availableTechs.value = [];
            unavailableTechs.value = [];
            selectedServiceInfo.value = null;
        }
    }
};

const handleSubmit = async () => {
    submitting.value = true;
    try {
        const payload = {
            member_id: Number(formData.value.member_id),
            tech_id: Number(formData.value.tech_id),
            service_id: Number(formData.value.service_id),
            start_time: new Date(formData.value.start_time).toISOString(),
            allow_waitlist: formData.value.allow_waitlist,
        };

        await createAppointment(payload);

        alert(
            formData.value.allow_waitlist
                ? "预约已加入候补队列！"
                : "预约创建成功！",
        );
        emit("success");
        closeModal();
    } catch (error) {
        console.error("Failed to create appointment:", error);
        alert(
            "创建预约失败: " +
            (error.response?.data?.msg || error.message || "未知错误"),
        );
    } finally {
        submitting.value = false;
    }
};

const closeModal = () => {
    // Reset form
    currentStep.value = 1;
    formData.value = {
        member_id: "",
        service_id: "",
        start_time: "",
        tech_id: "",
        allow_waitlist: false,
    };
    availableTechs.value = [];
    unavailableTechs.value = [];
    selectedServiceInfo.value = null;
    emit("close");
};

// Watch for modal open
watch(
    () => props.show,
    (newVal) => {
        if (newVal) {
            fetchInitialData();
        }
    },
);

// 设置默认开始时间为当前时间的下一个整点
const setDefaultStartTime = () => {
    const now = new Date();
    now.setHours(now.getHours() + 1, 0, 0, 0);
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, "0");
    const day = String(now.getDate()).padStart(2, "0");
    const hours = String(now.getHours()).padStart(2, "0");
    const minutes = String(now.getMinutes()).padStart(2, "0");
    formData.value.start_time = `${year}-${month}-${day}T${hours}:${minutes}`;
};

watch(
    () => props.show,
    (newVal) => {
        if (newVal && !formData.value.start_time) {
            setDefaultStartTime();
        }
    },
);
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl max-w-2xl p-0 overflow-hidden">
            <!-- Modal Header -->
            <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
                <div>
                    <h3 class="font-bold text-lg text-base-content">
                        新建预约
                    </h3>
                    <div class="text-sm text-base-content/60 mt-1">
                        步骤 {{ currentStep }} / 3
                    </div>
                </div>
                <button @click="closeModal" class="btn btn-ghost btn-sm btn-square">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                        stroke="currentColor" class="w-5 h-5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <!-- Progress Steps -->
            <div class="px-6 pt-4">
                <ul class="steps steps-horizontal w-full">
                    <li class="step" :class="currentStep >= 1 ? 'step-primary' : ''">
                        选择服务
                    </li>
                    <li class="step" :class="currentStep >= 2 ? 'step-primary' : ''">
                        选择技师
                    </li>
                    <li class="step" :class="currentStep >= 3 ? 'step-primary' : ''">
                        确认预约
                    </li>
                </ul>
            </div>

            <!-- Modal Body -->
            <div class="p-6" v-if="loading">
                <div class="flex justify-center py-12">
                    <span class="loading loading-spinner loading-lg"></span>
                </div>
            </div>

            <div class="p-6" v-else>
                <!-- Step 1: Service & Time Selection -->
                <div v-if="currentStep === 1" class="space-y-5">
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text font-medium">选择会员
                                <span class="text-error">*</span></span>
                        </label>
                        <select v-model="formData.member_id" class="select select-bordered w-full" required>
                            <option disabled value="">请选择会员</option>
                            <option v-for="m in members" :key="m.id" :value="m.id">
                                {{ m.name }} ({{ m.phone }})
                            </option>
                        </select>
                    </div>

                    <div class="form-control">
                        <label class="label">
                            <span class="label-text font-medium">选择服务项目
                                <span class="text-error">*</span></span>
                        </label>
                        <select v-model="formData.service_id" class="select select-bordered w-full" required>
                            <option disabled value="">请选择服务项目</option>
                            <option v-for="s in services" :key="s.id" :value="s.id">
                                {{ s.name }} - ¥{{ s.price }} ({{
                                    s.duration
                                }}分钟)
                            </option>
                        </select>
                    </div>

                    <div class="form-control">
                        <label class="label">
                            <span class="label-text font-medium">选择开始时间
                                <span class="text-error">*</span></span>
                        </label>
                        <input type="datetime-local" v-model="formData.start_time" class="input input-bordered w-full"
                            required />
                        <label class="label">
                            <span class="label-text-alt text-base-content/60">预计结束时间:
                                {{
                                    selectedService
                                        ? formatTime(endTime)
                                        : "请先选择服务项目"
                                }}</span>
                        </label>
                    </div>

                    <div class="alert alert-info">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            class="stroke-current shrink-0 w-6 h-6">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span class="text-sm">下一步将根据您选择的服务和时间，为您筛选具备相应技能的可用技师</span>
                    </div>
                </div>

                <!-- Step 2: Technician Selection -->
                <div v-if="currentStep === 2" class="space-y-5">
                    <div v-if="formData.allow_waitlist" class="alert alert-warning">
                        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                        <span>该时间段所有技师都不可用，您选择的预约将加入候补队列</span>
                    </div>

                    <div class="form-control">
                        <label class="label">
                            <span class="label-text font-medium">选择技师
                                <span class="text-error">*</span></span>
                        </label>
                    </div>

                    <!-- Available Technicians -->
                    <div v-if="availableTechs.length > 0">
                        <div class="text-sm font-medium text-success mb-3 flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-5 h-5">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            可用技师 ({{ availableTechs.length }})
                            <span class="text-xs font-normal text-base-content/60" v-if="selectedServiceInfo">
                                (均具备 {{ selectedServiceInfo.name }} 技能)
                            </span>
                        </div>
                        <div class="grid grid-cols-2 gap-3">
                            <label v-for="tech in availableTechs" :key="tech.id" class="cursor-pointer">
                                <input type="radio" name="tech" :value="tech.id" v-model="formData.tech_id"
                                    class="peer sr-only" />
                                <div
                                    class="card bg-base-100 border-2 border-base-300 hover:border-primary peer-checked:border-primary peer-checked:bg-primary/5 transition-all">
                                    <div class="card-body p-4">
                                        <div class="flex items-center gap-3">
                                            <div
                                                class="w-12 h-12 rounded-full bg-success/10 flex items-center justify-center text-lg font-bold text-success">
                                                {{ tech.name.charAt(0) }}
                                            </div>
                                            <div class="flex-1">
                                                <div class="font-semibold text-base-content">
                                                    {{ tech.name }}
                                                </div>
                                                <div class="text-xs text-success font-medium">
                                                    空闲可用
                                                </div>
                                                <!-- Skill matching indicator -->
                                                <div class="text-xs text-base-content/60 mt-1 flex items-center gap-1" v-if="selectedServiceInfo">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3 text-success">
                                                        <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                                                    </svg>
                                                    具备 {{ selectedServiceInfo.name }} 技能
                                                </div>
                                            </div>
                                            <svg v-if="
                                                formData.tech_id === tech.id
                                            " xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                                                fill="currentColor" class="w-6 h-6 text-primary">
                                                <path fill-rule="evenodd"
                                                    d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
                                                    clip-rule="evenodd" />
                                            </svg>
                                        </div>
                                    </div>
                                </div>
                            </label>
                        </div>
                    </div>

                    <!-- Unavailable Technicians -->
                    <div v-if="
                        unavailableTechs.length > 0 &&
                        formData.allow_waitlist
                    " class="mt-6">
                        <div class="text-sm font-medium text-base-content/60 mb-3 flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-5 h-5">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                            </svg>
                            不可用技师 ({{ unavailableTechs.length }})
                        </div>
                        <div class="grid grid-cols-2 gap-3">
                            <label v-for="tech in unavailableTechs" :key="tech.id" class="cursor-pointer">
                                <input type="radio" name="tech" :value="tech.id" v-model="formData.tech_id"
                                    class="peer sr-only" />
                                <div
                                    class="card bg-base-100 border-2 border-base-300 hover:border-primary peer-checked:border-primary peer-checked:bg-primary/5 transition-all opacity-60">
                                    <div class="card-body p-4">
                                        <div class="flex items-center gap-3">
                                            <div
                                                class="w-12 h-12 rounded-full bg-error/10 flex items-center justify-center text-lg font-bold text-error">
                                                {{ tech.name.charAt(0) }}
                                            </div>
                                            <div class="flex-1">
                                                <div class="font-semibold text-base-content">
                                                    {{ tech.name }}
                                                </div>
                                                <!-- Updated unavailable reason -->
                                                <div class="text-xs text-error font-medium mt-1">
                                                    忙碌/休息/或不具备该服务技能
                                                </div>
                                            </div>
                                            <svg v-if="
                                                formData.tech_id === tech.id
                                            " xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                                                fill="currentColor" class="w-6 h-6 text-primary">
                                                <path fill-rule="evenodd"
                                                    d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
                                                    clip-rule="evenodd" />
                                            </svg>
                                        </div>
                                    </div>
                                </div>
                            </label>
                        </div>
                    </div>

                    <div v-if="
                        availableTechs.length === 0 &&
                        unavailableTechs.length === 0
                    " class="text-center py-8 text-base-content/60">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                            stroke="currentColor" class="w-12 h-12 mx-auto mb-2 opacity-30">
                            <path stroke-linecap="round" stroke-linejoin="round"
                                d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
                        </svg>
                        <p>没有可用的技师</p>
                    </div>
                </div>

                <!-- Step 3: Confirmation -->
                <div v-if="currentStep === 3" class="space-y-5">
                    <div class="alert alert-info">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            class="stroke-current shrink-0 w-6 h-6">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span>请确认以下预约信息</span>
                    </div>

                    <div class="bg-base-200 rounded-lg p-4 space-y-3">
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">会员</span>
                            <span class="font-medium">{{
                                selectedMember?.name
                            }}</span>
                        </div>
                        <div class="divider my-0"></div>
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">服务项目</span>
                            <span class="font-medium">{{
                                selectedService?.name
                            }}</span>
                        </div>
                        <div class="divider my-0"></div>
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">技师</span>
                            <span class="font-medium">{{
                                selectedTech?.name
                            }}</span>
                        </div>
                        <div class="divider my-0"></div>
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">开始时间</span>
                            <span class="font-medium">{{
                                formatTime(formData.value.start_time)
                            }}</span>
                        </div>
                        <div class="divider my-0"></div>
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">预计结束</span>
                            <span class="font-medium">{{
                                formatTime(endTime)
                            }}</span>
                        </div>
                        <div class="divider my-0"></div>
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">服务费用</span>
                            <span class="font-medium">¥{{ selectedService?.price }}</span>
                        </div>
                        <div class="divider my-0"></div>
                        <div class="flex justify-between items-center">
                            <span class="text-base-content/60">预约状态</span>
                            <span class="font-medium" :class="{
                                'text-warning': formData.value.allow_waitlist,
                                'text-success': !formData.value.allow_waitlist
                            }">
                                {{
                                    formData.value.allow_waitlist
                                        ? "候补中"
                                        : "待服务"
                                }}
                            </span>
                        </div>
                    </div>

                    <div v-if="formData.value.allow_waitlist" class="alert alert-warning">
                        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                        <span class="text-sm">预约将加入候补队列，一旦有技师可用，系统将自动为您安排。</span>
                    </div>
                </div>
            </div>

            <!-- Modal Footer -->
            <div class="px-6 py-4 border-t border-base-200 bg-base-200/30 flex justify-between">
                <button @click="goBack" class="btn btn-ghost" v-if="currentStep > 1">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                        stroke="currentColor" class="w-4 h-4 mr-2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
                    </svg>
                    上一步
                </button>
                <div class="flex gap-2">
                    <button @click="closeModal" class="btn btn-ghost">
                        取消
                    </button>
                    <button v-if="currentStep === 1" @click="goToStep2" class="btn btn-neutral"
                        :disabled="!canProceedStep1 || loading">
                        <span v-if="loading" class="loading loading-spinner loading-xs mr-2"></span>
                        下一步
                    </button>
                    <button v-if="currentStep === 2" @click="goToStep3" class="btn btn-neutral"
                        :disabled="!canProceedStep2">
                        下一步
                    </button>
                    <button v-if="currentStep === 3" @click="handleSubmit" class="btn btn-neutral"
                        :disabled="submitting">
                        <span v-if="submitting" class="loading loading-spinner loading-xs mr-2"></span>
                        {{
                            submitting
                                ? "提交中..."
                                : "确认预约"
                        }}
                    </button>
                </div>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
            <button @click="closeModal">close</button>
        </form>
    </dialog>
</template>
