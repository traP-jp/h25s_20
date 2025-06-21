<template>
  <div :class="$style.user">
    <div :class="$style.avatar"></div>
    <div :class="$style.info">
      <div :class="$style.name">{{ name }}</div>
      <div :class="$style.score">{{ score }}点</div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ icon: string; name: string; score: number }>();
const { icon, name, score } = props;
</script>

<style module>
.user {
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
</style>
