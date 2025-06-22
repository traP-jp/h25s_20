<template>
  <div
    :class="[$style.box, { [$style.highlighted]: isHighlighted }]"
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
const props = defineProps<{
  number: number;
  isHighlighted?: boolean;
}>();

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
</style>
