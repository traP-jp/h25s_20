import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: () => import("@/views/RouteView.vue"),
      meta: { showLayout: false },
    },
    {
      path: "/rooms",
      component: () => import("@/views/RoomSelectView.vue"),
      meta: { showLayout: false },
    },
    {
      path: "/signup",
      component: () => import("@/views/SignUpView.vue"),
      meta: { showLayout: true },
    },
    {
      path: "/play/:roomId?",
      name: "play",
      component: () => import("@/views/PlayView.vue"),
      meta: { showLayout: true },
    },
    {
      path: "/api-test",
      component: () => import("@/views/ApiTestView.vue"),
      meta: { showLayout: false },
    },
    {
      path: "/help",
      component: () => import("@/views/HelpView.vue"),
      meta: { showLayout: true },
    },
    {
      path: "/:pathMatch(.*)*",
      component: () => import("@/views/NotFoundView.vue"),
      meta: { showLayout: false },
    }
  ],
});

export default router;
