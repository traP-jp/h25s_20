<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";

const isConnected = ref(false);
const ws = ref<WebSocket | null>(null);
const messages = ref<string[]>([]);

onMounted(() => {
  ws.value = new WebSocket("ws://localhost:8080/ws?username=debug_user");

  ws.value.onopen = () => {
    console.log("WebSocket opened");
    isConnected.value = true;
  };

  ws.value.onerror = (event) => {
    console.error("WebSocket error", event);
  };

  ws.value.onmessage = (event) => {
    console.log(event.data);
    try {
      const data = JSON.parse(event.data);
      if (data.event === "connection") {
        messages.value.push(`Connected: ${JSON.stringify(data.content)}`);
      }
    } catch {
      messages.value.push(event.data);
    }
  };

  ws.value.onclose = () => {
    console.log("WebSocket closed");
    isConnected.value = false;
  };
});

onBeforeUnmount(() => {
  ws.value?.close();
});
</script>

<template>
  <div>
    <h1>WebSocket 接続テスト</h1>
    <p v-if="isConnected">接続済み</p>
    <p v-else>未接続</p>
    <ul>
      <li v-for="(msg, idx) in messages" :key="idx">{{ msg }}</li>
    </ul>
  </div>
</template>
