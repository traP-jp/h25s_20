<template>
  <span v-for="n in props.gold" :key="'gold-' + n" :class="$style.gold">
    <img :class="$style.gold" src="/icons/gold.svg" />
  </span>
  <span v-for="n in props.silver" :key="'silver-' + n" :class="$style.silver">
    <img :class="$style.silver" src="/icons/silver.svg" />
  </span>
  <span v-for="n in props.bronze" :key="'bronze-' + n" :class="$style.bronze">
    <img :class="$style.bronze" src="/icons/bronze.svg" />
  </span>
</template>

<script lang="ts" setup>
const props = defineProps<{
  gold: number;
  silver: number;
  bronze: number;
}>();
</script>

<style module>
.gold img {
  width: 20px;
  height: 20px;
  margin-right: -2px;
}

.silver img {
  width: 20px;
  height: 20px;
  margin-right: -2px;
}

.bronze img {
  width: 20px;
  height: 20px;
  margin-right: -2px;
}
</style>
