<script setup>
import { computed } from 'vue';
import Icon from './Icon.vue';

const props = defineProps({
  state: {
    type: String,
    default: 'default',
  },
  icon: {
    type: Boolean,
    default: false,
  },
});

const states = {
  rainbow: 'button--rainbow',
};

const currentState = computed(() => states[props.state]);
</script>

<template>
  <button class="button" :class="currentState">
    <Icon v-if="icon" icon="plus" class="icon" />
    <slot>Button Medium</slot>
  </button>
</template>

<style scoped>
.button {
  width: auto;
  width: auto;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  padding: 5px 10px;
  gap: 5px;
  font-weight: 600;
  font-size: 14px;
  border-radius: 8px;
  color: var(--color-heading);
  background-color: var(--color-background);
  border: none;
  transition: all 250ms ease;
  cursor: pointer;
}

.button--rainbow {
  background: linear-gradient(
    114.82deg,
    #3498db 2.44%,
    #2ecc71 32.23%,
    #f1c40f 80.73%,
    #e74c3c 118.14%
  );
  background-size: 500% 100%;
  background-repeat: repeat;
  animation: gradientText 2s linear infinite;
  animation-direction: alternate;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  -moz-background-clip: text;
  -moz-text-fill-color: transparent;
}

@keyframes gradientText {
  0% {
    background-position: 0%;
  }
  100% {
    background-position: 100%;
  }
}

.icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.button:hover {
  background-color: var(--color-background-soft);
}
</style>
