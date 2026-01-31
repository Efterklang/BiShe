<script setup>
import { ref, onMounted } from 'vue';
import { getAppointments } from '../api/appointments';
import { getTechnicians } from '../api/technicians';
import { getServices } from '../api/services';
import { getMembers } from '../api/members';
import { getProducts, getAllInventoryLogs } from '../api/products';
import Avatar from '../components/Avatar.vue';
import { History, Package, User, Clock, ShoppingBag } from 'lucide-vue-next';

const activeTab = ref('service'); // 'service' | 'product'
const appointments = ref([]);
const inventoryLogs = ref([]);
const technicians = ref([]);
const services = ref([]);
const members = ref([]);
const products = ref([]);
const loading = ref(true);

// Fetch data
const fetchData = async () => {
  loading.value = true;
  try {
    const [apptRes, techRes, serviceRes, memberRes, productRes] = await Promise.allSettled([
      getAppointments({ status: 'completed' }), // Filter for completed orders
      getTechnicians(),
      getServices(),
      getMembers(),
      getProducts()
    ]);

    if (apptRes.status === 'fulfilled') appointments.value = apptRes.value || [];
    if (techRes.status === 'fulfilled') technicians.value = techRes.value || [];
    if (serviceRes.status === 'fulfilled') services.value = serviceRes.value || [];
    if (memberRes.status === 'fulfilled') members.value = memberRes.value || [];
    if (productRes.status === 'fulfilled') products.value = productRes.value || [];

    // Fetch product sales (inventory logs with action_type='sale')
    const logsRes = await getAllInventoryLogs({ action_type: 'sale' });
    if (logsRes && logsRes.logs) {
      inventoryLogs.value = logsRes.logs;
    } else if (Array.isArray(logsRes)) {
      inventoryLogs.value = logsRes;
    }

  } catch (error) {
    console.error("Error fetching history:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchData);

// Helpers
const formatDate = (dateStr) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const getTechName = (id) => technicians.value.find(t => t.id === id)?.name || `技师#${id}`;
const getServiceName = (id) => services.value.find(s => s.id === id)?.name || `项目#${id}`;
const getMemberName = (id) => members.value.find(m => m.id === id)?.name || `会员#${id}`;
const getProductName = (id) => products.value.find(p => p.id === id)?.name || `商品#${id}`;

</script>

<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
      <div>
        <h1 class="text-3xl font-bold text-base-content flex items-center gap-3">
          <History class="w-8 h-8 text-primary" />
          历史订单
        </h1>
        <p class="text-base-content/60 mt-2">查看所有已完成的服务订单记录。</p>
      </div>
    </div>

    <!-- Tab Navigation -->
    <div class="tabs tabs-boxed bg-base-200 mb-4">
      <button
        class="tab"
        :class="{ 'tab-active': activeTab === 'service' }"
        @click="activeTab = 'service'"
      >
        <ShoppingBag class="w-4 h-4 mr-2" />
        服务订单 ({{ appointments.length }})
      </button>
      <button
        class="tab"
        :class="{ 'tab-active': activeTab === 'product' }"
        @click="activeTab = 'product'"
      >
        <Package class="w-4 h-4 mr-2" />
        商品订单 ({{ inventoryLogs.length }})
      </button>
    </div>

    <!-- History Table -->
    <div class="card bg-base-100 border border-base-300 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
            <tr v-if="activeTab === 'service'">
              <th class="px-6 py-4 font-semibold">订单号</th>
              <th class="px-6 py-4 font-semibold">会员</th>
              <th class="px-6 py-4 font-semibold">服务项目</th>
              <th class="px-6 py-4 font-semibold">技师</th>
              <th class="px-6 py-4 font-semibold">完成时间</th>
              <th class="px-6 py-4 font-semibold">实收金额</th>
            </tr>
            <tr v-else>
              <th class="px-6 py-4 font-semibold">记录号</th>
              <th class="px-6 py-4 font-semibold">购买会员</th>
              <th class="px-6 py-4 font-semibold">商品名称</th>
              <th class="px-6 py-4 font-semibold">完成时间</th>
              <th class="px-6 py-4 font-semibold">购买数量</th>
              <th class="px-6 py-4 font-semibold">实付金额</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-base-200">
            <tr v-if="loading">
              <td colspan="6" class="px-6 py-24 text-center">
                <span class="loading loading-spinner loading-lg text-primary"></span>
              </td>
            </tr>
            <tr v-else-if="appointments.length === 0">
              <td colspan="6" class="px-6 py-24 text-center">
                <div class="flex flex-col items-center justify-center gap-4">
                  <div class="w-16 h-16 bg-base-200 rounded-full flex items-center justify-center">
                    <History class="w-8 h-8 text-base-content/30" />
                  </div>
                  <div class="text-center">
                    <h3 class="text-lg font-medium text-base-content">暂无历史订单</h3>
                    <p class="text-base-content/60 mt-1">目前没有任何已完成的订单记录。</p>
                  </div>
                </div>
              </td>
            </tr>
            <!-- 服务订单列表 -->
            <template v-if="activeTab === 'service'">
              <tr v-for="appt in appointments" :key="appt.id" class="hover:bg-base-200/50 transition-colors">
                <td class="px-6 py-4 text-base-content/60 font-mono text-xs">
                  #{{ appt.id }}
                </td>
                <td class="px-6 py-4 font-medium text-base-content">
                  {{ getMemberName(appt.member_id || appt.MemberID) }}
                </td>
                <td class="px-6 py-4 text-base-content/80">
                  <div class="badge badge-ghost gap-1">
                    {{ getServiceName(appt.service_id || appt.ServiceID) }}
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <Avatar :name="appt.technician?.name || getTechName(appt.tech_id)" :src="appt.technician?.avatar_url" :size="32" class="ring-1 ring-base-300" />
                    <span class="text-sm font-medium">{{ appt.technician?.name || getTechName(appt.tech_id) }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-base-content/60 text-sm">
                  {{ formatDate(appt.end_time || appt.EndTime) }}
                </td>
                <td class="px-6 py-4">
                  <div class="flex flex-col gap-1">
                    <span class="font-medium text-success">
                      ¥{{ (appt.actual_price || appt.ActualPrice).toFixed(2) }}
                    </span>
                    <span v-if="appt.commission_amount > 0" class="text-xs text-warning">
                      佣金 ¥{{ appt.commission_amount.toFixed(2) }} → {{ appt.commission_to?.name || '未知' }}
                    </span>
                  </div>
                </td>
              </tr>
            </template>

            <!-- 商品订单列表 -->
            <template v-else>
              <tr v-for="log in inventoryLogs" :key="log.id" class="hover:bg-base-200/50 transition-colors">
                <td class="px-6 py-4 text-base-content/60 font-mono text-xs">
                  #{{ log.id }}
                </td>
                <td class="px-6 py-4 font-medium text-base-content">
                  <div v-if="log.member" class="flex items-center gap-2">
                    <User class="w-4 h-4 text-base-content/40" />
                    {{ log.member.name }}
                  </div>
                  <span v-else class="text-base-content/40">-</span>
                </td>
                <td class="px-6 py-4 text-base-content/80">
                  <div class="badge badge-ghost gap-1">
                    <Package class="w-3 h-3" />
                    {{ log.product?.name || '未知商品' }}
                  </div>
                </td>
                <td class="px-6 py-4 text-base-content/60 text-sm">
                  {{ formatDate(log.created_at) }}
                </td>
                <td class="px-6 py-4">
                  <span class="font-medium">{{ Math.abs(log.change_amount) }}</span>
                  <span class="text-xs text-base-content/60">件</span>
                </td>
                <td class="px-6 py-4">
                  <div class="flex flex-col gap-1">
                    <span class="font-medium text-success">
                      ¥{{ (log.sale_amount || (Math.abs(log.change_amount) * (log.product?.retail_price || 0))).toFixed(2) }}
                    </span>
                    <span v-if="log.commission_amount > 0" class="text-xs text-warning">
                      佣金 ¥{{ log.commission_amount.toFixed(2) }} → {{ log.commission_to?.name || '未知' }}
                    </span>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>