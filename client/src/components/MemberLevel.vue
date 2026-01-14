<script setup>
import { computed } from 'vue';

const props = defineProps({
    level: {
        type: String,
        default: '普通会员'
    },
    size: {
        type: String,
        default: 'sm', // xs, sm, md, lg
        validator: (value) => ['xs', 'sm', 'md', 'lg'].includes(value)
    }
});

// 根据等级返回 badge 颜色
const badgeColor = computed(() => {
    const levelStr = props.level || '普通会员';
    const colorMap = {
        'basic': 'badge-ghost',
        'vip': 'badge-success',
        'silver': 'badge-warning',
        'gold': 'badge-info',
        'platinum': 'badge-primary',
        '普通会员': 'badge-ghost'
    };
    return colorMap[levelStr] || 'badge-ghost';
});

// badge 大小类
const sizeClasses = computed(() => {
    const sizeMap = {
        'xs': 'badge-xs',
        'sm': 'badge-sm',
        'md': 'badge-md',
        'lg': 'badge-lg'
    };
    return sizeMap[props.size] || sizeMap['sm'];
});

// 显示的等级文本
const displayLevel = computed(() => {
    return props.level || '普通会员';
});
</script>

<template>
    <span :class="['badge', badgeColor, sizeClasses]">
        {{ displayLevel }}
    </span>
</template>