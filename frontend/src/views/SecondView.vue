<script setup lang="ts">
// https://medium.com/@erkalanhanife/why-use-vue-3-pinia-and-websockets-5018a470f245

import { ref, onMounted, onUnmounted } from "vue";
import { useNotificationStore } from "@/store";

const ws = ref<WebSocket | null>(null);
const isConnected = ref(false);
const retryInterval = 3000;
let retryTimeout: number | null = null; // retrtTimeout はタイマーの ID。割合小さな自然数

const notificationStore = useNotificationStore();

const connectWebSocket = () => {
  // ws.value = new WebSocket("ws://10ten.trap.show");
  // ws.value = new WebSocket("wss://echo.websocket.org");
  ws.value = new WebSocket("wss://kaitoyama-websocket-poc.trap.show/ws?room_id=room1&user_id=kitsne");

  ws.value.onopen = () => {
    console.log("WebSocket connected.");
    isConnected.value = true;
    if (retryTimeout) {
      clearTimeout(retryTimeout);
      retryTimeout = null;
    }
  };

  ws.value.onmessage = (event) => {
    // 単に pinia のストアにメッセージの内容を追加するだけ
    notificationStore.addNotification(event.data);
  };

  ws.value.onerror = (error) => {
    console.error("WebSocket error:", error);
    isConnected.value = false;
  };

  // 接続が閉じられたとき・接続に失敗したときにも実行される
  ws.value.onclose = (event) => {
    if (event.wasClean) {
      console.log("GoodBye!");
    } else {
      console.log("WebSocket closed. Retrying...");
      isConnected.value = false;
      attemptReconnect();
    }
  };
};

const attemptReconnect = () => {
  if (!isConnected.value) {
    retryTimeout = setTimeout(() => {
      console.log("Attempting to reconnect...");
      connectWebSocket();
    }, retryInterval);
  }
};

onMounted(connectWebSocket);

onUnmounted(() => {
  if (ws.value) {
    ws.value.close();
  }
  if (retryTimeout) {
    clearTimeout(retryTimeout);
  }
});
</script>

<template>
  <div class="app">
    <h1>Notification System</h1>
    <button @click="ws?.close()">Close</button>
    <p v-if="!isConnected">Reconnecting...</p>
    <ul>
      <li v-for="(notification, index) in notificationStore.notifications" :key="index">
        {{ notification }}
      </li>
    </ul>
  </div>
</template>
