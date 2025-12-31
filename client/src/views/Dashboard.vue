<script setup>
import { ref, onMounted } from 'vue';
import { getDashboardStats, getFissionRanking } from '../api/dashboard';

const stats = ref({
  dailyRevenue: 0,
  newMembers: 0,
  activeTechs: 0,
  occupancyRate: 0
});

const fissionRanking = ref([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const [statsRes, rankingRes] = await Promise.all([
      getDashboardStats().catch(err => console.warn("Stats API failed", err)),
      getFissionRanking().catch(err => console.warn("Ranking API failed", err))
    ]);

    if (statsRes) {
        stats.value = { ...stats.value, ...statsRes };
    }

    if (rankingRes) {
        fissionRanking.value = rankingRes;
    }
  } catch (error) {
    console.error("Failed to load dashboard data:", error);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="max-w-7xl mx-auto space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-base-content">经营概览</h1>
      <p class="mt-2 text-base-content/60">
        实时监控店铺运营数据，掌握核心业务指标。
      </p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Stat Card 1 -->
      <div class="p-6 bg-base-100 rounded-xl border border-base-300 shadow-sm">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-base-content/60">今日营收</p>
          <span class="p-2 bg-success/10 text-success rounded-lg">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </span>
        </div>
        <div class="mt-4">
          <h3 class="text-3xl font-bold text-base-content">¥{{ stats.dailyRevenue }}</h3>
          <div class="flex items-center mt-1 text-sm">
            <span class="text-success font-medium flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 mr-0.5">
                <path fill-rule="evenodd" d="M12 7a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0V8.414l-4.293 4.293a1 1 0 01-1.414 0L8 10.414l-4.293 4.293a1 1 0 01-1.414-1.414l5-5a1 1 0 011.414 0L11 10.586 14.586 7H12z" clip-rule="evenodd" />
              </svg>
              12%
            </span>
            <span class="text-base-content/40 ml-2">较昨日</span>
          </div>
        </div>
      </div>

      <!-- Stat Card 2 -->
      <div class="p-6 bg-base-100 rounded-xl border border-base-300 shadow-sm">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-base-content/60">新增会员</p>
          <span class="p-2 bg-info/10 text-info rounded-lg">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 0110.374 21c-2.331 0-4.512-.645-6.374-1.766z" />
            </svg>
          </span>
        </div>
        <div class="mt-4">
          <h3 class="text-3xl font-bold text-base-content">{{ stats.newMembers }}</h3>
          <div class="flex items-center mt-1 text-sm">
            <span class="text-base-content/60">本月累计: <span class="font-medium text-base-content">128</span></span>
          </div>
        </div>
      </div>

      <!-- Stat Card 3 -->
      <div class="p-6 bg-base-100 rounded-xl border border-base-300 shadow-sm">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-base-content/60">技师负载率</p>
          <span class="p-2 bg-secondary/10 text-secondary rounded-lg">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
            </svg>
          </span>
        </div>
        <div class="mt-4">
          <h3 class="text-3xl font-bold text-base-content">{{ stats.occupancyRate }}%</h3>
          <div class="flex items-center mt-1 text-sm">
            <span class="text-base-content/60">活跃技师: <span class="font-medium text-base-content">{{ stats.activeTechs }}</span></span>
          </div>
        </div>
      </div>

      <!-- Stat Card 4 (Placeholder) -->
      <div class="p-6 bg-base-100 rounded-xl border border-base-300 shadow-sm">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-base-content/60">待处理预约</p>
          <span class="p-2 bg-warning/10 text-warning rounded-lg">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </span>
        </div>
        <div class="mt-4">
          <h3 class="text-3xl font-bold text-base-content">8</h3>
          <div class="flex items-center mt-1 text-sm">
            <span class="text-warning font-medium">需要关注</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Revenue Trend -->
      <div class="p-6 bg-base-100 rounded-xl border border-base-300 shadow-sm">
        <h3 class="text-lg font-semibold text-base-content mb-6">近30天营收趋势</h3>
        <div class="h-64 flex items-end justify-between gap-2">
          <!-- Simple CSS Bar Chart Mockup -->
          <div v-for="i in 15" :key="i"
               class="bg-primary w-full rounded-t-sm hover:opacity-80 transition-opacity"
               :style="{ height: `${Math.floor(Math.random() * 70 + 30)}%`, opacity: 0.3 + (i/25) }">
          </div>
        </div>
        <div class="flex justify-between text-xs text-base-content/40 mt-4 font-medium">
          <span>30天前</span>
          <span>今天</span>
        </div>
      </div>

      <!-- Service Ranking -->
      <div class="p-6 bg-base-100 rounded-xl border border-base-300 shadow-sm">
        <h3 class="text-lg font-semibold text-base-content mb-6">热门项目排行</h3>
        <div class="space-y-6">
          <div class="relative">
            <div class="flex justify-between mb-2 text-sm">
              <span class="font-medium text-base-content">1. 全身精油SPA</span>
              <span class="text-base-content/60">128单</span>
            </div>
            <div class="w-full bg-base-200 rounded-full h-2">
              <div class="bg-primary h-2 rounded-full" style="width: 85%"></div>
            </div>
          </div>
          <div class="relative">
            <div class="flex justify-between mb-2 text-sm">
              <span class="font-medium text-base-content">2. 中式推拿</span>
              <span class="text-base-content/60">96单</span>
            </div>
            <div class="w-full bg-base-200 rounded-full h-2">
              <div class="bg-primary h-2 rounded-full opacity-80" style="width: 65%"></div>
            </div>
          </div>
          <div class="relative">
            <div class="flex justify-between mb-2 text-sm">
              <span class="font-medium text-base-content">3. 足底按摩</span>
              <span class="text-base-content/60">85单</span>
            </div>
            <div class="w-full bg-base-200 rounded-full h-2">
              <div class="bg-primary h-2 rounded-full opacity-60" style="width: 55%"></div>
            </div>
          </div>
          <div class="relative">
            <div class="flex justify-between mb-2 text-sm">
              <span class="font-medium text-base-content">4. 艾灸护理</span>
              <span class="text-base-content/60">42单</span>
            </div>
            <div class="w-full bg-base-200 rounded-full h-2">
              <div class="bg-primary h-2 rounded-full opacity-40" style="width: 30%"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Fission Ranking Table -->
    <div class="bg-base-100 rounded-xl border border-base-300 shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-base-200 flex justify-between items-center">
        <h3 class="text-lg font-semibold text-base-content">🏆 裂变达人榜</h3>
        <button class="text-sm text-base-content/60 hover:text-base-content font-medium transition-colors">查看全部</button>
      </div>
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead class="bg-base-200 text-base-content/60 uppercase text-xs">
            <tr>
              <th class="px-6 py-3 font-medium">排名</th>
              <th class="px-6 py-3 font-medium">会员姓名</th>
              <th class="px-6 py-3 font-medium">邀请人数</th>
              <th class="px-6 py-3 font-medium">累计佣金</th>
              <th class="px-6 py-3 font-medium">等级</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-base-200">
            <tr v-for="(item, index) in fissionRanking" :key="item.id" class="hover:bg-base-200/50 transition-colors">
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold"
                  :class="{
                    'bg-warning/20 text-warning': index === 0,
                    'bg-base-300 text-base-content': index === 1,
                    'bg-error/20 text-error': index === 2,
                    'text-base-content/60': index > 2
                  }"
                >
                  {{ index + 1 }}
                </span>
              </td>
              <td class="px-6 py-4 font-medium text-base-content">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-base-300 flex items-center justify-center text-xs font-bold text-base-content/60">
                    {{ item.name ? item.name.charAt(0) : '?' }}
                  </div>
                  {{ item.name }}
                </div>
              </td>
              <td class="px-6 py-4 text-base-content/60">{{ item.inviteCount }} 人</td>
              <td class="px-6 py-4 font-medium text-success">+¥{{ item.totalCommission }}</td>
              <td class="px-6 py-4">
                <span class="badge badge-ghost badge-sm">
                  {{ item.level || '普通会员' }}
                </span>
              </td>
            </tr>
            <tr v-if="fissionRanking.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-base-content/60">暂无数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
