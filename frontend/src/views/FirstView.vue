<template>
  <div class="websocket-view">
    <h2>Vue.js WebSocket Tutorial</h2>
    <button @click="sendMessage('hello')">Send Message</button>
  </div>
</template>

<script setup lang="ts">
// https://tutorialedge.net/javascript/vuejs/vuejs-websocket-tutorial/

import { ref, onMounted } from "vue";

const connection = ref<WebSocket | null>(null);

const sendMessage = (message: string): void => {
  console.log(connection.value);
  if (connection.value) {
    connection.value.send(message);
  }
};

onMounted(() => {
  console.log("Starting connection to WebSocket Server");
  connection.value = new WebSocket("wss://echo.websocket.org");

  connection.value.onmessage = (event: MessageEvent) => {
    console.log(event);
    console.log(event.data);
  };

  connection.value.onopen = (event: Event) => {
    console.log(event);
    console.log("Successfully connected to the echo websocket server...");
  };

  // 他に onclose, onerror イベントも設定できる
});
</script>

<style scoped>
.websocket-view {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
