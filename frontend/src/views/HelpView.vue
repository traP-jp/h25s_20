<template>
  <div :class="$style.container">
    <img :class="$style.logo" src="/logo.svg" alt="Logo" />
    <div :class="$style.slide">
      <img
        v-for="(slide, index) in availableSlides"
        :key="index"
        :src="slide"
        :alt="`Help slide ${index}`"
        :class="$style.slideImage"
        @error="onImageError(index)"
      />
    </div>

    <button :class="$style.button" @click="onClick">
      ゲームに戻る
    </button>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const availableSlides = ref<string[]>([]);

const onClick = () => {
  router.push("/rooms");
};

const checkImageExists = async (url: string): Promise<boolean> => {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => resolve(true);
    img.onerror = () => resolve(false);
    img.src = url;
  });
};

const loadSlides = async () => {
  const slides: string[] = [];
  let index = 0;

  while (true) {
    const slideUrl = `/slides/${index}.png`;
    const exists = await checkImageExists(slideUrl);

    if (!exists) {
      break;
    }

    slides.push(slideUrl);
    index++;
  }

  availableSlides.value = slides;
};

const onImageError = (index: number) => {
  console.warn(`Failed to load slide ${index}`);
};

onMounted(() => {
  loadSlides();
});
</script>

<style module>
.container {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.logo {
  width: 300px;
  margin-bottom: 20px;
}

.slide {
  width: 100%;
  max-width: 1200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 30px;
  margin-bottom: 40px;
  overflow-y: auto;
  padding: 20px;
}

.slideImage {
  width: 90%;
  max-width: 1000px;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  margin: 10px 0;
}

.button {
  padding: 12px 24px;
  font-size: 16px;
  font-weight: 500;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  margin-top: 20px;
}

.button:hover {
  background-color: #0056b3;
}
</style>
