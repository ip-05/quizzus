<template>
  <button class="button" :class="{ 'button--minimalistic': minimalistic }">
    <div class="img" :class="{ 'img--small': smallIcon }">
      <nuxt-img :src="src" />
    </div>
    <span class="text">
      <slot>Medium Button</slot>
    </span>
  </button>
</template>

<script setup>
import { defineProps } from 'vue';

defineProps({
  src: {
    type: String,
    default: 'svg/icon-settings.svg',
  },
  minimalistic: {
    type: Boolean,
    default: false,
  },
  smallIcon: {
    type: Boolean,
    default: false,
  },
});
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

.img--small {
  width: 16px;
  height: 16px;
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

.button--minimalistic {
  width: auto;
  padding: 0;
  display: inline-flex;
}

.button--minimalistic:hover {
  background: none;
}
</style>
