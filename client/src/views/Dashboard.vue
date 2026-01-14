<script setup>
import { ref, onMounted, computed } from "vue";
import VChart from "vue-echarts";
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LineChart } from "echarts/charts";
import {
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    GridComponent,
} from "echarts/components";
import {
    getDashboardStats,
    getFissionRanking,
    getRevenueTrend,
    getServiceRanking,
    getProductSales,
} from "../api/dashboard";
import Avatar from "../components/Avatar.vue";
import MemberLevel from "../components/MemberLevel.vue";

// 注册 ECharts 组件
use([
    CanvasRenderer,
    LineChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    GridComponent,
]);

const stats = ref({
    dailyRevenue: 0,
    revenueGrowth: 0,
    newMembers: 0,
    activeTechs: 0,
    occupancyRate: 0,
});

const fissionRanking = ref([]);
const revenueTrend = ref([]);
const serviceRanking = ref([]);
const productSales = ref({
    topProducts: [],
    totalRevenue: 0,
    totalSales: 0,
    lowStockCount: 0,
});
const loading = ref(true);
const trendLoading = ref(false);
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
        const [statsRes, rankingRes, trendRes, serviceRes, productRes] =
            await Promise.all([
                getDashboardStats().catch((err) =>
                    console.warn("Stats API failed", err),
                ),
                getFissionRanking().catch((err) =>
                    console.warn("Ranking API failed", err),
                ),
                getRevenueTrend({ days: trendPeriod.value }).catch((err) =>
                    console.warn("Trend API failed", err),
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

        if (trendRes) {
            revenueTrend.value = trendRes;
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

const changeTrendPeriod = async (days) => {
    trendPeriod.value = days;
    trendLoading.value = true;
    try {
        const trendRes = await getRevenueTrend({ days });
        if (trendRes) {
            revenueTrend.value = trendRes;
        }
    } catch (error) {
        console.error("Failed to load revenue trend:", error);
    } finally {
        trendLoading.value = false;
    }
};

onMounted(fetchData);

// ECharts 配置选项
const chartOption = computed(() => {
    if (revenueTrend.value.length === 0) {
        return {
            title: {
                text: "暂无数据",
                left: "center",
                top: "center",
            },
        };
    }

    const dates = revenueTrend.value.map((item) => item.date.substring(5));
    const serviceData = revenueTrend.value.map(
        (item) => item.service_revenue || 0,
    );
    const productData = revenueTrend.value.map(
        (item) => item.product_revenue || 0,
    );

    return {
        tooltip: {
            trigger: "axis",
            backgroundColor: "rgba(31, 41, 55, 0.95)",
            borderColor: "transparent",
            textStyle: {
                color: "#fff",
            },
            formatter: (params) => {
                let result = `<div style="font-weight: 600; margin-bottom: 6px;">${params[0].axisValue}</div>`;
                params.forEach((param) => {
                    const color = param.color;
                    result += `<div style="display: flex; align-items: center; gap: 8px; margin-top: 4px;">
                        <span style="display: inline-block; width: 10px; height: 10px; border-radius: 50%; background: ${color};"></span>
                        <span>${param.seriesName}: ¥${param.value.toFixed(2)}</span>
                    </div>`;
                });
                return result;
            },
        },
        legend: {
            data: ["服务营收", "商品营收"],
            top: 0,
            right: 0,
        },
        grid: {
            left: "3%",
            right: "4%",
            bottom: "3%",
            top: "15%",
            containLabel: true,
        },
        xAxis: {
            type: "category",
            boundaryGap: false,
            data: dates,
            axisLabel: {
                color: "#9ca3af",
                fontSize: 11,
            },
        },
        yAxis: {
            type: "value",
            axisLine: {
                show: false,
            },
            axisTick: {
                show: false,
            },
            axisLabel: {
                color: "#9ca3af",
                fontSize: 11,
                formatter: (value) => {
                    if (value >= 1000) {
                        return (value / 1000).toFixed(1) + "k";
                    }
                    return value;
                },
            },
            splitLine: {
                lineStyle: {
                    color: "#f3f4f6",
                },
            },
        },
        series: [
            {
                name: "服务营收",
                type: "line",
                smooth: true,
                lineStyle: {
                    width: 3,
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 1,
                        y2: 0,
                        colorStops: [
                            { offset: 0, color: "#3b82f6" },
                            { offset: 1, color: "#60a5fa" },
                        ],
                    },
                },
                itemStyle: {
                    color: "#3b82f6",
                    borderWidth: 2,
                    borderColor: "#fff",
                },
                areaStyle: {
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [
                            { offset: 0, color: "rgba(59, 130, 246, 0.2)" },
                            { offset: 1, color: "rgba(59, 130, 246, 0.02)" },
                        ],
                    },
                },
                data: serviceData,
            },
            {
                name: "商品营收",
                type: "line",
                smooth: true,
                lineStyle: {
                    width: 3,
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 1,
                        y2: 0,
                        colorStops: [
                            { offset: 0, color: "#8b5cf6" },
                            { offset: 1, color: "#a78bfa" },
                        ],
                    },
                },
                itemStyle: {
                    color: "#8b5cf6",
                    borderWidth: 2,
                    borderColor: "#fff",
                },
                areaStyle: {
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [
                            { offset: 0, color: "rgba(139, 92, 246, 0.2)" },
                            { offset: 1, color: "rgba(139, 92, 246, 0.02)" },
                        ],
                    },
                },
                data: productData,
            },
        ],
    };
});

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
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-8 h-8">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                        </div>
                        <div class="stat-title">今日营收</div>
                        <div class="stat-value">¥{{ formatNumber(stats.dailyRevenue) }}</div>
                        <div class="stat-desc flex items-center">
                            <span :class="stats.revenueGrowth >= 0
                                ? 'text-success'
                                : 'text-error'
                                " class="flex items-center mr-1">
                                <svg v-if="stats.revenueGrowth >= 0" xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 mr-0.5">
                                    <path fill-rule="evenodd"
                                        d="M12 7a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0V8.414l-4.293 4.293a1 1 0 01-1.414 0L8 10.414l-4.293 4.293a1 1 0 01-1.414-1.414l5-5a1 1 0 011.414 0L11 10.586 14.586 7H12z"
                                        clip-rule="evenodd" />
                                </svg>
                                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
                                    class="w-4 h-4 mr-0.5">
                                    <path fill-rule="evenodd"
                                        d="M12 13a1 1 0 100 2h5a1 1 0 001-1V9a1 1 0 10-2 0v2.586l-4.293-4.293a1 1 0 00-1.414 0L8 9.586 3.707 5.293a1 1 0 00-1.414 1.414l5 5a1 1 0 001.414 0L11 9.414 14.586 13H12z"
                                        clip-rule="evenodd" />
                                </svg>
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
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-8 h-8">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 0110.374 21c-2.331 0-4.512-.645-6.374-1.766z" />
                            </svg>
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
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-8 h-8">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
                            </svg>
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
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                stroke="currentColor" class="w-8 h-8">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
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
                <div v-if="loading || trendLoading" class="h-80 flex items-center justify-center">
                    <span class="loading loading-spinner loading-lg"></span>
                </div>
                <div v-else class="h-80">
                    <v-chart :option="chartOption" autoresize />
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
                        class="flex items-center justify-center py-12 text-base-content/40">
                        暂无数据
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
                            class="flex items-center justify-center py-8 text-base-content/40">
                            暂无商品销售数据
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
                    <h3 class="card-title">
                        🏆 裂变达人榜
                    </h3>
                    <button class="btn btn-ghost btn-sm">
                        查看全部
                    </button>
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
                                    暂无数据
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>