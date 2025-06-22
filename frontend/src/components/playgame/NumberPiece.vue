<template>
  <div :class="$style.box" :style="{ backgroundColor: colorMap[number] }">
    {{ props.number }}
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ number: number }>();

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
}
</style>
