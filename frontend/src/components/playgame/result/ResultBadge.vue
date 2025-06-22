<template>

  <span :class="$style.rank">
    <img :src="`/icons/${props.rank}.svg`" :alt="'Rank ' + props.rank" />
  </span>

</template>

  <script lang="ts" setup>

  const props = defineProps<{
    rank: number;
  }>();

  </script>

  <style module>

  .rank img {
    width: 30px;
    height: 30px;
  }

  </style>
