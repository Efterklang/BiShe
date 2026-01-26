<script setup>
import { ref, computed, watch } from "vue";
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

import { getRevenueTrend } from "../api/dashboard";

const props = defineProps({
    days: {
        type: Number,
        default: 30,
    },
});

const revenueTrend = ref([]);
const loading = ref(false);

const fetchRevenueTrend = async (days) => {
    loading.value = true;
    try {
        const res = await getRevenueTrend({ days });
        revenueTrend.value = res || [];
    } catch (error) {
        console.error("Failed to load revenue trend:", error);
        revenueTrend.value = [];
    } finally {
        loading.value = false;
    }
};

// 注册 ECharts 组件
use([
    CanvasRenderer,
    LineChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    GridComponent,
]);

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
    const totalData = serviceData.map((s, i) => s + productData[i]);

    return {
        tooltip: {
            trigger: "axis",
            backgroundColor: "rgba(255, 255, 255, 0.95)",
            borderColor: "#e5e7eb",
            borderWidth: 1,
            textStyle: {
                color: "#1f2937",
                fontSize: 12
            },
            padding: [8, 12],
            extraCssText: 'box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);',
            formatter: (params) => {
                let result = `<div style="font-weight: 600; margin-bottom: 6px; color: #111827;">${params[0].axisValue}</div>`;
                params.forEach((param) => {
                    const color = param.color.colorStops ? param.color.colorStops[0].color : param.color;
                    result += `<div style="display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 4px;">
                        <div style="display: flex; align-items: center; gap: 6px;">
                            <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: ${color};"></span>
                            <span style="color: #4b5563;">${param.seriesName}</span>
                        </div>
                        <span style="font-weight: 500; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;">¥${param.value.toFixed(2)}</span>
                    </div>`;
                });
                return result;
            },
        },
        legend: {
            data: ["服务营收", "商品营收", "总营收"],
            top: 0,
            right: 0,
            icon: "circle",
            itemGap: 16,
            textStyle: {
                color: "#6b7280"
            }
        },
        grid: {
            left: "1%",
            right: "1%",
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
                margin: 12
            },
            axisLine: {
                show: false
            },
            axisTick: {
                show: false
            }
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
                    type: "dashed"
                },
            },
        },
        series: [
            {
                name: "服务营收",
                type: "line",
                smooth: true,
                showSymbol: false,
                symbolSize: 8,
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
                    shadowColor: 'rgba(59, 130, 246, 0.3)',
                    shadowBlur: 10,
                    shadowOffsetY: 5
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
                            { offset: 1, color: "rgba(59, 130, 246, 0.0)" },
                        ],
                    },
                },
                data: serviceData,
            },
            {
                name: "商品营收",
                type: "line",
                smooth: true,
                showSymbol: false,
                symbolSize: 8,
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
                    shadowColor: 'rgba(139, 92, 246, 0.3)',
                    shadowBlur: 10,
                    shadowOffsetY: 5
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
                            { offset: 1, color: "rgba(139, 92, 246, 0.0)" },
                        ],
                    },
                },
                data: productData,
            },
            {
                name: "总营收",
                type: "line",
                smooth: true,
                showSymbol: false,
                symbolSize: 8,
                z: 10,
                lineStyle: {
                    width: 3,
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 1,
                        y2: 0,
                        colorStops: [
                            { offset: 0, color: "#10b981" },
                            { offset: 1, color: "#34d399" },
                        ],
                    },
                    shadowColor: 'rgba(16, 185, 129, 0.3)',
                    shadowBlur: 10,
                    shadowOffsetY: 5
                },
                itemStyle: {
                    color: "#10b981",
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
                            { offset: 0, color: "rgba(16, 185, 129, 0.2)" },
                            { offset: 1, color: "rgba(16, 185, 129, 0.0)" },
                        ],
                    },
                },
                data: totalData,
            },
        ],
    };
});

watch(() => props.days, (newDays) => {
    fetchRevenueTrend(newDays);
}, { immediate: true });

</script>

<template>
    <div class="h-80 w-full">
        <div v-if="loading" class="h-full flex items-center justify-center">
            <span class="loading loading-spinner loading-lg"></span>
        </div>
        <v-chart v-else :option="chartOption" autoresize class="w-full h-full" />
    </div>
</template>

<style scoped>
/* 确保图表容器大小正确 */
:deep(.echarts) {
    width: 100%;
    height: 100%;
}
</style>