<template>
  <div :class="$style.container">
    <button :class="$style.mainBtn" @click="onClickMain">Got it!</button>
  </div>
</template>
<script setup lang="ts">
import { inject } from "vue";

const showStartModal = inject("showStartModal") as any;
const showResultModal = inject("showResultModal") as any;

const onClickMain = () => {
  // send an event to backend
  // close this modal and open start modal
  showResultModal.value = false;
  showStartModal.value = true;
};
</script>

<style module>
.mainBtn {
  transform: scale(1.5);
  margin: 20px 0;
}

.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}
</style>
