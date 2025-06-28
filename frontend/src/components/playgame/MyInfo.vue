<template>
  <div :class="$style.container">
    <div :class="$style.left">
      <div
        :class="$style.avatar"
        :style="{ backgroundImage: `url(${icon})` }"
      ></div>
      <div :class="$style.info">
        <div :class="$style.name">{{ name }}</div>
        <div :class="$style.score">{{ score }}点</div>
      </div>
    </div>
    <div :class="$style.right">
      <div
        :class="[
          $style.time,
          getTimeClass(time) && $style[getTimeClass(time)],
          { [$style.placeholder]: isTimePlaceholder(time) },
        ]"
      >
        LeftTime: {{ isTimePlaceholder(time) ? "xx:xx" : time }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  icon: string;
  name: string;
  score: number;
  time: string;
}>();
const { icon, name, score, time } = props;

// 時間文字列から数値を抽出して残り時間を判定
const getTimeClass = (timeStr: string) => {
  // プレースホルダー状態の場合はクラスを適用しない
  if (isTimePlaceholder(timeStr)) return "";

  let timeNum: number;

  // mm:ss形式の場合（例: "01:30" -> 90秒）
  if (timeStr.includes(":")) {
    const [minutes, seconds] = timeStr.split(":").map(Number);
    if (!isNaN(minutes) && !isNaN(seconds)) {
      timeNum = minutes * 60 + seconds;
    } else {
      return "";
    }
  } else {
    // 単純な数値の場合
    timeNum = parseInt(timeStr);
    if (isNaN(timeNum)) return "";
  }

  console.log(`Time: ${timeStr}, Seconds: ${timeNum}`); // デバッグ用

  if (timeNum <= 10) {
    console.log("Applying urgent-pulse class"); // デバッグ用
    return "urgent-pulse";
  } else if (timeNum <= 30) {
    console.log("Applying warning class"); // デバッグ用
    return "warning";
  }
  return "";
};

// 時間がプレースホルダー状態かどうかを判定
const isTimePlaceholder = (timeStr: string) => {
  return (
    !timeStr ||
    timeStr === "" ||
    timeStr === "0" ||
    timeStr === "xx:xx" ||
    isNaN(parseInt(timeStr))
  );
};
</script>

<style module>
.container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid var(--border-color, #ccc);
}
.left {
  display: flex;
  align-items: center;
}
.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  border: 2px dashed var(--border-color, #ccc);
  /* placeholder for icon */
}
.info {
  display: flex;
  flex-direction: column;
}
.name {
  font-weight: bold;
}
.score {
  color: var(--accent-color, #fff);
}
.right {
  display: flex;
  align-items: center;
}
.time {
  color: var(--accent-color, #fff);
  transition: color 0.3s ease;
}

/* プレースホルダー状態（データ未受信時） */
.placeholder {
  color: #cccccc !important;
  animation: none !important;
}

/* 30秒以下で赤色に */
.warning {
  color: #ff4757 !important;
}

/* 10秒以下で強調アニメーション（NumberPieceと同様の効果） */
.urgent-pulse {
  color: #ff3838 !important;
  animation: time-urgent 0.8s ease-in-out infinite !important;
  font-weight: bold !important;
}

@keyframes time-urgent {
  0% {
    transform: scale(1);
    opacity: 1;
    text-shadow: 0 0 0px rgba(255, 56, 56, 0);
  }
  25% {
    transform: scale(1.3);
    opacity: 0.9;
    text-shadow: 0 0 15px rgba(255, 56, 56, 0.8);
  }
  50% {
    transform: scale(1.2);
    opacity: 0.7;
    text-shadow: 0 0 25px rgba(255, 56, 56, 1);
  }
  75% {
    transform: scale(1.25);
    opacity: 0.9;
    text-shadow: 0 0 10px rgba(255, 56, 56, 0.6);
  }
  100% {
    transform: scale(1);
    opacity: 1;
    text-shadow: 0 0 0px rgba(255, 56, 56, 0);
  }
}

@keyframes subtle-blink {
  0%,
  80% {
    transform: scale(1);
  }
  40% {
    transform: scale(1.1);
  }
}
</style>
