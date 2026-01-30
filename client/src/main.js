import { createApp } from "vue";
import { createPinia } from "pinia";
import "./style.css";
import App from "./App.vue";
import router from "./router";
import { useAppStore } from "./stores/app";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

// Initialize auth state from localStorage
const appStore = useAppStore();
appStore.initAuth();

app.mount("#app");
