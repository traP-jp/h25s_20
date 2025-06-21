<template>
  <div :class="$style.container">
    <div :class="$style.left">
      <div
        :class="$style.avatar"
        :style="{ backgroundImage: `url(${icon})` }"
      ></div>
      <div :class="$style.info">
        <div :class="$style.name">{{ name }}</div>
        <div :class="$style.score">{{ score }}点</div>
      </div>
    </div>
    <div :class="$style.right">
      <div :class="$style.time">LeftTime: {{ time }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  icon: string;
  name: string;
  score: number;
  time: string;
}>();
const { icon, name, score, time } = props;
</script>

<style module>
.container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid var(--border-color, #ccc);
}
.left {
  display: flex;
  align-items: center;
}
.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  border: 2px dashed var(--border-color, #ccc);
  /* placeholder for icon */
}
.info {
  display: flex;
  flex-direction: column;
}
.name {
  font-weight: bold;
}
.score {
  color: var(--accent-color, #fff);
}
.right {
  display: flex;
  align-items: center;
}
.time {
  color: var(--accent-color, #fff);
}
</style>
