<script setup lang="ts">
import { computed, ref } from 'vue';
import Icon from '../icons/Icon.vue';

import { Icons } from '@/types';

interface Props {
  color: 'red' | 'blue' | 'green' | 'yellow';
  active: boolean;
  disabled: boolean;
}

const props = withDefaults(defineProps<Props>(), { color: 'red', disabled: false, active: false });
const hover = ref(false);

const icons = {
  red: 'triangle',
  blue: 'circle',
  green: 'square',
  yellow: 'diamond',
};

const currentStyle = computed(() => `button--${props.color}`);
const defaultIcon = computed((): Icons => (icons[props.color] + `-${props.color}`) as Icons);
const hoveredIcon = computed((): Icons => (icons[props.color] + `-white`) as Icons);
</script>

<template>
  <button
    class="button"
    :class="[currentStyle]"
    :disabled="disabled"
    :active="active"
    @mouseover="hover = true"
    @mouseleave="hover = false"
  >
    <Icon v-show="!active" :icon="defaultIcon" class="icon" />
    <Icon v-show="active" :icon="hoveredIcon" class="icon" />
    <slot>Primary Button</slot>
  </button>
</template>

<style scoped>
.button {
  --color: var(--color-heading);
  --color-hover: var(--c-white);
  --color-disabled: var(--color-text);
  --background: var(--color-background-soft);
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
  padding: 20px;
  gap: 20px;
  font-weight: 600;
  font-size: 16px;
  border-radius: 20px;
  line-height: 24px;
  transition: all 250ms ease;
  cursor: pointer;
  border: none;
  outline: 3px solid transparent;
  outline-offset: -3px;
  text-align: left;
  overflow-wrap: anywhere;
}

.button--red {
  --main-color: var(--c-red);
}

.button--blue {
  --main-color: var(--c-blue);
}

.button--green {
  --main-color: var(--c-green);
}

.button--yellow {
  --main-color: var(--c-yellow);
}

.button {
  color: var(--color-heading);
  background-color: var(--color-background-soft);
  outline-color: transparent;
}

.button:focus-visible {
  outline-color: var(--c-black);
}

.button:hover {
  outline-color: var(--main-color);
}

.button[disabled] {
  color: var(--color-text);
  outline-color: transparent;
}

.button[active='true'] {
  color: var(--c-white);
  background-color: var(--main-color);
}

.icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
  align-self: self-start;
}
</style>
