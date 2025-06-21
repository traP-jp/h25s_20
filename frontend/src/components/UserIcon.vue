<script setup lang="ts">
import { computed, ref, watch } from "vue";

const props = defineProps<{ id?: string; size?: number }>();

const currentImageUrl = ref("https://api.dicebear.com/9.x/thumbs/svg?seed=");
const imageError = ref(false);

const preloadImage = (url: string): Promise<void> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve();
    img.onerror = () => reject();
    img.src = url;
  });
};

watch(
  () => props.id,
  async (newId) => {
    imageError.value = false;
    const primaryUrl = `https://q.trap.jp/api/v3/public/icon/${newId}`;
    const fallbackUrl = `https://api.dicebear.com/9.x/thumbs/svg?seed=${newId}`;

    try {
      await preloadImage(primaryUrl);
      currentImageUrl.value = primaryUrl;
    } catch {
      currentImageUrl.value = fallbackUrl;
    }
  },
  { immediate: true }
);

const imageStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  objectFit: "contain" as const,
  borderRadius: `${props.size || 0}px`,
}));
</script>

<template>
  <img :style="imageStyle" :src="currentImageUrl" />
</template>
