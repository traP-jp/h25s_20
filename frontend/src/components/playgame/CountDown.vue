<template>
  <div v-if="time !== -1 && isStartCountdown" class="countdown-overlay">
    <!-- 背景オーバーレイ -->
    <div class="backdrop"></div>
    <!-- カウントダウン表示 -->
    <div class="countdown-display">
      <div class="count">{{ time }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, toRef } from "vue";
const props = defineProps<{
  time: number;
  isStartCountdown?: boolean; // スタート時のカウントダウンかどうか
}>();
const time = toRef(props, "time");
const { isStartCountdown = true } = props; // デフォルトはスタート時
</script>

<style scoped>
.countdown-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000;
  pointer-events: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.backdrop {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  z-index: 1001;
}

.countdown-display {
  position: relative;
  z-index: 1002;
  display: flex;
  align-items: center;
  justify-content: center;
}

.count {
  color: #fff;
  font-size: 8rem;
  font-weight: bold;
  text-align: center;
}
</style>
