<script setup>
import { ref, onMounted } from 'vue';
import { getTechnicians, createTechnician, updateTechnician } from '../api/technicians';
import TechnicianSchedule from '../components/TechnicianSchedule.vue';

const activeTab = ref('overview');
const technicians = ref([]);
const loading = ref(true);
const showModal = ref(false);
const submitting = ref(false);
const editingId = ref(null);

const formData = ref({
  name: '',
  skills: '', // Comma separated string for input
  status: 0
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
  formData.value = { name: '', skills: '', status: 0 };
  showModal.value = true;
};

const handleEdit = (tech) => {
  editingId.value = tech.id;
  const skills = parseSkills(tech.skills || tech.Skills);
  formData.value = {
    name: tech.name,
    skills: Array.isArray(skills) ? skills.join(', ') : '',
    status: tech.status
  };
  showModal.value = true;
};

const handleSchedule = (tech) => {
  activeTab.value = 'schedule';
};

const handleSubmit = async () => {
  submitting.value = true;
  try {
    // Convert comma-separated skills string to array
    const skillsArray = formData.value.skills
      .split(/[,，]/) // Split by comma (English or Chinese)
      .map(s => s.trim())
      .filter(s => s);

    const payload = {
      name: formData.value.name,
      skills: JSON.stringify(skillsArray),
      status: Number(formData.value.status)
    };

    if (editingId.value) {
      await updateTechnician(editingId.value, payload);
    } else {
      await createTechnician(payload);
    }

    showModal.value = false;
    await fetchTechnicians();
  } catch (error) {
    alert((editingId.value ? '更新' : '创建') + '失败: ' + (error.message || '未知错误'));
  } finally {
    submitting.value = false;
  }
};

const getStatusInfo = (status) => {
  switch (Number(status)) {
    case 0: return { text: '空闲', class: 'badge-success' };
    case 1: return { text: '忙碌', class: 'badge-warning' };
    case 2: return { text: '请假', class: 'badge-ghost' };
    default: return { text: '未知', class: 'badge-ghost' };
  }
};

// Helper to parse skills
const parseSkills = (skills) => {
  if (Array.isArray(skills)) return skills;
  if (typeof skills === 'string') {
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
        <h1 class="text-3xl font-bold tracking-tight text-base-content">技师管理</h1>
        <p class="mt-2 text-base-content/60">
          管理店内技师团队，查看实时状态与技能分布。
        </p>
      </div>
      <button
        v-if="activeTab === 'overview'"
        @click="openCreateModal"
        class="btn btn-neutral"
      >
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 mr-2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        添加技师
      </button>
    </div>

    <!-- Tabs -->
    <div role="tablist" class="tabs tabs-bordered mb-6">
      <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'overview' }" @click="activeTab = 'overview'">技师总览</a>
      <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'schedule' }" @click="activeTab = 'schedule'">排班管理</a>
    </div>

    <div v-if="activeTab === 'overview'">
    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div v-for="i in 4" :key="i" class="h-64 rounded-xl border border-base-300 bg-base-200 animate-pulse"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="technicians.length === 0" class="flex flex-col items-center justify-center py-20 border border-dashed border-base-300 rounded-xl bg-base-200/50">
      <div class="p-4 bg-base-300 rounded-full mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8 text-base-content/40">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-base-content">暂无技师</h3>
      <p class="text-base-content/60 mt-1">点击右上角按钮添加第一位技师。</p>
    </div>

    <!-- Technicians Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div
        v-for="tech in technicians"
        :key="tech.id"
        class="group relative flex flex-col bg-base-100 border border-base-300 rounded-xl p-6 hover:border-base-content/20 transition-all duration-200 hover:shadow-sm"
      >
        <!-- Status Badge -->
        <div class="absolute top-4 right-4">
          <span
            class="badge badge-sm"
            :class="getStatusInfo(tech.status).class"
          >
            {{ getStatusInfo(tech.status).text }}
          </span>
        </div>

        <!-- Avatar & Info -->
        <div class="flex flex-col items-center text-center mb-4">
          <div class="w-20 h-20 rounded-full bg-base-200 flex items-center justify-center text-2xl font-semibold text-base-content/60 mb-4 ring-4 ring-base-100 shadow-sm">
            {{ tech.name ? tech.name.charAt(0) : '?' }}
          </div>
          <h3 class="text-lg font-semibold text-base-content">{{ tech.name }}</h3>

          <!-- Rating -->
          <div class="flex items-center gap-1 mt-1 text-warning text-sm font-medium">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-4 h-4">
              <path fill-rule="evenodd" d="M10.788 3.21c.448-1.077 1.976-1.077 2.424 0l2.082 5.007 5.404.433c1.164.093 1.636 1.545.749 2.305l-4.117 3.527 1.257 5.273c.271 1.136-.964 2.033-1.96 1.425L12 18.354 7.373 21.18c-.996.608-2.231-.29-1.96-1.425l1.257-5.273-4.117-3.527c-.887-.76-.415-2.212.749-2.305l5.404-.433 2.082-5.006z" clip-rule="evenodd" />
            </svg>
            <span>{{ tech.average_rating || tech.AverageRating || 5.0 }}</span>
          </div>
        </div>

        <!-- Skills -->
        <div class="flex-1">
          <div class="flex flex-wrap gap-2 justify-center">
            <span
              v-for="(skill, idx) in parseSkills(tech.skills || tech.Skills)"
              :key="idx"
              class="badge badge-outline text-xs"
            >
              {{ skill }}
            </span>
          </div>
        </div>

        <!-- Actions -->
        <div class="mt-6 pt-4 border-t border-base-200 flex gap-2">
          <button
            @click="handleSchedule(tech)"
            class="btn btn-sm btn-outline flex-1"
          >
            排班
          </button>
          <button
            @click="handleEdit(tech)"
            class="btn btn-sm btn-outline flex-1"
          >
            编辑
          </button>
        </div>
      </div>
    </div>
    </div>

    <div v-else-if="activeTab === 'schedule'">
      <TechnicianSchedule />
    </div>

    <!-- Create Modal -->
    <dialog class="modal" :class="{ 'modal-open': showModal }">
      <div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
          <h3 class="font-semibold text-lg text-base-content">{{ editingId ? '编辑技师' : '添加新技师' }}</h3>
          <button @click="showModal = false" class="btn btn-ghost btn-sm btn-square text-base-content/60">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6">
          <form @submit.prevent="handleSubmit" class="space-y-5">
            <div>
              <label class="block text-sm font-medium text-base-content/80 mb-1">姓名</label>
              <input
                type="text"
                v-model="formData.name"
                placeholder="请输入技师姓名"
                class="input input-bordered w-full bg-base-100"
                required
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-base-content/80 mb-1">
                技能标签 <span class="text-base-content/40 font-normal">(用逗号分隔)</span>
              </label>
              <textarea
                v-model="formData.skills"
                placeholder="例如: 精油SPA, 足疗, 中式推拿"
                rows="3"
                class="textarea textarea-bordered w-full bg-base-100 resize-none"
              ></textarea>
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
              <button
                type="submit"
                class="btn btn-neutral w-full"
                :disabled="submitting"
              >
                <span v-if="submitting" class="loading loading-spinner loading-xs mr-2"></span>
                {{ submitting ? '保存中...' : (editingId ? '确认修改' : '确认添加') }}
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
