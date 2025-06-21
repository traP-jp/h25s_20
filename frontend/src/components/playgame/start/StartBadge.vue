<template>

<span v-for="n in props.gold" :key="n" :class="$style.gold">●</span>
<span v-for="n in props.silver" :key="n" :class="$style.silver">●</span>
<span v-for="n in props.bronze" :key="n" :class="$style.bronze">●</span>

</template>

<script lang="ts" setup>

const props = defineProps<{
  gold: number;
  silver: number;
  bronze: number;
}>();

</script>

<style module>

.gold {
  color: gold;
}

.silver {
  color: silver;
}

.bronze {
  color: #cd7f32; /* ブロンズの色 */
}


</style>
