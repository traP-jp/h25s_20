<template>
  <button :class="$style.mainBtn" @click="onClickMain" >Got it!</button>
</template>
<script setup lang="ts">

import { ref, inject } from "vue";

const showStartModal = inject('showStartModal');
const showResultModal = inject('showResultModal');

const onClickMain = () => {
  // send an event to backend
  // close this modal and open start modal
  showResultModal.value = false;
  showStartModal.value = true;
};

</script>

<style module>
.container {
  ;
}
</style>
