<template>
  <!-- start -->
  <OverlayModal v-if="showStartModal" @close="showStartModal = false">
    <!-- ここに Component を差し込む -->
    <div>start modal</div>
    <StartUserList />
    <StartControl @quit="showStartModal = false" />
  </OverlayModal>
  <button @click="showStartModal= true">startモーダルのテスト</button>
</template>

<script setup lang="ts">

import OverlayModal from "@/components/OverlayModal.vue";
import StartUserList from "@/components/playgame/start/StartUserList.vue";
import StartControl from "@/components/playgame/start/StartControl.vue";
import { ref } from "vue";
const showStartModal = ref(false);

</script>
