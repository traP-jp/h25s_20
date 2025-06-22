<template>
  <div :class="$style.wrapper">
    <!-- プレビュー -->
    <div :class="$style.preview">
      {{ viewExpression }}
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
          <MathInputButton
            icon="mdi:backspace-outline"
            @click="backspace"
            @long-press="clearAll"
            :show-long-press-style="true"
          />
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
import { computed, ref, onMounted, onUnmounted, inject } from "vue";
import MathInputButton from "./MathInputButton.vue";
import { defineModel, watch } from "vue";
import { checkMath } from "@/lib/board-update";
import { apiClient } from "@/api";

const expression = ref("");
const board = defineModel<number[]>("board");
const version = ref(0); // フォーミュラのバージョン管理

// 親コンポーネントからroomを注入
const currentRoom = inject<any>("currentRoom");

const handleKeydown = (event: KeyboardEvent) => {
  // 数字キー (1-9)
  if (/^[1-9]$/.test(event.key)) {
    addSymbol(event.key);
    event.preventDefault();
    return;
  }

  // 演算子キー
  if (/^[+\-*/]$/.test(event.key)) {
    addSymbol(event.key);
    event.preventDefault();
    return;
  }

  // 括弧キー
  if (event.key === "(" || event.key === ")") {
    addParentheses();
    event.preventDefault();
    return;
  }

  // バックスペースキー
  if (event.key === "Backspace") {
    backspace();
    event.preventDefault();
    return;
  }

  // ESCキーで全消し
  if (event.key === "Escape") {
    clearAll();
    event.preventDefault();
    return;
  }
};

onMounted(() => {
  window.addEventListener("keydown", handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeydown);
});

watch(expression, async (newValue) => {
  if (!board.value) return;
  // 一時的に checkMath で検証のみ行う（board は更新しない）
  const tempBoard = [...board.value]; // 元のboardのコピーを作成
  const result = checkMath(tempBoard, newValue);

  // 結果が10で、inputが空文字列（消せる列が存在）の場合、自動提出
  if (result["input"] === "" && newValue.trim() !== "" && currentRoom.value) {
    console.log("10に到達し、消せる列が存在します。自動提出します。");

    try {
      const submission = {
        version: version.value,
        formula: newValue,
      };

      const response = await apiClient.submitFormula(currentRoom.value.roomId, submission);

      if (response.success) {
        console.log("数式が正常に提出されました:", response.data);
        version.value++; // バージョンを更新
        // 提出成功時は入力をクリアし、バックエンドからのWebSocket更新を待つ
        expression.value = "";
      } else {
        console.error("数式の提出に失敗しました:", response.data);
        // 提出失敗時は入力をそのまま残す
      }
    } catch (error) {
      console.error("数式提出時にエラーが発生しました:", error);
      // エラー時も入力をそのまま残す
    }
  }
  // 提出に該当しない場合は何もしない（expressionの値はそのまま）
});

const viewExpression = computed(() => {
  return expression.value.replace("-", "−").replace("*", "×").replace("/", "÷");
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
      if (expression.value.length !== 0) {
        expression.value = expression.value.slice(0, -1) + value;
      }
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

const clearAll = () => {
  expression.value = "";
};
</script>

<style module>
.wrapper {
  width: 100%;
  margin-bottom: 20px;
}

.preview {
  font-family: "M PLUS Code Latin", sans-serif;
  letter-spacing: 5px;
  border-radius: 8px;
  padding: 5px;
  margin-bottom: 5px;
  font-size: 30px;
  font-weight: 500;
  text-align: center;
  min-height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: black;
}

.container {
  margin: 0 auto;
  width: 300px;
  display: flex;
  justify-content: space-between;
}

.tens {
  height: 170px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.tensRow {
  width: 170px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.symbolRow {
  width: 110px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}
</style>
