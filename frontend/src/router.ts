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
      component: () => import("@/views/SignUpView.vue"),
      meta: { showLayout: true },
    },
    {
      path: "/play",
      component: () => import("@/views/PlayView.vue"),
      meta: { showLayout: true },
    },
  ],
});

export default router;
