<script setup>
import { ref, onMounted, computed } from "vue";
import {
    CreditCard,
    TrendingUp,
    TrendingDown,
    UserPlus,
    Activity,
    AlertCircle,
    Trophy,
    Package,
    LayoutList
} from 'lucide-vue-next';

import Avatar from "../components/Avatar.vue";
import MemberLevel from "../components/MemberLevel.vue";

import {
    getDashboardStats,
    getFissionRanking,
    getServiceRanking,
    getProductSales,
} from "../api/dashboard";
import RevenueTrend from "../components/RevenueTrend.vue";

const stats = ref({
    dailyRevenue: 0,
    revenueGrowth: 0,
    newMembers: 0,
    activeTechs: 0,
    occupancyRate: 0,
});

const fissionRanking = ref([]);
const serviceRanking = ref([]);
const productSales = ref({
    topProducts: [],
    totalRevenue: 0,
    totalSales: 0,
    lowStockCount: 0,
});
const loading = ref(true);
const trendPeriod = ref(30); // 默认30天

// 时间段选项
const periodOptions = [
    { label: "7天", value: 7 },
    { label: "30天", value: 30 },
    { label: "90天", value: 90 },
];

const fetchData = async () => {
    loading.value = true;
    try {
        const [statsRes, rankingRes, serviceRes, productRes] =
            await Promise.all([
                getDashboardStats().catch((err) =>
                    console.warn("Stats API failed", err),
                ),
                getFissionRanking().catch((err) =>
                    console.warn("Ranking API failed", err),
                ),
                getServiceRanking().catch((err) =>
                    console.warn("Service API failed", err),
                ),
                getProductSales({ days: 30 }).catch((err) =>
                    console.warn("Product sales API failed", err),
                ),
            ]);

        if (statsRes) {
            stats.value = { ...stats.value, ...statsRes };
        }

        if (rankingRes) {
            fissionRanking.value = rankingRes;
        }

        if (serviceRes) {
            serviceRanking.value = serviceRes;
        }

        if (productRes) {
            productSales.value = productRes;
        }
    } catch (error) {
        console.error("Failed to load dashboard data:", error);
    } finally {
        loading.value = false;
    }
};

const changeTrendPeriod = (days) => {
    trendPeriod.value = days;
};

onMounted(fetchData);

// 计算服务排行的最大值，用于进度条宽度
const getBarWidth = (count) => {
    if (serviceRanking.value.length === 0) return "30%";
    const max = Math.max(...serviceRanking.value.map((s) => s.order_count));
    if (max === 0) return "30%";
    return `${Math.max((count / max) * 100, 10)}%`;
};

// 计算商品销售的最大值，用于进度条宽度
const getProductBarWidth = (count) => {
    if (productSales.value.topProducts.length === 0) return "30%";
    const max = Math.max(
        ...productSales.value.topProducts.map((p) => p.sales_count),
    );
    if (max === 0) return "30%";
    return `${Math.max((count / max) * 100, 10)}%`;
};

// 格式化数字
const formatNumber = (num) => {
    return typeof num === "number" ? num.toFixed(2) : "0.00";
};

// 格式化日期显示
const formatDateLabel = computed(() => {
    if (trendPeriod.value === 7) return "近7天";
    if (trendPeriod.value === 30) return "近30天";
    if (trendPeriod.value === 90) return "近90天";
    return `近${trendPeriod.value}天`;
});
</script>

