<script setup>
import { ref, computed, reactive, watch } from "vue";
import { RouterLink, RouterView, useRoute } from "vue-router";
import {
    LayoutDashboard,
    Calendar,
    History,
    ClipboardList,
    Package,
    Users,
    UserCircle,
    Menu,
    Bot,
    Sun,
    Moon,
    X,
    ChevronRight
} from 'lucide-vue-next';
import { useTheme } from "../composables/useTheme";
import { usePermission } from "../composables/usePermission";
import UserMenu from "./UserMenu.vue";
import AIReport from "./AIReport.vue";

const { themePreference, setThemePreference } = useTheme();
const { canViewAI, canManageUsers } = usePermission();

const route = useRoute();

const showMobileMenu = ref(false);
const aiReportRef = ref(null);

const openAIAdvisor = () => {
    aiReportRef.value?.open();
};

const menuItems = computed(() => {
    const items = [
        { name: "Dashboard", path: "/", icon: LayoutDashboard },
        {
            name: "订单管理",
            icon: Calendar,
            children: [
                { name: "预约管理", path: "/appointments", icon: Calendar },
                { name: "历史订单", path: "/history", icon: History },
            ]
        },
        {
            name: "业务管理",
            icon: Package,
            children: [
                { name: "服务项目", path: "/services", icon: ClipboardList },
                { name: "实体商品", path: "/products", icon: Package },
            ]
        },
        { name: "技师管理", path: "/technicians", icon: Users },
        { name: "会员管理", path: "/members", icon: UserCircle },
    ];

    // Add user management for managers only
    if (canManageUsers.value) {
        items.push({ name: "用户管理", path: "/users", icon: Users });
    }

    return items;
});

const menuOpenState = reactive({});
watch(
    () => route.path,
    () => {
        menuItems.value.forEach((item) => {
            if (item.children && item.children.some((child) => route.path === child.path)) {
                menuOpenState[item.name] = true;
            }
        });
    },
    { immediate: true } // 初始化时立即执行一次
)

const handleMenuToggle = (event, itemName) => {
    // 同步 DOM 的 open 状态到我们的变量中
    menuOpenState[itemName] = event.target.open;
};


</script>

