<script setup>
import { ref, computed } from "vue";
import MarkdownIt from "markdown-it";
import { generateAIReport } from "../api/ai";

const showAIModal = ref(false);
const aiReport = ref("");
const displayedReport = ref("");
const aiLoading = ref(false);
const md = new MarkdownIt();
const STORAGE_KEY = "spa_ai_report";

const renderedReport = computed(() => md.render(displayedReport.value));

const typeWriter = (text, index = 0) => {
    if (index < text.length && showAIModal.value) {
        displayedReport.value += text.charAt(index);
        // Randomize typing speed slightly for realism
        const delay = Math.random() * 20 + 10;
        setTimeout(() => typeWriter(text, index + 1), delay);
    }
};

const open = async () => {
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
    open();
};

defineExpose({
    open
});
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': showAIModal }">
        <div class="modal-box w-11/12 max-w-4xl bg-base-100 border border-base-300 shadow-2xl h-[80vh] flex flex-col">
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
</template>

<style scoped>
@reference "../style.css";

.markdown-body :deep(h1) {
    @apply text-2xl font-bold my-4;
}

.markdown-body :deep(h2) {
    @apply text-xl font-bold my-3;
}

.markdown-body :deep(h3) {
    @apply text-lg font-bold my-2;
}

.markdown-body :deep(p) {
    @apply my-2 leading-relaxed;
}

.markdown-body :deep(ul) {
    @apply list-disc list-inside my-2;
}

.markdown-body :deep(ol) {
    @apply list-decimal list-inside my-2;
}

.markdown-body :deep(li) {
    @apply my-1;
}

.markdown-body :deep(strong) {
    @apply font-bold;
}

.markdown-body :deep(blockquote) {
    @apply border-l-4 border-base-300 pl-4 italic my-4;
}
</style>
