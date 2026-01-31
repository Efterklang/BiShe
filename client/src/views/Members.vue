<script setup>
import { ref, onMounted, computed } from 'vue';
import MarkdownIt from "markdown-it";
import { getMembers, createMember, updateMemberBalance, deleteMember } from '../api/members';
import { listOrders } from "../api/orders";
import { generateMemberProfile } from "../api/ai";
import Avatar from '../components/Avatar.vue';
import MemberLevel from '../components/MemberLevel.vue';
import { Plus, X, UserPlus, Users, Wallet, Ban, ReceiptText, Sparkles } from 'lucide-vue-next';

const members = ref([]);
const loading = ref(true);
const createModalRef = ref(null);
const submitting = ref(false);
const selectedMember = ref(null);
const balanceModalRef = ref(null);
const balanceSubmitting = ref(false);
const balanceValue = ref(0);
const deactivateModalRef = ref(null);
const deactivateSubmitting = ref(false);
const ordersModalRef = ref(null);
const ordersLoading = ref(false);
const orders = ref([]);
const ordersPage = ref(1);
const ordersTotal = ref(0);
const aiModalRef = ref(null);
const aiLoading = ref(false);
const aiProfile = ref("");
const displayedProfile = ref("");
const md = new MarkdownIt();
const renderedProfile = computed(() => md.render(displayedProfile.value));

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

// 根据 referrer_id 获取推荐人信息
const getReferrer = (referrerId) => {
  if (!referrerId) return null;
  return members.value.find(m => m.id === referrerId || m.ID === referrerId);
};

const openBalanceModal = (member) => {
  selectedMember.value = member;
  balanceValue.value = Number(member.balance ?? member.Balance ?? 0);
  balanceModalRef.value?.showModal();
};

const submitBalance = async () => {
  if (!selectedMember.value) return;
  balanceSubmitting.value = true;
  try {
    await updateMemberBalance(selectedMember.value.id, Number(balanceValue.value));
    balanceModalRef.value?.close();
    await fetchMembers();
    alert('余额已更新');
  } catch (error) {
    alert('更新失败: ' + (error.message || '未知错误'));
  } finally {
    balanceSubmitting.value = false;
  }
};

const openDeactivateModal = (member) => {
  selectedMember.value = member;
  deactivateModalRef.value?.showModal();
};

const submitDeactivate = async () => {
  if (!selectedMember.value) return;
  deactivateSubmitting.value = true;
  try {
    await deleteMember(selectedMember.value.id);
    deactivateModalRef.value?.close();
    await fetchMembers();
    alert('会员已注销');
  } catch (error) {
    alert('操作失败: ' + (error.message || '未知错误'));
  } finally {
    deactivateSubmitting.value = false;
  }
};

const openOrdersModal = async (member) => {
  selectedMember.value = member;
  ordersPage.value = 1;
  orders.value = [];
  ordersTotal.value = 0;
  ordersModalRef.value?.showModal();
  await fetchOrders();
};

const fetchOrders = async () => {
  if (!selectedMember.value) return;
  ordersLoading.value = true;
  try {
    const res = await listOrders({
      member_id: selectedMember.value.id,
      page: ordersPage.value,
      page_size: 20
    });
    orders.value = res.orders || [];
    ordersTotal.value = res.total || 0;
  } catch (error) {
    orders.value = [];
    ordersTotal.value = 0;
  } finally {
    ordersLoading.value = false;
  }
};

