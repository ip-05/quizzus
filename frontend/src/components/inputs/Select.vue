<script setup lang="ts">
/**
 *
 * This is a wrapper component for v-select
 * https://github.com/sagalbot/vue-select
 *
 */

import { defineProps, defineEmits, ref, watch } from 'vue';

interface Props {
  placeholder?: string;
  options: (string | number)[];
  modelValue: string | number;
}

const emit = defineEmits(['update:modelValue']);
withDefaults(defineProps<Props>(), {
  placeholder: 'Please select one',
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
