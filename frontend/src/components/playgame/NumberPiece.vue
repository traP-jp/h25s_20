<template>
  <div :class="$style.box">
    {{ props.number }}
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ number: number }>();
</script>

<style module>
.box {
  height: 50px;
  width: 50px;
  border-radius: 10px;
  background-color: #cccccc;
  border: 2px solid #000000;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}
</style>
