<script setup lang="ts">
import UserIcon from "@/components/UserIcon.vue";
import { Icon } from "@iconify/vue";

const props = withDefaults(
  defineProps<{
    id?: string;
    size?: number;
    border?: string;
    isReady: boolean;
  }>(),
  {
    isReady: false, // デフォルト値を設定
  }
);

console.log(props.isReady);
</script>

<template>
  <div :class="$style.iconContainer">
    <UserIcon
      :id="props.id"
      :size="props.size"
      :style="{
        // borderColor: props.borderColor || 'transparent',
        // borderWidth: props.borderWidth || '0px',
        border: props.border || '0px solid transparent',
      }"
      :class="$style.userIcon"
    />
    <Icon
      v-show="props.isReady"
      :icon="'fe:check'"
      :class="$style.checkIcon"
      :style="{ border: props.border || '0px solid transparent' }"
    />
  </div>
</template>

<style module>
.iconContainer {
  position: relative;
  display: inline-block;
  margin-right: 10px;
}

.checkIcon {
  color: green;
  font-weight: bold;
  position: absolute;
  top: 0;
  left: 0;
  transform: translate(105%, 105%); /* 左上に微調整 */
  font-size: 1rem; /* サイズ調整 */
  background: white; /* 必要なら背景色を追加 */
  border-radius: 50%; /* 丸くしたい場合 */
  width: 20px; /* バッジのサイズ */
  height: 20px; /* バッジのサイズ */
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
