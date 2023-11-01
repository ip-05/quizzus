<script setup lang="ts">
import { computed, ref } from 'vue';
import Icon from '../icons/Icon.vue';

interface Props {
  value: string | number | (string | number)[];
  modelValue: string | number | boolean | (string | number)[] | undefined;
  id?: string;
  disabled: boolean;
  checked: boolean;
  color: 'standart' | 'red' | 'blue' | 'green' | 'yellow';
}

const emit = defineEmits(['update:modelValue']);
const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  checked: false,
  id: 'default_id',
  color: 'standart',
});

const computedModel = computed({
  get(): string | number | boolean | (string | number)[] | undefined {
    return props.modelValue;
  },
  set(value: boolean | string | number | (string | number)[] | undefined) {
    emit('update:modelValue', value);
  },
});

const isChecked = ref(props.checked);
</script>

<template>
  <label>
    <input
      :id="id"
      v-model="computedModel"
      type="checkbox"
      :disabled="disabled"
      :value="value"
      :checked="isChecked"
      class="input"
      @change="isChecked = !isChecked"
    />
    <span class="checkbox">
      <template v-if="isChecked">
        <Icon v-if="color === 'standart'" icon="checkbox-checked" />
        <Icon v-if="color === 'blue'" icon="checkbox-blue-checked" />
        <Icon v-if="color === 'green'" icon="checkbox-green-checked" />
        <Icon v-if="color === 'red'" icon="checkbox-red-checked" />
        <Icon v-if="color === 'yellow'" icon="checkbox-yellow-checked" />
      </template>
      <Icon v-else icon="checkbox" />
    </span>
  </label>
</template>

<style scoped>
.input {
  position: absolute;
  -webkit-appearance: none;
  appearance: none;
}

.checkbox {
  cursor: pointer;
  display: inline-flex;
  width: 24px;
  height: 24px;
  justify-content: center;
  align-items: center;
  border-radius: 4px;
}

.input:focus-visible + .checkbox {
  outline: 2px solid var(--focus-color);
}
</style>
