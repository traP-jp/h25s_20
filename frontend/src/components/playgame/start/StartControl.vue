<template>
  <button :class="$style.mainBtn" @click="onClickMain" :disabled="isBtnDisabled">{{ BtnMsg }}</button>
  <button v-show="!isBtnDisabled" :class="$style.quitBtn" @click="onClickQuit">部屋から抜ける</button>
  <button v-show="isBtnDisabled" :class="$style.cancelBtn" @click="onClickCancel">キャンセル</button>
</template>
<script setup lang="ts">
import { ref } from "vue";

const emit = defineEmits(['quit']);

const BtnMsg = ref("準備OK!");
const isBtnDisabled = ref(false);

const onClickMain = () => {
  // send an event to backend
  // change message
  BtnMsg.value = "他のプレイヤーを待っています...";
  isBtnDisabled.value = true;
};

const onClickQuit = () => {
  // send an event to backend
  // close modal
  emit('quit');
};

const onClickCancel = () => {
  // send an event to backend
  // change message
  BtnMsg.value = "準備OK!";
  isBtnDisabled.value = false;
};
</script>

<style module>
.container {
  ;
}
</style>
