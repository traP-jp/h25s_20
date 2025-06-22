<template>

  <span :class="$style.rank">{{ props.rank }}</span>

</template>

  <script lang="ts" setup>

  const props = defineProps<{
    rank: number;
  }>();

  </script>

  <style module>

  .rank {
    display: inline-block;
    background-color: blue;
    color: white;
    border-radius: 50%;
    width: 2em;
    height: 2em;
    line-height: 2em;
    text-align: center;
    font-weight: bold;
  }

  </style>
