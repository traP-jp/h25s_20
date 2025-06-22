<template>
  <div :class="$style.container">
    <ResultBadge :rank="props.rank" :class="$style.rank" />
    <UserIcon :id="props.name" :size="40" />
    <div :class="$style.userName">{{ props.name }}</div>
    <div :class="$style.score">{{ props.score }}点</div>
  </div>
</template>
<script setup lang="ts">
// import { ref } from "vue";
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
  width: 100%;
  display: flex;
  align-items: center;
}

.rank {
  margin-top: 5px;
  margin-right: 10px;
}

.userName {
  margin: 10px;
  font-weight: bold;
}
</style>
