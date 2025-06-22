<template>
  <!-- start -->
  <OverlayModal v-if="showStartModal" @close="showStartModal = false">
    <StartUserList :class="$style.userList" />
    <StartControl v-model:showResultModal="showStartModal" />
  </OverlayModal>
</template>

<script setup lang="ts">
import OverlayModal from "@/components/OverlayModal.vue";
import StartUserList from "@/components/playgame/start/StartUserList.vue";
import StartControl from "@/components/playgame/start/StartControl.vue";
import { defineModel } from "vue";

const showStartModal = defineModel<boolean>("showStartModal");
</script>

<style module>
.userList {
  margin-top: 60px;
  margin-bottom: 40px;
}
</style>
