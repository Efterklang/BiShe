<script setup>
import { ref, computed, onMounted, watch } from "vue";
import {
    getSchedules,
    batchSetSchedule,
    getTechnicianScheduleDetail,
} from "../api/schedules";
import { getTechnicians } from "../api/technicians";
import { getAppointments } from "../api/appointments";

// --- Props ---
const props = defineProps({
    selectedTechnician: {
        type: Object,
        default: null,
    },
});

// --- State ---
const viewMode = ref("month"); // 'month' | 'day'
const currentDate = ref(new Date());
const technicians = ref([]);
const selectedTechIds = ref([]);
const schedules = ref([]);
const appointments = ref([]);
const loading = ref(false);
const processing = ref(false);

// Batch Edit Mode (Month View)
const isEditMode = ref(false);
const selectedDates = ref(new Set());

// Schedule Detail Modal
const showDetailModal = ref(false);
const detailLoading = ref(false);
const scheduleDetail = ref(null);

// --- Computed ---
const currentYear = computed(() => currentDate.value.getFullYear());
const currentMonth = computed(() => currentDate.value.getMonth());

const currentLabel = computed(() => {
    if (viewMode.value === "month") {
        return currentDate.value.toLocaleDateString("zh-CN", {
            year: "numeric",
            month: "long",
        });
    }
    return currentDate.value.toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "long",
        day: "numeric",
        weekday: "short",
    });
});

const activeTechnicians = computed(() => {
    if (selectedTechIds.value.length === 0) return technicians.value;
    return technicians.value.filter((t) =>
        selectedTechIds.value.includes(t.id),
    );
});

// Month Grid
const calendarGrid = computed(() => {
    const year = currentYear.value;
    const month = currentMonth.value;
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);

    const days = [];

    // Monday start padding
    let startPadding = firstDay.getDay() - 1;
    if (startPadding < 0) startPadding = 6;

    for (let i = 0; i < startPadding; i++) {
        days.push({ date: null, id: `pad-start-${i}` });
    }

    for (let d = 1; d <= lastDay.getDate(); d++) {
        const date = new Date(year, month, d);
        const dateStr = formatDate(date);
        days.push({
            date: date,
            dateStr: dateStr,
            id: dateStr,
            isToday: dateStr === formatDate(new Date()),
            isSelected: selectedDates.value.has(dateStr),
        });
    }

    return days;
});

// Day View Time Slots (10:00 - 22:00)
const timeSlots = Array.from({ length: 13 }, (_, i) => i + 10); // 10, 11, ... 22

// --- Helpers ---
const formatDate = (date) => {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
};

const getMinutesFrom10AM = (dateStr) => {
    const date = new Date(dateStr);
    const hours = date.getHours();
    const minutes = date.getMinutes();
    // 10:00 is 0 minutes
    return hours * 60 + minutes - 10 * 60;
};

// --- API ---
const fetchTechniciansList = async () => {
    try {
        const res = await getTechnicians();
        technicians.value = res || [];
    } catch (e) {
        console.error("Failed to fetch technicians", e);
    }
};

const fetchSchedules = async () => {
    loading.value = true;
    try {
        let start, end;
        if (viewMode.value === "month") {
            start = formatDate(
                new Date(currentYear.value, currentMonth.value, 1),
            );
            end = formatDate(
                new Date(currentYear.value, currentMonth.value + 1, 0),
            );
        } else {
            start = formatDate(currentDate.value);
            end = start;
        }

        // Add timestamp to prevent caching
        const params = {
            start_date: start,
            end_date: end,
            _t: new Date().getTime(),
        };
        const data = await getSchedules(params);
        schedules.value = data || [];

        if (viewMode.value === "day") {
            await fetchAppointments(start);
        }
    } catch (e) {
        console.error("Failed to fetch schedules", e);
    } finally {
        loading.value = false;
    }
};

