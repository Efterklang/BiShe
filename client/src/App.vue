<script setup>
import { RouterView } from "vue-router";
import { useTheme } from "./composables/useTheme";

// Initialize theme support
useTheme();
</script>

<template>
    <RouterView />
</template>
