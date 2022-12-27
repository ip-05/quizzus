<template>
  <form class="form" :class="{ 'form--minimalistic': minimalistic }">
    <div class="icon">
      <nuxt-img :src="img" alt="Icon" />
    </div>
    <input
      class="input"
      type="text"
      :placeholder="placeholder"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
    />
  </form>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

/**
 * Reference:
 * https://vuejs.org/guide/components/events.html#usage-with-v-model
 */

defineEmits(['update:modelValue']);
defineProps({
  placeholder: {
    type: String,
    default: 'Input',
  },
  img: {
    type: String,
    default: 'svg/icon-pen.svg',
  },
  minimalistic: {
    type: Boolean,
    default: false,
  },
  modelValue: {
    type: String,
    default: null,
  },
});
</script>

<style scoped>
.form {
  display: grid;
  align-items: center;
  grid-template-columns: 16px auto;
  gap: 10px;
  padding: 10px 12px;
  background-color: var(--background-secondary-color);
  border-radius: 10px;
  transition: all 250ms ease;
  outline: 3px solid transparent;
  outline-offset: -3px;
}

.form:focus-within {
  outline-color: var(--green-color);
}

.icon {
  width: 16px;
  height: 16px;
}

.icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.input {
  width: 100%;
  height: auto;
  line-height: 20px;
  border: none;
  background: transparent;
  font-size: var(--font-quaternary-size);
  color: var(--font-primary-color);
  font-family: 'Inter-SemiBold', sans-serif;
  outline: none;
}

.input::placeholder {
  color: var(--font-secondary-color);
}

.form--minimalistic {
  background-color: transparent;
  outline: none;
  border-radius: none;
  padding: 0;
  gap: 5px;
}

.form--minimalistic .input::placeholder {
  color: var(--font-primary-color);
}
</style>
