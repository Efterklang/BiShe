import { defineShikiSetup } from "@slidev/types";

export default defineShikiSetup(() => {
	return {
		themes: {
			light: "catppuccin-latte",
			dark: "catppuccin-mocha",
		},
		transformers: [
			// ...
		],
	};
});
