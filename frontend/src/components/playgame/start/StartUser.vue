<template>
  <div :class="$style.container">
    <UserIcon :id="props.name" :size="40" />
    <span>{{ props.name }}</span>
    <StartBadge :gold="props.gold" :silver="props.silver" :bronze="props.bronze" />
  </div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import StartBadge from "@/components/playgame/start/StartBadge.vue";
import UserIcon from "@/components/UserIcon.vue";

const props = defineProps<{
  name: string;
  gold: number;
  silver: number;
  bronze: number;
}>();

</script>

<style module>
.container {
  border: 1px solid #ccc;
}
</style>
