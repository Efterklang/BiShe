<script setup>
import { ref, computed, onMounted, watch } from "vue";
import {
    getSchedules,
    batchSetSchedule,
} from "../api/schedules";
import { getTechnicians } from "../api/technicians";

// --- Props ---
const props = defineProps({
    selectedTechnician: {
        type: Object,
        default: null,
    },
});

// --- State ---
const currentDate = ref(new Date());
const technicians = ref([]);
const selectedTechIds = ref([]);
const schedules = ref([]);
const loading = ref(false);
const processing = ref(false);

// Batch Edit Mode (Month View)
const isEditMode = ref(false);
const selectedDates = ref(new Set());



// --- Computed ---
const currentYear = computed(() => currentDate.value.getFullYear());
const currentMonth = computed(() => currentDate.value.getMonth());

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

const currentLabel = computed(() => {
    return currentDate.value.toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "long",
    });
});







// --- Helpers ---
const formatDate = (date) => {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
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
        const start = formatDate(
            new Date(currentYear.value, currentMonth.value, 1),
        );
        const end = formatDate(
            new Date(currentYear.value, currentMonth.value + 1, 0),
        );

        // Add timestamp to prevent caching
        const params = {
            start_date: start,
            end_date: end,
            _t: new Date().getTime(),
        };
        const data = await getSchedules(params);
        schedules.value = data || [];
    } catch (e) {
        console.error("Failed to fetch schedules", e);
    } finally {
        loading.value = false;
    }
};



// --- Actions ---


const changeDate = (delta) => {
    const newDate = new Date(currentDate.value);
    newDate.setMonth(newDate.getMonth() + delta);
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
    <div class="bg-base-100 rounded-xl border border-base-300 shadow-sm overflow-hidden">
        <!-- Toolbar -->
        <div class="p-4 border-b border-base-200 flex flex-col md:flex-row justify-between items-center gap-4">
            <!-- Left: Navigation -->
            <div class="flex items-center gap-2">
                <div class="join">
                    <button @click="changeDate(-1)" class="btn btn-sm join-item">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                            stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
                        </svg>
                    </button>
                    <button @click="goToToday" class="btn btn-sm join-item">
                        今天
                    </button>
                    <button @click="changeDate(1)" class="btn btn-sm join-item">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                            stroke="currentColor" class="w-4 h-4">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                        </svg>
                    </button>
                </div>
                <span class="text-lg font-bold ml-2 min-w-35">{{
                    currentLabel
                }}</span>
            </div>

            <!-- Right: Controls -->
            <div class="flex items-center gap-3 flex-wrap justify-end">


                <!-- Tech Filter -->
                <div class="dropdown dropdown-end">
                    <label tabindex="0" class="btn btn-sm btn-outline">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                            stroke="currentColor" class="w-4 h-4 mr-1">
                            <path stroke-linecap="round" stroke-linejoin="round"
                                d="M12 3c2.755 0 5.455.232 8.083.678.533.09.917.556.917 1.096v1.044a2.25 2.25 0 01-.659 1.591l-5.432 5.432a2.25 2.25 0 00-.659 1.591v2.927a2.25 2.25 0 01-1.244 2.013L9.75 21v-6.568a2.25 2.25 0 00-.659-1.591L3.659 7.409A2.25 2.25 0 013 5.818V4.774c0-.54.384-1.006.917-1.096A48.32 48.32 0 0112 3z" />
                        </svg>
                        筛选技师 ({{ selectedTechIds.length || "全" }})
                    </label>
                    <div tabindex="0"
                        class="dropdown-content z-1 menu p-2 shadow bg-base-100 rounded-box w-52 max-h-64 overflow-y-auto border border-base-200">
                        <li v-for="tech in technicians" :key="tech.id">
                            <label class="label cursor-pointer justify-start gap-2">
                                <input type="checkbox" class="checkbox checkbox-xs checkbox-primary" :value="tech.id"
                                    v-model="selectedTechIds" />
                                <span class="label-text">{{ tech.name }}</span>
                            </label>
                        </li>
                    </div>
                </div>

                <!-- Batch Edit Toggle -->
                <div class="flex gap-2">
                    <button class="btn btn-sm" :class="isEditMode ? 'btn-neutral' : 'btn-ghost'"
                        @click="isEditMode = !isEditMode">
                        {{ isEditMode ? "退出排班" : "批量排班" }}
                    </button>

                    <div v-if="isEditMode" class="join">
                        <button class="btn btn-sm btn-success join-item text-white" @click="handleBatchSet(true)"
                            :disabled="processing">
                            设为在岗
                        </button>
                        <button class="btn btn-sm btn-error join-item text-white" @click="handleBatchSet(false)"
                            :disabled="processing">
                            设为休息
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Content Area -->
        <div class="relative min-h-125">
            <div v-if="loading" class="absolute inset-0 bg-base-100/50 z-10 flex items-center justify-center">
                <span class="loading loading-spinner loading-lg text-primary"></span>
            </div>

            <!-- Month View -->
            <div class="p-4">
                <div class="grid grid-cols-7 gap-2 mb-2 text-center text-sm font-semibold text-base-content/60">
                    <div>周一</div>
                    <div>周二</div>
                    <div>周三</div>
                    <div>周四</div>
                    <div>周五</div>
                    <div>周六</div>
                    <div>周日</div>
                </div>
                <div class="grid grid-cols-7 gap-2">
                    <div v-for="day in calendarGrid" :key="day.id"
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
                        ]" @click="handleDayClick(day)">
                        <div v-if="day.date">
                            <span class="text-sm font-medium" :class="day.isToday ? 'text-primary' : ''">{{
                                day.date.getDate() }}</span>

                            <!-- Status Dots -->
                            <div class="mt-2 space-y-1">
                                <div v-if="getDayStatus(day.dateStr).leave > 0"
                                    class="flex items-center gap-1 text-xs text-error">
                                    <span class="w-2 h-2 rounded-full bg-error"></span>
                                    <span>{{
                                        getDayStatus(day.dateStr).leave
                                        }}
                                        休</span>
                                </div>
                                <div v-if="getDayStatus(day.dateStr).work > 0"
                                    class="flex items-center gap-1 text-xs text-success">
                                    <span class="w-2 h-2 rounded-full bg-success"></span>
                                    <span>{{
                                        getDayStatus(day.dateStr).work
                                        }}
                                        班</span>
                                </div>
                            </div>

                            <!-- Selection Checkmark -->
                            <div v-if="isEditMode && day.isSelected" class="absolute top-2 right-2 text-primary">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
                                    class="w-5 h-5">
                                    <path fill-rule="evenodd"
                                        d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
                                        clip-rule="evenodd" />
                                </svg>
                            </div>
                        </div>
                    </div>
                </div>
            </div>


        </div>


    </div>
</template>