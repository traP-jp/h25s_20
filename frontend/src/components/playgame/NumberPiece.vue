<template>
  <div 
    :class="[$style.box, { [$style.changing]: isChanging }]" 
    :style="{ backgroundColor: colorMap[number] }"
  >
    {{ props.number }}
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{ number: number }>();

const isChanging = ref(false);

// 値が変わった時のアニメーション
watch(() => props.number, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    isChanging.value = true;
    setTimeout(() => {
      isChanging.value = false;
    }, 500); // アニメーション時間に合わせる
  }
});

const colorMap: Record<number, string> = {
  1: "#DF3E3E",
  2: "#F49F37",
  3: "#F2D500",
  4: "#277A33",
  5: "#5CA333",
  6: "#BD611F",
  7: "#F66D3B",
  8: "#A1E825",
  9: "#764040",
};
</script>

<style module>
.box {
  height: 40px;
  width: 40px;
  border-radius: 10px;
  color: white;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  transition: all 0.3s ease;
}

.changing {
  animation: valueChange 0.5s ease;
}

@keyframes valueChange {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.2);
    box-shadow: 0 0 20px rgba(255, 255, 255, 0.6);
  }
  100% {
    transform: scale(1);
  }
}
</style>
