<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import Icon from '../icons/Icon.vue';
import { Icons } from '@/types';

/**
 * Reference:
 * https://vuejs.org/guide/components/events.html#usage-with-v-model
 */

interface Props {
  placeholder?: string;
  icon?: Icons;
  modelValue: string | number;
}

defineEmits(['update:modelValue']);
withDefaults(defineProps<Props>(), {
  placeholder: 'Input',
});
</script>

<template>
  <div class="form">
    <div v-if="icon" class="icon">
      <Icon :icon="icon" />
    </div>
    <input
      class="input"
      :class="{ 'input--icon': icon }"
      type="text"
      :placeholder="placeholder"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    />
  </div>
</template>

<style scoped>
.form {
  position: relative;
  display: grid;
}

.icon {
  position: absolute;
  left: 15px;
  display: flex;
  align-self: center;
}

.icon * {
  width: 16px;
  height: 16px;
}

.input {
  padding: 10px 15px;
  font-size: 14px;
  font-weight: 600;
  background: var(--color-background);
  color: var(--color-heading);
  border: none;
  outline: 2px solid var(--color-border);
  outline-offset: -2px;
  border-radius: 8px;
  transition: all 250ms ease;
}

.input--icon {
  padding-left: 40px;
}

.input:focus {
  outline-color: var(--color-background-dark);
}

.input::placeholder {
  color: var(--color-text);
}
</style>
