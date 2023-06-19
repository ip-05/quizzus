<script setup>
import { defineProps, defineEmits } from 'vue';

defineEmits(['update:modelValue']);
defineProps({
  placeholder: {
    type: String,
    default: 'Input',
  },
  modelValue: {
    type: [String, Number],
    default: null,
  },
});
</script>

<template>
  <div class="form">
    <textarea
      class="textarea"
      :placeholder="placeholder"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
    ></textarea>
  </div>
</template>

<style scoped>
.form {
  position: relative;
  display: grid;
}

.textarea {
  height: 150%;
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
  resize: none;
}

.textarea:focus {
  outline-color: var(--color-background-dark);
}

.textarea::placeholder {
  color: var(--color-text);
}
</style>
