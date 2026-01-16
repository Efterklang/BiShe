<script setup>
import { ref, computed } from "vue";
import { RouterLink, RouterView, useRoute } from "vue-router";
import MarkdownIt from "markdown-it";
import { generateAIReport } from "../api/ai";
import { useTheme } from "../composables/useTheme";
import { useAppStore } from "../stores/app";
import { usePermission } from "../composables/usePermission";
import UserMenu from "./UserMenu.vue";

const { themePreference, setThemePreference } = useTheme();
const { canViewAI, canManageUsers } = usePermission();

const md = new MarkdownIt();
const STORAGE_KEY = "spa_ai_report";

const route = useRoute();

const showAIModal = ref(false);
const aiReport = ref("");
const displayedReport = ref("");
const aiLoading = ref(false);

const renderedReport = computed(() => md.render(displayedReport.value));

const typeWriter = (text, index = 0) => {
    if (index < text.length && showAIModal.value) {
        displayedReport.value += text.charAt(index);
        // Randomize typing speed slightly for realism
        const delay = Math.random() * 20 + 10;
        setTimeout(() => typeWriter(text, index + 1), delay);
    }
};

const openAIAdvisor = async () => {
    showAIModal.value = true;

    // Try load from storage if empty
    if (!aiReport.value) {
        const cached = localStorage.getItem(STORAGE_KEY);
        if (cached) {
            aiReport.value = cached;
            displayedReport.value = cached; // Show immediately without typing effect
            return;
        }
    }

    if (!aiReport.value) {
        aiLoading.value = true;
        displayedReport.value = "";
        try {
            const data = await generateAIReport();
            aiReport.value = data.report;
            localStorage.setItem(STORAGE_KEY, data.report);
            typeWriter(aiReport.value);
        } catch (error) {
            aiReport.value =
                "⚠️ 获取分析报告失败，请检查网络或服务器状态。\n\n错误详情: " +
                (error.message || "未知错误");
            displayedReport.value = aiReport.value;
        } finally {
            aiLoading.value = false;
        }
    } else if (displayedReport.value.length < aiReport.value.length) {
        // Restart typing if previously interrupted
        displayedReport.value = "";
        typeWriter(aiReport.value);
    }
};

const regenerateReport = () => {
    localStorage.removeItem(STORAGE_KEY);
    aiReport.value = "";
    displayedReport.value = "";
    openAIAdvisor();
};

const menuItems = computed(() => {
    const items = [
        { name: "Dashboard", path: "/", icon: "📊" },
        { name: "预约管理", path: "/appointments", icon: "📅" },
        { name: "技师管理", path: "/technicians", icon: "💆" },
        { name: "服务项目", path: "/services", icon: "📋" },
        { name: "实体商品", path: "/products", icon: "📦" },
        { name: "会员管理", path: "/members", icon: "👥" },
        { name: "历史订单", path: "/history", icon: "📜" },
    ];

    // Add user management for managers only
    if (canManageUsers.value) {
        items.push({ name: "用户管理", path: "/users", icon: "👤" });
    }

    return items;
});


</script>

