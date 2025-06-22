<template>
  <!-- start -->
  <OverlayModal v-if="showResultModal" @close="showResultModal = false">
    <!-- ここに Component を差し込む -->
    <div>result modal</div>
    <ResultUserList />
    <ResultControl />
  </OverlayModal>
  <button @click="showResultModal = true">resultモーダルのテスト</button>
</template>

<script setup lang="ts">
import OverlayModal from "@/components/OverlayModal.vue";
import ResultUserList from "@/components/playgame/result/ResultUserList.vue";
import ResultControl from "@/components/playgame/result/ResultControl.vue";
import { ref, inject } from "vue";

const showStartModal = inject("showStartModal");
const showResultModal = inject("showResultModal");
</script>