const fetchAppointments = async (dateStr) => {
    try {
        // Note: In a real app, backend should support date range filtering for appointments.
        // Here we fetch list and filter client side as per current backend capabilities.
        const allAppts = await getAppointments();
        appointments.value = (allAppts || []).filter((app) =>
            app.start_time.startsWith(dateStr),
        );
    } catch (e) {
        console.error("Failed to fetch appointments", e);
    }
};

// --- Actions ---
const switchView = (mode) => {
    viewMode.value = mode;
    selectedDates.value.clear();
    isEditMode.value = false;
    fetchSchedules();
};

const changeDate = (delta) => {
    const newDate = new Date(currentDate.value);
    if (viewMode.value === "month") {
        newDate.setMonth(newDate.getMonth() + delta);
    } else {
        newDate.setDate(newDate.getDate() + delta);
    }
    currentDate.value = newDate;
    fetchSchedules();
};

const goToToday = () => {
    currentDate.value = new Date();
    fetchSchedules();
};

const handleDayClick = (day) => {
    if (!day.date) return;

    if (isEditMode.value) {
        if (selectedDates.value.has(day.dateStr)) {
            selectedDates.value.delete(day.dateStr);
        } else {
            selectedDates.value.add(day.dateStr);
        }
    } else {
        currentDate.value = day.date;
        switchView("day");
    }
};

const handleBatchSet = async (isAvailable) => {
    if (selectedTechIds.value.length === 0) {
        alert("请先在右上角筛选中选择至少一位技师");
        return;
    }
    if (selectedDates.value.size === 0) {
        alert("请先点击日历选择日期");
        return;
    }

    processing.value = true;
    try {
        const dates = Array.from(selectedDates.value);

        await batchSetSchedule({
            tech_ids: selectedTechIds.value,
            dates: dates,
            is_available: isAvailable,
        });

        await fetchSchedules();
        selectedDates.value.clear();
        isEditMode.value = false;
        alert("排班设置成功");
    } catch (e) {
        alert("设置失败: " + (e.response?.data?.msg || e.message));
    } finally {
        processing.value = false;
    }
};

// --- Data Helpers ---
const getDayStatus = (dateStr) => {
    const daySchedules = schedules.value.filter((s) =>
        s.date.startsWith(dateStr),
    );

    // If filtering specific techs
    if (selectedTechIds.value.length > 0) {
        const statuses = selectedTechIds.value.map((tid) => {
            const s = daySchedules.find((ds) => ds.tech_id === tid);
            // Default to available if no record exists
            return s ? (s.is_available ? "work" : "leave") : "work";
        });
        return {
            total: selectedTechIds.value.length,
            leave: statuses.filter((s) => s === "leave").length,
            work: statuses.filter((s) => s === "work").length,
        };
    }

    // Global view
    const leaveCount = daySchedules.filter((s) => !s.is_available).length;
    const totalTechs = technicians.value.length;
    return {
        total: totalTechs,
        leave: leaveCount,
        work: totalTechs - leaveCount,
    };
};

const getTechDaySchedule = (techId) => {
    const dateStr = formatDate(currentDate.value);
    const s = schedules.value.find(
        (s) => s.date.startsWith(dateStr) && s.tech_id === techId,
    );
    return s ? s.is_available : true; // Default available
};

const getTechAppointments = (techId) => {
    return appointments.value
        .filter((a) => a.tech_id === techId)
        .map((a) => {
            const startMins = getMinutesFrom10AM(a.start_time);
            const endMins = getMinutesFrom10AM(a.end_time);
            const duration = endMins - startMins;

            // Calculate position and width percentage relative to 12 hours (720 mins)
            const left = (startMins / 720) * 100;
            const width = (duration / 720) * 100;

            return {
                ...a,
                style: {
                    left: `${Math.max(0, left)}%`,
                    width: `${Math.min(100 - left, width)}%`,
                },
            };
        });
};

// --- Detail Modal ---
const showScheduleDetail = async (tech, dateStr) => {
    detailLoading.value = true;
    showDetailModal.value = true;
    scheduleDetail.value = null;

    try {
        const data = await getTechnicianScheduleDetail({
            tech_id: tech.id,
            date: dateStr,
        });
        scheduleDetail.value = {
            ...data,
            technician: tech,
        };
    } catch (e) {
        console.error("Failed to fetch schedule detail", e);
        alert("获取排班详情失败");
        showDetailModal.value = false;
    } finally {
        detailLoading.value = false;
    }
};

