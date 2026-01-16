<script setup>
import { ref, onMounted } from 'vue';
import { getAppointments } from '../api/appointments';
import { getTechnicians } from '../api/technicians';
import { getServices } from '../api/services';
import { getMembers } from '../api/members';
import Avatar from '../components/Avatar.vue';

const appointments = ref([]);
const technicians = ref([]);
const services = ref([]);
const members = ref([]);
const loading = ref(true);

// Fetch data
const fetchData = async () => {
  loading.value = true;
  try {
    const [apptRes, techRes, serviceRes, memberRes] = await Promise.allSettled([
      getAppointments({ status: 'completed' }), // Filter for completed orders
      getTechnicians(),
      getServices(),
      getMembers()
    ]);

    if (apptRes.status === 'fulfilled') appointments.value = apptRes.value || [];
    if (techRes.status === 'fulfilled') technicians.value = techRes.value || [];
    if (serviceRes.status === 'fulfilled') services.value = serviceRes.value || [];
    if (memberRes.status === 'fulfilled') members.value = memberRes.value || [];

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

</script>

<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between mb-10 gap-4">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-base-content">历史订单</h1>
      </div>
    </div>

    <!-- History Table -->
    <div class="bg-base-100 rounded-xl border border-base-300 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead class="bg-base-200 text-base-content/60 uppercase text-xs">
            <tr>
              <th class="px-6 py-3 font-medium">订单号</th>
              <th class="px-6 py-3 font-medium">会员</th>
              <th class="px-6 py-3 font-medium">服务项目</th>
              <th class="px-6 py-3 font-medium">技师</th>
              <th class="px-6 py-3 font-medium">完成时间</th>
              <th class="px-6 py-3 font-medium">实收金额</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-base-200">
            <tr v-if="loading">
              <td colspan="7" class="px-6 py-12 text-center">
                <span class="loading loading-spinner loading-lg text-base-content/40"></span>
              </td>
            </tr>
            <tr v-else-if="appointments.length === 0">
              <td colspan="7" class="px-6 py-12 text-center text-base-content/60">暂无历史订单</td>
            </tr>
            <tr v-else v-for="appt in appointments" :key="appt.id" class="hover:bg-base-200/50 transition-colors">
              <td class="px-6 py-4 text-base-content/60 font-mono text-xs">#{{ appt.id }}</td>
              <td class="px-6 py-4 font-medium text-base-content">
                {{ getMemberName(appt.member_id || appt.MemberID) }}
              </td>
              <td class="px-6 py-4 text-base-content/80">
                {{ getServiceName(appt.service_id || appt.ServiceID) }}
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <Avatar :name="getTechName(appt.tech_id)" :size="15" />
                  <span class="text-base-content/80">{{ getTechName(appt.tech_id) }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-base-content/60 text-xs">
                {{ formatDate(appt.end_time || appt.EndTime) }}
              </td>
              <td class="px-6 py-4 font-medium text-success">¥{{ appt.actual_price || appt.ActualPrice }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>