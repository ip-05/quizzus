<script setup lang="ts">
import { computed } from 'vue';
import Icon from '../icons/Icon.vue';

import { Icons } from '@/types';

interface Props {
  style?: 'default' | 'bordered' | 'action';
  disabled: boolean;
  icon?: Icons;
}

const props = withDefaults(defineProps<Props>(), { style: 'default', disabled: false });

const currentStyle = computed(() => `button--${props.style}`);
</script>

<template>
  <button class="button" :class="currentStyle" :disabled="disabled">
    <Icon v-if="icon" :icon="icon" class="icon" />
    <slot>Primary Button</slot>
  </button>
</template>

<style scoped>
.button {
  width: auto;
  min-height: 36px;
  height: auto;
  display: inline-flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 10px 20px;
  gap: 10px;
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

  --background: var(--color-background-soft);
  --background-hover: var(--color-background-dark);
  --background-disabled: var(--color-background-superlight);

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

.button--action {
  --color: var(--c-white);
  --color-hover: var(--c-white);
  --color-disabled: var(--c-white);

  --background: var(--c-green);
  --background-hover: var(--c-green-dark);
  --background-disabled: var(--c-green-light);

  --border-color: transparent;
  --border-color-hover: transparent;
  --border-color-disabled: transparent;
}

.button {
  color: var(--color);
  background-color: var(--background);
  outline-color: var(--border-color);
}

.button:focus-visible {
  outline-color: var(--c-black);
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

.icon {
  width: 16px;
  height: 16px;
  object-fit: contain;
}
</style>
