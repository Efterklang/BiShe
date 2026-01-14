<script setup>
import { computed } from 'vue';

const props = defineProps({
    name: {
        type: String,
        required: true,
        default: '?'
    },
    size: {
        type: String,
        default: 'md', // xs, sm, md, lg, xl
        validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(value)
    }
});

// 简单的字符串 hash 函数
const hashString = (str) => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // Convert to 32bit integer
    }
    return Math.abs(hash);
};

// 颜色配对（背景色 + 文字颜色）
const colorPairs = [
    { bg: 'bg-red-100', text: 'text-red-700' },
    { bg: 'bg-blue-100', text: 'text-blue-700' },
    { bg: 'bg-green-100', text: 'text-green-700' },
    { bg: 'bg-yellow-100', text: 'text-yellow-700' },
    { bg: 'bg-purple-100', text: 'text-purple-700' },
    { bg: 'bg-pink-100', text: 'text-pink-700' },
    { bg: 'bg-indigo-100', text: 'text-indigo-700' },
    { bg: 'bg-cyan-100', text: 'text-cyan-700' },
    { bg: 'bg-orange-100', text: 'text-orange-700' },
    { bg: 'bg-lime-100', text: 'text-lime-700' }
];

const colorIndex = computed(() => {
    const hash = hashString(props.name || '?');
    return hash % colorPairs.length;
});

const avatarColors = computed(() => colorPairs[colorIndex.value]);

const sizeClasses = computed(() => {
    const sizeMap = {
        'xs': 'w-6 h-6 text-xs',
        'sm': 'w-8 h-8 text-sm',
        'md': 'w-10 h-10 text-base',
        'lg': 'w-16 h-16 text-xl',
        'xl': 'w-20 h-20 text-2xl'
    };
    return sizeMap[props.size] || sizeMap['md'];
});

const initial = computed(() => {
    return props.name ? props.name.charAt(0).toUpperCase() : '?';
});
</script>

<template>
    <div :class="[
        avatarColors.bg,
        avatarColors.text,
        sizeClasses,
        'rounded-full flex items-center justify-center font-bold shrink-0'
    ]" :title="name">
        {{ initial }}
    </div>
</template>