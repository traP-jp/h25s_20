<template>
  <!-- start -->
  <OverlayModal v-if="showStartModal" @close="showStartModal = false">
    <StartUserList :class="$style.userList" />
    <StartControl />
  </OverlayModal>
  <button @click="showStartModal = true">startモーダルのテスト</button>
</template>

<script setup lang="ts">
import OverlayModal from "@/components/OverlayModal.vue";
import StartUserList from "@/components/playgame/start/StartUserList.vue";
import StartControl from "@/components/playgame/start/StartControl.vue";

import { inject } from "vue";

const showStartModal = inject("showStartModal");
</script>

<style module>
.userList {
  margin-top: 60px;
  margin-bottom: 40px;
}
</style>
