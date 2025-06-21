<script setup lang="ts">
import { computed, ref, watch } from "vue";

const props = defineProps<{ id?: string; size?: number }>();

const imageError = ref(false);

watch(
  () => props.id,
  (newId) => {
    if (newId) {
      imageError.value = false; // Reset error state when id changes
    }
  },
  { immediate: true }
);

const imageUrl = computed(() => {
  return imageError.value
    ? `https://api.dicebear.com/9.x/thumbs/svg?seed=${props.id}`
    : `https://q.trap.jp/api/v3/public/icon/${props.id}`;
});

const imageStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  objectFit: "contain" as const,
  borderRadius: `${props.size || 0}px`,
}));

const handleImageError = () => {
  imageError.value = true;
};
</script>

<template>
  <img :style="imageStyle" :src="imageUrl" @error="handleImageError" />
</template>
