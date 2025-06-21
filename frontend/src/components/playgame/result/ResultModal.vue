<template>
  <!-- start -->
  <OverlayModal v-if="showResultModal" @close="showResultModal = false">
    <!-- ここに Component を差し込む -->
    <div>result modal</div>
  </OverlayModal>
  <button @click="showResultModal = true">resultモーダルのテスト</button>
</template>

<script setup lang="ts">
import OverlayModal from "@/components/OverlayModal.vue";
import { ref } from "vue";
const showResultModal = ref(false);
</script>
