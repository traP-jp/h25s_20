import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: () => import("@/views/RoomSelectView.vue"),
      meta: { showLayout: false },
    },
    {
      path: "/signup",
      name: "サインアップ",
      component: () => import("@/views/SignUpView.vue"),
      meta: { showLayout: true },
    },
    {
      path: "/play",
      name: "サインアップ",
      component: () => import("@/views/PlayView.vue"),
      meta: { showLayout: true },
    },
  ],
});

export default router;
