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
        'silver': 'badge-info',
        'gold': 'badge-warning',
        'platinum': 'badge-primary',
        '普通会员': 'badge-ghost'
    };
    return colorMap[levelStr] || 'badge-ghost';
});

// 显示的等级文本
const displayLevel = computed(() => {
    const levelMap = {
        'basic': '普通',
        'vip': '会员',
        'silver': '白银',
        'gold': '黄金',
        'platinum': '白金'
    };
    return levelMap[props.level] || '普通会员';
});
</script>

<template>
    <span :class="['badge', badgeColor]">
        {{ displayLevel }}
    </span>
</template>