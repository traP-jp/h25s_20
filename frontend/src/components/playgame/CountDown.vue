<template>
  <div v-if="time !== -1" class="countdown-display" :class="{ 'urgent': time <= 3 }">
    <div class="count" :data-time="time">{{ time }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, toRef } from "vue";
const props = defineProps<{ time: number }>();
const time = toRef(props, "time");
</script>

<style scoped>
.countdown-display {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #ff6b6b, #ee5a24);
  border: 3px solid #fff;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  z-index: 1005;
  pointer-events: none; /* タップ操作を通すため */
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 15px 25px;
  backdrop-filter: blur(10px);
  animation: countdown-pulse 1s ease-in-out infinite alternate;
}

.count {
  color: #fff;
  font-size: 4rem;
  font-weight: bold;
  text-shadow: 0 4px 8px rgba(0, 0, 0, 0.5);
  min-width: 80px;
  text-align: center;
}

@keyframes countdown-pulse {
  0% {
    transform: translateX(-50%) scale(1);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }
  100% {
    transform: translateX(-50%) scale(1.05);
    box-shadow: 0 12px 40px rgba(255, 107, 107, 0.4);
  }
}

/* 残り時間が少ない時の緊急感のあるスタイル */
.countdown-display.urgent {
  background: linear-gradient(135deg, #ff3838, #c44569);
  animation: countdown-urgent 0.5s ease-in-out infinite alternate;
}

@keyframes countdown-urgent {
  0% {
    transform: translateX(-50%) scale(1);
  }
  100% {
    transform: translateX(-50%) scale(1.1);
  }
}
</style>
