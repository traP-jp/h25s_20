<template>
  <div :class="$style.container">
    <ReadyIcon :id="props.name" :size="40" :isReady="props.isReady" />
    <span>{{ props.name }}</span>
    <StartBadge :gold="props.gold" :silver="props.silver" :bronze="props.bronze" />
  </div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import StartBadge from "@/components/playgame/start/StartBadge.vue";
import ReadyIcon from "@/components/ReadyIcon.vue";

const props = defineProps<{
  name: string;
  gold: number;
  silver: number;
  bronze: number;
  isReady: boolean;
}>();

</script>

<style module>
.container {
  border: 1px solid #ccc;
}
</style>
