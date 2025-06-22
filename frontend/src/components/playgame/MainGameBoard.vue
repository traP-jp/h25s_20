<template>
  <div :class="$style.gameboard">
    <div :class="$style.grid">
      <TransitionGroup name="piece" tag="div" :class="$style.gridContent">
        <NumberPiece
          v-for="(num, idx) in board"
          :key="`${idx}-${num}`"
          :number="num"
          :class="$style.piece"
        />
      </TransitionGroup>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineModel } from "vue";
import NumberPiece from "@/components/playgame/NumberPiece.vue";

const board = defineModel<number[]>("board");
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

.gridContent {
  display: contents;
}

.piece {
  transition: all 0.3s ease;
}

/* アニメーション定義 */
:global(.piece-enter-active),
:global(.piece-leave-active) {
  transition: all 0.3s ease;
}

:global(.piece-enter-from) {
  opacity: 0;
  transform: scale(0.8);
}

:global(.piece-leave-to) {
  opacity: 0;
  transform: scale(0.8);
}

:global(.piece-move) {
  transition: transform 0.3s ease;
}
</style>