const closeDetailModal = () => {
    showDetailModal.value = false;
    scheduleDetail.value = null;
};

const formatTime = (timeStr) => {
    if (!timeStr) return "";
    return timeStr.substring(11, 16);
};

// Watch for selectedTechnician prop changes
watch(
    () => props.selectedTechnician,
    (newTech) => {
        if (newTech) {
            selectedTechIds.value = [newTech.id];
            fetchSchedules();
        }
    },
    { immediate: true },
);

onMounted(async () => {
    await fetchTechniciansList();
    await fetchSchedules();
});
</script>

<template>
    <div
        class="bg-base-100 rounded-xl border border-base-300 shadow-sm overflow-hidden"
    >
        <!-- Toolbar -->
        <div
            class="p-4 border-b border-base-200 flex flex-col md:flex-row justify-between items-center gap-4"
        >
            <!-- Left: Navigation -->
            <div class="flex items-center gap-2">
                <div class="join">
                    <button
                        @click="changeDate(-1)"
                        class="btn btn-sm join-item"
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
                                d="M15.75 19.5L8.25 12l7.5-7.5"
                            />
                        </svg>
                    </button>
                    <button @click="goToToday" class="btn btn-sm join-item">
                        今天
                    </button>
                    <button @click="changeDate(1)" class="btn btn-sm join-item">
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
                                d="M8.25 4.5l7.5 7.5-7.5 7.5"
                            />
                        </svg>
                    </button>
                </div>
                <span class="text-lg font-bold ml-2 min-w-[140px]">{{
                    currentLabel
                }}</span>
            </div>

            <!-- Right: Controls -->
            <div class="flex items-center gap-3 flex-wrap justify-end">
                <!-- View Switcher -->
                <div class="tabs tabs-boxed tabs-sm bg-base-200">
                    <a
                        class="tab"
                        :class="{ 'tab-active': viewMode === 'month' }"
                        @click="switchView('month')"
                        >月视图</a
                    >
                    <a
                        class="tab"
                        :class="{ 'tab-active': viewMode === 'day' }"
                        @click="switchView('day')"
                        >日视图</a
                    >
                </div>

                <!-- Tech Filter -->
                <div class="dropdown dropdown-end">
                    <label tabindex="0" class="btn btn-sm btn-outline">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="1.5"
                            stroke="currentColor"
                            class="w-4 h-4 mr-1"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M12 3c2.755 0 5.455.232 8.083.678.533.09.917.556.917 1.096v1.044a2.25 2.25 0 01-.659 1.591l-5.432 5.432a2.25 2.25 0 00-.659 1.591v2.927a2.25 2.25 0 01-1.244 2.013L9.75 21v-6.568a2.25 2.25 0 00-.659-1.591L3.659 7.409A2.25 2.25 0 013 5.818V4.774c0-.54.384-1.006.917-1.096A48.32 48.32 0 0112 3z"
                            />
                        </svg>
                        筛选技师 ({{ selectedTechIds.length || "全" }})
                    </label>
                    <div
                        tabindex="0"
                        class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52 max-h-64 overflow-y-auto border border-base-200"
                    >
                        <li v-for="tech in technicians" :key="tech.id">
                            <label
                                class="label cursor-pointer justify-start gap-2"
                            >
                                <input
                                    type="checkbox"
                                    class="checkbox checkbox-xs checkbox-primary"
                                    :value="tech.id"
                                    v-model="selectedTechIds"
                                />
                                <span class="label-text">{{ tech.name }}</span>
                            </label>
                        </li>
                    </div>
                </div>

                <!-- Batch Edit Toggle (Month Only) -->
                <div v-if="viewMode === 'month'" class="flex gap-2">
                    <button
                        class="btn btn-sm"
                        :class="isEditMode ? 'btn-neutral' : 'btn-ghost'"
                        @click="isEditMode = !isEditMode"
                    >
                        {{ isEditMode ? "退出排班" : "批量排班" }}
                    </button>

                    <div v-if="isEditMode" class="join">
                        <button
                            class="btn btn-sm btn-success join-item text-white"
                            @click="handleBatchSet(true)"
                            :disabled="processing"
                        >
                            设为在岗
                        </button>
                        <button
                            class="btn btn-sm btn-error join-item text-white"
                            @click="handleBatchSet(false)"
                            :disabled="processing"
                        >
                            设为休息
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Content Area -->
        <div class="relative min-h-[500px]">
            <div
                v-if="loading"
                class="absolute inset-0 bg-base-100/50 z-10 flex items-center justify-center"
            >
                <span
                    class="loading loading-spinner loading-lg text-primary"
                ></span>
            </div>

            <!-- Month View -->
            <div v-if="viewMode === 'month'" class="p-4">
                <div
                    class="grid grid-cols-7 gap-2 mb-2 text-center text-sm font-semibold text-base-content/60"
                >
                    <div>周一</div>
                    <div>周二</div>
                    <div>周三</div>
                    <div>周四</div>
                    <div>周五</div>
                    <div>周六</div>
                    <div>周日</div>
                </div>
                <div class="grid grid-cols-7 gap-2">
                    <div
                        v-for="day in calendarGrid"
                        :key="day.id"
                        class="aspect-square rounded-lg border p-2 relative transition-all cursor-pointer flex flex-col justify-between"
                        :class="[
                            day.date
                                ? 'hover:border-primary/50 bg-base-100'
                                : 'border-transparent bg-transparent pointer-events-none',
                            day.isToday
                                ? 'ring-2 ring-primary ring-offset-2'
                                : 'border-base-200',
                            day.isSelected
                                ? 'bg-primary/10 border-primary'
                                : '',
                        ]"
                        @click="handleDayClick(day)"
                    >
                        <div v-if="day.date">
                            <span
                                class="text-sm font-medium"
                                :class="day.isToday ? 'text-primary' : ''"
                                >{{ day.date.getDate() }}</span
                            >

                            <!-- Status Dots -->
                            <div class="mt-2 space-y-1">
                                <div
                                    v-if="getDayStatus(day.dateStr).leave > 0"
                                    class="flex items-center gap-1 text-xs text-error"
                                >
                                    <span
                                        class="w-2 h-2 rounded-full bg-error"
                                    ></span>
                                    <span
                                        >{{
                                            getDayStatus(day.dateStr).leave
                                        }}
                                        休</span
                                    >
                                </div>
                                <div
                                    v-if="getDayStatus(day.dateStr).work > 0"
                                    class="flex items-center gap-1 text-xs text-success"
                                >
                                    <span
                                        class="w-2 h-2 rounded-full bg-success"
                                    ></span>
                                    <span
                                        >{{
                                            getDayStatus(day.dateStr).work
                                        }}
                                        班</span
                                    >
                                </div>
                            </div>

                            <!-- Selection Checkmark -->
                            <div
                                v-if="isEditMode && day.isSelected"
                                class="absolute top-2 right-2 text-primary"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 24 24"
                                    fill="currentColor"
                                    class="w-5 h-5"
                                >
                                    <path
                                        fill-rule="evenodd"
                                        d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
                                        clip-rule="evenodd"
                                    />
                                </svg>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Day View (Swimlane) -->
            <div v-else class="flex flex-col h-full overflow-x-auto">
                <!-- Header Time Scale -->
                <div
                    class="flex border-b border-base-200 bg-base-200/30 min-w-[800px]"
                >
                    <div
                        class="w-32 flex-shrink-0 p-3 font-semibold text-sm border-r border-base-200 bg-base-100 sticky left-0 z-20"
                    >
                        技师
                    </div>
                    <div class="flex-1 flex relative h-10">
                        <div
                            v-for="hour in timeSlots"
                            :key="hour"
                            class="flex-1 border-l border-base-300/50 text-xs text-base-content/50 p-1"
                        >
                            {{ hour }}:00
                        </div>
                    </div>
                </div>

                <!-- Rows -->
                <div class="flex-1 overflow-y-auto min-w-[800px]">
                    <div
                        v-for="tech in activeTechnicians"
                        :key="tech.id"
                        class="flex border-b border-base-200 hover:bg-base-100/50 transition-colors h-16"
                    >
                        <!-- Tech Name -->
                        <div
                            class="w-32 flex-shrink-0 p-3 flex items-center gap-2 border-r border-base-200 bg-base-100 sticky left-0 z-10"
                        >
                            <div
                                class="w-8 h-8 rounded-full bg-base-300 flex items-center justify-center text-xs font-bold"
                            >
                                {{ tech.name.charAt(0) }}
                            </div>
                            <div class="flex flex-col">
                                <span
                                    class="text-sm font-medium truncate w-16"
                                    >{{ tech.name }}</span
                                >
                                <span
                                    class="text-[10px] px-1.5 py-0.5 rounded-full w-fit"
                                    :class="
                                        getTechDaySchedule(tech.id)
                                            ? 'bg-success/10 text-success'
                                            : 'bg-error/10 text-error'
                                    "
                                >
                                    {{
                                        getTechDaySchedule(tech.id)
                                            ? "在岗"
                                            : "休息"
                                    }}
                                </span>
                            </div>
                            <button
                                @click="
                                    showScheduleDetail(
                                        tech,
                                        formatDate(currentDate),
                                    )
                                "
                                class="btn btn-xs btn-ghost ml-auto"
                                title="查看详情"
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
                                        d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                                    />
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                    />
                                </svg>
                            </button>
                        </div>

                        <!-- Timeline -->
                        <div class="flex-1 relative bg-base-100">
                            <!-- Grid Lines -->
                            <div
                                class="absolute inset-0 flex pointer-events-none"
                            >
                                <div
                                    v-for="hour in timeSlots"
                                    :key="hour"
                                    class="flex-1 border-l border-base-200 border-dashed"
                                ></div>
                            </div>

                            <!-- Unavailable Overlay -->
                            <div
                                v-if="!getTechDaySchedule(tech.id)"
                                class="absolute inset-0 bg-base-300/30 flex items-center justify-center"
                            >
                                <span
                                    class="text-xs text-base-content/40 font-medium tracking-widest"
                                    >RESTING</span
                                >
                                <div
                                    class="absolute inset-0"
                                    style="
                                        background-image: repeating-linear-gradient(
                                            45deg,
                                            transparent,
                                            transparent 10px,
                                            rgba(0, 0, 0, 0.03) 10px,
                                            rgba(0, 0, 0, 0.03) 20px
                                        );
                                    "
                                ></div>
                            </div>

                            <!-- Appointments -->
                            <div v-else class="absolute inset-0 top-2 bottom-2">
                                <div
                                    v-for="appt in getTechAppointments(tech.id)"
                                    :key="appt.id"
                                    class="absolute h-full rounded-md bg-primary text-primary-content text-xs p-1 overflow-hidden shadow-sm border border-primary-focus hover:z-10 hover:shadow-md transition-all cursor-pointer"
                                    :style="appt.style"
                                    :title="`${appt.start_time.substring(11, 16)} - ${appt.end_time.substring(11, 16)}`"
                                >
                                    <div class="font-bold truncate">
                                        {{ appt.member?.name || "Guest" }}
                                    </div>
                                    <div
                                        class="opacity-80 truncate scale-90 origin-top-left"
                                    >
                                        {{ appt.service_item?.name }}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Schedule Detail Modal -->
        <dialog class="modal" :class="{ 'modal-open': showDetailModal }">
            <div
                class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl max-w-2xl"
            >
                <!-- Modal Header -->
                <div
                    class="flex justify-between items-center mb-4 pb-4 border-b border-base-200"
                >
                    <h3 class="font-bold text-lg">
                        排班详情
                        <span
                            v-if="scheduleDetail"
                            class="text-base-content/60 font-normal ml-2"
                        >
                            {{ scheduleDetail.technician?.name }} -
                            {{ scheduleDetail.date }}
                        </span>
                    </h3>
                    <button
                        @click="closeDetailModal"
                        class="btn btn-ghost btn-sm btn-square"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="2"
                            stroke="currentColor"
                            class="w-5 h-5"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <!-- Modal Body -->
                <div v-if="detailLoading" class="flex justify-center py-12">
                    <span class="loading loading-spinner loading-lg"></span>
                </div>

                <div v-else-if="scheduleDetail" class="space-y-4">
                    <!-- Status Badge -->
                    <div class="flex items-center gap-2">
                        <span class="text-sm font-medium">状态:</span>
                        <span
                            class="badge"
                            :class="
                                scheduleDetail.is_available
                                    ? 'badge-success'
                                    : 'badge-error'
                            "
                        >
                            {{ scheduleDetail.is_available ? "在岗" : "休息" }}
                        </span>
                    </div>

                    <!-- Appointments List -->
                    <div>
                        <h4 class="font-semibold mb-3 flex items-center gap-2">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke-width="1.5"
                                stroke="currentColor"
                                class="w-5 h-5"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5"
                                />
                            </svg>
                            预约列表 ({{
                                scheduleDetail.appointments?.length || 0
                            }})
                        </h4>

                        <div
                            v-if="
                                !scheduleDetail.appointments ||
                                scheduleDetail.appointments.length === 0
                            "
                            class="text-center py-8 text-base-content/60"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke-width="1.5"
                                stroke="currentColor"
                                class="w-12 h-12 mx-auto mb-2 opacity-30"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5"
                                />
                            </svg>
                            <p>暂无预约</p>
                        </div>

                        <div v-else class="space-y-3 max-h-96 overflow-y-auto">
                            <div
                                v-for="appt in scheduleDetail.appointments"
                                :key="appt.id"
                                class="p-4 border border-base-200 rounded-lg hover:border-primary/50 transition-colors"
                            >
                                <div
                                    class="flex justify-between items-start mb-2"
                                >
                                    <div>
                                        <div class="font-semibold text-base">
                                            {{
                                                appt.member?.name || "未知客户"
                                            }}
                                        </div>
                                        <div
                                            class="text-sm text-base-content/60"
                                        >
                                            {{
                                                appt.service_item?.name ||
                                                "未知服务"
                                            }}
                                        </div>
                                    </div>
                                    <span
                                        class="badge badge-sm"
                                        :class="{
                                            'badge-warning':
                                                appt.status === 'pending',
                                            'badge-success':
                                                appt.status === 'completed',
                                            'badge-info':
                                                appt.status === 'waitlist' ||
                                                appt.status === 'waiting',
                                            'badge-error':
                                                appt.status === 'cancelled',
                                        }"
                                    >
                                        {{
                                            appt.status === "pending"
                                                ? "待服务"
                                                : appt.status === "completed"
                                                  ? "已完成"
                                                  : appt.status ===
                                                          "waitlist" ||
                                                      appt.status === "waiting"
                                                    ? "候补中"
                                                    : appt.status ===
                                                        "cancelled"
                                                      ? "已取消"
                                                      : appt.status
                                        }}
                                    </span>
                                </div>

                                <div
                                    class="flex items-center gap-4 text-sm text-base-content/70"
                                >
                                    <div class="flex items-center gap-1">
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
                                                d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                        {{ formatTime(appt.start_time) }} -
                                        {{ formatTime(appt.end_time) }}
                                    </div>
                                    <div class="flex items-center gap-1">
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
                                                d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                        ¥{{ appt.actual_price }}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Modal Footer -->
                <div class="modal-action">
                    <button @click="closeDetailModal" class="btn btn-neutral">
                        关闭
                    </button>
                </div>
            </div>
            <form
                method="dialog"
                class="modal-backdrop bg-base-content/20 backdrop-blur-sm"
            >
                <button @click="closeDetailModal">close</button>
            </form>
        </dialog>
    </div>
</template>
