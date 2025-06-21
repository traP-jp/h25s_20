<template>
  <div :class="$style.wrapper">
    <!-- プレビュー -->
    <div :class="[$style.preview, { [$style.correct]: isCorrect }]">
      {{ expression }}
    </div>

    <div :class="$style.container">
      <div :class="$style.tens">
        <div :class="$style.tensRow">
          <MathInputButton icon="mdi:numeric-1" @click="addToExpression('1')" />
          <MathInputButton icon="mdi:numeric-2" @click="addToExpression('2')" />
          <MathInputButton icon="mdi:numeric-3" @click="addToExpression('3')" />
        </div>
        <div :class="$style.tensRow">
          <MathInputButton icon="mdi:numeric-4" @click="addToExpression('4')" />
          <MathInputButton icon="mdi:numeric-5" @click="addToExpression('5')" />
          <MathInputButton icon="mdi:numeric-6" @click="addToExpression('6')" />
        </div>
        <div :class="$style.tensRow">
          <MathInputButton icon="mdi:numeric-7" @click="addToExpression('7')" />
          <MathInputButton icon="mdi:numeric-8" @click="addToExpression('8')" />
          <MathInputButton icon="mdi:numeric-9" @click="addToExpression('9')" />
        </div>
      </div>
      <div :class="$style.tens">
        <div :class="$style.symbolRow">
          <MathInputButton text="( )" @click="addParentheses" />
          <MathInputButton icon="mdi:backspace-outline" @click="backspace" />
        </div>
        <div :class="$style.symbolRow">
          <MathInputButton icon="mdi:plus" @click="addToExpression('+')" />
          <MathInputButton icon="mdi:minus" @click="addToExpression('-')" />
        </div>
        <div :class="$style.symbolRow">
          <MathInputButton icon="mdi:close" @click="addToExpression('*')" />
          <MathInputButton icon="mdi:division" @click="addToExpression('/')" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import MathInputButton from "./MathInputButton.vue";

const expression = ref("");

import { solvePoland } from "@/lib/solver/solvePoland";
import { encodePoland } from "@/lib/solver/encodePoland";

const isCorrect = computed(() => {
  try {
    console.log(solvePoland(encodePoland(expression.value)));
    return solvePoland(encodePoland(expression.value)) === "10";
  } catch {
    return false;
  }
});

const addToExpression = (value: string) => {
  // 数字の場合のみ制限をチェック
  if (/\d/.test(value)) {
    // 現在の式から数字のみを抽出してカウント
    const digitCount = (expression.value.match(/\d/g) || []).length;

    // すでに4つの数字がある場合は追加しない
    if (digitCount >= 4) {
      return;
    }
  }

  expression.value += value;
};

const addParentheses = () => {
  const openParens = (expression.value.match(/\(/g) || []).length;
  const closeParens = (expression.value.match(/\)/g) || []).length;

  if (openParens > closeParens) {
    expression.value += ")";
  } else {
    expression.value += "(";
  }
};

const backspace = () => {
  expression.value = expression.value.slice(0, -1);
};
</script>

<style module>
.wrapper {
  width: 100%;
}

.preview {
  letter-spacing: 5px;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
  font-size: 30px;
  font-weight: 500;
  text-align: center;
  min-height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.correct {
  color: #3b82f6; /* 青色 */
}

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
