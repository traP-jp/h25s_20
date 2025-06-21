<template>
  <div :class="$style.container">
    <StartBadge :rank="props.rank"/>
    <UserIcon :id="props.name" :size="40" />
    <span>{{ props.name }}</span>
    <span>{{ props.score }}点</span>
  </div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import ResultBadge from "@/components/playgame/result/ResultBadge.vue";
import UserIcon from "@/components/UserIcon.vue";

const props = defineProps<{
  name: string;
  rank: number;
  score: number;
}>();

</script>

<style module>
.container {
  border: 1px solid #ccc;
}
</style>
