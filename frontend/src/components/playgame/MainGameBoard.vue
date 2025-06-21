<template>
  <div :class="$style.gameboard">
    <div :class="$style.grid">
      <NumberPiece v-for="(num, idx) in numbers" :key="idx" :number="num" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import NumberPiece from "@/components/playgame/NumberPiece.vue";

const numbers = ref<number[]>([2, 8, 5, 3, 1, 6, 9, 4, 7, 2, 5, 8, 3, 1, 6, 9]);
</script>

<style module>
.gameboard {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}
.grid {
  display: grid;
  grid-template-columns: repeat(4, auto);
  grid-gap: 10px;
}
</style>
