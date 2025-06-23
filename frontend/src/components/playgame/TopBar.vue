<template>
  <div :class="$style.container">
    <div :class="$style.content">
      <UserIcon :class="$style.myIcon" :id="currentUsername" :size="40" />
      <div :class="$style.myName">{{ currentUsername }}</div>
      <div :class="$style.right">
        <TextMark :text="room?.roomName || ''" :bgColor="`#008800`" />
        <img src="/logo.svg" alt="Logo" :class="$style.logo" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Room } from "@/lib/types.ts";
import UserIcon from "@/components/UserIcon.vue";
import TextMark from "@/components/TextMark.vue";
import { useWebSocketStore } from "@/store";
import { computed } from "vue";

// WebSocketストアからユーザー名を取得
const webSocketStore = useWebSocketStore();

// WebSocketStoreのcurrentUsernameが空の場合はlocalStorageから取得
const currentUsername = computed(() => {
  return webSocketStore.currentUsername || sessionStorage.getItem("username") || "Unknown";
});

// propsからroomを受け取る
defineProps<{
  room: Room | null;
}>();
</script>

<style module>
.container {
  background-color: #ff88bb;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 10;
  width: 100vw;
  height: 40px;
}

.content {
  position: relative;
  width: min(360px, 100%);
  height: 100%;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.myIcon {
  z-index: 1;
  border: 3px solid #ff88bb;
  position: absolute;
  left: 15px;
  top: 10px;
}

.myName {
  margin-top: 6px;
  margin-left: 70px;
  font-weight: bold;
  font-size: 15px;
}

.right {
  height: 100%;
  display: flex;
  align-items: center;
}

.logo {
  height: 18px;
  margin: 0 10px;
}
</style>
