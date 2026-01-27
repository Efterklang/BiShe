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
                textStyle: {
                    color: "#2e3440" // Nord polar night (dark text for light bg)
                }
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
            backgroundColor: "#eceff4", // Nord snow storm 1 (light bg)
            borderColor: "#d8dee9", // Nord snow storm 2 (border)
            borderWidth: 1,
            textStyle: {
                color: "#2e3440", // Nord polar night (text)
                fontSize: 12
            },
            padding: [8, 12],
            extraCssText: 'box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.08), 0 2px 4px -1px rgba(0, 0, 0, 0.04);',
            formatter: (params) => {
                let result = `<div style="font-weight: 600; margin-bottom: 6px; color: #2e3440;">${params[0].axisValue}</div>`;
                params.forEach((param) => {
                    const color = param.color.colorStops ? param.color.colorStops[0].color : param.color;
                    result += `<div style="display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 4px;">
                        <div style="display: flex; align-items: center; gap: 6px;">
                            <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: ${color};"></span>
                            <span style="color: #4c566a;">${param.seriesName}</span>
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
                color: "#4c566a" // Nord polar night 3 (legend text)
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
                color: "#616e88", // Nord snow storm 3 (axis text)
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
                color: "#616e88", // Nord snow storm 3 (axis text)
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
                    color: "#e5e9f0", // Nord snow storm 2 (split line)
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
                            { offset: 0, color: "#5e81ac" }, // Nord frost 3 (blue)
                            { offset: 1, color: "#81a1c1" }, // Nord frost 4 (light blue)
                        ],
                    },
                    shadowColor: 'rgba(94, 129, 172, 0.2)',
                    shadowBlur: 10,
                    shadowOffsetY: 5
                },
                itemStyle: {
                    color: "#5e81ac", // Nord frost 3
                    borderWidth: 2,
                    borderColor: "#eceff4", // Nord snow storm 1 (border)
                },
                areaStyle: {
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [
                            { offset: 0, color: "rgba(94, 129, 172, 0.15)" },
                            { offset: 1, color: "rgba(94, 129, 172, 0.0)" },
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
                            { offset: 0, color: "#b48ead" }, // Nord aurora 4 (purple)
                            { offset: 1, color: "#d08770" }, // Nord aurora 3 (orange)
                        ],
                    },
                    shadowColor: 'rgba(180, 142, 173, 0.2)',
                    shadowBlur: 10,
                    shadowOffsetY: 5
                },
                itemStyle: {
                    color: "#b48ead", // Nord aurora 4
                    borderWidth: 2,
                    borderColor: "#eceff4", // Nord snow storm 1
                },
                areaStyle: {
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [
                            { offset: 0, color: "rgba(180, 142, 173, 0.15)" },
                            { offset: 1, color: "rgba(180, 142, 173, 0.0)" },
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
                            { offset: 0, color: "#8fbcbb" }, // Nord frost 1 (teal)
                            { offset: 1, color: "#88c0d0" }, // Nord frost 2 (light teal)
                        ],
                    },
                    shadowColor: 'rgba(143, 188, 187, 0.2)',
                    shadowBlur: 10,
                    shadowOffsetY: 5
                },
                itemStyle: {
                    color: "#8fbcbb", // Nord frost 1
                    borderWidth: 2,
                    borderColor: "#eceff4", // Nord snow storm 1
                },
                areaStyle: {
                    color: {
                        type: "linear",
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [
                            { offset: 0, color: "rgba(143, 188, 187, 0.15)" },
                            { offset: 1, color: "rgba(143, 188, 187, 0.0)" },
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
        <VChart v-else :option="chartOption" autoresize class="w-full h-full" />
    </div>
</template>

<style scoped>
/* 确保图表容器大小正确 */
:deep(.echarts) {
    width: 100%;
    height: 100%;
}
</style>
