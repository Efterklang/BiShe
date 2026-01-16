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
            data: ["服务营收", "商品营收", "总营收"],
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
                    color: "#9ca3af",
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
            {
                name: "总营收",
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
                            { offset: 0, color: "#10b981" },
                            { offset: 1, color: "#34d399" },
                        ],
                    },
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
                            { offset: 1, color: "rgba(16, 185, 129, 0.02)" },
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
    <div v-if="loading" class="h-80 flex items-center justify-center">
        <span class="loading loading-spinner loading-lg"></span>
    </div>
    <v-chart v-else :option="chartOption" autoresize />
</template>