<template>
    <div
        class="min-h-screen bg-base-100 font-sans text-base-content selection:bg-primary selection:text-primary-content drawer lg:drawer-open">
        <input id="drawer-toggle" type="checkbox" class="drawer-toggle" />
        <div class="drawer-content flex-1 min-w-0 flex flex-col">
            <!-- Mobile Header -->
            <div
                class="lg:hidden sticky top-0 z-40 flex items-center justify-between px-4 py-3 bg-base-100/80 backdrop-blur-sm border-b border-base-200">
                <div class="flex items-center gap-2 font-bold text-lg">
                    <span class="text-primary">Smart</span>Spa
                </div>
                <div class="flex items-center gap-2">
                    <UserMenu />
                </div>
            </div>

            <!-- Main Content -->
            <main class="flex-1 min-w-0 pb-20 lg:pb-10">
                <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 lg:py-10">
                    <RouterView />
                </div>
            </main>

            <!-- Mobile Bottom Navigation -->
            <div class="lg:hidden fixed bottom-0 left-0 right-0 z-40 bg-base-100 border-t border-base-200 pb-safe">
                <div class="grid grid-cols-4 h-16">
                    <RouterLink to="/" class="flex flex-col items-center justify-center gap-1 text-base-content/60"
                        :class="{ 'text-primary font-medium': route.path === '/' }">
                        <LayoutDashboard class="w-6 h-6" />
                        <span class="text-[10px]">概览</span>
                    </RouterLink>
                    <RouterLink to="/appointments"
                        class="flex flex-col items-center justify-center gap-1 text-base-content/60"
                        :class="{ 'text-primary font-medium': route.path === '/appointments' }">
                        <Calendar class="w-6 h-6" />
                        <span class="text-[10px]">预约</span>
                    </RouterLink>
                    <RouterLink to="/technicians"
                        class="flex flex-col items-center justify-center gap-1 text-base-content/60"
                        :class="{ 'text-primary font-medium': route.path === '/technicians' }">
                        <Users class="w-6 h-6" />
                        <span class="text-[10px]">技师</span>
                    </RouterLink>
                    <button @click="showMobileMenu = true"
                        class="flex flex-col items-center justify-center gap-1 text-base-content/60">
                        <Menu class="w-6 h-6" />
                        <span class="text-[10px]">菜单</span>
                    </button>
                </div>
            </div>

            <!-- Mobile Menu Sheet -->
            <dialog class="modal modal-bottom lg:hidden" :class="{ 'modal-open': showMobileMenu }">
                <div class="modal-box max-h-[85vh] flex flex-col p-0 bg-base-100">
                    <!-- Sheet Header -->
                    <div
                        class="flex items-center justify-between p-4 border-b border-base-200 sticky top-0 bg-base-100 z-10">
                        <span class="font-bold text-lg">功能导航</span>
                        <button @click="showMobileMenu = false" class="btn btn-sm btn-circle btn-ghost">
                            <X class="w-5 h-5" />
                        </button>
                    </div>

                    <!-- Sheet Content -->
                    <div class="p-4 overflow-y-auto space-y-6">
                        <template v-for="item in menuItems" :key="item.name">
                            <div v-if="item.children" class="space-y-3">
                                <h3 class="font-bold text-sm text-base-content/40 uppercase tracking-wider px-1">{{
                                    item.name }}</h3>
                                <div class="grid grid-cols-2 gap-3">
                                    <RouterLink v-for="child in item.children" :key="child.path" :to="child.path"
                                        @click="showMobileMenu = false"
                                        class="flex flex-col items-center justify-center gap-2 p-4 bg-base-200/50 rounded-xl active:scale-95 transition-transform"
                                        :class="{ 'bg-primary/10 text-primary ring-1 ring-primary/20': route.path === child.path }">
                                        <component :is="child.icon" class="w-8 h-8" />
                                        <span class="text-sm font-medium">{{ child.name }}</span>
                                    </RouterLink>
                                </div>
                            </div>
                            <RouterLink v-else :to="item.path" @click="showMobileMenu = false"
                                class="flex items-center gap-4 p-4 bg-base-200/50 rounded-xl active:scale-95 transition-transform"
                                :class="{ 'bg-primary/10 text-primary ring-1 ring-primary/20': route.path === item.path }">
                                <component :is="item.icon" class="w-8 h-8" />
                                <span class="text-sm font-medium">{{ item.name }}</span>
                            </RouterLink>
                        </template>

                        <!-- Mobile Footer Actions -->
                        <div class="pt-4 border-t border-base-200 space-y-4">
                            <label class="flex items-center justify-between p-4 bg-base-200/30 rounded-xl">
                                <div class="flex items-center gap-3">
                                    <Sun v-if="themePreference === 'light'" class="w-5 h-5" />
                                    <Moon v-else class="w-5 h-5" />
                                    <span class="font-medium">深色模式</span>
                                </div>
                                <input type="checkbox" value="dark" class="toggle toggle-primary"
                                    :checked="themePreference === 'dark'"
                                    @change="themePreference === 'dark' ? setThemePreference('light') : setThemePreference('dark')" />
                            </label>

                            <button v-if="canViewAI" @click="showMobileMenu = false; openAIAdvisor()"
                                class="w-full btn btn-primary h-12 text-lg font-medium shadow-lg shadow-primary/20">
                                <Bot class="w-6 h-6 mr-2" />
                                AI 经营顾问
                            </button>
                        </div>
                    </div>
                </div>
                <form method="dialog" class="modal-backdrop bg-base-content/20 backdrop-blur-sm">
                    <button @click="showMobileMenu = false">close</button>
                </form>
            </dialog>
        </div>

        <div class="drawer-side">
            <label for="drawer-toggle" aria-label="close sidebar" class="drawer-overlay"></label>
            <!-- Sidebar -->
            <aside class="min-h-full w-80 bg-base-100 text-base-content flex flex-col">
                <!-- Logo and User Menu -->
                <div class="h-16 flex items-center justify-between px-6 border-b border-base-200">
                    <div class="flex items-center gap-2 font-bold text-xl">
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
                    <template v-for="item in menuItems" :key="item.name">
                        <!-- Group with Children -->
                        <div v-if="item.children" class="space-y-1">
                            <details class="group" :open="menuOpenState[item.name]"
                                @toggle="handleMenuToggle($event, item.name)">
                                <summary
                                    class="flex items-center justify-between w-full gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors text-base-content/70 hover:bg-base-200 hover:text-base-content cursor-pointer select-none list-none marker:content-none">
                                    <div class="flex items-center gap-3">
                                        <component :is="item.icon" class="w-5 h-5" />
                                        {{ item.name }}
                                    </div>
                                    <ChevronRight class="size-4 transition-transform group-open:rotate-90 opacity-50" />
                                </summary>
                                <div class="mt-1 pl-4 space-y-1 border-l-2 border-base-200 ml-4">
                                    <RouterLink v-for="child in item.children" :key="child.path" :to="child.path"
                                        class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors"
                                        :class="[
                                            route.path === child.path
                                                ? 'bg-primary text-primary-content shadow-sm'
                                                : 'text-base-content/70 hover:bg-base-200 hover:text-base-content',
                                        ]">
                                        <component :is="child.icon" class="w-5 h-5" />
                                        {{ child.name }}
                                    </RouterLink>
                                </div>
                            </details>
                        </div>

                        <!-- Single Item -->
                        <RouterLink v-else :to="item.path"
                            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors"
                            :class="[
                                route.path === item.path
                                    ? 'bg-primary text-primary-content shadow-sm'
                                    : 'text-base-content/70 hover:bg-base-200 hover:text-base-content',
                            ]">
                            <component :is="item.icon" class="w-5 h-5" />
                            {{ item.name }}
                        </RouterLink>
                    </template>
                </nav>

                <!-- Footer Actions -->
                <div class="mt-auto p-4 border-t border-base-200 space-y-3">
                    <label class="flex cursor-pointer gap-2 bg-base-200/50 p-2 rounded-lg justify-center">
                        <Sun class="w-5 h-5" />
                        <input type="checkbox" value="dark" class="toggle theme-controller"
                            :checked="themePreference === 'dark'"
                            @change="themePreference === 'dark' ? setThemePreference('light') : setThemePreference('dark')" />
                        <Moon class="w-5 h-5" />
                    </label>
                    <button v-if="canViewAI" @click="openAIAdvisor"
                        class="w-full btn btn-primary btn-sm h-10 font-medium">
                        <Bot class="w-5 h-5" />
                        AI 经营顾问
                    </button>
                </div>
            </aside>
        </div>

        <!-- AI Advisor Modal -->
        <AIReport ref="aiReportRef" />
    </div>
</template>

<style>
@reference "../style.css";
</style>