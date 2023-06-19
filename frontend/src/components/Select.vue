<script setup>
/**
 *
 * This is a wrapper component for v-select
 * https://github.com/sagalbot/vue-select
 *
 */

import { defineProps, defineEmits, ref, watch } from 'vue';

const emit = defineEmits(['update:modelValue']);
defineProps({
  placeholder: {
    type: String,
    default: 'Please select one',
  },
  options: {
    type: Array,
  },
  modelValue: {
    type: [String, Number],
    default: null,
  },
});

const selected = ref('');
watch(selected, (value) => emit('update:modelValue', value));
</script>

<template>
  <div class="form">
    <v-select v-model="selected" :options="options" class="select"></v-select>
  </div>
</template>

<style scoped></style>
