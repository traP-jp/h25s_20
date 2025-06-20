<template>
  <div :class="$style.wrapper">
    <!-- プレビュー -->
    <div :class="[$style.preview]">
      {{ expression }}
    </div>

    <div :class="$style.container">
      <div :class="$style.tens">
        <div :class="$style.tensRow">
          <MathInputButton icon="mdi:numeric-1" @click="addSymbol('1')" />
          <MathInputButton icon="mdi:numeric-2" @click="addSymbol('2')" />
          <MathInputButton icon="mdi:numeric-3" @click="addSymbol('3')" />
        </div>
        <div :class="$style.tensRow">
          <MathInputButton icon="mdi:numeric-4" @click="addSymbol('4')" />
          <MathInputButton icon="mdi:numeric-5" @click="addSymbol('5')" />
          <MathInputButton icon="mdi:numeric-6" @click="addSymbol('6')" />
        </div>
        <div :class="$style.tensRow">
          <MathInputButton icon="mdi:numeric-7" @click="addSymbol('7')" />
          <MathInputButton icon="mdi:numeric-8" @click="addSymbol('8')" />
          <MathInputButton icon="mdi:numeric-9" @click="addSymbol('9')" />
        </div>
      </div>
      <div :class="$style.tens">
        <div :class="$style.symbolRow">
          <MathInputButton icon="mdi:code-parentheses" @click="addParentheses" />
          <MathInputButton icon="mdi:backspace-outline" @click="backspace" />
        </div>
        <div :class="$style.symbolRow">
          <MathInputButton icon="mdi:plus" @click="addSymbol('+')" />
          <MathInputButton icon="mdi:minus" @click="addSymbol('-')" />
        </div>
        <div :class="$style.symbolRow">
          <MathInputButton icon="mdi:close" @click="addSymbol('*')" />
          <MathInputButton icon="mdi:division" @click="addSymbol('/')" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import MathInputButton from "./MathInputButton.vue";
import { defineModel, watch } from "vue";
import { checkMath } from "@/lib/board-update";

const expression = ref("");
const board = defineModel<number[]>("board");

watch(expression, (newValue) => {
  if (!board.value) return;
  const result = checkMath(board.value, newValue);
  board.value = result["board"];
  expression.value = result["input"];
});

const addSymbol = (value: string) => {
  const last = expression.value.length > 0 ? expression.value[expression.value.length - 1] : "+";

  if (/[1-9]/.test(last)) {
    if (/[1-9]/.test(value)) {
      expression.value = expression.value.slice(0, -1) + value;
      return;
    } else if (/[+\-*/]/.test(value)) {
      const numberCount = (expression.value.match(/[1-9]/g) || []).length;
      if (numberCount >= 4) {
        return;
      }
      expression.value += value;
      return;
    }
  } else if (/[+\-*/]/.test(last)) {
    if (/[1-9]/.test(value)) {
      expression.value += value;
      return;
    } else if (/[+\-*/]/.test(value)) {
      expression.value = expression.value.slice(0, -1) + value;
      return;
    }
  } else if (last === "(") {
    if (/[1-9]/.test(value)) {
      expression.value += value;
      return;
    } else if (/[+\-*/]/.test(value)) {
      return;
    }
  } else if (last === ")") {
    if (/[1-9]/.test(value)) {
      return;
    } else if (/[+\-*/]/.test(value)) {
      expression.value += value;
      return;
    }
  }
};

const addParentheses = () => {
  const openParens = (expression.value.match(/\(/g) || []).length;
  const closeParens = (expression.value.match(/\)/g) || []).length;

  console.log("openParens: ", openParens);
  console.log("closeParens: ", closeParens);

  const last = expression.value.length > 0 ? expression.value[expression.value.length - 1] : "+";
  if (/[1-9]/.test(last) || last === ")") {
    if (openParens > closeParens) {
      expression.value += ")";
    }
    return;
  } else if (/[+\-*/]/.test(last) || last === "(") {
    if (openParens < 2) {
      expression.value += "(";
    }
    return;
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
  color: black;
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
