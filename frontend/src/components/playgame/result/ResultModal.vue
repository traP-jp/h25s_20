<template>
  <!-- start -->
  <OverlayModal v-if="showResultModal" @close="showResultModal = false">
    <!-- ここに Component を差し込む -->
    <div>result modal</div>
    <ResultUserList :class="$style.userList" />
    <ResultControl v-model:showResultModal="showResultModal" v-model:showStartModal="showStartModal" />
  </OverlayModal>
</template>

<script setup lang="ts">
import OverlayModal from "@/components/OverlayModal.vue";
import ResultUserList from "@/components/playgame/result/ResultUserList.vue";
import ResultControl from "@/components/playgame/result/ResultControl.vue";
import { defineModel } from "vue";

const showStartModal = defineModel<boolean>("showStartModal");
const showResultModal = defineModel<boolean>("showResultModal");
</script>

<style module>
.userList {
  margin: 20px auto;
}
</style>