<template>
    <div
        class="min-h-screen bg-base-100 font-sans text-base-content selection:bg-primary selection:text-primary-content drawer lg:drawer-open">
        <input id="drawer-toggle" type="checkbox" class="drawer-toggle" />
        <div class="drawer-content flex-1 min-w-0">
            <!-- Mobile Header -->
            <div
                class="lg:hidden sticky top-0 z-40 flex items-center justify-between px-4 py-3 bg-base-100/80 backdrop-blur-sm border-b border-base-200"
            >
                <div class="flex items-center gap-2 font-bold text-lg tracking-tight">
                    <span class="text-primary">Smart</span>Spa
                </div>
                <div class="flex items-center gap-2">
                    <UserMenu />
                    <label for="drawer-toggle" class="p-2 text-base-content/70 hover:bg-base-200 rounded-md">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="1.5"
                            stroke="currentColor"
                            class="w-6 h-6"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                            />
                        </svg>
                    </label>
                </div>
            </div>




            <!-- Main Content -->
            <main class="flex-1 min-w-0">
                <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 lg:py-10">
                    <RouterView />
                </div>
            </main>
        </div>

        <div class="drawer-side">
            <label for="drawer-toggle" aria-label="close sidebar" class="drawer-overlay"></label>
            <!-- Sidebar -->
            <aside class="min-h-full w-80 bg-base-100 text-base-content">
                <div class="flex flex-col h-full">
                    <!-- Logo and User Menu -->
                    <div class="h-16 flex items-center justify-between px-6 border-b border-base-200">
                        <div class="flex items-center gap-2 font-bold text-xl tracking-tight">
                            <div
                                class="w-8 h-8 bg-primary text-primary-content rounded-lg flex items-center justify-center text-sm font-bold">
                                S
                            </div>
                            <span>XX养生店</span>
                        </div>
                        <div class="hidden lg:block">
                            <UserMenu />
                        </div>
                    </div>

                    <!-- Navigation -->
                    <nav class="flex-1 px-4 py-6 space-y-1">
                        <p class="px-2 text-xs font-semibold text-base-content/50 uppercase tracking-wider mb-4">
                            Menu
                        </p>
                        <RouterLink v-for="item in menuItems" :key="item.path" :to="item.path"
                            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors"
                            :class="[
                                route.path === item.path
                                    ? 'bg-primary text-primary-content shadow-sm'
                                    : 'text-base-content/70 hover:bg-base-200 hover:text-base-content',
                            ]">
                            <span class="text-lg">{{ item.icon }}</span>
                            {{ item.name }}
                        </RouterLink>
                    </nav>

                    <!-- Footer Actions -->
                    <div class="p-4 border-t border-base-200 space-y-3">
                        <label class="flex cursor-pointer gap-2 bg-base-200/50 p-2 rounded-lg justify-center">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="20"
                                height="20"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round">
                                <circle cx="12" cy="12" r="5" />
                                <path
                                    d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4" />
                            </svg>
                            <input type="checkbox" value="dark" class="toggle theme-controller" :checked="themePreference === 'dark'" @change="themePreference === 'dark' ? setThemePreference('light') : setThemePreference('dark')" />
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="20"
                                height="20"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round">
                                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
                            </svg>
                        </label>
                        <button
                            v-if="canViewAI"
                            @click="openAIAdvisor"
                            class="w-full btn btn-primary btn-sm h-10 font-medium"
                        >
                            <span>🤖</span>
                            AI 经营顾问
                        </button>
                    </div>
                </div>
            </aside>
        </div>

        <!-- AI Advisor Modal -->
        <dialog class="modal" :class="{ 'modal-open': showAIModal }">
            <div
                class="modal-box w-11/12 max-w-4xl bg-base-100 border border-base-300 shadow-2xl h-[80vh] flex flex-col">
                <div class="flex justify-between items-center mb-4 pb-2 border-b border-base-200">
                    <h3 class="font-bold text-lg flex items-center gap-2 text-base-content">
                        <span>🤖</span> 智能经营顾问
                    </h3>
                    <button @click="showAIModal = false" class="btn btn-sm btn-circle btn-ghost text-base-content/60">
                        ✕
                    </button>
                </div>

                <div class="flex-1 overflow-y-auto p-6 bg-base-200/30 rounded-xl text-base-content">
                    <div v-if="aiLoading" class="flex flex-col items-center justify-center h-full gap-4">
                        <span class="loading loading-dots loading-lg text-primary"></span>
                        <p class="text-base-content/60 animate-pulse">
                            正在分析经营数据...
                        </p>
                    </div>
                    <div v-else class="markdown-body">
                        <div v-html="renderedReport"></div>
                        <span class="animate-pulse inline-block w-2 h-4 bg-primary ml-1 align-middle"
                            v-if="displayedReport.length < aiReport.length"></span>
                    </div>
                </div>

                <div class="modal-action mt-4 flex justify-between items-center">
                    <div class="text-xs text-base-content/40">
                        基于近30天运营数据生成
                    </div>
                    <div class="flex gap-2">
                        <button @click="regenerateReport" class="btn btn-outline btn-sm" :disabled="aiLoading">
                            🔄 重新生成
                        </button>
                        <button @click="showAIModal = false" class="btn btn-primary btn-sm">
                            关闭
                        </button>
                    </div>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
                <button @click="showAIModal = false">close</button>
            </form>
        </dialog>
    </div>
</template>

<style>
@reference "../style.css";

.markdown-body h1 {
    @apply text-2xl font-bold my-4;
}

.markdown-body h2 {
    @apply text-xl font-bold my-3;
}

.markdown-body h3 {
    @apply text-lg font-bold my-2;
}

.markdown-body p {
    @apply my-2 leading-relaxed;
}

.markdown-body ul {
    @apply list-disc list-inside my-2;
}

.markdown-body ol {
    @apply list-decimal list-inside my-2;
}

.markdown-body li {
    @apply my-1;
}

.markdown-body strong {
    @apply font-bold;
}

.markdown-body blockquote {
    @apply border-l-4 border-base-300 pl-4 italic my-4;
}
</style>
