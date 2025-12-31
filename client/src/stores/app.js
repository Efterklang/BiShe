import { defineStore } from "pinia";
import { ref } from "vue";

export const useAppStore = defineStore("app", () => {
	const isLoading = ref(false);
	const theme = ref("emerald");

	function setLoading(status) {
		isLoading.value = status;
	}

	function setTheme(newTheme) {
		theme.value = newTheme;
		document.documentElement.setAttribute("data-theme", newTheme);
	}

	return { isLoading, theme, setLoading, setTheme };
});
