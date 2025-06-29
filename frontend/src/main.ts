import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";

const app = createApp(App);

import router from "./router";
import { createPinia } from "pinia";

app.use(router);
app.use(createPinia());

app.mount("#app");
