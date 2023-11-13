<script setup lang="ts">
import { computed } from 'vue';
import Icon from '../icons/Icon.vue';

import { Icons } from '@/types';

interface Props {
  style?: 'default' | 'bordered' | 'action' | 'active';
  disabled: boolean;
  icon?: Icons;
}

const props = withDefaults(defineProps<Props>(), { style: 'default', disabled: false });

const currentStyle = computed(() => `button--${props.style}`);
</script>

<template>
  <button class="button" :class="currentStyle" :disabled="disabled">
    <Icon v-if="icon" :icon="icon" class="icon" />
    <slot>Secondary Button</slot>
  </button>
</template>

<style scoped>
.button {
  width: auto;
  width: auto;
  display: inline-flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 5px 10px;
  gap: 5px;
  font-weight: 600;
  font-size: 14px;
  border-radius: 8px;
  line-height: 100%;
  transition: all 250ms ease;
  cursor: pointer;
  border: none;
  outline: 2px solid transparent;
  outline-offset: -2px;
}

.button--default {
  --color: var(--color-heading);
  --color-hover: var(--color-heading);
  --color-disabled: var(--color-text);

  --background: var(--color-background);
  --background-hover: var(--color-background-soft);
  --background-disabled: var(--color-background);

  --border-color: transparent;
  --border-color-hover: transparent;
  --border-color-disabled: transparent;
}

.button--active {
  --color: var(--color-heading);
  --color-hover: var(--color-heading);
  --color-disabled: var(--color-text);

  --background: var(--color-background-soft);
  --background-hover: var(--color-background-soft);
  --background-disabled: var(--color-background-soft);

  --border-color: transparent;
  --border-color-hover: transparent;
  --border-color-disabled: transparent;
}

.button--bordered {
  --color: var(--color-heading);
  --color-hover: var(--color-heading);
  --color-disabled: var(--color-text);

  --background: var(--color-background);
  --background-hover: var(--color-background);
  --background-disabled: var(--color-background);

  --border-color: var(--color-background-soft);
  --border-color-hover: var(--color-background-dark);
  --border-color-disabled: var(--color-background-soft);
}

.button {
  color: var(--color);
  background-color: var(--background);
  outline-color: var(--border-color);
}

.button:hover {
  color: var(--color-hover);
  background-color: var(--background-hover);
  outline-color: var(--border-color-hover);
}

.button[disabled] {
  color: var(--color-disabled);
  background-color: var(--background-disabled);
  outline-color: var(--border-color-disabled);
}

/* Exception */
.button--action {
  background: linear-gradient(
    114.82deg,
    #3498db 2.44%,
    #2ecc71 32.23%,
    #f1c40f 80.73%,
    #e74c3c 118.14%
  );
  background-size: 500% 100%;
  background-repeat: repeat;
  animation: gradientText 2s linear infinite;
  animation-direction: alternate;
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  -moz-background-clip: text;
  -moz-text-fill-color: transparent;
  outline-color: transparent;
}

.button--action:hover {
  color: var(--c-white);
  animation-duration: 4s;
  background-clip: border-box;
  -webkit-background-clip: border-box;
  -webkit-text-fill-color: var(--c-white);
  outline-color: transparent;
}

@keyframes gradientText {
  0% {
    background-position: 0%;
  }
  100% {
    background-position: 100%;
  }
}

.icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}
</style>
