import { defineStore } from "pinia";
import { ref } from "vue";

export const useNotificationStore = defineStore("notificationStore", () => {
  const notifications = ref<string[]>([]);

  const addNotification = (message: string) => {
    notifications.value.push(message);
  };

  return {
    notifications,
    addNotification,
  };
});
