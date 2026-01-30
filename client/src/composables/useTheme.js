import { onMounted, ref } from "vue";

const THEME_KEY = "spa-admin-theme-preference";

export function useTheme() {
	// 'nord' is light, 'mocha' is dark
	const lightTheme = "nord";
	const darkTheme = "mocha";

	// Current active theme (applied to DOM)
	const currentTheme = ref(lightTheme);

	// User preference: 'light', 'dark', or 'system'
	const themePreference = ref(localStorage.getItem(THEME_KEY) || "system");

	const applyTheme = (theme) => {
		document.documentElement.setAttribute("data-theme", theme);
		currentTheme.value = theme;
	};

	const handleSystemChange = (e) => {
		if (themePreference.value === "system") {
			applyTheme(e.matches ? darkTheme : lightTheme);
		}
	};

	const updateTheme = () => {
		const pref = themePreference.value;

		if (pref === "system") {
			const systemDark = window.matchMedia(
				"(prefers-color-scheme: dark)",
			).matches;
			applyTheme(systemDark ? darkTheme : lightTheme);
		} else if (pref === "dark") {
			applyTheme(darkTheme);
		} else {
			applyTheme(lightTheme);
		}
	};

	const setThemePreference = (pref) => {
		themePreference.value = pref;
		localStorage.setItem(THEME_KEY, pref);
		updateTheme();
	};

	onMounted(() => {
		// Initial check
		updateTheme();

		// Listen for system changes
		const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
		mediaQuery.addEventListener("change", handleSystemChange);
	});

	return {
		themePreference,
		currentTheme,
		setThemePreference,
		lightTheme,
		darkTheme,
	};
}
