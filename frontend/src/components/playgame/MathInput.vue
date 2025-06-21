<template>
  <div :class="$style.container">
    <div :class="$style.tens">
      <div :class="$style.tensRow">
        <MathInputButton text="1" />
        <MathInputButton text="2" />
        <MathInputButton text="3" />
      </div>
      <div :class="$style.tensRow">
        <MathInputButton text="4" />
        <MathInputButton text="5" />
        <MathInputButton text="6" />
      </div>
      <div :class="$style.tensRow">
        <MathInputButton text="7" />
        <MathInputButton text="8" />
        <MathInputButton text="9" />
      </div>
    </div>
    <div :class="$style.tens">
      <div :class="$style.symbolRow">
        <MathInputButton icon="ic:outline-plus" />
        <MathInputButton icon="ic:outline-minus" />
      </div>
      <div :class="$style.symbolRow">
        <MathInputButton icon="ic:sharp-close" />
        <MathInputButton icon="ic:sharp-divide" />
      </div>
      <div :class="$style.symbolRow">
        <MathInputButton text="( )" />
        <MathInputButton icon="ic:outline-backspace" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import MathInputButton from "./MathInputButton.vue";
</script>

<style module>
.container {
  margin: 30px auto 0 auto;
  width: 350px;
  display: flex;
  justify-content: space-between;
}

.tens {
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.tensRow {
  width: 200px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.symbolRow {
  width: 130px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}
</style>
