<template>
  <div :class="$style.gameboard">
    <div :class="$style.grid">
      <NumberPiece
        v-for="(num, idx) in board"
        :key="idx"
        :number="num"
        :is-highlighted="highlightedNumbers.includes(num)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineModel, defineProps } from "vue";
import NumberPiece from "@/components/playgame/NumberPiece.vue";

const board = defineModel<number[]>("board");

defineProps<{
  highlightedNumbers: number[];
}>();
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
