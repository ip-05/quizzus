<script setup lang="ts">
import Icon from '../icons/Icon.vue';

import { Icons } from '@/types';

interface Props {
  icon?: Icons;
  iconRight?: Icons;
  disabled: boolean;
}

withDefaults(defineProps<Props>(), { style: 'default', disabled: false });
</script>

<template>
  <button class="button button--default" :disabled="disabled">
    <Icon v-if="icon" :icon="icon" class="icon" />
    <slot>Large Button</slot>
    <Icon v-if="iconRight" :icon="iconRight" class="icon" />
  </button>
</template>

<style scoped>
.button {
  width: auto;
  min-height: 51px;
  height: auto;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 13.5px 15px;
  gap: 10px;
  font-weight: 500;
  font-size: 22px;
  border-radius: 15px;
  line-height: 100%;
  transition: all 250ms ease;
  cursor: pointer;
  border: none;
  outline: 3px solid transparent;
  outline-offset: -3px;
}

.button--default {
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
  width: 24px;
  height: 24px;
  object-fit: contain;
}
</style>
