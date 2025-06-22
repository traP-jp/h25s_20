<template>
  <div :class="$style.container">
    <button :class="$style.mainBtn" @click="onClickMain">Got it!</button>
    <!-- Placeholder buttons for layout consistency -->
    <div :class="$style.placeholder"></div>
  </div>
</template>
<script setup lang="ts">
import { defineModel } from "vue";

const showStartModal = defineModel<boolean>("showStartModal");
const showResultModal = defineModel<boolean>("showResultModal");

const onClickMain = () => {
  // send an event to backend
  // close this modal and open start modal
  showResultModal.value = false;
  showStartModal.value = true;
};
</script>

<style module>
.container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  background-color: white;
  padding: 20px;
  z-index: 1000;
}

.mainBtn {
  transform: scale(1.5);
  margin: 20px 0;
}

.placeholder {
  height: 40px; /* Same height as quit/cancel buttons */
  align-self: flex-start;
}
</style>
