<script setup lang="ts">
import { computed, useSlots } from 'vue';
import Icon from '../icons/Icon.vue';

const slots = useSlots();

interface Props {
  label: string;
  value: string;
  modelValue: string | number;
  color: 'standart' | 'red' | 'blue' | 'green' | 'yellow';
}

defineEmits(['update:modelValue']);
const props = withDefaults(defineProps<Props>(), {
  color: 'standart',
});
const isChecked = computed(() => props.modelValue === props.value);
</script>

<template>
  <label>
    <template v-if="slots.content">
      <span @click="() => $el.querySelector('.input').click()">
        <slot name="content"></slot>
      </span>
    </template>
    <input
      class="input"
      type="radio"
      :checked="isChecked"
      :value="value"
      @change="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-show="!slots.content" class="radio">
      <template v-if="isChecked">
        <Icon v-if="color === 'standart'" icon="radio-checked" />
        <Icon v-if="color === 'blue'" icon="radio-blue-checked" />
        <Icon v-if="color === 'green'" icon="radio-green-checked" />
        <Icon v-if="color === 'red'" icon="radio-red-checked" />
        <Icon v-if="color === 'yellow'" icon="radio-yellow-checked" />
      </template>
      <Icon v-else icon="radio" />
    </span>
  </label>
</template>

<style scoped>
.input {
  position: absolute;
  -webkit-appearance: none;
  appearance: none;
}

.radio {
  border-radius: 9999px;
  cursor: pointer;
  display: inline-flex;
  width: 24px;
  height: 24px;
  justify-content: center;
  align-items: center;
}

.input:focus-visible + .radio {
  outline: 2px solid var(--focus-color);
}
</style>
