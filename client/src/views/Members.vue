<script setup>
import { ref, onMounted } from 'vue';
import { getMembers, createMember } from '../api/members';
import Avatar from '../components/Avatar.vue';
import MemberLevel from '../components/MemberLevel.vue';
import { Plus, X, UserPlus, Users } from 'lucide-vue-next';

const members = ref([]);
const loading = ref(true);
const createModalRef = ref(null);
const submitting = ref(false);
const formData = ref({
  name: '',
  phone: '',
  invitation_code: ''
});

const fetchMembers = async () => {
  loading.value = true;
  try {
    const res = await getMembers();
    members.value = res || [];
  } catch (error) {
    console.error("Failed to load members:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchMembers);

const openCreateModal = () => {
  formData.value = { name: '', phone: '', invitation_code: '' };
  createModalRef.value?.showModal();
};

const closeCreateModal = () => {
  createModalRef.value?.close();
};

const handleCreateMember = async () => {
  submitting.value = true;
  try {
    await createMember({
      name: formData.value.name,
      phone: formData.value.phone,
      invitation_code: formData.value.invitation_code || undefined
    });
    closeCreateModal();
    formData.value = { name: '', phone: '', invitation_code: '' };
    await fetchMembers();
    alert('会员注册成功');
  } catch (error) {
    alert('注册失败: ' + (error.message || '未知错误'));
  } finally {
    submitting.value = false;
  }
};

// 根据 member id 生成一致的随机背景色
const getAvatarBgColor = (memberId) => {
  const colors = [
    'bg-red-100',
    'bg-blue-100',
    'bg-green-100',
    'bg-yellow-100',
    'bg-purple-100',
    'bg-pink-100',
    'bg-indigo-100',
    'bg-cyan-100',
    'bg-orange-100',
    'bg-lime-100'
  ];
  const index = (memberId % colors.length);
  return colors[index];
};

// 获取对应字体颜色（与背景色匹配）
const getAvatarTextColor = (memberId) => {
  const textColors = [
    'text-red-700',
    'text-blue-700',
    'text-green-700',
    'text-yellow-700',
    'text-purple-700',
    'text-pink-700',
    'text-indigo-700',
    'text-cyan-700',
    'text-orange-700',
    'text-lime-700'
  ];
  const index = (memberId % textColors.length);
  return textColors[index];
};


</script>

<template>
  <div>
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-base-content">会员管理</h1>
        <p class="mt-1 text-base-content/60">
          查看会员列表、等级及消费记录，管理客户关系。
        </p>
      </div>
      <button @click="openCreateModal" class="btn btn-primary">
        <Plus class="w-4 h-4 mr-1" />
        注册会员
      </button>
    </div>

    <!-- Members Table -->
    <div class="bg-base-100 rounded-xl border border-base-300 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
            <tr>
              <th class="px-6 py-3 font-medium">ID</th>
              <th class="px-6 py-3 font-medium">姓名</th>
              <th class="px-6 py-3 font-medium">手机号</th>
              <th class="px-6 py-3 font-medium">等级</th>
              <th class="px-6 py-3 font-medium">年消费总额</th>
              <th class="px-6 py-3 font-medium">余额</th>
              <th class="px-6 py-3 font-medium">邀请码</th>
              <th class="px-6 py-3 font-medium">推荐人ID</th>
              <th class="px-6 py-3 font-medium text-right">操作</th>
            </tr>
          </thead>
          <tbody class="text-sm divide-y divide-base-200">
            <tr v-if="loading">
              <td colspan="9" class="px-6 py-12 text-center">
                <span class="loading loading-spinner loading-lg text-primary"></span>
              </td>
            </tr>
            <tr v-else-if="members.length === 0">
              <td colspan="9" class="px-6 py-16 text-center">
                <div class="flex flex-col items-center justify-center">
                  <div class="w-16 h-16 bg-base-200 rounded-full flex items-center justify-center mb-4">
                    <Users class="w-8 h-8 text-base-content/40" />
                  </div>
                  <h3 class="text-lg font-bold text-base-content">暂无会员数据</h3>
                  <p class="text-base-content/60 mt-1">点击右上角按钮注册新会员</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="member in members" :key="member.id" class="hover:bg-base-50/50 transition-colors">
              <td class="px-6 py-4 text-base-content/50 font-mono text-xs">#{{ member.id }}</td>
              <td class="px-6 py-4 font-medium text-base-content">
                <div class="flex items-center gap-3">
                  <Avatar :name="member.name" size="sm" />
                  {{ member.name }}
                </div>
              </td>
              <td class="px-6 py-4 text-base-content/80 font-mono">{{ member.phone }}</td>
              <td class="px-6 py-4">
                <MemberLevel :level="member.level || member.Level" />
              </td>
              <td class="px-6 py-4 font-mono text-base-content">¥{{ member.yearly_total_consumption ||
                member.YearlyTotalConsumption || 0 }}</td>
              <td class="px-6 py-4 font-mono text-success font-medium">¥{{ member.balance || member.Balance || 0 }}</td>
              <td class="px-6 py-4">
                <code class="badge badge-neutral badge-outline font-mono text-xs">
                  {{ member.invitation_code || member.InvitationCode }}
                </code>
              </td>
              <td class="px-6 py-4 text-base-content/50 font-mono text-xs">{{ member.referrer_id || member.ReferrerID ||
                '-' }}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn btn-ghost btn-xs">详情</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- Pagination or count could go here -->
      <div
        class="bg-base-50 px-6 py-3 border-t border-base-200 text-xs text-base-content/60 flex justify-between items-center"
        v-if="members.length > 0">
        <span>共 {{ members.length }} 位会员</span>
      </div>
    </div>

    <!-- Create Modal -->
    <dialog ref="createModalRef" class="modal">
      <div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl p-0 overflow-hidden max-w-md">
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/50">
          <h3 class="font-semibold text-lg text-base-content flex items-center gap-2">
            <UserPlus class="w-5 h-5 text-primary" />
            注册新会员
          </h3>
          <button @click="closeCreateModal" class="btn btn-ghost btn-sm btn-square text-base-content/60">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6">
          <form @submit.prevent="handleCreateMember" class="space-y-5">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">姓名</span>
              </label>
              <input type="text" v-model="formData.name" placeholder="请输入会员姓名"
                class="input input-bordered w-full bg-base-100" required />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">手机号</span>
              </label>
              <input type="tel" v-model="formData.phone" placeholder="请输入手机号"
                class="input input-bordered w-full bg-base-100" required />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">
                  邀请码 <span class="text-base-content/40 font-normal">(选填)</span>
                </span>
              </label>
              <input type="text" v-model="formData.invitation_code" placeholder="如有推荐人请填写"
                class="input input-bordered w-full bg-base-100" />
            </div>

            <div class="pt-2">
              <button type="submit" class="btn btn-primary w-full" :disabled="submitting">
                <span v-if="submitting" class="loading loading-spinner loading-xs"></span>
                {{ submitting ? '注册中...' : '确认注册' }}
              </button>
            </div>
          </form>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
        <button @click="closeCreateModal">close</button>
      </form>
    </dialog>
  </div>
</template>