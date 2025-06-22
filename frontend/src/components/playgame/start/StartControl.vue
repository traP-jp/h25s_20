<template>
  <div :class="$style.container">
    <button
      :class="$style.mainBtn"
      @click="onClickMain"
      :disabled="isBtnDisabled"
    >
      {{ BtnMsg }}
    </button>
    <button
      v-show="!isBtnDisabled"
      :class="$style.quitBtn"
      @click="onClickQuit"
    >
      部屋から抜ける
    </button>
    <button
      v-show="isBtnDisabled"
      :class="$style.cancelBtn"
      @click="onClickCancel"
    >
      キャンセル
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, inject } from "vue";

const showStartModal = inject("showStartModal");
const showResultModal = inject("showResultModal");

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
  showStartModal.value = false;
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

.quitBtn {
  align-self: flex-start;
}

.cancelBtn {
  align-self: flex-start;
}
</style>
