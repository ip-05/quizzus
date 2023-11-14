<script setup lang="ts">
import Icon from '../icons/Icon.vue';

import { Icons } from '@/types';

interface Props {
  icon: Icons;
}

defineProps<Props>();
</script>

<template>
  <button class="button button--default">
    <Icon :icon="icon" class="icon" />
  </button>
</template>

<style scoped>
.button {
  width: auto;
  height: auto;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  border-radius: 9999px;
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
  width: 30px;
  height: 30px;
  object-fit: contain;
}
</style>
