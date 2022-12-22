<template>
  <button class="button">
    <div class="img">
      <nuxt-img :src="img" :alt="alt" />
    </div>
    <span class="text">
      <slot>Medium Button</slot>
    </span>
  </button>
</template>

<script setup>
import { computed, defineProps } from 'vue';

const props = defineProps({
  mode: {
    type: String,
    default: 'default',
  },
});

const imgs = {
  logout: { src: 'svg/icon-logout.svg', alt: 'log out' },
  settings: { src: 'svg/icon-settings.svg', alt: 'Settings' },
  google: { src: 'svg/icon-google-small.svg', alt: 'Sign In' },
};
const img = computed(() => imgs[props.mode].src);
const alt = computed(() => imgs[props.mode].alt);
</script>

<style scoped>
.button {
  width: 100%;
  height: 100%;
  display: flex;
  gap: 5px;
  align-items: center;
  padding: 8px;
  border: none;
  border-radius: 10px;
  background: var(--background-main-color);
  transition: all 250ms ease;
  cursor: pointer;
}

.button:hover {
  background: var(--background-secondary-color);
}

.img {
  width: 24px;
  height: 24px;
}

.img img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.text {
  color: var(--font-primary-color);
  font-family: 'Inter-SemiBold', sans-serif;
  font-size: var(--font-quaternary-size);
  line-height: 20px;
}
</style>
