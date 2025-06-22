<template>
  <div
    :class="[
      $style.box,
      {
        [$style.highlighted]: isHighlighted,
        [$style.changed]: isAnimating,
      },
    ]"
    :style="{
      backgroundColor: colorMap[number],
      borderColor: isHighlighted ? colorMap[number] : 'transparent',
      boxShadow: isHighlighted
        ? `0 0 5px ${colorMap[number]}cc, 0 0 8px ${colorMap[number]}`
        : 'none',
    }"
  >
    {{ props.number }}
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

const props = defineProps<{
  number: number;
  isHighlighted?: boolean;
}>();

const isAnimating = ref(false);
const previousNumber = ref(props.number);

// 数字が変更された時にアニメーションを実行
watch(
  () => props.number,
  (newNumber, oldNumber) => {
    if (newNumber !== oldNumber) {
      isAnimating.value = true;
      setTimeout(() => {
        isAnimating.value = false;
      }, 600); // アニメーション時間と合わせる
    }
    previousNumber.value = newNumber;
  }
);

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
  transition: all 0.2s ease;
  box-sizing: border-box;
  border: 3px solid transparent;
}

.changed {
  animation: numberChange 0.6s ease-in-out;
}

@keyframes numberChange {
  0% {
    transform: scale(1) rotate(0deg);
    opacity: 1;
  }
  25% {
    transform: scale(1.1) rotate(5deg);
    opacity: 0.8;
    box-shadow: 0 0 20px rgba(255, 255, 255, 0.8);
  }
  50% {
    transform: scale(1.1) rotate(-5deg);
    opacity: 0.6;
    box-shadow: 0 0 30px rgba(255, 255, 255, 1);
  }
  75% {
    transform: scale(1.1) rotate(3deg);
    opacity: 0.9;
    box-shadow: 0 0 15px rgba(255, 255, 255, 0.6);
  }
  100% {
    transform: scale(1) rotate(0deg);
    opacity: 1;
    box-shadow: none;
  }
}
</style>
