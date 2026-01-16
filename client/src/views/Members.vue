<script setup>
import { ref, onMounted } from 'vue';
import { getMembers, createMember } from '../api/members';
import Avatar from '../components/Avatar.vue';
import MemberLevel from '../components/MemberLevel.vue';

const members = ref([]);
const loading = ref(true);
const showModal = ref(false);
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

const handleCreateMember = async () => {
  submitting.value = true;
  try {
    await createMember({
      name: formData.value.name,
      phone: formData.value.phone,
      invitation_code: formData.value.invitation_code || undefined
    });
    showModal.value = false;
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
  <div class="max-w-7xl mx-auto">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between mb-10 gap-4">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-base-content">会员管理</h1>
        <p class="mt-2 text-base-content/60">
          查看会员列表、等级及消费记录，管理客户关系。
        </p>
      </div>
      <button @click="showModal = true" class="btn btn-primary btn-sm">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"
          class="w-4 h-4 mr-2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        注册会员
      </button>
    </div>

    <!-- Members Table -->
    <div class="bg-base-100 rounded-box border border-base-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table table-zebra w-full">
          <thead class="bg-base-200 text-base-content/70 uppercase text-xs">
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
          <tbody class="text-sm">
            <tr v-if="loading">
              <td colspan="8" class="px-6 py-12 text-center">
                <span class="loading loading-spinner loading-lg text-base-content/30"></span>
              </td>
            </tr>
            <tr v-else-if="members.length === 0">
              <td colspan="8" class="px-6 py-12 text-center text-base-content/50">暂无会员数据</td>
            </tr>
            <tr v-else v-for="member in members" :key="member.id" class="hover:bg-base-200/50 transition-colors">
              <td class="px-6 py-4 text-base-content/50 font-mono text-xs">#{{ member.id }}</td>
              <td class="px-6 py-4 font-medium text-base-content">
                <div class="flex items-center gap-3">
                  <Avatar :name="member.name" size="sm" />
                  {{ member.name }}
                </div>
              </td>
              <td class="px-6 py-4 text-base-content/80">{{ member.phone }}</td>
              <td class="px-6 py-4">
                <MemberLevel :level="member.level || member.Level" />
              </td>
              <td class="px-6 py-4 font-mono text-base-content">¥{{ member.yearly_total_consumption ||
                member.YearlyTotalConsumption || 0 }}</td>
              <td class="px-6 py-4 font-mono text-success">¥{{ member.balance || member.Balance || 0 }}</td>
              <td class="px-6 py-4">
                <code class="badge badge-neutral badge-outline font-mono text-xs">
                  {{ member.invitation_code || member.InvitationCode }}
                </code>
              </td>
              <td class="px-6 py-4 text-base-content/50">{{ member.referrer_id || member.ReferrerID || '-' }}</td>
              <td class="px-6 py-4 text-right">
                <button class="btn btn-ghost btn-xs">详情</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <dialog class="modal" :class="{ 'modal-open': showModal }">
      <div class="modal-box bg-base-100 border border-base-200 shadow-2xl rounded-box p-0 overflow-hidden max-w-md">
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center bg-base-200/30">
          <h3 class="font-semibold text-lg text-base-content">注册新会员</h3>
          <button @click="showModal = false" class="btn btn-ghost btn-sm btn-circle">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
              stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6">
          <form @submit.prevent="handleCreateMember" class="space-y-5">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">姓名</span>
              </label>
              <input type="text" v-model="formData.name" placeholder="请输入会员姓名" class="input input-bordered w-full"
                required />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">手机号</span>
              </label>
              <input type="tel" v-model="formData.phone" placeholder="请输入手机号" class="input input-bordered w-full"
                required />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">
                  邀请码 <span class="text-base-content/40 font-normal">(选填)</span>
                </span>
              </label>
              <input type="text" v-model="formData.invitation_code" placeholder="如有推荐人请填写"
                class="input input-bordered w-full" />
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
        <button @click="showModal = false">close</button>
      </form>
    </dialog>
  </div>
</template>