const openAIModal = async (member) => {
  selectedMember.value = member;
  aiProfile.value = "";
  displayedProfile.value = "";
  aiModalRef.value?.showModal();

  aiLoading.value = true;
  try {
    const data = await generateMemberProfile(member.id);
    aiProfile.value = data.profile || data.report || "";
    displayedProfile.value = aiProfile.value;
  } catch (error) {
    aiProfile.value = "⚠️ 获取分析失败： " + (error.message || "未知错误");
    displayedProfile.value = aiProfile.value;
  } finally {
    aiLoading.value = false;
  }
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
              <th class="px-6 py-3 font-medium">推荐人</th>
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
              <td class="px-6 py-4">
                <div v-if="getReferrer(member.referrer_id || member.ReferrerID)" class="flex items-center gap-2">
                  <Avatar :name="getReferrer(member.referrer_id || member.ReferrerID).name" size="sm" />
                  <span class="text-sm text-base-content">{{ getReferrer(member.referrer_id || member.ReferrerID).name
                  }}</span>
                </div>
                <span v-else class="text-base-content/40">-</span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-1">
                  <div class="tooltip tooltip-bottom" data-tip="修改余额">
                    <button class="btn btn-ghost btn-xs btn-square" aria-label="修改余额" title="修改余额"
                      @click="openBalanceModal(member)">
                      <Wallet class="w-4 h-4" />
                    </button>
                  </div>
                  <div class="tooltip tooltip-bottom" data-tip="注销">
                    <button class="btn btn-ghost btn-xs btn-square text-warning" aria-label="注销" title="注销"
                      @click="openDeactivateModal(member)">
                      <Ban class="w-4 h-4" />
                    </button>
                  </div>
                  <div class="tooltip tooltip-bottom" data-tip="查看消费记录">
                    <button class="btn btn-ghost btn-xs btn-square" aria-label="查看消费记录" title="查看消费记录"
                      @click="openOrdersModal(member)">
                      <ReceiptText class="w-4 h-4" />
                    </button>
                  </div>
                  <div class="tooltip tooltip-bottom" data-tip="AI分析用户画像">
                    <button class="btn btn-ghost btn-xs btn-square" aria-label="AI分析用户画像" title="AI分析用户画像"
                      @click="openAIModal(member)">
                      <Sparkles class="w-4 h-4" />
                    </button>
                  </div>
                </div>
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

    <dialog ref="balanceModalRef" class="modal">
      <div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl max-w-md">
        <h3 class="font-semibold text-lg text-base-content">修改余额</h3>
        <p class="text-sm text-base-content/60 mt-1">
          {{ selectedMember?.name }}（#{{ selectedMember?.id }}）
        </p>
        <div class="form-control mt-4">
          <label class="label">
            <span class="label-text font-medium">余额（元）</span>
          </label>
          <input type="number" min="0" step="0.01" v-model="balanceValue"
            class="input input-bordered w-full bg-base-100" />
        </div>
        <div class="modal-action">
          <button class="btn btn-ghost" @click="balanceModalRef?.close()" :disabled="balanceSubmitting">取消</button>
          <button class="btn btn-primary" @click="submitBalance" :disabled="balanceSubmitting">
            保存
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog ref="deactivateModalRef" class="modal">
      <div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl max-w-md">
        <h3 class="font-semibold text-lg text-base-content text-warning">注销会员</h3>
        <p class="text-sm text-base-content/60 mt-2">
          确认注销 {{ selectedMember?.name }}（#{{ selectedMember?.id }}）？
        </p>
        <div class="modal-action">
          <button class="btn btn-ghost" @click="deactivateModalRef?.close()"
            :disabled="deactivateSubmitting">取消</button>
          <button class="btn btn-warning" @click="submitDeactivate" :disabled="deactivateSubmitting">
            确认注销
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog ref="ordersModalRef" class="modal">
      <div class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl w-11/12 max-w-4xl">
        <div class="flex justify-between items-center">
          <div>
            <h3 class="font-semibold text-lg text-base-content">消费记录</h3>
            <p class="text-sm text-base-content/60 mt-1">
              {{ selectedMember?.name }}（#{{ selectedMember?.id }}）
            </p>
          </div>
          <button class="btn btn-ghost btn-sm" @click="ordersModalRef?.close()">关闭</button>
        </div>

        <div class="mt-4">
          <div v-if="ordersLoading" class="flex items-center justify-center py-10">
            <span class="loading loading-spinner loading-lg"></span>
          </div>
          <div v-else class="overflow-x-auto">
            <table class="table w-full">
              <thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
                <tr>
                  <th class="px-4 py-3 font-medium">时间</th>
                  <th class="px-4 py-3 font-medium">类型</th>
                  <th class="px-4 py-3 font-medium">支付金额</th>
                  <th class="px-4 py-3 font-medium">佣金</th>
                </tr>
              </thead>
              <tbody class="text-sm divide-y divide-base-200">
                <tr v-if="orders.length === 0">
                  <td colspan="4" class="px-4 py-10 text-center text-base-content/50">
                    暂无消费记录
                  </td>
                </tr>
                <tr v-else v-for="o in orders" :key="o.id">
                  <td class="px-4 py-3 font-mono text-xs text-base-content/70">
                    {{ (o.created_at || o.createdAt || '').toString().slice(0, 19).replace('T', ' ') }}
                  </td>
                  <td class="px-4 py-3">
                    <span class="badge badge-outline">
                      {{ o.order_type || o.orderType }}
                    </span>
                  </td>
                  <td class="px-4 py-3 font-mono text-success">
                    ¥{{ Number(o.paid_amount || o.paidAmount || 0).toFixed(2) }}
                  </td>
                  <td class="px-4 py-3 font-mono text-base-content">
                    ¥{{ Number(o.commission_amount || o.commissionAmount || 0).toFixed(2) }}
                  </td>
                </tr>
              </tbody>
            </table>
            <div class="mt-3 text-xs text-base-content/50 flex justify-between items-center" v-if="ordersTotal > 0">
              <span>共 {{ ordersTotal }} 条</span>
              <div class="join">
                <button class="btn btn-xs join-item" :disabled="ordersPage <= 1 || ordersLoading"
                  @click="ordersPage--; fetchOrders()">上一页</button>
                <button class="btn btn-xs join-item" :disabled="ordersLoading"
                  @click="ordersPage++; fetchOrders()">下一页</button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog ref="aiModalRef" class="modal">
      <div
        class="modal-box bg-base-100 border border-base-300 shadow-2xl rounded-xl w-11/12 max-w-4xl h-[80vh] flex flex-col">
        <div class="flex justify-between items-center mb-3 pb-2 border-b border-base-200">
          <div>
            <h3 class="font-semibold text-lg text-base-content">AI 用户画像</h3>
            <p class="text-sm text-base-content/60 mt-1">
              {{ selectedMember?.name }}（#{{ selectedMember?.id }}）
            </p>
          </div>
          <button class="btn btn-ghost btn-sm" @click="aiModalRef?.close()">关闭</button>
        </div>
        <div class="flex-1 overflow-y-auto p-4 bg-base-200/30 rounded-xl text-base-content">
          <div v-if="aiLoading" class="flex flex-col items-center justify-center h-full gap-4">
            <span class="loading loading-dots loading-lg text-primary"></span>
            <p class="text-base-content/60">正在生成画像...</p>
          </div>
          <div v-else class="markdown-body">
            <div v-html="renderedProfile"></div>
          </div>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>