<template>
    <div class="max-w-7xl mx-auto space-y-8">
        <!-- Header -->
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-base-content">
                经营概览
            </h1>
            <p class="mt-2 text-base-content/60">
                实时监控店铺运营数据，掌握核心业务指标。
            </p>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <!-- Stat Card 1 -->
            <div class="card bg-base-100 border border-base-300 shadow-sm">
                <div class="card-body">
                    <div class="stat">
                        <div class="stat-figure text-success">
                            <CreditCard class="w-8 h-8" />
                        </div>
                        <div class="stat-title">今日营收</div>
                        <div class="stat-value">¥{{ formatNumber(stats.dailyRevenue) }}</div>
                        <div class="stat-desc flex items-center">
                            <span :class="stats.revenueGrowth >= 0
                                ? 'text-success'
                                : 'text-error'
                                " class="flex items-center mr-1">
                                <TrendingUp v-if="stats.revenueGrowth >= 0" class="w-4 h-4 mr-0.5" />
                                <TrendingDown v-else class="w-4 h-4 mr-0.5" />
                                {{ formatNumber(Math.abs(stats.revenueGrowth)) }}%
                            </span>
                            较昨日
                        </div>
                    </div>
                </div>
            </div>

            <!-- Stat Card 2 -->
            <div class="card bg-base-100 border border-base-300 shadow-sm">
                <div class="card-body">
                    <div class="stat">
                        <div class="stat-figure text-info">
                            <UserPlus class="w-8 h-8" />
                        </div>
                        <div class="stat-title">新增会员</div>
                        <div class="stat-value">{{ stats.newMembers }}</div>
                        <div class="stat-desc">本月累计: 128</div>
                    </div>
                </div>
            </div>

            <!-- Stat Card 3 -->
            <div class="card bg-base-100 border border-base-300 shadow-sm">
                <div class="card-body">
                    <div class="stat">
                        <div class="stat-figure text-secondary">
                            <Activity class="w-8 h-8" />
                        </div>
                        <div class="stat-title">技师负载率</div>
                        <div class="stat-value">{{ formatNumber(stats.occupancyRate) }}%</div>
                        <div class="stat-desc">活跃技师: {{ stats.activeTechs }}</div>
                    </div>
                </div>
            </div>

            <!-- Stat Card 4 -->
            <div class="card bg-base-100 border border-base-300 shadow-sm">
                <div class="card-body">
                    <div class="stat">
                        <div class="stat-figure text-warning">
                            <AlertCircle class="w-8 h-8" />
                        </div>
                        <div class="stat-title">待处理预约</div>
                        <div class="stat-value">8</div>
                        <div class="stat-desc text-warning">需要关注</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Revenue Trend - Full Width -->
        <div class="card bg-base-100 border border-base-300 shadow-sm">
            <div class="card-body">
                <div class="flex justify-between items-center mb-6">
                    <h3 class="card-title">
                        {{ formatDateLabel }}营收趋势
                    </h3>
                    <div class="tabs tabs-boxed tabs-sm">
                        <button v-for="option in periodOptions" :key="option.value" class="tab"
                            :class="{ 'tab-active': trendPeriod === option.value }"
                            @click="changeTrendPeriod(option.value)">
                            {{ option.label }}
                        </button>
                    </div>
                </div>
                <div class="h-80">
                    <RevenueTrend :days="trendPeriod" />
                </div>
            </div>
        </div>

        <!-- Service Ranking and Product Sales -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Service Ranking -->
            <div class="card bg-base-100 border border-base-300 shadow-sm">
                <div class="card-body">
                    <h3 class="card-title">
                        热门项目排行
                    </h3>
                    <div v-if="loading" class="flex items-center justify-center py-12">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>
                    <div v-else-if="serviceRanking.length === 0"
                        class="flex flex-col items-center justify-center py-12 text-base-content/60 gap-3">
                        <LayoutList class="w-10 h-10 text-base-content/20" />
                        <p>暂无热门项目数据</p>
                    </div>
                    <div v-else class="space-y-6">
                        <div v-for="(service, index) in serviceRanking.slice(0, 5)" :key="service.service_id"
                            class="relative">
                            <div class="flex justify-between mb-2 text-sm">
                                <span class="font-medium text-base-content">{{ index + 1 }}.
                                    {{ service.service_name }}</span>
                                <span class="text-base-content/60">{{ service.order_count }}单 / ¥{{
                                    formatNumber(service.total_revenue)
                                    }}</span>
                            </div>
                            <div class="w-full bg-base-200 rounded-full h-2">
                                <div class="bg-primary h-2 rounded-full transition-all" :style="{
                                    width: getBarWidth(service.order_count),
                                    opacity: 1 - index * 0.15,
                                }"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Product Sales Overview -->
            <div class="card bg-base-100 border border-base-300 shadow-sm">
                <div class="card-body">
                    <h3 class="card-title">
                        实体商品销售概览
                    </h3>
                    <div v-if="loading" class="flex items-center justify-center py-12">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>
                    <div v-else>
                        <!-- 统计卡片 -->
                        <div class="grid grid-cols-3 gap-4 mb-6">
                            <div class="stat bg-base-200/50 rounded-lg p-3 text-center">
                                <div class="stat-title text-xs text-base-content/60 font-medium">
                                    总销售额
                                </div>
                                <div class="stat-value text-lg font-bold text-success">
                                    ¥{{ formatNumber(productSales.totalRevenue) }}
                                </div>
                            </div>
                            <div class="stat bg-base-200/50 rounded-lg p-3 text-center">
                                <div class="stat-title text-xs text-base-content/60 font-medium">
                                    销售订单
                                </div>
                                <div class="stat-value text-lg font-bold text-info">
                                    {{ productSales.totalSales }}
                                </div>
                            </div>
                            <div class="stat bg-base-200/50 rounded-lg p-3 text-center">
                                <div class="stat-title text-xs text-base-content/60 font-medium">
                                    库存预警
                                </div>
                                <div class="stat-value text-lg font-bold text-warning">
                                    {{ productSales.lowStockCount }}
                                </div>
                            </div>
                        </div>

                        <!-- 热销商品排行 -->
                        <div v-if="productSales.topProducts.length === 0"
                            class="flex flex-col items-center justify-center py-8 text-base-content/60 gap-3">
                            <Package class="w-10 h-10 text-base-content/20" />
                            <p>暂无商品销售数据</p>
                        </div>
                        <div v-else class="space-y-4">
                            <div v-for="(product, index) in productSales.topProducts" :key="product.product_id"
                                class="relative">
                                <div class="flex justify-between mb-2 text-sm">
                                    <span class="font-medium text-base-content">{{ index + 1 }}.
                                        {{ product.product_name }}</span>
                                    <span class="text-base-content/60">{{ product.sales_count }}件 / ¥{{
                                        formatNumber(product.total_revenue)
                                        }}</span>
                                </div>
                                <div class="w-full bg-base-200 rounded-full h-2">
                                    <div class="bg-secondary h-2 rounded-full transition-all" :style="{
                                        width: getProductBarWidth(
                                            product.sales_count,
                                        ),
                                        opacity: 1 - index * 0.15,
                                    }"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Fission Ranking Table -->
        <div class="card bg-base-100 border border-base-300 shadow-sm overflow-hidden">
            <div class="card-body">
                <div class="flex justify-between items-center">
                    <h3 class="card-title flex items-center gap-2">
                        <Trophy class="w-6 h-6 text-warning" />
                        裂变达人榜
                    </h3>
                    <button class="btn btn-ghost btn-sm">
                        查看全部
                    </button>
                </div>
                <div class="overflow-x-auto">
                    <table class="table w-full">
                        <thead class="bg-base-200/50 text-base-content/60 uppercase text-xs">
                            <tr>
                                <th class="px-6 py-4 font-semibold">排名</th>
                                <th class="px-6 py-4 font-semibold">会员姓名</th>
                                <th class="px-6 py-4 font-semibold">邀请人数</th>
                                <th class="px-6 py-4 font-semibold">累计佣金</th>
                                <th class="px-6 py-4 font-semibold">等级</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-base-200">
                            <tr v-for="(item, index) in fissionRanking" :key="item.id"
                                class="hover:bg-base-200/50 transition-colors">
                                <td class="px-6 py-4">
                                    <span
                                        class="inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold"
                                        :class="{
                                            'bg-warning/20 text-warning':
                                                index === 0,
                                            'bg-base-300 text-base-content':
                                                index === 1,
                                            'bg-error/20 text-error': index === 2,
                                            'text-base-content/60': index > 2,
                                        }">
                                        {{ index + 1 }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 font-medium text-base-content">
                                    <div class="flex items-center gap-3">
                                        <Avatar :name="item.name" size="sm" />
                                        {{ item.name }}
                                    </div>
                                </td>
                                <td class="px-6 py-4 text-base-content/60">
                                    {{ item.inviteCount }} 人
                                </td>
                                <td class="px-6 py-4 font-medium text-success">
                                    +¥{{ item.totalCommission }}
                                </td>
                                <td class="px-6 py-4">
                                    <MemberLevel :level="item.level" size="sm" />
                                </td>
                            </tr>
                            <tr v-if="fissionRanking.length === 0">
                                <td colspan="5" class="px-6 py-12 text-center text-base-content/60">
                                    <div class="flex flex-col items-center justify-center gap-3">
                                        <Trophy class="w-10 h-10 text-base-content/20" />
                                        <p>暂无裂变排行数据</p>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>