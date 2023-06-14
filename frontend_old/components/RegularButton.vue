<template>
  <button class="button" :class="{ 'button--active': active, 'button--disabled': disabled }">
    <span class="text">
      <slot>Medium Button</slot>
    </span>
  </button>
</template>

<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  active: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});
</script>

<style scoped>
.button {
  border: none;
  border-radius: 8px;
  width: auto;
  height: auto;
  font-size: var(--font-quaternary-size);
  color: var(--font-primary-color);
  padding: 10px 20px;
  line-height: 20px;
  font-family: 'Inter-SemiBold', sans-serif;
  background: var(--background-secondary-color);
  cursor: pointer;
  transition: all 250ms ease;
}

.button--active {
  color: white;
  background: var(--green-color);
}

.button--disabled {
  opacity: 0.5;
}
</style>
