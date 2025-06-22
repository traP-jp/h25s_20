<template>
  <div
    :class="[
      $style.container,
      { [$style.longPressing]: isLongPressing && showLongPressStyle },
    ]"
    @click="handleClick"
    @mousedown="startLongPress"
    @mouseup="endLongPress"
    @mouseleave="endLongPress"
    @touchstart="startLongPress"
    @touchend="endLongPress"
  >
    <Icon v-if="icon" :icon="icon" :class="$style.icon" />
    <span v-else :class="$style.text">{{ text }}</span>
  </div>
</template>

<script setup lang="ts">
import { Icon } from "@iconify/vue";
import { ref } from "vue";

defineProps<{
  text?: string;
  icon?: string;
  showLongPressStyle?: boolean;
}>();

const emit = defineEmits<{
  click: [];
  longPress: [];
}>();

const isLongPressing = ref(false);
let longPressTimer: number | null = null;

const startLongPress = () => {
  isLongPressing.value = true;
  longPressTimer = window.setTimeout(() => {
    emit("longPress");
  }, 500);
};

const endLongPress = () => {
  isLongPressing.value = false;
  if (longPressTimer) {
    clearTimeout(longPressTimer);
    longPressTimer = null;
  }
};

const handleClick = () => {
  if (!isLongPressing.value) {
    emit("click");
  }
};
</script>

<style module>
.container {
  width: 50px;
  height: 50px;
  border-radius: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f0f0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.container:hover {
  background-color: #88aaff;
  transform: scale(1.05);
}

.container:active {
  transform: scale(0.95);
}

.longPressing {
  background-color: #ff6b6b !important;
  transform: scale(1.1) !important;
  box-shadow: 0 0 15px rgba(255, 107, 107, 0.5);
}

.text {
  color: black;
  font-size: 28px;
  font-weight: bold;
}

.icon {
  color: black;
  font-size: 32px;
}
</style>
