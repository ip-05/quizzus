<script setup lang="ts">
interface Props {
  checked: boolean;
  size: 'small' | 'medium' | 'large';
}

withDefaults(defineProps<Props>(), { checked: false, size: 'medium' });
</script>

<template>
  <label class="container">
    <h3 v-if="size === 'large'" class="label"><slot>Switch</slot></h3>
    <h4 v-if="size === 'medium'" class="label"><slot>Switch</slot></h4>
    <h5 v-if="size === 'small'" class="label"><slot>Switch</slot></h5>
    <input
      v-bind="$attrs"
      class="input"
      type="checkbox"
      :checked="checked"
      @change="$emit('update:checked', ($event.target as HTMLInputElement).checked)"
    />
    <span class="switch" :class="`switch--${size}`"></span>
  </label>
</template>

<style scoped>
.container {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.input {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

.switch--small {
  --container-width: 24px;
  --conteiner-height: 16px;
  --switch-width: 12px;
  --switch-height: 12px;
  --switch-padding: 3px;
}

.switch--medium {
  --container-width: 36px;
  --conteiner-height: 24px;
  --switch-width: 18px;
  --switch-height: 18px;
  --switch-padding: 4px;
}

.switch--large {
  --container-width: 48px;
  --conteiner-height: 32px;
  --switch-width: 25px;
  --switch-height: 24px;
  --switch-padding: 4px;
}

.switch {
  --height: var(--conteiner-height);
  display: flex;
  position: relative;
  height: var(--height);
  flex-basis: var(--container-width);
  border-radius: var(--container-width);
  background-color: var(--color-background-dark);
  flex-shrink: 0;
  transition: background-color 0.25s ease-in-out;
}

.switch::before {
  content: '';
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: var(--switch-padding);
  width: var(--switch-width);
  height: var(--switch-height);
  border-radius: 9999px;
  background-color: white;
  transition: all 0.25s ease-in-out;
}

.input:checked + .switch {
  background-color: var(--color-background-active);
}

.input:checked + .switch::before {
  left: calc(var(--container-width) - var(--switch-width) - var(--switch-padding));
}

/* .input:focus + .switch {
  outline: 2px black solid;
} */

/* .input:focus:checked + .switch {
  outline: 2px black solid;
} */
</style>
